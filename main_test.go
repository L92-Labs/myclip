package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"

	"github.com/L92-Labs/myclip/internal/settings"
)

func TestTrayIconsArePNG(t *testing.T) {
	for _, name := range settings.MenuIcons {
		for _, paused := range []bool{false, true} {
			img, err := png.Decode(bytes.NewReader(trayIcon(name, paused)))
			if err != nil || img.Bounds().Dx() != 32 {
				t.Fatalf("%s paused=%v: %v", name, paused, err)
			}
			if opaque(img) == 0 {
				t.Errorf("%s paused=%v: blank icon", name, paused)
			}
		}
	}
}

// opaque counts the pixels drawn at least half opaque.
func opaque(img image.Image) int {
	n := 0
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA).A > 0 {
				n++
			}
		}
	}
	return n
}

func TestRecentTitle(t *testing.T) {
	if got := recentTitle("short"); got != "short" {
		t.Errorf("got %q", got)
	}
	long := strings.Repeat("é", 50)
	if got := []rune(recentTitle(long)); len(got) != recentTitleRunes || got[len(got)-1] != '…' {
		t.Errorf("got %q", string(got))
	}
}

func TestScaleDown(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 400, 200))
	if b := scaleDown(img, 40, 1000).Bounds(); b.Dx() != 80 || b.Dy() != 40 {
		t.Errorf("by height: %v", b)
	}
	if b := scaleDown(img, 1000, 100).Bounds(); b.Dx() != 100 || b.Dy() != 50 {
		t.Errorf("by width: %v", b)
	}
	if scaleDown(img, 500, 500) != image.Image(img) {
		t.Error("scaled an image that fits")
	}
}
