package main

import (
	"errors"
	"log"
	"sync"

	"github.com/L92-Labs/myclip/internal/engine"
	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/native"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// ocrMaxBytes bounds the images read for text: larger ones stay unread.
const ocrMaxBytes = 16 << 20

// ocrQueue reads the text in image clips, one at a time off the main
// thread, so search finds a screenshot by its words (Maccy's Vision titles).
// Nothing leaves the computer: Vision on macOS, tesseract on Linux.
type ocrQueue struct {
	mu      sync.Mutex
	pending []*history.Item // waiting, oldest first; never dropped
	running bool
	// save stores a reading for the image with that digest, refusing it
	// atomically when the ID now holds another image.
	save func(id int64, image [32]byte, text string) error
	// recognize reads an image's text; nil uses the system's recognizer.
	recognize func([]byte) (string, error)
	stopped   bool // no recognizer on this system
}

func (q *ocrQueue) add(it *history.Item) {
	if !needsOCR(it) {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.stopped {
		return
	}
	q.pending = append(q.pending, it)
	if !q.running {
		q.running = true
		go q.run()
	}
}

// next takes the next job, or ends the worker when none is left.
func (q *ocrQueue) next() (*history.Item, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.pending) == 0 || q.stopped {
		q.running = false
		return nil, false
	}
	it := q.pending[0]
	q.pending = q.pending[1:]
	return it, true
}

// needsOCR reports whether it is an image not read yet, small enough.
func needsOCR(it *history.Item) bool {
	data, _, ok := it.Image()
	return ok && len(data) <= ocrMaxBytes && len(it.FilePaths()) == 0 && !it.TextRecognized() && it.Text() == ""
}

func (q *ocrQueue) run() {
	for {
		it, ok := q.next()
		if !ok {
			return
		}
		data, typ, _ := it.Image()
		// Go cannot decode HEIC to check it, but Vision reads it.
		if typ != pasteboard.TypeHEIC && !decodable(data) {
			continue
		}
		sum := engine.ImageDigest(it)
		recognize := q.recognize
		if recognize == nil {
			recognize = native.RecognizeText
		}
		text, err := recognize(data)
		if err != nil {
			log.Print("myclip: text in image: ", err)
			if errors.Is(err, native.ErrNoOCR) {
				q.mu.Lock()
				q.stopped, q.running, q.pending = true, false, nil // no recognizer here
				q.mu.Unlock()
				return
			}
			continue
		}
		// The item may be gone meanwhile, its ID taken by a newer clip: save
		// checks the image under its lock. An empty reading is saved too, so
		// the image is not read again at the next launch.
		if err := q.save(it.ID, sum, text); errors.Is(err, engine.ErrNoItem) {
			continue
		} else if err != nil {
			log.Print("myclip: text in image: ", err)
		}
	}
}
