package main

import (
	"testing"

	"github.com/L92-Labs/myclip/internal/native"
)

func TestFlags(t *testing.T) {
	if got := flags("Cmd+Shift+C"); got != native.ModCommand|native.ModShift {
		t.Fatalf("%x", got)
	}
	if got := flags("Ctrl+Alt+V"); got != native.ModControl|native.ModOption {
		t.Fatalf("%x", got)
	}
	if flags("nonsense") != 0 {
		t.Fatal("parsed nonsense")
	}
}

func TestHoldToCycle(t *testing.T) {
	var held uint
	c := cycler{held: func() uint { return held }}
	var mods uint = native.ModCommand | native.ModShift

	// Press, press again while holding, release: select.
	held = mods
	if !c.start(mods) {
		t.Fatal("hold not started")
	}
	if done, _ := c.poll(); done {
		t.Fatal("done while held")
	}
	if !c.press() {
		t.Fatal("second press not taken")
	}
	held = 0
	if done, sel := c.poll(); !done || !sel {
		t.Fatalf("release: done=%v select=%v", done, sel)
	}

	// Press and release at once: the popup stays open.
	held = mods
	c.start(mods)
	held = native.ModShift // ⌘ let go first
	if done, sel := c.poll(); !done || sel {
		t.Fatalf("quick press: done=%v select=%v", done, sel)
	}
	if c.press() {
		t.Fatal("a press after the hold ended cycled")
	}

	// No modifier state (Linux): no hold at all.
	held = 0
	if c.start(mods) {
		t.Fatal("hold started without modifiers")
	}
}
