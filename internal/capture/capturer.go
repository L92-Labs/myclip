package capture

import (
	"sort"
	"sync"
	"time"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

const (
	// DefaultInterval is Maccy's clipboardCheckInterval default.
	DefaultInterval = 500 * time.Millisecond
	// MinInterval is the D7 floor.
	MinInterval = 100 * time.Millisecond
	// DefaultMaxItemBytes is the D6 per-item cap.
	DefaultMaxItemBytes int64 = 32 << 20
)

// ClampInterval applies the D7 floor; zero means the default.
func ClampInterval(d time.Duration) time.Duration {
	if d == 0 {
		return DefaultInterval
	}
	return max(d, MinInterval)
}

// Outcome says what one Check did.
type Outcome int

const (
	Unchanged Outcome = iota // changeCount as last seen
	Ignored                  // events paused, or app/type rules matched: zero data reads
	Raced                    // the pasteboard changed mid-check: nothing kept, retried next tick
	Empty                    // every item was filtered out
	Oversized                // over MaxItemBytes (D6): skipped and recorded
	Captured
)

func (o Outcome) String() string {
	return [...]string{"unchanged", "ignored", "raced", "empty", "oversized", "captured"}[o]
}

// Result is a captured pasteboard, ready for history.
type Result struct {
	Contents    []history.Content
	Application string
	ChangeCount int
	OwnWrite    bool // carried TypeOwnWrite: a MyClip copy-back
	Size        int64
}

// Capturer runs the N6 capture sequence against a pasteboard. Check must not
// run concurrently with itself; Prefs methods are safe from any goroutine.
type Capturer struct {
	Board        pasteboard.Reader
	Frontmost    func() string // bundle ID of the source app; nil ⇒ ""
	MaxItemBytes int64         // 0 ⇒ DefaultMaxItemBytes

	mu    sync.Mutex
	prefs Prefs
	last  int
}

// NewCapturer starts at the pasteboard's current changeCount, so contents
// present before launch are not captured (as Maccy).
func NewCapturer(board pasteboard.Reader, prefs Prefs) *Capturer {
	return &Capturer{Board: board, prefs: prefs, last: board.ChangeCount()}
}

// Prefs returns a copy of the current prefs.
func (c *Capturer) Prefs() Prefs {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.prefs
}

// UpdatePrefs changes prefs under the lock.
func (c *Capturer) UpdatePrefs(fn func(*Prefs)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fn(&c.prefs)
}

// Check runs one capture step. Order (N6 sandwich): changeCount c0 → early
// exit on ignoreEvents → types-only snapshot → app + type ignore chain →
// re-check c1 == c0 before any data read → read → re-check c2 == c1, else
// discard the bytes. A raced check leaves the change unseen so the next
// tick reprocesses it.
func (c *Capturer) Check() (*Result, Outcome) {
	c0 := c.Board.ChangeCount()
	c.mu.Lock()
	if c0 == c.last {
		c.mu.Unlock()
		return nil, Unchanged
	}
	if c.prefs.IgnoreEvents {
		// ignore-next skips exactly one change: clear both flags.
		if c.prefs.IgnoreOnlyNextEvent {
			c.prefs.IgnoreEvents, c.prefs.IgnoreOnlyNextEvent = false, false
		}
		c.last = c0
		c.mu.Unlock()
		return nil, Ignored
	}
	prefs := c.prefs
	c.mu.Unlock()

	snap, err := c.Board.Snapshot()
	if err != nil || snap.ChangeCount != c0 {
		return nil, Raced
	}
	app := ""
	if c.Frontmost != nil {
		app = c.Frontmost()
	}
	if ShouldIgnoreApp(app, prefs) || ShouldIgnoreTypes(snap.Types, prefs) {
		c.seen(c0)
		return nil, Ignored
	}
	if c.Board.ChangeCount() != c0 {
		return nil, Raced // c1 path: no data read happened
	}
	items := c.Board.ReadItems()
	if c.Board.ChangeCount() != c0 {
		return nil, Raced // c2 path: bytes are dropped here, never stored
	}
	c.seen(c0)

	var contents []history.Content
	for _, it := range items {
		types := StripDynamicAndOLE(DropWordLinkSet(it.Types))
		if ItemMatchesRegexp(string(types[pasteboard.TypeString]), prefs.IgnoreRegexp) {
			continue
		}
		names := make([]string, 0, len(types))
		for t := range types {
			names = append(names, t)
		}
		sort.Strings(names)
		for _, t := range names {
			contents = append(contents, history.Content{Type: t, Value: types[t]})
		}
	}
	if len(contents) == 0 {
		return nil, Empty
	}
	r := &Result{
		Contents:    contents,
		Application: app,
		ChangeCount: c0,
		OwnWrite:    history.HasOwnWriteMarker(contents),
	}
	for _, ct := range contents {
		r.Size += int64(len(ct.Value))
	}
	limit := c.MaxItemBytes
	if limit == 0 {
		limit = DefaultMaxItemBytes
	}
	if r.Size > limit {
		return &Result{Application: app, ChangeCount: c0, Size: r.Size}, Oversized
	}
	return r, Captured
}

func (c *Capturer) seen(n int) {
	c.mu.Lock()
	c.last = n
	c.mu.Unlock()
}
