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
	pending []ocrJob // waiting, oldest first; never dropped, only IDs held
	running bool
	// item returns the clip stored under id now (Engine.Item).
	item func(id int64) *history.Item
	// save stores a reading for the image with that digest, refusing it
	// atomically when the ID now holds another image.
	save func(id int64, image [32]byte, text string) error
	// recognize reads an image's text; nil uses the system's recognizer.
	recognize func([]byte) (string, error)
	stopped   bool // no recognizer on this system
}

// ocrJob names an image to read without holding its bytes, so a backlog
// never keeps evicted images in memory.
type ocrJob struct {
	id     int64
	digest [32]byte
}

// add queues an image kept from before, behind the others.
func (q *ocrQueue) add(it *history.Item) { q.enqueue(it, false) }

// addLive queues an image just copied ahead of the startup backlog: the
// screenshot you just took is the one you will search for.
func (q *ocrQueue) addLive(it *history.Item) { q.enqueue(it, true) }

func (q *ocrQueue) enqueue(it *history.Item, live bool) {
	if !needsOCR(it) {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.stopped {
		return
	}
	job := ocrJob{it.ID, engine.ImageDigest(it)}
	if live {
		q.pending = append([]ocrJob{job}, q.pending...)
	} else {
		q.pending = append(q.pending, job)
	}
	if !q.running {
		q.running = true
		go q.run()
	}
}

// next takes the next job, or ends the worker when none is left.
func (q *ocrQueue) next() (ocrJob, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.pending) == 0 || q.stopped {
		q.running = false
		return ocrJob{}, false
	}
	job := q.pending[0]
	q.pending = q.pending[1:]
	return job, true
}

// needsOCR reports whether it is an image not read yet, small enough.
func needsOCR(it *history.Item) bool {
	data, _, ok := it.Image()
	return ok && len(data) <= ocrMaxBytes && len(it.FilePaths()) == 0 && !it.TextRecognized() && it.Text() == ""
}

func (q *ocrQueue) run() {
	for {
		job, ok := q.next()
		if !ok {
			return
		}
		// Gone (trimmed, deleted) or another image under the ID now: skip
		// without reading it.
		it := q.item(job.id)
		if it == nil || engine.ImageDigest(it) != job.digest || it.TextRecognized() {
			continue
		}
		data, typ, _ := it.Image()
		// Go cannot decode HEIC to check it, but Vision reads it.
		if typ != pasteboard.TypeHEIC && !decodable(data) {
			continue
		}
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
		if err := q.save(job.id, job.digest, text); errors.Is(err, engine.ErrNoItem) {
			continue
		} else if err != nil {
			log.Print("myclip: text in image: ", err)
		}
	}
}
