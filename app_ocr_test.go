package main

import (
	"bytes"
	"image"
	"image/png"
	"sync"
	"testing"
	"time"

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

// A reading saves only onto the image read: not onto a newer clip that
// took a deleted image's ID meanwhile.
func TestOCRSavesOnlyOntoTheImageRead(t *testing.T) {
	old, reused := pngItem(t, 7, 10), pngItem(t, 7, 20)
	var mu sync.Mutex
	var saved []string // every save, in order
	read := make(chan struct{}, 2)
	oldData, _, _ := old.Image()
	q := ocrQueue{
		recognize: func(data []byte) (string, error) {
			read <- struct{}{}
			if bytes.Equal(data, oldData) {
				return "old words", nil
			}
			return "new words", nil
		},
		current: func(int64) *history.Item { return reused }, // id 7 is another image now
		save: func(id int64, text string) error {
			mu.Lock()
			defer mu.Unlock()
			saved = append(saved, text)
			return nil
		},
	}
	q.add(old)    // read, then dropped: id 7 holds another image
	q.add(reused) // read and saved
	<-read
	<-read
	deadline := time.Now().Add(2 * time.Second)
	for {
		mu.Lock()
		got := append([]string(nil), saved...)
		mu.Unlock()
		if len(got) == 1 && got[0] == "new words" {
			return
		}
		if len(got) > 0 && got[0] == "old words" {
			t.Fatalf("the deleted image's text went onto the new clip: %v", got)
		}
		if time.Now().After(deadline) {
			t.Fatalf("saved %v", saved)
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
		recognize: func([]byte) (string, error) { return "w", nil },
		current:   func(id int64) *history.Item { return all[id] },
		save: func(int64, string) error {
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
