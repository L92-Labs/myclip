package native

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"slices"
	"strings"
	"testing"
	"unicode/utf16"
)

func TestCFHTMLParse(t *testing.T) {
	frag := []byte("<p>hi & bye</p>")
	payload := cfhtmlBuild(frag)
	if got := cfhtmlParse(payload); !bytes.Equal(got, frag) {
		t.Fatalf("round trip: got %q", got)
	}
	// A payload with no readable header comes back whole.
	plain := []byte("<div>raw</div>")
	if got := cfhtmlParse(plain); !bytes.Equal(got, plain) {
		t.Fatalf("no header: got %q", got)
	}
	// Garbage offsets fall back to the whole payload.
	if got := cfhtmlParse([]byte("Version:0.9\r\nStartFragment:9999\r\nEndFragment:12\r\n")); !bytes.HasPrefix(got, []byte("Version:0.9")) {
		t.Fatalf("garbage offsets: got %q", got)
	}
}

func TestCFHTMLBuildOffsets(t *testing.T) {
	frag := []byte("x")
	payload := cfhtmlBuild(frag)
	startHTML := readField(t, payload, "StartHTML")
	endHTML := readField(t, payload, "EndHTML")
	startFrag := readField(t, payload, "StartFragment")
	endFrag := readField(t, payload, "EndFragment")
	if startFrag != int(startHTML)+strings.Index(string(payload[startHTML:endHTML]), "<!--StartFragment-->")+len("<!--StartFragment-->") {
		t.Fatalf("StartFragment %d does not follow the mark at %d", startFrag, startHTML)
	}
	if !bytes.Equal(payload[startFrag:endFrag], frag) {
		t.Fatalf("fragment bytes %q", payload[startFrag:endFrag])
	}
	if endHTML != endFrag+len("<!--EndFragment--></body></html>") {
		t.Fatalf("EndHTML %d, EndFragment %d", endHTML, endFrag)
	}
}

func readField(t *testing.T, payload []byte, name string) int {
	t.Helper()
	i := bytes.Index(payload, []byte(name+":"))
	if i < 0 {
		t.Fatalf("no %s", name)
	}
	n := 0
	for _, c := range payload[i+len(name)+1:] {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	if n == 0 {
		t.Fatalf("%s is 0", name)
	}
	return n
}

func TestDIBRoundTrip(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 3, 2))
	img.Set(0, 0, color.NRGBA{255, 0, 0, 255})
	img.Set(1, 0, color.NRGBA{0, 255, 0, 128})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	dib, ok := pngToDIB(buf.Bytes())
	if !ok {
		t.Fatal("pngToDIB failed")
	}
	back, ok := dibToPNG(dib)
	if !ok {
		t.Fatal("dibToPNG failed")
	}
	decoded, err := png.Decode(bytes.NewReader(back))
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, a := decoded.At(0, 0).RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Fatalf("red pixel: %d %d %d", r>>8, g>>8, b>>8)
	}
	if a>>8 == 0 {
		t.Fatal("alpha lost")
	}
}

