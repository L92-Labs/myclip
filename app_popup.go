package main

import (
	"log"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/native"
	"github.com/L92-Labs/myclip/internal/popup"
	"github.com/L92-Labs/myclip/internal/settings"
	"github.com/L92-Labs/myclip/internal/shortcut"
	"github.com/L92-Labs/myclip/locales"
)

// The popup's Host: the panel, the pasteboard and the history.

// registerShortcut registers the popup's global shortcut, replacing the
// one registered before.
func (a *app) registerShortcut() {
	if a.shortcut == a.settings.PopupShortcut {
		return
	}
	if a.shortcut != "" {
		mygo.GlobalShortcut.Unregister(a.shortcut)
		a.shortcut = ""
	}
	accel := a.settings.PopupShortcut
	if accel == "" {
		return
	}
	if err := mygo.GlobalShortcut.Register(accel, a.togglePopup); err != nil {
		log.Printf("myclip: shortcut %s: %v", accel, err)
		go mygo.Dialog.Message(mygo.MessageOptions{
			Type:    mygo.MessageWarning,
			Message: locales.ShortcutTaken(shortcut.Symbols(accel)),
		})
		return
	}
	a.shortcut = accel
}

func (a *app) togglePopup() {
	if a.panel.IsVisible() {
		a.Close()
		return
	}
	a.showPopup()
}

// showPopup opens the popup where the settings say (PopupPosition).
func (a *app) showPopup() {
	a.popup.SetItems(a.sortedItems())
	a.popup.Reset()
	w, h := a.popup.Fit()
	a.panel.SetSize(w, h)
	x, y := a.popupOrigin(w, h)
	a.panel.SetPosition(x, y)
	a.panel.Show()
}

// popupOrigin places a w×h popup (PopupPosition.origin @ eb03eba), kept
// on its screen.
func (a *app) popupOrigin(w, h int) (int, int) {
	s := a.settings
	cursor := mygo.Screen.CursorScreenPoint()
	area := a.popupScreen(cursor).WorkArea
	center := func(r mygo.Rectangle) (int, int) {
		return r.X + (r.Width-w)/2, r.Y + (r.Height-h)/2
	}
	var x, y int
	switch s.PopupPosition {
	case settings.PopupCenter:
		x, y = center(area)
	case settings.PopupWindow:
		if f, ok := native.FrontmostWindowFrame(); ok {
			r := mygo.Rectangle{X: int(f.X), Y: int(f.Y), Width: int(f.Width), Height: int(f.Height)}
			area = mygo.Screen.DisplayMatching(r).WorkArea
			x, y = center(r)
		} else {
			x, y = center(area)
		}
	case settings.PopupStatusItem:
		if a.tray != nil {
			if b := a.tray.Bounds(); b.Width > 0 {
				area = mygo.Screen.DisplayNearestPoint(mygo.Point{X: b.X, Y: b.Y}).WorkArea
				x, y = b.X, b.Y+b.Height
				break
			}
		}
		x, y = center(area)
	case settings.PopupLastPosition:
		// The anchor is the popup's top middle, as a fraction of the screen.
		x = area.X + int(float64(area.Width)*s.WindowPositionX) - w/2
		y = area.Y + int(float64(area.Height)*(1-s.WindowPositionY))
	default:
		x, y = cursor.X, cursor.Y
		area = mygo.Screen.DisplayNearestPoint(cursor).WorkArea
	}
	x = max(area.X, min(x, area.X+area.Width-w))
	y = max(area.Y, min(y, area.Y+area.Height-h))
	return x, y
}

// popupScreen is the display of popupScreen, 0 for the one with the
// pointer.
func (a *app) popupScreen(cursor mygo.Point) mygo.Display {
	if n := a.settings.PopupScreen; n > 0 {
		if ds := mygo.Screen.Displays(); n <= len(ds) {
			return ds[n-1]
		}
	}
	return mygo.Screen.DisplayNearestPoint(cursor)
}

// Close hides the popup, keeping where it was for lastPosition.
func (a *app) Close() {
	if !a.panel.IsVisible() {
		return
	}
	if a.settings.PopupPosition == settings.PopupLastPosition {
		b := a.panel.Bounds()
		area := mygo.Screen.DisplayMatching(b).WorkArea
		if area.Width > 0 && area.Height > 0 {
			a.settings.WindowPositionX = float64(b.X+b.Width/2-area.X) / float64(area.Width)
			a.settings.WindowPositionY = 1 - float64(b.Y-area.Y)/float64(area.Height)
			a.settings.Normalize()
		}
	}
	a.panel.Hide()
	a.saveSettings()
}

