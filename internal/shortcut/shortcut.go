// Package shortcut converts between MyGo accelerators ("Cmd+Shift+C"), the
// keys of native UI, and the symbols macOS shows (⇧⌘C).
package shortcut

import (
	"strings"

	"github.com/egoist/mygo/ui"
)

// Chord is a key with exactly these modifiers.
type Chord struct {
	Mods ui.Modifiers
	Key  ui.Key
}

var names = map[string]ui.Key{
	"backspace": ui.KeyBackspace, "delete": ui.KeyDelete, "del": ui.KeyDelete,
	"space": ui.KeySpace, "enter": ui.KeyEnter, "return": ui.KeyEnter,
	"escape": ui.KeyEscape, "esc": ui.KeyEscape, "tab": ui.KeyTab,
	"up": ui.KeyUp, "down": ui.KeyDown, "left": ui.KeyLeft, "right": ui.KeyRight,
	"home": ui.KeyHome, "end": ui.KeyEnd, "pageup": ui.KeyPageUp, "pagedown": ui.KeyPageDown,
	",": ui.KeyComma, ".": ui.KeyPeriod, "/": ui.KeySlash, ";": ui.KeySemicolon,
	"'": ui.KeyQuote, "[": ui.KeyBracketLeft, "]": ui.KeyBracketRight,
	"\\": ui.KeyBackslash, "`": ui.KeyBackquote, "-": ui.KeyMinus, "=": ui.KeyEqual,
	"f1": ui.KeyF1, "f2": ui.KeyF2, "f3": ui.KeyF3, "f4": ui.KeyF4, "f5": ui.KeyF5, "f6": ui.KeyF6,
	"f7": ui.KeyF7, "f8": ui.KeyF8, "f9": ui.KeyF9, "f10": ui.KeyF10, "f11": ui.KeyF11, "f12": ui.KeyF12,
}

// accelerator names, the first of each key as Format writes it.
var accel = map[ui.Key]string{
	ui.KeyBackspace: "Backspace", ui.KeyDelete: "Delete", ui.KeySpace: "Space",
	ui.KeyEnter: "Enter", ui.KeyEscape: "Escape", ui.KeyTab: "Tab",
	ui.KeyUp: "Up", ui.KeyDown: "Down", ui.KeyLeft: "Left", ui.KeyRight: "Right",
	ui.KeyHome: "Home", ui.KeyEnd: "End", ui.KeyPageUp: "PageUp", ui.KeyPageDown: "PageDown",
	ui.KeyComma: ",", ui.KeyPeriod: ".", ui.KeySlash: "/", ui.KeySemicolon: ";",
	ui.KeyQuote: "'", ui.KeyBracketLeft: "[", ui.KeyBracketRight: "]",
	ui.KeyBackslash: "\\", ui.KeyBackquote: "`", ui.KeyMinus: "-", ui.KeyEqual: "=",
	ui.KeyF1: "F1", ui.KeyF2: "F2", ui.KeyF3: "F3", ui.KeyF4: "F4", ui.KeyF5: "F5", ui.KeyF6: "F6",
	ui.KeyF7: "F7", ui.KeyF8: "F8", ui.KeyF9: "F9", ui.KeyF10: "F10", ui.KeyF11: "F11", ui.KeyF12: "F12",
}

var symbols = map[ui.Key]string{
	ui.KeyBackspace: "⌫", ui.KeyDelete: "⌦", ui.KeySpace: "Space", ui.KeyEnter: "↩",
	ui.KeyEscape: "⎋", ui.KeyTab: "⇥", ui.KeyUp: "↑", ui.KeyDown: "↓", ui.KeyLeft: "←",
	ui.KeyRight: "→", ui.KeyHome: "↖", ui.KeyEnd: "↘", ui.KeyPageUp: "⇞", ui.KeyPageDown: "⇟",
}

// Letter returns the key typing a letter or digit.
func Letter(r rune) (ui.Key, bool) {
	switch {
	case r >= 'a' && r <= 'z':
		return ui.KeyA + ui.Key(r-'a'), true
	case r >= 'A' && r <= 'Z':
		return ui.KeyA + ui.Key(r-'A'), true
	case r >= '0' && r <= '9':
		return ui.Key0 + ui.Key(r-'0'), true
	}
	return 0, false
}

// Parse reads an accelerator. Unknown keys fail.
func Parse(s string) (Chord, bool) {
	var ch Chord
	parts := strings.Split(s, "+")
	if len(parts) == 0 || s == "" {
		return ch, false
	}
	for _, p := range parts[:len(parts)-1] {
		switch strings.ToLower(p) {
		// ⌘ is Super: ui.Cmd is Super on macOS but Ctrl elsewhere, where
		// "Cmd" and "Ctrl" would collapse into one modifier.
		case "cmd", "command", "super", "meta":
			ch.Mods |= ui.Super
		case "cmdorctrl", "commandorcontrol":
			ch.Mods |= ui.Cmd
		case "ctrl", "control":
			ch.Mods |= ui.Ctrl
		case "alt", "option":
			ch.Mods |= ui.Alt
		case "shift":
			ch.Mods |= ui.Shift
		default:
			return ch, false
		}
	}
	last := parts[len(parts)-1]
	if last == "" { // "Cmd++"
		last = "+"
	}
	if k, ok := names[strings.ToLower(last)]; ok {
		ch.Key = k
		return ch, true
	}
	if r := []rune(last); len(r) == 1 {
		if k, ok := Letter(r[0]); ok {
			ch.Key = k
			return ch, true
		}
	}
	return ch, false
}

// String writes c as an accelerator.
func (c Chord) String() string {
	var b strings.Builder
	for _, m := range []struct {
		mod  ui.Modifiers
		name string
	}{{ui.Super, "Cmd"}, {ui.Ctrl, "Ctrl"}, {ui.Alt, "Alt"}, {ui.Shift, "Shift"}} {
		if c.Mods&m.mod != 0 {
			b.WriteString(m.name + "+")
		}
	}
	b.WriteString(keyName(c.Key, accel))
	return b.String()
}

// Symbols writes c as macOS menus do: ⌃⌥⇧⌘ then the key.
func (c Chord) Symbols() string {
	return Mods(c.Mods) + keyName(c.Key, symbols)
}

// Mods writes modifiers as macOS menus do.
func Mods(m ui.Modifiers) string {
	var b strings.Builder
	for _, x := range []struct {
		mod ui.Modifiers
		sym string
	}{{ui.Ctrl, "⌃"}, {ui.Alt, "⌥"}, {ui.Shift, "⇧"}, {ui.Super, "⌘"}} {
		if m&x.mod != 0 {
			b.WriteString(x.sym)
		}
	}
	return b.String()
}

func keyName(k ui.Key, table map[ui.Key]string) string {
	if s, ok := table[k]; ok {
		return s
	}
	if s, ok := accel[k]; ok {
		return s
	}
	switch {
	case k >= ui.KeyA && k <= ui.KeyZ:
		return string(rune('A' + k - ui.KeyA))
	case k >= ui.Key0 && k <= ui.Key9:
		return string(rune('0' + k - ui.Key0))
	}
	return "?"
}

// Symbols formats an accelerator for display, or "" when it is invalid.
func Symbols(accelerator string) string {
	c, ok := Parse(accelerator)
	if !ok {
		return ""
	}
	return c.Symbols()
}

// KeyRune returns the lowercase letter of a letter key, or 0.
func KeyRune(k ui.Key) rune {
	if k >= ui.KeyA && k <= ui.KeyZ {
		return rune('a' + k - ui.KeyA)
	}
	return 0
}
