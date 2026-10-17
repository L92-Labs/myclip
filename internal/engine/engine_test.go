package engine

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/pasteboard"
	"github.com/L92-Labs/myclip/internal/store"
)

func text(s string) []pasteboard.Item {
	return []pasteboard.Item{{Types: map[string][]byte{pasteboard.TypeString: []byte(s)}}}
}

type rig struct {
	e     *Engine
	f     *pasteboard.Fake
	s     *store.Store
	path  string
	clock time.Time
}

func newRig(t *testing.T, lim Limits) *rig {
	t.Helper()
	path := filepath.Join(t.TempDir(), "history.db")
	s, err := store.Open(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	r := &rig{f: pasteboard.NewFake(), s: s, path: path, clock: time.Unix(1_700_000_000, 0)}
	r.e, err = New(context.Background(), Config{
		Board: r.f, Store: s, Limits: lim,
		Now: func() time.Time { r.clock = r.clock.Add(time.Second); return r.clock },
	})
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func (r *rig) copy(t *testing.T, s string) capture.Outcome {
	t.Helper()
	r.f.SetItems(text(s), nil)
	return r.e.Tick(context.Background())
}

func TestCaptureIsStored(t *testing.T) {
	r := newRig(t, Limits{})
	if o := r.copy(t, "hello"); o != capture.Captured {
		t.Fatalf("outcome %v", o)
	}
	got, err := r.s.Load(context.Background())
	if err != nil || len(got) != 1 || got[0].Text() != "hello" || got[0].ID == 0 {
		t.Fatalf("stored %+v, %v", got, err)
	}
}

// Copy-back bump: writing an item back merges into it (one row, count 2)
// and moves it to the top.
func TestCopyBackBumps(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "one")
	r.copy(t, "two")
	one := r.e.Items()[1]
	if err := r.e.CopyBack(one.ID); err != nil {
		t.Fatal(err)
	}
	if o := r.e.Tick(context.Background()); o != capture.Captured {
		t.Fatalf("copy-back tick %v", o)
	}
	items := r.e.Items()
	if len(items) != 2 || items[0].Text() != "one" || items[0].NumberOfCopies != 2 || items[0].ID != one.ID {
		t.Fatalf("after copy-back: %+v", items[0])
	}
	if _, ok := items[0].Value(pasteboard.TypeOwnWrite); ok {
		t.Fatal("own-write marker stored as content")
	}
	got, _ := r.s.Load(context.Background())
	if len(got) != 2 || got[0].ID != one.ID || got[0].NumberOfCopies != 2 {
		t.Fatalf("stored after copy-back: %+v", got)
	}
}

func TestTrimEvictsOldestUnpinnedFromStore(t *testing.T) {
	r := newRig(t, Limits{MaxItems: 3})
	for i := range 5 {
		r.copy(t, fmt.Sprint("clip ", i))
	}
	items := r.e.Items()
	if len(items) != 3 || items[0].Text() != "clip 4" || items[2].Text() != "clip 2" {
		t.Fatalf("memory: %d items, first %q", len(items), items[0].Text())
	}
	got, _ := r.s.Load(context.Background())
	if len(got) != 3 || got[2].Text() != "clip 2" {
		t.Fatalf("store: %d rows", len(got))
	}
}

func TestRestartReloads(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "persisted")
	r.s.Close()
	s, err := store.Open(context.Background(), r.path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	e, err := New(context.Background(), Config{Board: r.f, Store: s})
	if err != nil {
		t.Fatal(err)
	}
	if items := e.Items(); len(items) != 1 || items[0].Text() != "persisted" {
		t.Fatalf("reloaded %+v", items)
	}
	// The pasteboard still holds "persisted": no recapture after restart.
	if o := e.Tick(context.Background()); o != capture.Unchanged {
		t.Fatalf("tick after restart %v", o)
	}
}

func TestPauseStopsCapture(t *testing.T) {
	r := newRig(t, Limits{})
	r.e.SetPaused(true)
	if o := r.copy(t, "private"); o != capture.Ignored {
		t.Fatalf("paused outcome %v", o)
	}
	r.e.SetPaused(false)
	if o := r.copy(t, "public"); o != capture.Captured {
		t.Fatalf("resumed outcome %v", o)
	}
	if n := len(r.e.Items()); n != 1 {
		t.Fatalf("%d items", n)
	}
}
