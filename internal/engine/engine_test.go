package engine

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/history"
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
	if err := r.e.CopyBack(one.ID, false); err != nil {
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

func TestPinDeleteClear(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "a")
	r.copy(t, "b")
	r.copy(t, "c")
	items := r.e.Items()
	pin, err := r.e.TogglePin(items[1].ID, "b")
	if err != nil || pin != "c" {
		t.Fatalf("pin %q %v", pin, err)
	}
	if err := r.e.SetTitle(items[1].ID, "Bee"); err != nil {
		t.Fatal(err)
	}
	if items[1].Pin != "" || items[1].Title != "" {
		t.Fatal("handed-out item was mutated")
	}
	if err := r.e.Delete(items[0].ID); err != nil {
		t.Fatal(err)
	}
	if err := r.e.Clear(false); err != nil {
		t.Fatal(err)
	}
	got, _ := r.s.Load(context.Background())
	if len(got) != 1 || got[0].Pin != "c" || got[0].Title != "Bee" || got[0].Text() != "b" {
		t.Fatalf("after clear %+v", got)
	}
	if err := r.e.Clear(true); err != nil {
		t.Fatal(err)
	}
	if got, _ := r.s.Load(context.Background()); len(got) != 0 || len(r.e.Items()) != 0 {
		t.Fatalf("after clear all %d", len(got))
	}
}

func TestSetPinTaken(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "a")
	r.copy(t, "b")
	items := r.e.Items()
	r.e.TogglePin(items[0].ID)
	r.e.TogglePin(items[1].ID)
	if err := r.e.SetPin(items[1].ID, r.e.Item(items[0].ID).Pin); err == nil {
		t.Fatal("took a used key")
	}
}

func TestCopyBackPlainAndFiles(t *testing.T) {
	r := newRig(t, Limits{})
	r.f.SetItems([]pasteboard.Item{
		{Types: map[string][]byte{pasteboard.TypeString: []byte("x"), pasteboard.TypeRTF: []byte("{\\rtf1 x}")}},
	}, nil)
	r.e.Tick(context.Background())
	id := r.e.Items()[0].ID
	if err := r.e.CopyBack(id, true); err != nil {
		t.Fatal(err)
	}
	got := r.f.ReadItems()
	if len(got) != 1 || got[0].Types[pasteboard.TypeRTF] != nil || string(got[0].Types[pasteboard.TypeString]) != "x" {
		t.Fatalf("plain write %+v", got)
	}

	items := boardItems([]history.Content{
		{Type: pasteboard.TypeFileURL, Value: []byte("file:///a")},
		{Type: pasteboard.TypeFileURL, Value: []byte("file:///b")},
	}, false)
	if len(items) != 2 || string(items[1].Types[pasteboard.TypeFileURL]) != "file:///b" {
		t.Fatalf("files %+v", items)
	}
}

func TestSetLimitsTrims(t *testing.T) {
	r := newRig(t, Limits{})
	for _, s := range []string{"a", "b", "c"} {
		r.copy(t, s)
	}
	r.e.SetLimits(Limits{MaxItems: 1})
	if got, _ := r.s.Load(context.Background()); len(got) != 1 || got[0].Text() != "c" {
		t.Fatalf("stored %+v", got)
	}
}

// Text read in an image is stored and searchable, never pasted, and does
// not stop the same image from merging when copied again.
func TestRecognizedText(t *testing.T) {
	r := newRig(t, Limits{})
	img := []pasteboard.Item{{Types: map[string][]byte{pasteboard.TypePNG: []byte("png-bytes")}}}
	r.f.SetItems(img, nil)
	r.e.Tick(context.Background())
	id := r.e.Items()[0].ID
	if err := r.e.SetRecognizedText(id, ImageDigest(r.e.Item(id)), "Invoice 42"); err != nil {
		t.Fatal(err)
	}
	if got := r.e.Item(id).RecognizedText(); got != "Invoice 42" {
		t.Fatalf("text %q", got)
	}
	if err := r.e.CopyBack(id, false); err != nil {
		t.Fatal(err)
	}
	for _, it := range r.f.ReadItems() {
		if _, ok := it.Types[pasteboard.TypeRecognizedText]; ok {
			t.Fatal("recognized text written to the pasteboard")
		}
	}
	r.f.SetItems(img, nil) // the same image, copied again
	r.e.Tick(context.Background())
	if items := r.e.Items(); len(items) != 1 || items[0].RecognizedText() != "Invoice 42" {
		t.Fatalf("after recopy: %d items", len(items))
	}
	got, _ := r.s.Load(context.Background())
	if len(got) != 1 || got[0].RecognizedText() != "Invoice 42" {
		t.Fatalf("stored %+v", got)
	}
}