// Select copies an item back, and pastes it into the app the popup opened
// over: the panel never activated MyClip, so that app has the keyboard
// again once the panel hides.
func (a *app) Select(id int64, act popup.Action) {
	a.Close()
	if err := a.engine.CopyBack(id, act.Plain); err != nil {
		log.Print("myclip: copy: ", err)
		return
	}
	if !act.Paste {
		return
	}
	if !native.AccessibilityTrusted(true) {
		go mygo.Dialog.Message(mygo.MessageOptions{Type: mygo.MessageInfo, Message: locales.PasteNeedsAccessibility()})
		return
	}
	if native.SecureInputEnabled() {
		go mygo.Dialog.Message(mygo.MessageOptions{Type: mygo.MessageInfo, Message: locales.PasteSecureInput()})
		return
	}
	// Let the pasteboard write and the focus settle before ⌘V.
	time.AfterFunc(60*time.Millisecond, func() {
		mygo.RunOnMain(func() {
			if err := native.Paste(); err != nil {
				log.Print("myclip: paste: ", err)
			}
		})
	})
}

// pinReserved are the letters of the in-popup shortcuts, which pins
// cannot take (HistoryItem.supportedPins).
func (a *app) pinReserved() []string {
	var out []string
	for _, accel := range []string{a.settings.PinShortcut, a.settings.DeleteShortcut, a.settings.PreviewShortcut} {
		if ch, ok := shortcut.Parse(accel); ok {
			if r := shortcut.KeyRune(ch.Key); r != 0 {
				out = append(out, string(r))
			}
		}
	}
	return out
}

func (a *app) TogglePin(id int64) {
	if _, err := a.engine.TogglePin(id, a.pinReserved()...); err != nil {
		log.Print("myclip: pin: ", err)
	}
}

func (a *app) Delete(id int64) {
	if err := a.engine.Delete(id); err != nil {
		log.Print("myclip: delete: ", err)
	}
}

// Clear clears the history from the popup, which closes (History.clear).
func (a *app) Clear(all bool) {
	a.Close()
	a.clearHistory(all)
	a.saveSettings() // suppressClearAlert
}

func (a *app) OpenSettings() {
	a.Close()
	a.openPrefs()
}

func (a *app) About() {
	a.Close()
	mygo.App.Focus()
	go mygo.Dialog.Message(mygo.MessageOptions{
		Type:    mygo.MessageInfo,
		Message: locales.AboutMessage(appVersion()),
		Detail:  locales.AboutDetail(),
	})
}

func (a *app) Quit() { mygo.App.Quit() }

func (a *app) Thumbnail(it *history.Item, maxHeight int) *ui.Bitmap {
	return a.images.thumbnail(it, maxHeight)
}

func (a *app) Preview(it *history.Item) *ui.Bitmap { return a.images.preview(it) }

// PreviewSize is the size of the image Preview shows, as copied.
func (a *app) PreviewSize(it *history.Item) (int, int, bool) { return a.images.previewSize(it) }

func (a *app) AppIcon(bundleID string) *ui.Bitmap { return a.images.appIcon(bundleID) }

func (a *app) AppName(bundleID string) string { return a.images.appName(bundleID) }

// Resize takes the size the popup asks for after the frame, keeping its
// top left where it is unless that puts it off its screen, as the preview
// widening it near the right edge would.
func (a *app) Resize(w, h int) {
	a.panel.Update(func() {
		b := a.panel.Bounds()
		if b.Width == w && b.Height == h {
			return
		}
		area := mygo.Screen.DisplayMatching(b).WorkArea
		x := max(area.X, min(b.X, area.X+area.Width-w))
		y := max(area.Y, min(b.Y, area.Y+area.Height-h))
		a.panel.SetSize(w, h)
		if x != b.X || y != b.Y {
			a.panel.SetPosition(x, y)
		}
	})
}

func appVersion() string {
	if v := mygo.App.Version(); v != "" {
		return v
	}
	return "dev"
}
