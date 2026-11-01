package transfer

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
	"github.com/L92-Labs/myclip/internal/settings"
)

func sample(now time.Time) []*history.Item {
	return []*history.Item{
		{Application: "Notes", FirstCopiedAt: now, LastCopiedAt: now, NumberOfCopies: 2,
			Pin: "b", Title: "renamed",
			Contents: []history.Content{
				{Type: pasteboard.TypeString, Value: []byte("hello")},
				{Type: pasteboard.TypePNG, Value: []byte{0x89, 'P', 'N', 'G'}},
			}},
		{FirstCopiedAt: now, LastCopiedAt: now, NumberOfCopies: 1,
			Contents: []history.Content{{Type: pasteboard.TypeString, Value: []byte("ünïcode ✓")}}},
	}
}

func TestExportImportRoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	var buf bytes.Buffer
	def := settings.Default()
	if err := Export(&buf, &def, sample(now)); err != nil {
		t.Fatal(err)
	}
	a, err := Read(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if a.Format != FormatVersion || a.Exported.IsZero() {
		t.Fatalf("header: %d %v", a.Format, a.Exported)
	}
	if a.Settings == nil || a.Settings.SearchMode != settings.SearchExact {
		t.Fatal("settings lost")
	}
	items := a.HistoryItems()
	if len(items) != 2 {
		t.Fatalf("got %d items", len(items))
	}
	x := items[0]
	if x.Application != "Notes" || x.Pin != "b" || x.Title != "renamed" || x.NumberOfCopies != 2 {
		t.Fatalf("fields: %+v", x)
	}
	if len(x.Contents) != 2 || string(x.Contents[0].Value) != "hello" || string(x.Contents[1].Value) != "\x89PNG" {
		t.Fatalf("contents: %+v", x.Contents)
	}
	if string(items[1].Contents[0].Value) != "ünïcode ✓" {
		t.Fatalf("unicode: %q", items[1].Contents[0].Value)
	}
}

func TestReadRejectsJunk(t *testing.T) {
	for name, in := range map[string]string{
		"not json":     "hello",
		"no format":    `{"items":[]}`,
		"newer format": `{"myclip":99,"items":[]}`,
		"no clips":     `{"myclip":1}`,
		"empty list":   `{"myclip":1,"items":[]}`,
	} {
		if _, err := Read(bytes.NewReader([]byte(in))); err == nil {
			t.Fatalf("%s: accepted", name)
		}
	}
}

func TestReadSanitizesRows(t *testing.T) {
	in := `{"myclip":1,"items":[
		{"contents":[]},
		{"contents":[{"value":"aGVsbG8="}]},
		{"numberOfCopies":-4,"contents":[{"type":"public.utf8-plain-text","value":"aGk="}]},
		{"contents":[{"type":"","value":"eCJ="},{"type":"public.html","value":"PGI+"}]},
		{"pin":"b","contents":[{"type":"public.utf8-plain-text","value":"cGlraw=="}]}
	]}`
	a, err := Read(bytes.NewReader([]byte(in)))
	if err != nil {
		t.Fatal(err)
	}
	items := a.HistoryItems()
	if len(items) != 3 {
		t.Fatalf("got %d items: %+v", len(items), items)
	}
	// Countless and zero-timestamped rows are repaired.
	if items[0].NumberOfCopies != 1 || items[0].FirstCopiedAt.IsZero() || string(items[0].Contents[0].Value) != "hi" {
		t.Fatalf("not sanitized: %+v", items[0])
	}
	// Typeless contents are dropped, keeping the typed one.
	if len(items[1].Contents) != 1 || items[1].Contents[0].Type != "public.html" {
		t.Fatalf("typeless content kept: %+v", items[1])
	}
	if items[2].Pin != "b" {
		t.Fatalf("pin lost: %+v", items[2])
	}
}

func TestExportIsPlainJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, nil, sample(time.Now())); err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	if _, ok := v["settings"]; ok {
		t.Fatal("nil settings marshaled")
	}
}

// A settings-only backup — the app's own export after clearing history —
// must round-trip instead of being refused on read.
func TestSettingsOnlyBackupRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	def := settings.Default()
	if err := Export(&buf, &def, nil); err != nil {
		t.Fatal(err)
	}
	a, err := Read(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("settings-only backup rejected: %v", err)
	}
	if a.Settings == nil || a.Settings.SearchMode != settings.SearchExact || len(a.Items) != 0 {
		t.Fatalf("archive: %+v", a)
	}
	if items := a.HistoryItems(); len(items) != 0 {
		t.Fatalf("items out of an empty archive: %+v", items)
	}
	// History-only junk stays rejected.
	if _, err := Read(bytes.NewReader([]byte(`{"myclip":1,"items":[]}`))); err == nil {
		t.Fatal("empty history-only backup accepted")
	}
}

// The snippet flag survives export and import.
func TestSnippetRoundTrip(t *testing.T) {
	now := time.Now().UTC()
	items := []*history.Item{
		{FirstCopiedAt: now, LastCopiedAt: now, NumberOfCopies: 1, Snippet: true,
			Contents: []history.Content{{Type: pasteboard.TypeString, Value: []byte("on {date}")}}},
	}
	var buf bytes.Buffer
	if err := Export(&buf, nil, items); err != nil {
		t.Fatal(err)
	}
	a, err := Read(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	got := a.HistoryItems()
	if len(got) != 1 || !got[0].Snippet {
		t.Fatalf("snippet flag lost: %+v", got)
	}
}
