package history

import (
	"net/url"
	"path"
	"strings"
	"unicode"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// maxTitleRunes bounds a derived title, enough for a menu row.
const maxTitleRunes = 80

// Size is the stored byte size of the item's contents (D6 caps).
func (it *Item) Size() int64 {
	var n int64
	for _, c := range it.Contents {
		n += int64(len(c.Value))
	}
	return n
}

// Value returns the bytes of the first content of type t.
func (it *Item) Value(t string) ([]byte, bool) {
	for _, c := range it.Contents {
		if c.Type == t {
			return c.Value, true
		}
	}
	return nil, false
}

// Text is the item's plain text, if it has any.
func (it *Item) Text() string {
	v, _ := it.Value(pasteboard.TypeString)
	return string(v)
}

// DisplayTitle is the user's title, else one derived from the contents:
// files by name, then text (whitespace collapsed, truncated), then images.
func (it *Item) DisplayTitle() string {
	if it.Title != "" {
		return it.Title
	}
	var files []string
	for _, c := range it.Contents {
		if c.Type != pasteboard.TypeFileURL {
			continue
		}
		if u, err := url.Parse(string(c.Value)); err == nil && u.Path != "" {
			files = append(files, path.Base(u.Path))
		}
	}
	if len(files) > 0 {
		return truncate(strings.Join(files, ", "))
	}
	if s := collapse(it.Text()); s != "" {
		return truncate(s)
	}
	for _, t := range []string{pasteboard.TypePNG, pasteboard.TypeTIFF, pasteboard.TypeJPEG, pasteboard.TypeHEIC} {
		if _, ok := it.Value(t); ok {
			return "Image"
		}
	}
	if len(it.Contents) > 0 {
		return it.Contents[0].Type
	}
	return "(empty)"
}

func collapse(s string) string {
	return strings.Join(strings.FieldsFunc(s, unicode.IsSpace), " ")
}

func truncate(s string) string {
	r := []rune(s)
	if len(r) <= maxTitleRunes {
		return s
	}
	return string(r[:maxTitleRunes-1]) + "…"
}
