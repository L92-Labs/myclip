package history

import (
	"bytes"
	"time"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// Content is one representation (UTI + bytes).
type Content struct {
	Type  string
	Value []byte
}

// Item is one history row (Maccy HistoryItem shape).
type Item struct {
	Application    string
	FirstCopiedAt  time.Time
	LastCopiedAt   time.Time
	NumberOfCopies int
	Pin            string
	Title          string
	Contents       []Content
}

// NewItem builds an item from pasteboard contents, stripping the own-write
// marker from stored representations (marker is metadata-only).
func NewItem(contents []Content, application string, now time.Time) *Item {
	filtered := make([]Content, 0, len(contents))
	for _, c := range contents {
		if c.Type == pasteboard.TypeOwnWrite {
			continue
		}
		filtered = append(filtered, Content{Type: c.Type, Value: append([]byte(nil), c.Value...)})
	}
	return &Item{
		Application:    application,
		FirstCopiedAt:  now,
		LastCopiedAt:   now,
		NumberOfCopies: 1,
		Contents:       filtered,
	}
}

// HasOwnWriteMarker reports whether the incoming pasteboard carried the marker
// (checked before stripping into Contents).
func HasOwnWriteMarker(contents []Content) bool {
	for _, c := range contents {
		if c.Type == pasteboard.TypeOwnWrite {
			return true
		}
	}
	return false
}

// Modified returns the changeCount encoded in TypeModified, if present.
func (it *Item) Modified() (int, bool) {
	for _, c := range it.Contents {
		if c.Type != pasteboard.TypeModified {
			continue
		}
		// Maccy stores the prior changeCount as UTF-8 decimal string.
		var n int
		for _, b := range c.Value {
			if b < '0' || b > '9' {
				return 0, false
			}
			n = n*10 + int(b-'0')
		}
		return n, true
	}
	return 0, false
}

// Supersedes reports whether it contains every non-transient content of other
// (HistoryItem.supersedes @ eb03eba).
func (it *Item) Supersedes(other *Item) bool {
	transient := map[string]struct{}{}
	for _, t := range pasteboard.TransientTypes {
		transient[t] = struct{}{}
	}
	for _, oc := range other.Contents {
		if _, skip := transient[oc.Type]; skip {
			continue
		}
		found := false
		for _, ic := range it.Contents {
			if ic.Type == oc.Type && bytes.Equal(ic.Value, oc.Value) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
