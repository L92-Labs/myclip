package capture

import (
	"regexp"
	"strings"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// Prefs are the capture-time ignore settings (subset needed for M1 tests).
type Prefs struct {
	IgnoreEvents              bool
	IgnoreOnlyNextEvent       bool
	EnabledPasteboardTypes    []string
	IgnoredPasteboardTypes    []string // user set; defaults = pasteboard.DefaultIgnoredPasteboardTypes
	IgnoredApps               []string
	IgnoreAllAppsExceptListed bool
	IgnoreRegexp              []string
}

// ShouldIgnoreTypes implements Clipboard.shouldIgnore(types) @ eb03eba:
// ignore if disjoint from enabled, OR any intersection with (always-ignore ∪ user ignored).
func ShouldIgnoreTypes(types []string, prefs Prefs) bool {
	enabled := toSet(prefs.EnabledPasteboardTypes)
	if len(enabled) == 0 {
		enabled = toSet(pasteboard.DefaultEnabledPasteboardTypes)
	}
	ignored := toSet(pasteboard.AlwaysIgnoreTypes)
	userIgn := prefs.IgnoredPasteboardTypes
	if userIgn == nil {
		userIgn = pasteboard.DefaultIgnoredPasteboardTypes
	}
	for _, t := range userIgn {
		ignored[t] = struct{}{}
	}

	typeSet := toSet(types)
	if disjoint(typeSet, enabled) {
		return true
	}
	if intersects(typeSet, ignored) {
		return true
	}
	return false
}

// ShouldIgnoreApp implements Clipboard.shouldIgnore(sourceAppBundle) @ eb03eba.
//
// A source that cannot be identified (no frontmost app, or an unnamed
// window on Linux) is captured normally, but ignored in allowlist mode: an
// unknown app is not on the list.
func ShouldIgnoreApp(bundleID string, prefs Prefs) bool {
	if bundleID == "" {
		return prefs.IgnoreAllAppsExceptListed
	}
	listed := false
	for _, a := range prefs.IgnoredApps {
		if a == bundleID {
			listed = true
			break
		}
	}
	if prefs.IgnoreAllAppsExceptListed {
		return !listed
	}
	return listed
}

// ItemMatchesRegexp reports whether any regexp matches the item's plain-text
// (or any string-ish) content. Maccy evaluates per pasteboard item inside the
// merge loop — matching item dropped, non-matching survive.
func ItemMatchesRegexp(text string, patterns []string) bool {
	for _, p := range patterns {
		if p == "" {
			continue
		}
		re, err := regexp.Compile(p)
		if err != nil {
			continue
		}
		if re.MatchString(text) {
			return true
		}
	}
	return false
}

// StripDynamicAndOLE removes dyn.* and com.microsoft.ole.source.* types
// (Clipboard.swift content filter).
func StripDynamicAndOLE(types map[string][]byte) map[string][]byte {
	out := make(map[string][]byte, len(types))
	for t, v := range types {
		if strings.HasPrefix(t, "dyn.") {
			continue
		}
		if strings.HasPrefix(t, "com.microsoft.ole.source.") {
			continue
		}
		out[t] = v
	}
	return out
}

// DropWordLinkSet applies the Word ObjectLink/Link-Source + PDF subtraction
// when both Word link types are present (issues #613/#770).
func DropWordLinkSet(types map[string][]byte) map[string][]byte {
	_, hasObj := types[pasteboard.TypeMicrosoftObjectLink]
	_, hasSrc := types[pasteboard.TypeMicrosoftLinkSource]
	if !(hasObj && hasSrc) {
		return types
	}
	out := make(map[string][]byte, len(types))
	for t, v := range types {
		switch t {
		case pasteboard.TypeMicrosoftObjectLink, pasteboard.TypeMicrosoftLinkSource, pasteboard.TypePDF:
			continue
		default:
			out[t] = v
		}
	}
	return out
}

func toSet(xs []string) map[string]struct{} {
	m := make(map[string]struct{}, len(xs))
	for _, x := range xs {
		m[x] = struct{}{}
	}
	return m
}

func disjoint(a, b map[string]struct{}) bool {
	for t := range a {
		if _, ok := b[t]; ok {
			return false
		}
	}
	return true
}

func intersects(a, b map[string]struct{}) bool {
	for t := range a {
		if _, ok := b[t]; ok {
			return true
		}
	}
	return false
}
