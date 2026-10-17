package search

import (
	"reflect"
	"testing"
)

func indices(rs []Result) []int {
	var out []int
	for _, r := range rs {
		out = append(out, r.Index)
	}
	return out
}

var titles = []string{"Hello World", "help me", "World peace", "héllo ÉTÉ", "yellow"}

func TestEmptyQueryMatchesAll(t *testing.T) {
	if got := indices(Search(Exact, "", titles)); len(got) != len(titles) {
		t.Fatalf("got %v", got)
	}
}

func TestExactCaseInsensitive(t *testing.T) {
	got := Search(Exact, "WORLD", titles)
	if !reflect.DeepEqual(indices(got), []int{0, 2}) {
		t.Fatalf("got %v", indices(got))
	}
	if got[0].Ranges[0] != (Range{6, 11}) {
		t.Fatalf("range %v", got[0].Ranges)
	}
	if got := Search(Exact, "été", titles); len(got) != 1 || got[0].Index != 3 {
		t.Fatalf("unicode fold %v", got)
	}
	if got := Search(Exact, "l.o", titles); len(got) != 0 {
		t.Fatalf("exact treated as regexp: %v", got)
	}
}

func TestRegexp(t *testing.T) {
	if got := indices(Search(Regexp, "^[Hh]el", titles)); !reflect.DeepEqual(got, []int{0, 1}) {
		t.Fatalf("got %v", got)
	}
	if got := Search(Regexp, "(", titles); got != nil {
		t.Fatalf("bad pattern matched %v", got)
	}
	if got := Search(Regexp, "world", titles); len(got) != 0 {
		t.Fatalf("regexp should be case-sensitive: %v", got)
	}
}

func TestFuzzyOrderAndRanges(t *testing.T) {
	got := Search(Fuzzy, "hlo", titles)
	if !reflect.DeepEqual(indices(got), []int{0, 3}) {
		t.Fatalf("got %v", indices(got))
	}
	if !reflect.DeepEqual(got[0].Ranges, []Range{{0, 1}, {2, 3}, {4, 5}}) {
		t.Fatalf("ranges %v", got[0].Ranges)
	}
	// A contiguous match beats a spread one.
	got = Search(Fuzzy, "low", []string{"l o w", "yellow"})
	if got[0].Index != 1 {
		t.Fatalf("got %v", indices(got))
	}
}

func TestMixedFallsThrough(t *testing.T) {
	if got := indices(Search(Mixed, "peace", titles)); !reflect.DeepEqual(got, []int{2}) {
		t.Fatalf("exact step %v", got)
	}
	if got := indices(Search(Mixed, "^y", titles)); !reflect.DeepEqual(got, []int{4}) {
		t.Fatalf("regexp step %v", got)
	}
	if got := indices(Search(Mixed, "wpc", titles)); !reflect.DeepEqual(got, []int{2}) {
		t.Fatalf("fuzzy step %v", got)
	}
}
