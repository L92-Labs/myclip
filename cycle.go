package main

import (
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/L92-Labs/myclip/internal/native"
	"github.com/L92-Labs/myclip/internal/shortcut"
)

// cycler is Maccy's hold-to-cycle: the popup shortcut opens the popup; with
// its modifiers still held, pressing its key again moves the choice down,
// and releasing the modifiers then selects the item chosen. Releasing them
// without a second press leaves the popup open, as a plain press does.
type cycler struct {
	held   func() uint // the modifier flags held now (NSEvent.modifierFlags)
	mods   uint        // the shortcut's modifiers, as flags
	active bool        // the popup was opened by the shortcut, modifiers held
	cycled bool        // the key was pressed again in this hold
}

// start begins a hold after the shortcut opened the popup; it reports
// whether the modifiers are held (only macOS reports them).
func (c *cycler) start(mods uint) bool {
	c.mods, c.cycled = mods, false
	c.active = mods != 0 && c.held()&mods == mods
	return c.active
}

// press is the shortcut pressed while the popup shows: true when it
// belongs to the hold and should move the choice down.
func (c *cycler) press() bool {
	if !c.active {
		return false
	}
	c.cycled = true
	return true
}

// poll checks the modifiers: done once released, with selectNow when the
// hold cycled.
func (c *cycler) poll() (done, selectNow bool) {
	if !c.active {
		return true, false
	}
	if c.held()&c.mods == c.mods {
		return false, false
	}
	c.active = false
	return true, c.cycled
}

func (c *cycler) stop() { c.active = false }

// flags turns the shortcut's modifiers into NSEvent modifier flags.
func flags(accel string) uint {
	ch, ok := shortcut.Parse(accel)
	if !ok {
		return 0
	}
	var f uint
	for m, flag := range map[ui.Modifiers]uint{ui.Shift: native.ModShift, ui.Ctrl: native.ModControl, ui.Alt: native.ModOption, ui.Super: native.ModCommand} {
		if ch.Mods&m != 0 {
			f |= flag
		}
	}
	return f
}

// shortcutPressed is the popup shortcut's handler.
func (a *app) shortcutPressed() {
	if a.panel.IsVisible() {
		if a.cycle.press() {
			a.popup.ChooseNext()
			a.panel.Invalidate()
			return
		}
		a.Close()
		return
	}
	a.showPopup()
	if a.cycle.start(flags(a.settings.PopupShortcut)) {
		go a.watchHold()
	}
}

// watchHold polls the modifiers while a hold lasts (key-up events do not
// reach a global shortcut), every 15 ms on the main thread.
func (a *app) watchHold() {
	for {
		time.Sleep(15 * time.Millisecond)
		var done, sel bool
		mygo.RunOnMain(func() {
			if !a.panel.IsVisible() {
				a.cycle.stop()
			}
			done, sel = a.cycle.poll()
			if sel {
				a.popup.SelectChosen()
			}
		})
		if done {
			return
		}
	}
}
