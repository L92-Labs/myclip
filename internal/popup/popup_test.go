package popup

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
	"github.com/L92-Labs/myclip/internal/settings"
)

type call struct {
	what string
	id   int64
	a    Action
}

type fakeHost struct {
	calls  []call
	width  int
	height int
}

func (h *fakeHost) Select(id int64, a Action) { h.calls = append(h.calls, call{"select", id, a}) }
func (h *fakeHost) TogglePin(id int64)        { h.calls = append(h.calls, call{what: "pin", id: id}) }
func (h *fakeHost) Delete(id int64)           { h.calls = append(h.calls, call{what: "delete", id: id}) }
func (h *fakeHost) Clear(all bool) {
	h.calls = append(h.calls, call{what: fmt.Sprint("clear ", all)})
}
func (h *fakeHost) OpenSettings()                           { h.calls = append(h.calls, call{what: "settings"}) }
func (h *fakeHost) About()                                  { h.calls = append(h.calls, call{what: "about"}) }
func (h *fakeHost) Quit()                                   { h.calls = append(h.calls, call{what: "quit"}) }
func (h *fakeHost) Close()                                  { h.calls = append(h.calls, call{what: "close"}) }
func (h *fakeHost) Thumbnail(*history.Item, int) *ui.Bitmap { return nil }
func (h *fakeHost) Preview(*history.Item) *ui.Bitmap        { return nil }
func (h *fakeHost) AppIcon(string) *ui.Bitmap               { return nil }
func (h *fakeHost) AppName(id string) string                { return id }
func (h *fakeHost) Resize(w, hh int)                        { h.width, h.height = w, hh }
func (h *fakeHost) last() call {
	if len(h.calls) == 0 {
		return call{}
	}
	return h.calls[len(h.calls)-1]
}

func item(id int64, text, pin string) *history.Item {
	it := history.NewItem([]history.Content{{Type: pasteboard.TypeString, Value: []byte(text)}}, "com.apple.TextEdit", time.Unix(id, 0))
	it.ID, it.Pin = id, pin
	return it
}

func setup(t *testing.T, mutate func(*settings.Settings)) (*Popup, *fakeHost, *ui.Tester) {
	t.Helper()
	s := settings.Default()
	s.OpenPreviewAutomatically = false
	if mutate != nil {
		mutate(&s)
	}
	h := &fakeHost{}
	p := &Popup{Host: h, Settings: &s}
	p.SetItems([]*history.Item{
		item(9, "pinned thing", "b"),
		item(1, "apple pie", ""),
		item(2, "banana split", ""),
		item(3, "cherry tart", ""),
	})
	p.Reset()
	tt := ui.NewTester(p.View, 450, 600)
	tt.Frame()
	return p, h, tt
}

func TestResolveMatrix(t *testing.T) {
	copyOnly, paste, plain := Action{}, Action{Paste: true}, Action{Paste: true, Plain: true}
	for _, tc := range []struct {
		mods     ui.Modifiers
		pbd, rfd bool
		want     Action
		ok       bool
	}{
		{0, false, false, copyOnly, true},
		{0, true, true, plain, true},
		{ui.Cmd, false, false, copyOnly, true},
		{ui.Cmd, true, false, paste, true},
		{ui.Cmd, true, true, plain, true},
		{ui.Alt, false, false, paste, true},
		{ui.Alt, false, true, plain, true},
		{ui.Alt, true, false, copyOnly, true},
		{ui.Alt | ui.Shift, false, false, plain, true},
		{ui.Alt | ui.Shift, false, true, paste, true},
		{ui.Alt | ui.Shift, true, false, Action{}, false},
		{ui.Cmd | ui.Shift, true, false, plain, true},
		{ui.Cmd | ui.Shift, true, true, paste, true},
		{ui.Cmd | ui.Shift, false, false, Action{}, false},
		{ui.Ctrl, false, false, Action{}, false},
	} {
		got, ok := Resolve(tc.mods, tc.pbd, tc.rfd)
		if got != tc.want || ok != tc.ok {
			t.Errorf("%v pbd=%v rfd=%v: %+v %v, want %+v %v", tc.mods, tc.pbd, tc.rfd, got, ok, tc.want, tc.ok)
		}
	}
}

func TestFirstUnpinnedChosenAndEnterCopies(t *testing.T) {
	_, h, tt := setup(t, nil)
	tt.Key(0, ui.KeyEnter)
	if got := h.last(); got.what != "select" || got.id != 1 || got.a != (Action{}) {
		t.Fatalf("got %+v", got)
	}
}

func TestNavigationKeys(t *testing.T) {
	p, h, tt := setup(t, nil)
	tt.Key(0, ui.KeyDown)
	tt.Key(ui.Ctrl, ui.KeyN)
	if p.chosen != 3 {
		t.Fatalf("chosen %d after Down, ⌃N", p.chosen)
	}
	tt.Key(ui.Ctrl, ui.KeyP)
	tt.Key(ui.Cmd, ui.KeyUp)
	if p.chosen != 0 {
		t.Fatalf("chosen %d after ⌃P, ⌘↑", p.chosen)
	}
	tt.Key(ui.Alt, ui.KeyDown)
	tt.Key(ui.Alt, ui.KeyEnter)
	if got := h.last(); got.id != 3 || got.a != (Action{Paste: true}) {
		t.Fatalf("⌥↩ %+v", got)
	}
}

