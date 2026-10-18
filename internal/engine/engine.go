// Package engine runs capture → history merge → store, and copy-back.
// Pasteboard calls go through RunOnMain (AppKit's thread); history and store
// work stays on the caller's goroutine.
package engine

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
	"github.com/L92-Labs/myclip/internal/store"
)

// Limits are the D6 caps plus Maccy's history length.
type Limits struct {
	MaxItems int   // unpinned items kept (Maccy `size`, default 200)
	MaxBytes int64 // total stored content
}

// DefaultLimits are the approved D6 defaults.
var DefaultLimits = Limits{MaxItems: 200, MaxBytes: 512 << 20}

// Config wires an Engine.
type Config struct {
	Board     pasteboard.Board
	Frontmost func() string
	Store     *store.Store
	Prefs     capture.Prefs
	Limits    Limits
	// RunOnMain runs fn on the UI thread and waits; nil runs it inline.
	RunOnMain func(fn func())
	// OnChange is called after the history list changes.
	OnChange func()
	// OnNew is called with each new unique copy (not merges), for
	// Maccy's copy notification.
	OnNew func(*history.Item)
	// OnCopy is called for each copy captured from another app, new or
	// merged into an existing item (not MyClip's own copy-backs).
	OnCopy func(*history.Item)
	// OnError reports store failures; capture continues in memory.
	OnError func(error)
	Now     func() time.Time
}

// Engine is safe from any goroutine.
type Engine struct {
	cfg Config
	cap *capture.Capturer

	mu        sync.Mutex
	hist      *history.History
	oversized int

	// persist orders the store writes as the changes they follow: held from
	// a change through its write, so a slower writer cannot save an older
	// state over a newer one. Callbacks run after it is released (they wait
	// on the main thread, which may be waiting for persist).
	persist sync.Mutex
}

// New loads stored history and starts capture at the current changeCount.
func New(ctx context.Context, cfg Config) (*Engine, error) {
	if cfg.Board == nil || cfg.Store == nil {
		return nil, errors.New("engine: Board and Store are required")
	}
	if cfg.RunOnMain == nil {
		cfg.RunOnMain = func(fn func()) { fn() }
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if cfg.Limits == (Limits{}) {
		cfg.Limits = DefaultLimits
	}
	items, err := cfg.Store.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("engine: load history: %w", err)
	}
	e := &Engine{cfg: cfg, hist: history.New()}
	e.hist.Load(items)
	cfg.RunOnMain(func() { e.cap = capture.NewCapturer(cfg.Board, cfg.Prefs) })
	e.cap.Frontmost = cfg.Frontmost
	return e, nil
}

// Run checks the pasteboard every interval (D7 floor applied) until ctx ends.
func (e *Engine) Run(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(capture.ClampInterval(interval))
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			e.Tick(ctx)
		}
	}
}

// Tick runs one capture check and records its result.
func (e *Engine) Tick(ctx context.Context) capture.Outcome {
	var r *capture.Result
	var o capture.Outcome
	e.cfg.RunOnMain(func() { r, o = e.cap.Check() })
	switch o {
	case capture.Captured:
		e.record(ctx, r)
	case capture.Oversized:
		e.mu.Lock()
		e.oversized++
		e.mu.Unlock()
	}
	return o
}

func (e *Engine) record(ctx context.Context, r *capture.Result) {
	now := e.cfg.Now()
	e.persist.Lock()
	e.mu.Lock()
	it := history.NewItem(r.Contents, r.Application, now)
	res := e.hist.Add(it, r.ChangeCount, r.OwnWrite, now)
	evicted := e.hist.Trim(e.cfg.Limits.MaxItems, e.cfg.Limits.MaxBytes)
	e.mu.Unlock()

	err := e.cfg.Store.Save(ctx, res.Item)
	var ids []int64
	for _, x := range evicted {
		if x.ID != 0 {
			ids = append(ids, x.ID)
		}
	}
	err = errors.Join(err, e.cfg.Store.Delete(ctx, ids...))
	e.persist.Unlock()
	e.report(err)
	if !res.Merged && e.cfg.OnNew != nil {
		e.cfg.OnNew(res.Item)
	}
	if !r.OwnWrite && e.cfg.OnCopy != nil {
		e.cfg.OnCopy(res.Item) // every copy made in another app, merged or new
	}
	e.changed()
}

func (e *Engine) report(err error) {
	if err != nil && e.cfg.OnError != nil {
		e.cfg.OnError(err)
	}
}

func (e *Engine) changed() {
	if e.cfg.OnChange != nil {
		e.cfg.OnChange()
	}
}

// Items returns the history, newest first. Items are shared: read only.
func (e *Engine) Items() []*history.Item {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.hist.Items()
}

// Oversized counts clips skipped by the D6 per-item cap this session.
func (e *Engine) Oversized() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.oversized
}

// Paused reports whether capture is paused (Maccy ignoreEvents).
func (e *Engine) Paused() bool { return e.cap.Prefs().IgnoreEvents }

// SetPaused pauses or resumes capture. Resuming clears ignore-next.
func (e *Engine) SetPaused(paused bool) {
	e.cap.UpdatePrefs(func(p *capture.Prefs) {
		p.IgnoreEvents = paused
		if !paused {
			p.IgnoreOnlyNextEvent = false
		}
	})
	e.changed()
}

// IgnoreNext skips exactly the next pasteboard change.
func (e *Engine) IgnoreNext() {
	e.cap.UpdatePrefs(func(p *capture.Prefs) { p.IgnoreEvents, p.IgnoreOnlyNextEvent = true, true })
	e.changed()
}

// IgnoringNext reports whether only the next change will be skipped.
func (e *Engine) IgnoringNext() bool {
	p := e.cap.Prefs()
	return p.IgnoreEvents && p.IgnoreOnlyNextEvent
}

// SetCaptureSafe keeps capture off until cleared (settings unreadable).
func (e *Engine) SetCaptureSafe(on bool) { e.cap.SetCaptureSafe(on) }
