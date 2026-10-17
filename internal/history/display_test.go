package history

import (
	"slices"
	"testing"
	"time"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

func textItem(s string) *Item {
	return NewItem([]Content{{Type: pasteboard.TypeString, Value: []byte(s)}}, "", time.Unix(0, 0))
}

func TestListTitleSpecialSymbols(t *testing.T) {
	it := textItem("  a\tb\nc ")
	if got := it.ListTitle(true); got != "··a⇥b⏎c·" {
		t.Errorf("special: %q", got)
	}
	if got := it.ListTitle(false); got != "a\tb\nc" {
		t.Errorf("plain: %q", got)
	}
	it.Title = "Renamed"
	if got := it.ListTitle(true); got != "Renamed" {
		t.Errorf("renamed: %q", got)
	}
}

func TestListTitleImageEmpty(t *testing.T) {
	it := NewItem([]Content{{Type: pasteboard.TypePNG, Value: []byte{1}}}, "", time.Now())
	if got := it.ListTitle(true); got != "" {
		t.Errorf("image title %q", got)
	}
}

func TestPreviewTextFallbacks(t *testing.T) {
	files := NewItem([]Content{{Type: pasteboard.TypeFileURL, Value: []byte("file:///Users/me/My%20File.txt")}}, "", time.Now())
	if got := files.PreviewText(); got != "/Users/me/My File.txt" {
		t.Errorf("files %q", got)
	}
	h := NewItem([]Content{{Type: pasteboard.TypeHTML, Value: []byte("<p>Hello &amp; <b>bye</b></p>")}}, "", time.Now())
	if got := h.PreviewText(); got != "Hello & bye" {
		t.Errorf("html %q", got)
	}
	r := NewItem([]Content{{Type: pasteboard.TypeRTF, Value: []byte(`{\rtf1\ansi{\fonttbl\f0\fswiss Helvetica;}{\colortbl;\red255\green255\blue255;}\f0\fs24 Hello\par World \{x\}}`)}}, "", time.Now())
	if got := r.PreviewText(); got != "Hello\nWorld {x}" {
		t.Errorf("rtf %q", got)
	}
}

func TestSort(t *testing.T) {
	base := time.Unix(1000, 0)
	a := &Item{ID: 1, FirstCopiedAt: base, LastCopiedAt: base.Add(3), NumberOfCopies: 1}
	b := &Item{ID: 2, FirstCopiedAt: base.Add(1), LastCopiedAt: base.Add(1), NumberOfCopies: 5}
	c := &Item{ID: 3, FirstCopiedAt: base.Add(2), LastCopiedAt: base.Add(2), NumberOfCopies: 2, Pin: "b"}
	ids := func(items []*Item) []int64 {
		var out []int64
		for _, it := range items {
			out = append(out, it.ID)
		}
		return out
	}
	for _, tc := range []struct {
		by, pin string
		want    []int64
	}{
		{SortLastCopiedAt, "top", []int64{3, 1, 2}},
		{SortLastCopiedAt, PinBottom, []int64{1, 2, 3}},
		{SortFirstCopiedAt, "top", []int64{3, 2, 1}},
		{SortNumberOfCopies, "top", []int64{3, 2, 1}},
	} {
		if got := ids(Sort([]*Item{a, b, c}, tc.by, tc.pin)); !slices.Equal(got, tc.want) {
			t.Errorf("%s/%s: %v want %v", tc.by, tc.pin, got, tc.want)
		}
	}
}

func TestAvailablePins(t *testing.T) {
	got := AvailablePins([]*Item{{Pin: "b"}, {Pin: "d"}}, "p")
	if got[0] != "c" || slices.Contains(got, "b") || slices.Contains(got, "d") || slices.Contains(got, "p") {
		t.Errorf("pins %v", got)
	}
}
