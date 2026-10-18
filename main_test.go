package main

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
	"time"

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

func TestDecodableRefusesHugeImages(t *testing.T) {
	encode := func(w, h int) []byte {
		var b bytes.Buffer
		png.Encode(&b, image.NewGray(image.Rect(0, 0, w, h)))
		return b.Bytes()
	}
	if !decodable(encode(100, 100)) {
		t.Error("refused a small image")
	}
	// A header claiming 100k×100k pixels: patch IHDR's width and height,
	// then its CRC (over type + data, bytes 12–28), so the header still
	// parses and only the pixel cap can refuse it.
	huge := encode(1, 1)
	copy(huge[16:24], []byte{0, 1, 0x86, 0xa0, 0, 1, 0x86, 0xa0})
	binary.BigEndian.PutUint32(huge[29:33], crc32.ChecksumIEEE(huge[12:29]))
	if cfg, _, err := image.DecodeConfig(bytes.NewReader(huge)); err != nil || cfg.Width != 100000 {
		t.Fatalf("patched header does not parse: %v %+v", err, cfg)
	}
	if decodable(huge) {
		t.Error("accepted a 10-gigapixel header")
	}
	if decodable([]byte("not an image")) {
		t.Error("accepted garbage")
	}
}

// The preview reports the copied image's size, not the scaled bitmap's.
func TestDecodeScaledKeepsOriginalSize(t *testing.T) {
	var b bytes.Buffer
	png.Encode(&b, image.NewGray(image.Rect(0, 0, 4000, 3000)))
	bmp, size := decodeScaled(b.Bytes(), 1600, 1600)
	if bmp == nil || size != [2]int{4000, 3000} {
		t.Fatalf("size %v", size)
	}
	if w, _ := bmp.Size(); w != 1600 {
		t.Fatalf("bitmap width %d, want 1600", w)
	}
	if bmp, _ := decodeScaled([]byte("junk"), 10, 10); bmp != nil {
		t.Fatal("decoded junk")
	}
}

func TestCopySoundIsShortAndQuiet(t *testing.T) {
	s := copySound()
	if d := s.Duration(); d < 50*time.Millisecond || d > 200*time.Millisecond {
		t.Fatalf("duration %v", d)
	}
}
