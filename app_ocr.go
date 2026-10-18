package main

import (
	"errors"
	"log"
	"sync"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/native"
)

// ocrMaxBytes bounds the images read for text: larger ones stay unread.
const ocrMaxBytes = 16 << 20

// ocrQueue reads the text in image clips, one at a time off the main
// thread, so search finds a screenshot by its words (Maccy's Vision titles).
// Nothing leaves the computer: Vision on macOS, tesseract on Linux.
type ocrQueue struct {
	once sync.Once
	jobs chan *history.Item
	save func(id int64, text string) error
}

func (q *ocrQueue) add(it *history.Item) {
	if !needsOCR(it) {
		return
	}
	q.once.Do(func() {
		q.jobs = make(chan *history.Item, 256)
		go q.run()
	})
	select {
	case q.jobs <- it:
	default: // backlog full: this one stays unread
	}
}

// needsOCR reports whether it is an image not read yet, small enough.
func needsOCR(it *history.Item) bool {
	data, _, ok := it.Image()
	return ok && len(data) <= ocrMaxBytes && len(it.FilePaths()) == 0 && it.RecognizedText() == "" && it.Text() == ""
}

func (q *ocrQueue) run() {
	for it := range q.jobs {
		data, _, _ := it.Image()
		if !decodable(data) {
			continue
		}
		text, err := native.RecognizeText(data)
		if err != nil {
			log.Print("myclip: text in image: ", err)
			if errors.Is(err, native.ErrNoOCR) {
				return // no recognizer here: stop trying
			}
			continue
		}
		if text == "" {
			continue
		}
		if err := q.save(it.ID, text); err != nil {
			log.Print("myclip: text in image: ", err)
		}
	}
}
