package native

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"testing"
	"unicode/utf16"
)

// The clipboard's text allocation can hold more than the text: decoding
// must stop at the first NUL code unit, not the last.
func TestUTF16StopsAtFirstTerminator(t *testing.T) {
	// "hello", its NUL, and four padding bytes an over-reporting GlobalSize
	// would hand back.
	raw := append(utf16StringToBytes("hello"), 0, 0, 0, 0)
	if got := utf16BytesToString(raw); got != "hello" {
		t.Fatalf("padding kept: %q", got)
	}
	// Interior NULs are terminators too (they end the text).
	raw = []byte{'h', 0, 0, 0, 'i', 0}
	if got := utf16BytesToString(raw); got != "h" {
		t.Fatalf("interior NUL ignored: %q", got)
	}
	// No terminator at all: what is there still decodes.
	raw = []byte{'h', 0, 'i', 0}
	if got := utf16BytesToString(raw); got != "hi" {
		t.Fatalf("unterminated lost: %q", got)
	}
}

// A malformed wide CF_HDROP without a terminator must not panic the parser.
func TestDropFilesUnterminatedWidePath(t *testing.T) {
	payload := make([]byte, 20+2) // header + one nonzero UTF-16 unit, no NUL
	binary.LittleEndian.PutUint32(payload[0:4], 20)
	binary.LittleEndian.PutUint32(payload[16:20], 1)
	payload[20], payload[21] = 'x', 0
	if got := dropFiles(payload); len(got) != 0 {
		t.Fatalf("unterminated path kept: %q", got)
	}
	// A first path that ends, then an unterminated one: the first is kept,
	// the malformed tail is dropped.
	payload = make([]byte, 20+10+4)
	binary.LittleEndian.PutUint32(payload[0:4], 20)
	binary.LittleEndian.PutUint32(payload[16:20], 1)
	for i, v := range utf16.Encode([]rune("C:\\a\x00")) {
		binary.LittleEndian.PutUint16(payload[20+i*2:], v)
	}
	for i, v := range utf16.Encode([]rune("bx")) {
		binary.LittleEndian.PutUint16(payload[30+i*2:], v)
	}
	if got := dropFiles(payload); len(got) != 1 || got[0] != `C:\a` {
		t.Fatalf("got %q", got)
	}
}

// Height 0x80000000 must be rejected before it can overflow into a huge
// image.
func TestDIBRejectsExtremeHeight(t *testing.T) {
	dib := make([]byte, 40)
	binary.LittleEndian.PutUint32(dib[0:4], 40)
	binary.LittleEndian.PutUint32(dib[4:8], 1<<15)
	binary.LittleEndian.PutUint32(dib[8:12], 0x80000000)
	binary.LittleEndian.PutUint16(dib[12:14], 1)
	binary.LittleEndian.PutUint16(dib[14:16], 24)
	if _, ok := dibToPNG(dib); ok {
		t.Fatal("height 0x80000000 accepted")
	}
}

