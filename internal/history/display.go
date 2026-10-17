package history

import (
	"cmp"
	"html"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// Sort orders (Sorter.By) and pin positions, as the settings name them.
const (
	SortLastCopiedAt   = "lastCopiedAt"
	SortFirstCopiedAt  = "firstCopiedAt"
	SortNumberOfCopies = "numberOfCopies"
	PinBottom          = "bottom"
)

// Sort returns items ordered as Maccy's Sorter: by the sort order, then
// pinned items first (or last with pinTo "bottom"). Stable.
func Sort(items []*Item, by, pinTo string) []*Item {
	out := slices.Clone(items)
	slices.SortStableFunc(out, func(a, b *Item) int {
		switch by {
		case SortFirstCopiedAt:
			return b.FirstCopiedAt.Compare(a.FirstCopiedAt)
		case SortNumberOfCopies:
			return cmp.Compare(b.NumberOfCopies, a.NumberOfCopies)
		default:
			return b.LastCopiedAt.Compare(a.LastCopiedAt)
		}
	})
	slices.SortStableFunc(out, func(a, b *Item) int {
		ap, bp := a.Pin != "", b.Pin != ""
		if ap == bp {
			return 0
		}
		if ap != (pinTo == PinBottom) {
			return -1
		}
		return 1
	})
	return out
}

// SupportedPins are the keys a pin may take (HistoryItem.supportedPins):
// a, q, v, w and z are reserved for select all, quit, paste, close, undo.
// Keys of the in-popup shortcuts are removed by the caller.
var SupportedPins = []string{
	"b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l",
	"m", "n", "o", "p", "r", "s", "t", "u", "x", "y",
}

// AvailablePins returns the supported keys no item has, in order, leaving
// out reserved ones.
func AvailablePins(items []*Item, reserved ...string) []string {
	var out []string
	for _, k := range SupportedPins {
		if slices.Contains(reserved, k) || slices.ContainsFunc(items, func(it *Item) bool { return it.Pin == k }) {
			continue
		}
		out = append(out, k)
	}
	return out
}

// imageTypes are StorageType.images, in Maccy's lookup order.
var imageTypes = []string{pasteboard.TypePNG, pasteboard.TypeTIFF, pasteboard.TypeJPEG, pasteboard.TypeHEIC}

// Image returns the item's first image representation and its type.
func (it *Item) Image() (data []byte, typ string, ok bool) {
	for _, c := range it.Contents {
		if slices.Contains(imageTypes, c.Type) {
			return c.Value, c.Type, true
		}
	}
	return nil, "", false
}

// ResolveFileURL turns a file reference URL (file:///.file/id=…, as Finder
// copies) into a file path URL; nil leaves URLs as they are.
var ResolveFileURL func(string) string

// FilePaths are the paths of the item's file URLs, percent-decoded.
func (it *Item) FilePaths() []string {
	var out []string
	for _, c := range it.Contents {
		if c.Type != pasteboard.TypeFileURL {
			continue
		}
		s := string(c.Value)
		if ResolveFileURL != nil && strings.HasPrefix(s, "file:///.file/") {
			s = ResolveFileURL(s)
		}
		if u, err := url.Parse(s); err == nil && u.Path != "" {
			out = append(out, u.Path)
		}
	}
	return out
}

// PreviewText is HistoryItem.previewableText: files, else plain text, else
// the text of RTF or HTML, else the title.
func (it *Item) PreviewText() string {
	if files := it.FilePaths(); len(files) > 0 {
		return strings.Join(files, "\n")
	}
	if s := it.Text(); s != "" {
		return s
	}
	if v, ok := it.Value(pasteboard.TypeRTF); ok {
		if s := rtfText(string(v)); strings.TrimSpace(s) != "" {
			return s
		}
	}
	if v, ok := it.Value(pasteboard.TypeHTML); ok {
		if s := htmlText(string(v)); strings.TrimSpace(s) != "" {
			return s
		}
	}
	return it.Title
}

var (
	leadingSpaces  = regexp.MustCompile(`^ +`)
	trailingSpaces = regexp.MustCompile(` +$`)
)

// ListTitle is the row title: the item's own title (a renamed pin), else
// HistoryItem.generateTitle — empty for images, which show a thumbnail;
// with special symbols, edge spaces become ·, newlines ⏎ and tabs ⇥.
func (it *Item) ListTitle(specialSymbols bool) string {
	if it.Title != "" {
		return it.Title
	}
	if _, _, ok := it.Image(); ok && len(it.FilePaths()) == 0 {
		return ""
	}
	title := shorten(it.PreviewText(), 1000)
	if !specialSymbols {
		return strings.TrimSpace(title)
	}
	dots := func(s string) string { return strings.Repeat("·", len(s)) }
	title = leadingSpaces.ReplaceAllStringFunc(title, dots)
	title = trailingSpaces.ReplaceAllStringFunc(title, dots)
	return strings.NewReplacer("\n", "⏎", "\t", "⇥").Replace(title)
}

func shorten(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

var (
	htmlTags   = regexp.MustCompile(`(?s)<(script|style)[^>]*>.*?</(script|style)>|<[^>]+>`)
	htmlBlocks = regexp.MustCompile(`(?i)<(br|/p|/div|/li|/tr|/h[1-6])[^>]*>`)
)

// htmlText is the text of an HTML fragment, roughly as a browser shows it.
func htmlText(s string) string {
	s = htmlBlocks.ReplaceAllString(s, "\n$0")
	s = htmlTags.ReplaceAllString(s, "")
	return strings.TrimSpace(html.UnescapeString(s))
}

// rtfText is the text of an RTF document: groups of destinations
// (fonttbl, colortbl, …) dropped, control words removed, \par as newline.
func rtfText(s string) string {
	var b strings.Builder
	depth, skipAt := 0, -1
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch == '{':
			depth++
			if skipAt < 0 && strings.HasPrefix(s[i+1:], `\*`) {
				skipAt = depth
			}
		case ch == '}':
			if depth == skipAt {
				skipAt = -1
			}
			depth--
		case ch == '\\':
			j := i + 1
			if j < len(s) && (s[j] == '\\' || s[j] == '{' || s[j] == '}') {
				if skipAt < 0 {
					b.WriteByte(s[j])
				}
				i = j
				continue
			}
			if j < len(s) && s[j] == '\'' && j+2 < len(s) {
				// \'hh: a Windows-1252 byte; ASCII-range only.
				var v byte
				for _, h := range s[j+1 : j+3] {
					v <<= 4
					switch {
					case h >= '0' && h <= '9':
						v |= byte(h - '0')
					case h >= 'a' && h <= 'f':
						v |= byte(h-'a') + 10
					case h >= 'A' && h <= 'F':
						v |= byte(h-'A') + 10
					}
				}
				if skipAt < 0 && v < 0x80 {
					b.WriteByte(v)
				}
				i = j + 2
				continue
			}
			k := j
			for k < len(s) && (s[k] >= 'a' && s[k] <= 'z' || s[k] >= 'A' && s[k] <= 'Z') {
				k++
			}
			word := s[j:k]
			for k < len(s) && (s[k] == '-' || s[k] >= '0' && s[k] <= '9') {
				k++
			}
			if k < len(s) && s[k] == ' ' {
				k++
			}
			if skipAt < 0 {
				switch {
				case isRTFDestination(word):
					skipAt = depth
				case word == "par" || word == "line":
					b.WriteByte('\n')
				case word == "tab":
					b.WriteByte('\t')
				}
			}
			i = k - 1
		case ch == '\r' || ch == '\n':
		default:
			if skipAt < 0 && depth > 0 {
				b.WriteByte(ch)
			}
		}
	}
	return strings.TrimSpace(b.String())
}

func isRTFDestination(word string) bool {
	switch word {
	case "fonttbl", "colortbl", "stylesheet", "info", "pict", "header", "footer", "expandedcolortbl":
		return true
	}
	return false
}
