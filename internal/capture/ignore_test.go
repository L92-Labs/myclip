package capture_test

import (
	"testing"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

func TestShouldIgnoreConcealed(t *testing.T) {
	prefs := capture.Prefs{EnabledPasteboardTypes: pasteboard.DefaultEnabledPasteboardTypes}
	types := []string{pasteboard.TypeString, pasteboard.TypeConcealed}
	if !capture.ShouldIgnoreTypes(types, prefs) {
		t.Fatal("concealed must be ignored")
	}
}

func TestShouldIgnoreDisjointFromEnabled(t *testing.T) {
	prefs := capture.Prefs{EnabledPasteboardTypes: []string{pasteboard.TypeString}}
	types := []string{pasteboard.TypePNG}
	if !capture.ShouldIgnoreTypes(types, prefs) {
		t.Fatal("disjoint from enabled must ignore")
	}
}

func TestShouldCaptureDefaultText(t *testing.T) {
	prefs := capture.Prefs{
		EnabledPasteboardTypes: pasteboard.DefaultEnabledPasteboardTypes,
		IgnoredPasteboardTypes: pasteboard.DefaultIgnoredPasteboardTypes,
	}
	types := []string{pasteboard.TypeString}
	if capture.ShouldIgnoreTypes(types, prefs) {
		t.Fatal("plain text must capture under defaults")
	}
}

func TestDefaultIgnoredHasExactlyFive(t *testing.T) {
	if n := len(pasteboard.DefaultIgnoredPasteboardTypes); n != 5 {
		t.Fatalf("default ignored count=%d want 5", n)
	}
	for _, tname := range pasteboard.DefaultIgnoredPasteboardTypes {
		if tname == pasteboard.TypeUniversalClipboard {
			t.Fatal("com.apple.is-remote-clipboard must not be in defaults")
		}
	}
}

func TestShouldIgnoreAppAllowlist(t *testing.T) {
	prefs := capture.Prefs{
		IgnoreAllAppsExceptListed: true,
		IgnoredApps:               []string{"com.ok.app"},
	}
	if !capture.ShouldIgnoreApp("com.other", prefs) {
		t.Fatal("allowlist mode: unlisted must ignore")
	}
	if capture.ShouldIgnoreApp("com.ok.app", prefs) {
		t.Fatal("allowlist mode: listed must capture")
	}
}

func TestDropWordLinkSet(t *testing.T) {
	in := map[string][]byte{
		pasteboard.TypeMicrosoftObjectLink: []byte("o"),
		pasteboard.TypeMicrosoftLinkSource: []byte("s"),
		pasteboard.TypePDF:                 []byte("%PDF"),
		pasteboard.TypeString:              []byte("keep"),
	}
	out := capture.DropWordLinkSet(in)
	if _, ok := out[pasteboard.TypePDF]; ok {
		t.Fatal("pdf must drop with word link set")
	}
	if _, ok := out[pasteboard.TypeString]; !ok {
		t.Fatal("string must remain")
	}
}

func TestItemMatchesRegexp(t *testing.T) {
	if !capture.ItemMatchesRegexp("secret-token", []string{`secret-.*`}) {
		t.Fatal("expected match")
	}
	if capture.ItemMatchesRegexp("hello", []string{`secret-.*`}) {
		t.Fatal("expected no match")
	}
}