// BI_BITFIELDS channels follow their declared masks, not an assumed BGRA
// layout.
func TestDIBBitfieldsUseDeclaredMasks(t *testing.T) {
	const w, h = 1, 1
	dib := make([]byte, 40+12+w*4*h)
	binary.LittleEndian.PutUint32(dib[0:4], 40)
	binary.LittleEndian.PutUint32(dib[4:8], w)
	binary.LittleEndian.PutUint32(dib[8:12], h)
	binary.LittleEndian.PutUint16(dib[12:14], 1)
	binary.LittleEndian.PutUint16(dib[14:16], 32)
	binary.LittleEndian.PutUint32(dib[16:20], 3) // BI_BITFIELDS
	// Masks r=0x000000ff, g=0x0000ff00, b=0x00ff0000: little-endian RGB.
	binary.LittleEndian.PutUint32(dib[40:44], 0x000000ff)
	binary.LittleEndian.PutUint32(dib[44:48], 0x0000ff00)
	binary.LittleEndian.PutUint32(dib[48:52], 0x00ff0000)
	// One pure red pixel under those masks.
	binary.LittleEndian.PutUint32(dib[52:56], 0x000000ff)
	got, ok := dibToPNG(dib)
	if !ok {
		t.Fatal("bitfields DIB rejected")
	}
	img, err := png.Decode(bytes.NewReader(got))
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, _ := img.At(0, 0).RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Fatalf("red became %d %d %d", r>>8, g>>8, b>>8)
	}

	// A V4-style header keeps the masks inside its own 108 bytes.
	dib = make([]byte, 108+w*4*h)
	binary.LittleEndian.PutUint32(dib[0:4], 108)
	binary.LittleEndian.PutUint32(dib[4:8], w)
	binary.LittleEndian.PutUint32(dib[8:12], h)
	binary.LittleEndian.PutUint16(dib[12:14], 1)
	binary.LittleEndian.PutUint16(dib[14:16], 32)
	binary.LittleEndian.PutUint32(dib[16:20], 3)
	binary.LittleEndian.PutUint32(dib[40:44], 0x00ff0000)
	binary.LittleEndian.PutUint32(dib[44:48], 0x0000ff00)
	binary.LittleEndian.PutUint32(dib[48:52], 0x000000ff)
	binary.LittleEndian.PutUint32(dib[52:56], 0xff000000)   // the alpha mask
	binary.LittleEndian.PutUint32(dib[108:112], 0xff0000ff) // opaque blue under those masks
	got, ok = dibToPNG(dib)
	if !ok {
		t.Fatal("V4 bitfields DIB rejected")
	}
	img, err = png.Decode(bytes.NewReader(got))
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, a := img.At(0, 0).RGBA()
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 255 || a>>8 != 255 {
		t.Fatalf("blue became %d %d %d %d", r>>8, g>>8, b>>8, a>>8)
	}

	// A V2 header (52 bytes) declares RGB masks but no alpha mask: the
	// first pixel must not be mistaken for one.
	dib = make([]byte, 52+2*4)
	binary.LittleEndian.PutUint32(dib[0:4], 52)
	binary.LittleEndian.PutUint32(dib[4:8], 2)
	binary.LittleEndian.PutUint32(dib[8:12], 1)
	binary.LittleEndian.PutUint16(dib[12:14], 1)
	binary.LittleEndian.PutUint16(dib[14:16], 32)
	binary.LittleEndian.PutUint32(dib[16:20], 3)
	binary.LittleEndian.PutUint32(dib[40:44], 0x000000ff)
	binary.LittleEndian.PutUint32(dib[44:48], 0x0000ff00)
	binary.LittleEndian.PutUint32(dib[48:52], 0x00ff0000)
	binary.LittleEndian.PutUint32(dib[52:56], 0x000000ff) // opaque red
	binary.LittleEndian.PutUint32(dib[56:60], 0x00ff0000) // opaque blue
	if got, ok = dibToPNG(dib); !ok {
		t.Fatal("V2 bitfields DIB rejected")
	}
	img, err = png.Decode(bytes.NewReader(got))
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, a = img.At(0, 0).RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 || a>>8 != 255 {
		t.Fatalf("pixel 0: %d %d %d %d", r>>8, g>>8, b>>8, a>>8)
	}
	r, g, b, a = img.At(1, 0).RGBA()
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 255 || a>>8 != 255 {
		t.Fatalf("pixel 1 (read as an alpha mask?): %d %d %d %d", r>>8, g>>8, b>>8, a>>8)
	}

	// Non-contiguous masks are refused, not silently misdecoded.
	dib = make([]byte, 40+12+4)
	binary.LittleEndian.PutUint32(dib[0:4], 40)
	binary.LittleEndian.PutUint32(dib[8:12], 1)
	binary.LittleEndian.PutUint16(dib[14:16], 32)
	binary.LittleEndian.PutUint32(dib[16:20], 3)
	binary.LittleEndian.PutUint32(dib[40:44], 0x000000ff|0x0000f000) // gappy
	binary.LittleEndian.PutUint32(dib[44:48], 0x0000ff00)
	binary.LittleEndian.PutUint32(dib[48:52], 0x00ff0000)
	if _, ok := dibToPNG(dib); ok {
		t.Fatal("gappy mask accepted")
	}
}
