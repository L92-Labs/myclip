package shortcut

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestParseRoundTrip(t *testing.T) {
	for _, s := range []string{"Cmd+Shift+C", "Alt+P", "Alt+Backspace", "Ctrl+Space", "Cmd+Ctrl+Alt+Shift+F5", "Cmd+,"} {
		c, ok := Parse(s)
		if !ok {
			t.Fatalf("%q did not parse", s)
		}
		if got := c.String(); got != s {
			t.Errorf("%q → %q", s, got)
		}
	}
}

func TestParseAliasesAndErrors(t *testing.T) {
	c, ok := Parse("CmdOrCtrl+option+p")
	if !ok || c.Mods != ui.Cmd|ui.Alt || c.Key != ui.KeyP {
		t.Fatalf("got %+v %v", c, ok)
	}
	for _, s := range []string{"", "Hyper+P", "Cmd+Nope"} {
		if _, ok := Parse(s); ok {
			t.Errorf("%q parsed", s)
		}
	}
}

func TestSymbols(t *testing.T) {
	for in, want := range map[string]string{
		"Cmd+Shift+C":   "⇧⌘C",
		"Alt+Backspace": "⌥⌫",
		"Ctrl+Space":    "⌃Space",
		"bad+":          "",
	} {
		if got := Symbols(in); got != want {
			t.Errorf("%q: %q want %q", in, got, want)
		}
	}
}
