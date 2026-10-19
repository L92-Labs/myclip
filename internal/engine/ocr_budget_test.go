package engine

import (
	"context"
	"testing"
	"unicode/utf8"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

func TestRecognizedTextRespectsItemByteLimit(t *testing.T) {
	for _, tc := range []struct {
		name, html, reading, want string
	}{
		{"ASCII", "aa", "abcdef", "abc"},
		{"UTF8", "aa", "aéz", "aé"},
		{"partial rune", "aaa", "界", ""},
		{"full item", "aaaaa", "words", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newRig(t, Limits{})
			r.e.cap.MaxItemBytes = 8
			r.f.SetItems([]pasteboard.Item{{Types: map[string][]byte{
				pasteboard.TypePNG:  []byte("png"),
				pasteboard.TypeHTML: []byte(tc.html),
			}}}, nil)
			if got := r.e.Tick(context.Background()); got != capture.Captured {
				t.Fatalf("capture: %v", got)
			}
			id := r.e.Items()[0].ID
			digest := ImageDigest(r.e.Item(id))
			// A replacement reading can use the space occupied by the old one.
			for _, reading := range []string{"x", tc.reading} {
				if err := r.e.SetRecognizedText(id, digest, reading); err != nil {
					t.Fatal(err)
				}
			}
			it := r.e.Item(id)
			if it.Size() > 8 || it.RecognizedText() != tc.want || !utf8.ValidString(it.RecognizedText()) || !it.TextRecognized() {
				t.Fatalf("size=%d, OCR=%q, attempted=%v; want size<=8, OCR=%q, attempted=true", it.Size(), it.RecognizedText(), it.TextRecognized(), tc.want)
			}
			stored, err := r.s.Load(context.Background())
			if err != nil || len(stored) != 1 {
				t.Fatalf("load: %v, %d items", err, len(stored))
			}
			if stored[0].Size() > 8 || stored[0].RecognizedText() != tc.want || !stored[0].TextRecognized() {
				t.Fatalf("persisted size=%d, OCR=%q, attempted=%v", stored[0].Size(), stored[0].RecognizedText(), stored[0].TextRecognized())
			}
		})
	}
}
