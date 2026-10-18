package main

import (
	"bytes"
	"image"
	"image/png"
	"sync"
	"testing"
	"time"

	"github.com/L92-Labs/myclip/internal/engine"
	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

func pngItem(t *testing.T, id int64, w int) *history.Item {
	var b bytes.Buffer
	if err := png.Encode(&b, image.NewGray(image.Rect(0, 0, w, 4))); err != nil {
		t.Fatal(err)
	}
	it := history.NewItem([]history.Content{{Type: pasteboard.TypePNG, Value: b.Bytes()}}, "", time.Now())
	it.ID = id
	return it
}

// A queued image whose ID now holds another clip is skipped unread, and
// a reading lands only on the image read.
func TestOCRSkipsReplacedImages(t *testing.T) {
	old, reused := pngItem(t, 7, 10), pngItem(t, 7, 20)
	oldData, _, _ := old.Image()
	var mu sync.Mutex
	var saved []string
	readOld := false
	calls := 0 // the replaced job must not reach the recognizer at all
	stored := old
	q := ocrQueue{
		item: func(int64) *history.Item { mu.Lock(); defer mu.Unlock(); return stored },
		recognize: func(data []byte) (string, error) {
			mu.Lock()
			calls++
			mu.Unlock()
			if bytes.Equal(data, oldData) {
				mu.Lock()
				readOld = true
				mu.Unlock()
				return "old words", nil
			}
			return "new words", nil
		},
		save: func(id int64, image [32]byte, text string) error {
			mu.Lock()
			defer mu.Unlock()
			if engine.ImageDigest(stored) != image {
				return engine.ErrNoItem
			}
			saved = append(saved, text)
			return nil
		},
	}
	mu.Lock()
	q.pending = []ocrJob{{7, engine.ImageDigest(old)}} // queued, not started yet
	stored = reused                                    // then the old image went, id 7 reused
	mu.Unlock()
	q.add(reused)
	deadline := time.Now().Add(2 * time.Second)
	for {
		mu.Lock()
		got, sawOld := append([]string(nil), saved...), readOld
		mu.Unlock()
		if sawOld || len(got) > 0 && got[0] != "new words" {
			t.Fatalf("the replaced image was read or saved: %v", got)
		}
		if len(got) == 1 {
			time.Sleep(50 * time.Millisecond) // let a stray job show itself
			mu.Lock()
			defer mu.Unlock()
			if calls != 1 {
				t.Fatalf("recognizer ran %d times, want 1", calls)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("saved %v", got)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestNeedsOCR(t *testing.T) {
	it := pngItem(t, 1, 4)
	if !needsOCR(it) {
		t.Fatal("unread image skipped")
	}
	it.Contents = append(it.Contents, history.Content{Type: pasteboard.TypeRecognizedText})
	if needsOCR(it) {
		t.Fatal("an image read as empty is read again")
	}
}

// A backlog longer than any buffer is read in full, none dropped.
func TestOCRBacklogIsNotDropped(t *testing.T) {
	const n = 600
	var mu sync.Mutex
	saved := 0
	all := make(map[int64]*history.Item, n)
	for i := range n {
		all[int64(i+1)] = pngItem(t, int64(i+1), 4+i%3)
	}
	q := ocrQueue{
		item:      func(id int64) *history.Item { return all[id] },
		recognize: func([]byte) (string, error) { return "w", nil },
		save: func(int64, [32]byte, string) error {
			mu.Lock()
			saved++
			mu.Unlock()
			return nil
		},
	}
	for i := range n {
		q.add(all[int64(i+1)])
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		mu.Lock()
		got := saved
		mu.Unlock()
		if got == n {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("read %d of %d", got, n)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// A copy made now is read before the startup backlog.
func TestLiveOCRJumpsTheBacklog(t *testing.T) {
	block := make(chan struct{})
	var mu sync.Mutex
	var order []int64
	all := map[int64]*history.Item{}
	for id := int64(1); id <= 4; id++ {
		all[id] = pngItem(t, id, int(id)+3)
	}
	q := ocrQueue{
		item:      func(id int64) *history.Item { return all[id] },
		recognize: func([]byte) (string, error) { <-block; return "w", nil },
		save: func(id int64, _ [32]byte, _ string) error {
			mu.Lock()
			order = append(order, id)
			mu.Unlock()
			return nil
		},
	}
	q.add(all[1]) // the worker takes 1 and waits
	time.Sleep(20 * time.Millisecond)
	q.add(all[2])
	q.add(all[3])
	q.addLive(all[4])
	close(block)
	deadline := time.Now().Add(2 * time.Second)
	for {
		mu.Lock()
		got := append([]int64(nil), order...)
		mu.Unlock()
		if len(got) == 4 {
			if got[1] != 4 {
				t.Fatalf("order %v: the live copy waited", got)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("order %v", got)
		}
		time.Sleep(5 * time.Millisecond)
	}
}
