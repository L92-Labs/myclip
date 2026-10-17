package prefs

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
	"github.com/L92-Labs/myclip/internal/settings"
)

type fakeHost struct {
	changed   int
	recording []bool
	pins      []*history.Item
}

func (h *fakeHost) Changed()                     { h.changed++ }
func (h *fakeHost) Recording(on bool)            { h.recording = append(h.recording, on) }
func (h *fakeHost) PinnedItems() []*history.Item { return h.pins }
func (h *fakeHost) PinKeys(int64) []string       { return []string{"b", "c"} }
func (h *fakeHost) SetPin(int64, string) error   { return nil }
func (h *fakeHost) SetTitle(int64, string) error { return nil }
func (h *fakeHost) SetText(int64, string) error  { return nil }
func (h *fakeHost) Unpin(int64)                  {}
func (h *fakeHost) LaunchAtLogin() bool          { return false }
func (h *fakeHost) SetLaunchAtLogin(bool) error  { return nil }
func (h *fakeHost) AccessibilityTrusted() bool   { return false }
func (h *fakeHost) RequestAccessibility()        {}
func (h *fakeHost) StorageSize() string          { return "1.0 MB" }
func (h *fakeHost) PickApp()                     {}
func (h *fakeHost) AppName(id string) string     { return id }
func (h *fakeHost) AppIcon(string) *ui.Bitmap    { return nil }
func (h *fakeHost) Displays() []string           { return []string{"Built-in"} }

func setup(t *testing.T) (*View, *fakeHost, *ui.Tester) {
	t.Helper()
	s := settings.Default()
	s.IgnoredApps = []string{"com.example.secret"}
	h := &fakeHost{}
	pinned := history.NewItem([]history.Content{{Type: pasteboard.TypeString, Value: []byte("pinned text")}}, "", time.Now())
	pinned.ID, pinned.Pin = 1, "b"
	h.pins = []*history.Item{pinned}
	v := &View{Host: h, S: &s}
	v.Load()
	tt := ui.NewTester(v.View, 680, 900)
	tt.Frame()
	return v, h, tt
}

func TestEveryPaneBuilds(t *testing.T) {
	v, _, tt := setup(t)
	want := []string{"Launch at login", "Size:", "Popup at:", "pinned text", "com.example.secret", "Clear history on quit"}
	for tab, text := range want {
		v.SetTab(tab)
		tt.Frame()
		if !tt.HasText(text) {
			t.Errorf("tab %d: no %q in %q", tab, text, tt.Texts())
		}
	}
	v.SetTab(4)
	for i, text := range []string{"Pasteboard generator type", "Regular expression"} {
		v.ignoreTab = i + 1
		tt.Frame()
		if !tt.HasText(text) && !slices.Contains(tt.Texts(), text) {
			t.Errorf("ignore tab %d: no %q", i+1, text)
		}
	}
}

func TestRecordShortcut(t *testing.T) {
	v, h, tt := setup(t)
	if err := tt.Click("⇧⌘C"); err != nil {
		t.Fatal(err)
	}
	if v.recording != "open" || len(h.recording) != 1 || !h.recording[0] {
		t.Fatalf("not recording: %q %v", v.recording, h.recording)
	}
	tt.Key(0, ui.KeyA)               // no modifier: ignored
	tt.Key(ui.Super|ui.Alt, ui.KeyV) // ⌘⌥V
	if v.S.PopupShortcut != "Cmd+Alt+V" || v.recording != "" || h.changed != 1 {
		t.Fatalf("shortcut %q recording %q changed %d", v.S.PopupShortcut, v.recording, h.changed)
	}
	if !tt.HasText("⌥⌘V") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestStorageTypes(t *testing.T) {
	v, h, tt := setup(t)
	v.SetTab(1)
	tt.Frame()
	if err := tt.Click("Images"); err != nil {
		t.Fatal(err)
	}
	if slices.Contains(v.S.EnabledPasteboardTypes, pasteboard.TypePNG) || h.changed != 1 {
		t.Fatalf("types %v changed %d", v.S.EnabledPasteboardTypes, h.changed)
	}
	if !slices.Contains(v.S.EnabledPasteboardTypes, pasteboard.TypeString) {
		t.Fatal("lost text types")
	}
}
