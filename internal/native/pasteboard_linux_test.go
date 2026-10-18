//go:build linux

package native

import (
	"testing"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

func TestUTIMapping(t *testing.T) {
	for in, want := range map[string]string{
		"UTF8_STRING":                    pasteboard.TypeString,
		"text/plain;charset=utf-8":       pasteboard.TypeString,
		"text/html":                      pasteboard.TypeHTML,
		"image/png":                      pasteboard.TypePNG,
		"text/uri-list":                  pasteboard.TypeFileURL,
		"x-kde-passwordManagerHint":      pasteboard.TypeConcealed,
		"application/x-myclip-own-write": pasteboard.TypeOwnWrite,
		"application/x-custom":           "application/x-custom",
	} {
		if got := uti(in); got != want {
			t.Errorf("%s → %s, want %s", in, got, want)
		}
	}
}

// KeePassXC's hint alone makes the whole copy ignored, before any read.
func TestPasswordHintIsIgnored(t *testing.T) {
	types := []string{"UTF8_STRING", uti("UTF8_STRING"), "x-kde-passwordManagerHint", uti("x-kde-passwordManagerHint")}
	if !capture.ShouldIgnoreTypes(types, capture.Prefs{}) {
		t.Fatal("password-manager copy not ignored")
	}
	if capture.ShouldIgnoreTypes([]string{"UTF8_STRING", uti("UTF8_STRING")}, capture.Prefs{}) {
		t.Fatal("plain text ignored")
	}
}