func TestCopyManyJoinsInOrderAndIsNotRecorded(t *testing.T) {
	r := newRig(t, Limits{})
	for _, s := range []string{"first", "second", "third"} {
		r.copy(t, s)
	}
	items := r.e.Items() // third second first
	if err := r.e.CopyMany([]int64{items[2].ID, items[0].ID}); err != nil {
		t.Fatal(err)
	}
	if got := string(r.f.ReadItems()[0].Types[pasteboard.TypeString]); got != "first\nthird" {
		t.Fatalf("joined %q", got)
	}
	if o := r.e.Tick(context.Background()); o != capture.Ignored {
		t.Fatalf("combined copy recorded: %v", o)
	}
	if n := len(r.e.Items()); n != 3 {
		t.Fatalf("%d items", n)
	}
	// Files stay files.
	r.f.SetItems([]pasteboard.Item{{Types: map[string][]byte{pasteboard.TypeFileURL: []byte("file:///a")}}}, nil)
	r.e.Tick(context.Background())
	r.f.SetItems([]pasteboard.Item{{Types: map[string][]byte{pasteboard.TypeFileURL: []byte("file:///b")}}}, nil)
	r.e.Tick(context.Background())
	files := r.e.Items()
	if err := r.e.CopyMany([]int64{files[1].ID, files[0].ID}); err != nil {
		t.Fatal(err)
	}
	got := r.f.ReadItems()
	if len(got) != 2 || string(got[0].Types[pasteboard.TypeFileURL]) != "file:///a" || string(got[1].Types[pasteboard.TypeFileURL]) != "file:///b" {
		t.Fatalf("files %+v", got)
	}
}

type failingBoard struct{ *pasteboard.Fake }

func (failingBoard) Write([]pasteboard.Item) error { return errors.New("rejected") }

// A stack write the pasteboard rejects must not skip the user's next copy.
func TestCopyManyFailureDoesNotArmIgnoreNext(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "a")
	r.copy(t, "b")
	ids := []int64{r.e.Items()[0].ID, r.e.Items()[1].ID}
	s := r.s
	e, err := New(context.Background(), Config{Board: failingBoard{r.f}, Store: s})
	if err != nil {
		t.Fatal(err)
	}
	if err := e.CopyMany(ids); err == nil {
		t.Fatal("write error lost")
	}
	if e.IgnoringNext() {
		t.Fatal("ignore-next armed after a failed write")
	}
}

// An empty reading marks the image read; recognized text counts toward the
// byte budget.
func TestRecognizedTextEmptyAndLimits(t *testing.T) {
	r := newRig(t, Limits{MaxBytes: 64})
	r.f.SetItems([]pasteboard.Item{{Types: map[string][]byte{pasteboard.TypePNG: []byte("old-png")}}}, nil)
	r.e.Tick(context.Background())
	r.f.SetItems([]pasteboard.Item{{Types: map[string][]byte{pasteboard.TypePNG: []byte("new-png")}}}, nil)
	r.e.Tick(context.Background())
	newID, oldID := r.e.Items()[0].ID, r.e.Items()[1].ID
	if err := r.e.SetRecognizedText(oldID, ImageDigest(r.e.Item(oldID)), ""); err != nil {
		t.Fatal(err)
	}
	if it := r.e.Item(oldID); !it.TextRecognized() || it.RecognizedText() != "" {
		t.Fatal("empty reading not recorded")
	}
	if err := r.e.SetRecognizedText(newID, ImageDigest(r.e.Item(newID)), strings.Repeat("x", 52)); err != nil {
		t.Fatal(err)
	}
	if r.e.Item(oldID) != nil {
		t.Fatal("over the byte budget, the oldest item stayed")
	}
	if got, _ := r.s.Load(context.Background()); len(got) != 1 {
		t.Fatalf("stored %d items", len(got))
	}
}

// A copy the user makes right after a stack write, before the next check,
// is captured: only the stack's own change is skipped.
func TestCopyRightAfterAStackIsKept(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "a")
	r.copy(t, "b")
	ids := []int64{r.e.Items()[0].ID, r.e.Items()[1].ID}
	if err := r.e.CopyMany(ids); err != nil {
		t.Fatal(err)
	}
	r.f.SetItems(text("typed right after"), nil) // before any tick
	if o := r.e.Tick(context.Background()); o != capture.Captured {
		t.Fatalf("user copy %v", o)
	}
	if got := r.e.Items()[0].Text(); got != "typed right after" {
		t.Fatalf("top %q", got)
	}
}

// A reading for another image (same ID reused) is refused under the lock.
func TestRecognizedTextNeedsTheSameImage(t *testing.T) {
	r := newRig(t, Limits{})
	r.f.SetItems([]pasteboard.Item{{Types: map[string][]byte{pasteboard.TypePNG: []byte("this image")}}}, nil)
	r.e.Tick(context.Background())
	id := r.e.Items()[0].ID
	other := ImageDigest(&history.Item{Contents: []history.Content{{Type: pasteboard.TypePNG, Value: []byte("deleted image")}}})
	if err := r.e.SetRecognizedText(id, other, "wrong words"); err == nil {
		t.Fatal("reading of another image accepted")
	}
	if r.e.Item(id).TextRecognized() {
		t.Fatal("stored onto the wrong image")
	}
}

