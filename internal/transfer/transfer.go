// Package transfer exports and imports a MyClip backup: the whole history
// plus the settings in one JSON file. Binary representations (images, HTML)
// marshal as base64, so the file is plain JSON a person can inspect.
package transfer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/settings"
)

// FormatVersion is the backup schema version readers accept.
const FormatVersion = 1

// MaxFileSize bounds one backup read; the store's own budget is far
// smaller, so anything past this is corrupt, not a big library.
const MaxFileSize = 2 << 30

type Content struct {
	Type  string `json:"type"`
	Value []byte `json:"value,omitempty"`
}

type Item struct {
	Application    string    `json:"application,omitempty"`
	FirstCopiedAt  time.Time `json:"firstCopiedAt"`
	LastCopiedAt   time.Time `json:"lastCopiedAt"`
	NumberOfCopies int       `json:"numberOfCopies"`
	Pin            string    `json:"pin,omitempty"`
	Title          string    `json:"title,omitempty"`
	Contents       []Content `json:"contents"`
}

// Archive is one backup file. Settings is present when the exporter had
// them; readers treat a nil Settings as history-only.
type Archive struct {
	Format   int                `json:"myclip"`
	Exported time.Time          `json:"exported"`
	Settings *settings.Settings `json:"settings,omitempty"`
	Items    []Item             `json:"items"`
}

// Export writes s and items as a backup.
func Export(w io.Writer, s *settings.Settings, items []*history.Item) error {
	a := Archive{
		Format:   FormatVersion,
		Exported: time.Now().UTC(),
		Settings: s,
	}
	for _, it := range items {
		x := Item{
			Application:    it.Application,
			FirstCopiedAt:  it.FirstCopiedAt,
			LastCopiedAt:   it.LastCopiedAt,
			NumberOfCopies: it.NumberOfCopies,
			Pin:            it.Pin,
			Title:          it.Title,
		}
		for _, c := range it.Contents {
			x.Contents = append(x.Contents, Content{Type: c.Type, Value: c.Value})
		}
		a.Items = append(a.Items, x)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(a)
}

// Read parses a backup and sanitizes it: unknown fields are ignored, bad
// rows are dropped, and pins or counts that make no sense are cleaned.
func Read(r io.Reader) (*Archive, error) {
	var a Archive
	if err := json.NewDecoder(io.LimitReader(r, MaxFileSize)).Decode(&a); err != nil {
		return nil, fmt.Errorf("transfer: not a MyClip backup: %w", err)
	}
	if a.Format != FormatVersion {
		return nil, fmt.Errorf("transfer: backup format %d, want %d", a.Format, FormatVersion)
	}
	if a.Settings != nil {
		a.Settings.Normalize()
	}
	kept := a.Items[:0]
	for _, x := range a.Items {
		if len(x.Contents) == 0 {
			continue // nothing to paste: not a clip
		}
		contents := make([]Content, 0, len(x.Contents))
		for _, c := range x.Contents {
			if c.Type == "" {
				continue
			}
			contents = append(contents, c)
		}
		if len(contents) == 0 {
			continue
		}
		x.Contents = contents
		if x.NumberOfCopies < 1 {
			x.NumberOfCopies = 1
		}
		if x.LastCopiedAt.IsZero() {
			x.LastCopiedAt = time.Now().UTC()
		}
		if x.FirstCopiedAt.IsZero() {
			x.FirstCopiedAt = x.LastCopiedAt
		}
		kept = append(kept, x)
	}
	a.Items = kept
	if len(a.Items) == 0 {
		return nil, errors.New("transfer: backup has no clips")
	}
	return &a, nil
}

// HistoryItems converts the archive's rows to history items, IDs unset.
func (a *Archive) HistoryItems() []*history.Item {
	out := make([]*history.Item, 0, len(a.Items))
	for _, x := range a.Items {
		it := &history.Item{
			Application:    x.Application,
			FirstCopiedAt:  x.FirstCopiedAt,
			LastCopiedAt:   x.LastCopiedAt,
			NumberOfCopies: x.NumberOfCopies,
			Pin:            x.Pin,
			Title:          x.Title,
		}
		for _, c := range x.Contents {
			it.Contents = append(it.Contents, history.Content{Type: c.Type, Value: c.Value})
		}
		out = append(out, it)
	}
	return out
}