func TestNumberAndPinChords(t *testing.T) {
	_, h, tt := setup(t, nil)
	tt.Key(ui.Cmd, ui.Key2)
	if got := h.last(); got.id != 2 || got.a != (Action{}) {
		t.Fatalf("⌘2 %+v", got)
	}
	tt.Key(ui.Alt|ui.Shift, ui.Key1)
	if got := h.last(); got.id != 1 || got.a != (Action{Paste: true, Plain: true}) {
		t.Fatalf("⌥⇧1 %+v", got)
	}
	tt.Key(ui.Cmd, ui.KeyB)
	if got := h.last(); got.id != 9 {
		t.Fatalf("⌘B (pin) %+v", got)
	}
	if !tt.HasText("⌘1") || !tt.HasText("⌘B") {
		t.Fatalf("no key labels in %q", tt.Texts())
	}
}

func TestTypingFilters(t *testing.T) {
	p, h, tt := setup(t, nil)
	tt.Type("ban")
	if p.Query() != "ban" || len(p.rows) != 1 || p.rows[0].it.ID != 2 {
		t.Fatalf("query %q rows %d", p.Query(), len(p.rows))
	}
	tt.Key(0, ui.KeyEnter)
	if got := h.last(); got.id != 2 {
		t.Fatalf("enter after filter %+v", got)
	}
	tt.Key(ui.Ctrl, ui.KeyU)
	if p.Query() != "" || len(p.rows) != 4 {
		t.Fatalf("⌃U left %q", p.Query())
	}
}

func TestTypingWhileSearchHidden(t *testing.T) {
	p, _, tt := setup(t, func(s *settings.Settings) { s.SearchVisibility = settings.SearchDuringSearch })
	tt.Type("che")
	if p.Query() != "che" || len(p.rows) != 1 {
		t.Fatalf("query %q rows %d", p.Query(), len(p.rows))
	}
}

func TestPinDeleteAndCloseChords(t *testing.T) {
	p, h, tt := setup(t, nil)
	tt.Key(ui.Alt, ui.KeyP)
	if got := h.last(); got.what != "pin" || got.id != 1 {
		t.Fatalf("⌥P %+v", got)
	}
	if p.Query() != "" {
		t.Fatalf("⌥P typed %q", p.Query())
	}
	tt.Key(ui.Alt, ui.KeyBackspace)
	if got := h.last(); got.what != "delete" || got.id != 1 {
		t.Fatalf("⌥⌫ %+v", got)
	}
	tt.Key(0, ui.KeyEscape)
	if got := h.last(); got.what != "close" {
		t.Fatalf("esc %+v", got)
	}
	tt.Key(ui.Cmd, ui.KeyComma)
	if got := h.last(); got.what != "settings" {
		t.Fatalf("⌘, %+v", got)
	}
}

func TestClearAsksFirst(t *testing.T) {
	p, h, tt := setup(t, nil)
	tt.Key(ui.Cmd|ui.Alt, ui.KeyBackspace)
	if !p.confirm || len(h.calls) != 0 {
		t.Fatalf("no confirmation: %+v", h.calls)
	}
	if !tt.HasText("Don't ask again") {
		t.Fatalf("no dialog in %q", tt.Texts())
	}
	p.dontAsk = true
	tt.Key(0, ui.KeyEnter) // the focused Clear button
	if got := h.last(); got.what != "clear false" || !p.Settings.SuppressClearAlert {
		t.Fatalf("after confirm %+v suppress=%v", got, p.Settings.SuppressClearAlert)
	}
	tt.Key(ui.Cmd|ui.Alt|ui.Shift, ui.KeyBackspace)
	if got := h.last(); got.what != "clear true" || p.confirm {
		t.Fatalf("suppressed clear all %+v", got)
	}
}

func TestPreviewToggleAndAuto(t *testing.T) {
	p, h, tt := setup(t, nil)
	tt.Key(ui.Ctrl, ui.KeySpace)
	if !p.previewOpen || !tt.HasText("Number of copies:") {
		t.Fatalf("no preview: %q", tt.Texts())
	}
	if h.width != 450+400+1 {
		t.Fatalf("width %d", h.width)
	}
	tt.Key(ui.Ctrl, ui.KeySpace)
	if p.previewOpen {
		t.Fatal("still open")
	}

	p, _, tt = setup(t, func(s *settings.Settings) { s.OpenPreviewAutomatically = true })
	tt.Advance(time.Second)
	if p.previewOpen {
		t.Fatal("opened before the delay")
	}
	tt.Advance(600 * time.Millisecond)
	if !p.previewOpen {
		t.Fatal("did not open after the delay")
	}
}

func TestHighlightAndSpecialSymbols(t *testing.T) {
	p, _, tt := setup(t, nil)
	p.SetItems([]*history.Item{item(5, " a\tb\n", "")})
	tt.Frame()
	if !tt.HasText("·a⇥b⏎") {
		t.Fatalf("texts %q", tt.Texts())
	}
}