// A renamed image adds nothing to a stack: titles are not clip content.
func TestCopyManySkipsTitles(t *testing.T) {
	r := newRig(t, Limits{})
	r.f.SetItems([]pasteboard.Item{{Types: map[string][]byte{pasteboard.TypePNG: []byte("img")}}}, nil)
	r.e.Tick(context.Background())
	img := r.e.Items()[0].ID
	if err := r.e.SetTitle(img, "My screenshot"); err != nil {
		t.Fatal(err)
	}
	r.copy(t, "words")
	if err := r.e.CopyMany([]int64{img, r.e.Items()[0].ID}); err != nil {
		t.Fatal(err)
	}
	if got := string(r.f.ReadItems()[0].Types[pasteboard.TypeString]); got != "words" {
		t.Fatalf("stack text %q", got)
	}
}

// Every copy from another app reports OnCopy, merged or new; MyClip's own
// copy-back does not.
func TestOnCopyForEveryOutsideCopy(t *testing.T) {
	r := newRig(t, Limits{})
	var mu sync.Mutex
	n := 0
	e, err := New(context.Background(), Config{Board: r.f, Store: r.s,
		OnCopy: func(*history.Item) { mu.Lock(); n++; mu.Unlock() }})
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"x", "y", "x"} { // the second x merges
		r.f.SetItems(text(s), nil)
		e.Tick(context.Background())
	}
	if err := e.CopyBack(e.Items()[1].ID, false); err != nil {
		t.Fatal(err)
	}
	e.Tick(context.Background())
	mu.Lock()
	defer mu.Unlock()
	if n != 3 {
		t.Fatalf("OnCopy %d times, want 3", n)
	}
}

// Concurrent changes reach the database in the order they happened: at
// rest, what is stored is what is in memory.
func TestStoreMatchesMemoryAfterConcurrentChanges(t *testing.T) {
	r := newRig(t, Limits{MaxItems: 20})
	for i := range 10 {
		r.copy(t, fmt.Sprint("seed ", i))
	}
	var wg sync.WaitGroup
	for g := range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range 30 {
				items := r.e.Items()
				it := items[(g+i)%len(items)]
				switch i % 3 {
				case 0:
					r.e.SetTitle(it.ID, fmt.Sprint("t", g, i))
				case 1:
					r.e.TogglePin(it.ID)
				default:
					r.e.SetText(it.ID, fmt.Sprint("text ", g, " ", i))
				}
			}
		}()
	}
	wg.Wait()
	mem := map[int64]string{}
	for _, it := range r.e.Items() {
		mem[it.ID] = fmt.Sprint(it.Title, "|", it.Pin, "|", it.Text())
	}
	stored, err := r.s.Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(stored) != len(mem) {
		t.Fatalf("stored %d, memory %d", len(stored), len(mem))
	}
	for _, it := range stored {
		if got := fmt.Sprint(it.Title, "|", it.Pin, "|", it.Text()); got != mem[it.ID] {
			t.Fatalf("item %d stored %q, memory %q", it.ID, got, mem[it.ID])
		}
	}
}

// A writer delayed between its change and its write cannot store an older
// state over a newer one.
func TestDelayedWriterDoesNotOverwriteNewer(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "clip")
	id := r.e.Items()[0].ID
	var calls atomic.Int32
	firstIn := make(chan struct{})
	beforeSave = func() {
		if calls.Add(1) == 1 { // only the first writer stalls (sync.Once would hold the second too)
			close(firstIn)
			time.Sleep(100 * time.Millisecond)
		}
	}
	defer func() { beforeSave = nil }()
	done := make(chan struct{})
	go func() { r.e.SetTitle(id, "older"); close(done) }()
	<-firstIn
	r.e.SetTitle(id, "newer") // changes after "older"
	<-done
	stored, _ := r.s.Load(context.Background())
	if got := r.e.Item(id).Title; got != "newer" || stored[0].Title != got {
		t.Fatalf("memory %q, stored %q", got, stored[0].Title)
	}
}

// The stack marker is one-shot: the same marker type replayed later by
// another app is an ordinary copy.
func TestReplayedStackMarkerIsCaptured(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "a")
	r.copy(t, "b")
	if err := r.e.CopyMany([]int64{r.e.Items()[0].ID, r.e.Items()[1].ID}); err != nil {
		t.Fatal(err)
	}
	written := r.f.ReadItems()
	if o := r.e.Tick(context.Background()); o != capture.Ignored {
		t.Fatalf("stack write %v", o)
	}
	r.copy(t, "something else")
	// Another clipboard manager restores the stack copy, marker and all.
	r.f.SetItems(written, nil)
	if o := r.e.Tick(context.Background()); o != capture.Captured {
		t.Fatalf("replayed marker %v, want captured", o)
	}
}
