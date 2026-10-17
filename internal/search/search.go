// Package search filters history titles as Maccy's Search does: exact
// (case-insensitive substring), regexp, fuzzy, or mixed (the first of
// those three that finds anything).
package search

import (
	"cmp"
	"regexp"
	"slices"
	"unicode"
	"unicode/utf8"
)

// Modes, as the settings name them.
const (
	Exact  = "exact"
	Fuzzy  = "fuzzy"
	Regexp = "regexp"
	Mixed  = "mixed"
)

// fuzzyLimit bounds how much of a title fuzzy search reads (Maccy: 5,000).
const fuzzyLimit = 5000

// Range is a match in a title, as byte offsets [Start, End).
type Range struct{ Start, End int }

// Result is a title that matched: its index in the input and where.
type Result struct {
	Index  int
	Ranges []Range
	score  int
}

// Search returns the titles matching query, in input order except for fuzzy
// results, which come best first. An empty query matches every title.
func Search(mode, query string, titles []string) []Result {
	if query == "" {
		out := make([]Result, len(titles))
		for i := range titles {
			out[i] = Result{Index: i}
		}
		return out
	}
	switch mode {
	case Regexp:
		return byRegexp(query, titles)
	case Fuzzy:
		return byFuzzy(query, titles)
	case Mixed:
		if r := exact(query, titles); len(r) > 0 {
			return r
		}
		if r := byRegexp(query, titles); len(r) > 0 {
			return r
		}
		return byFuzzy(query, titles)
	default:
		return exact(query, titles)
	}
}

func exact(query string, titles []string) []Result {
	return matchAll(regexp.MustCompile("(?i)"+regexp.QuoteMeta(query)), titles)
}

// byRegexp is case-sensitive, as NSString's .regularExpression option is;
// a pattern that does not compile matches nothing.
func byRegexp(query string, titles []string) []Result {
	re, err := regexp.Compile(query)
	if err != nil {
		return nil
	}
	return matchAll(re, titles)
}

func matchAll(re *regexp.Regexp, titles []string) []Result {
	var out []Result
	for i, t := range titles {
		if loc := re.FindStringIndex(t); loc != nil && loc[1] > loc[0] {
			out = append(out, Result{Index: i, Ranges: []Range{{loc[0], loc[1]}}})
		}
	}
	return out
}

// byFuzzy finds the query's letters in order in a title, ignoring case and
// spaces in the query; titles with closer, earlier matches come first.
func byFuzzy(query string, titles []string) []Result {
	var q []rune
	for _, r := range query {
		if !unicode.IsSpace(r) {
			q = append(q, unicode.ToLower(r))
		}
	}
	if len(q) == 0 {
		return Search(Exact, "", titles)
	}
	var out []Result
	for i, t := range titles {
		if len(t) > fuzzyLimit {
			t = t[:fuzzyLimit]
		}
		if r, ok := fuzzy(q, t); ok {
			r.Index = i
			out = append(out, r)
		}
	}
	slices.SortStableFunc(out, func(a, b Result) int { return cmp.Compare(a.score, b.score) })
	return out
}

// fuzzy matches q in t, preferring a contiguous run: it starts from each
// occurrence of the first letter and keeps the tightest match. The score is
// the gaps between letters plus where the match starts; lower is better.
func fuzzy(q []rune, t string) (Result, bool) {
	best := Result{score: -1}
	for start := 0; start < len(t); {
		r, size := utf8.DecodeRuneInString(t[start:])
		if unicode.ToLower(r) != q[0] {
			start += size
			continue
		}
		ranges, gaps, ok := follow(q, t, start)
		if !ok {
			break // later starts cannot match either
		}
		score := gaps*4 + utf8.RuneCountInString(t[:start])
		if best.score < 0 || score < best.score {
			best = Result{Ranges: ranges, score: score}
		}
		if gaps == 0 {
			break
		}
		start += size
	}
	return best, best.score >= 0
}

// follow matches q greedily from byte offset start, merging adjacent
// letters into ranges, and counts the runes skipped between them.
func follow(q []rune, t string, start int) (ranges []Range, gaps int, ok bool) {
	qi, skipped := 0, 0
	for i := start; i < len(t) && qi < len(q); {
		r, size := utf8.DecodeRuneInString(t[i:])
		if unicode.ToLower(r) == q[qi] {
			if n := len(ranges); n > 0 && ranges[n-1].End == i {
				ranges[n-1].End = i + size
			} else {
				ranges = append(ranges, Range{i, i + size})
				if qi > 0 {
					gaps += skipped
				}
			}
			skipped = 0
			qi++
		} else {
			skipped++
		}
		i += size
	}
	return ranges, gaps, qi == len(q)
}