func TestDIBToPNGFormats(t *testing.T) {
	// 24-bpp bottom-up BI_RGB.
	w, h := 2, 2
	stride := (w*3 + 3) / 4 * 4
	dib := make([]byte, 40+stride*h)
	binary.LittleEndian.PutUint32(dib[0:4], 40)
	binary.LittleEndian.PutUint32(dib[4:8], uint32(w))
	binary.LittleEndian.PutUint32(dib[8:12], uint32(h)) // positive: bottom-up
	binary.LittleEndian.PutUint16(dib[12:14], 1)
	binary.LittleEndian.PutUint16(dib[14:16], 24)
	// First stored row is the image's last: make it pure blue.
	for x := 0; x < w; x++ {
		copy(dib[40+x*3:], []byte{255, 0, 0})
	}
	got, ok := dibToPNG(dib)
	if !ok {
		t.Fatal("24-bpp DIB rejected")
	}
	img, err := png.Decode(bytes.NewReader(got))
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, _ := img.At(0, h-1).RGBA() // bottom-up: the first stored row is the image's last
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 255 {
		t.Fatalf("bottom row should be blue: %d %d %d", r>>8, g>>8, b>>8)
	}

	// 32-bpp BI_BITFIELDS with masks.
	dib32 := make([]byte, 40+12+w*4*h)
	binary.LittleEndian.PutUint32(dib32[0:4], 40)
	binary.LittleEndian.PutUint32(dib32[4:8], uint32(w))
	binary.LittleEndian.PutUint32(dib32[8:12], uint32(h))
	binary.LittleEndian.PutUint16(dib32[12:14], 1)
	binary.LittleEndian.PutUint16(dib32[14:16], 32)
	binary.LittleEndian.PutUint32(dib32[16:20], 3) // BI_BITFIELDS
	binary.LittleEndian.PutUint32(dib32[40:44], 0x00ff0000)
	binary.LittleEndian.PutUint32(dib32[44:48], 0x0000ff00)
	binary.LittleEndian.PutUint32(dib32[48:52], 0x000000ff)
	if got, ok := dibToPNG(dib32); !ok {
		t.Fatal("bitfields DIB rejected")
	} else if _, err := png.Decode(bytes.NewReader(got)); err != nil {
		t.Fatal(err)
	}

	// Top-down negative height.
	binary.LittleEndian.PutUint32(dib32[16:20], 0)
	binary.LittleEndian.PutUint32(dib32[8:12], uint32(int32(-h)))
	if got, ok := dibToPNG(dib32); !ok {
		t.Fatal("top-down DIB rejected")
	} else if _, err := png.Decode(bytes.NewReader(got)); err != nil {
		t.Fatal(err)
	}

	// Junk is refused.
	for _, junk := range [][]byte{nil, {1, 2, 3}, make([]byte, 40)} {
		if junk != nil {
			if _, ok := dibToPNG(junk); ok {
				t.Fatalf("junk accepted: % x", junk[:8])
			}
		}
	}
}

func TestDropFiles(t *testing.T) {
	paths := []string{`C:\Users\lucas\a.txt`, `D:\data\my clip.png`}
	payload := dropFilesBuild(paths)
	if got := dropFiles(payload); !slices.Equal(got, paths) {
		t.Fatalf("wide: got %q", got)
	}
	// The payload parses as wide from its header.
	if binary.LittleEndian.Uint32(payload[16:20]) != 1 {
		t.Fatal("fWide not set")
	}
	// A narrow (ANSI) payload parses too.
	list := []byte("C:\\a.txt\x00C:\\b.txt\x00\x00")
	narrow := make([]byte, 20+len(list))
	binary.LittleEndian.PutUint32(narrow[0:4], 20)
	copy(narrow[20:], list)
	if got := dropFiles(narrow); !slices.Equal(got, []string{`C:\a.txt`, `C:\b.txt`}) {
		t.Fatalf("narrow: got %q", got)
	}
	// Junk is refused.
	for _, junk := range [][]byte{nil, {1, 2, 3}, make([]byte, 20)} {
		if junk != nil && len(junk) >= 20 {
			binary.LittleEndian.PutUint32(junk[0:4], 999)
		}
		if got := dropFiles(junk); got != nil {
			t.Fatalf("junk accepted: %q", got)
		}
	}
}

func TestUTF16Helpers(t *testing.T) {
	for _, s := range []string{"hello", "héllo wörld ✓", "", "emoji 🎉 clip"} {
		if got := utf16BytesToString(utf16StringToBytes(s)); got != s {
			t.Fatalf("round trip %q: got %q", s, got)
		}
	}
	if b := utf16StringToBytes("hi"); !slices.Equal(b, []byte{'h', 0, 'i', 0, 0, 0}) {
		t.Fatalf("no NUL terminator: % x", b)
	}
	if s := utf16BytesToString([]byte{'h', 0, 'i', 0, 0, 0}); s != "hi" {
		t.Fatalf("terminator kept: %q", s)
	}
}

func TestPowerShellScriptEncoding(t *testing.T) {
	raw, err := base64.StdEncoding.DecodeString(powershellCommand())
	if err != nil {
		t.Fatal(err)
	}
	u := make([]uint16, len(raw)/2)
	for i := range u {
		u[i] = uint16(raw[i*2]) | uint16(raw[i*2+1])<<8
	}
	if got := string(utf16.Decode(u)); got != ocrScript {
		t.Fatal("round trip through -EncodedCommand encoding changed the script")
	}
	// The script must not contain a raw backtick: it would break the Go
	// raw string and PowerShell's single-quoted strings alike.
	if strings.Contains(ocrScript, "`") {
		t.Fatal("script contains a backtick")
	}
}
