package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/L92-Labs/myclip/internal/darwin"
	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/prefs"
	"github.com/L92-Labs/myclip/internal/settings"
	"github.com/L92-Labs/myclip/locales"
)

// saveSettings writes the settings when they changed since the last save.
func (a *app) saveSettings() {
	if a.captureSafe {
		return // never overwrite the unreadable file: reset is the way out
	}
	if err := settings.Save(a.settingsPath, a.settings); err != nil {
		log.Print("myclip: save settings: ", err)
	}
}

// Changed applies edited settings: capture rules, limits, the poll
// interval, the shortcut, the menu icon and the popup.
func (a *app) Changed() {
	a.settings.Normalize()
	old, s := a.applied, a.settings
	a.applied = s
	a.saveSettings()

	if !slices.Equal(old.EnabledPasteboardTypes, s.EnabledPasteboardTypes) ||
		!slices.Equal(old.IgnoredPasteboardTypes, s.IgnoredPasteboardTypes) ||
		!slices.Equal(old.IgnoredApps, s.IgnoredApps) ||
		!slices.Equal(old.IgnoreRegexp, s.IgnoreRegexp) ||
		old.IgnoreAllAppsExceptListed != s.IgnoreAllAppsExceptListed {
		p := s.CapturePrefs()
		p.IgnoreEvents = a.engine.Paused()
		a.engine.SetPrefs(p)
	}
	if old.IgnoreEvents != s.IgnoreEvents {
		a.engine.SetPaused(s.IgnoreEvents)
	}
	if old.HistorySize != s.HistorySize {
		a.engine.SetLimits(a.limits())
	}
	if old.ClipboardCheckInterval != s.ClipboardCheckInterval {
		a.restartPoll()
	}
	if old.PopupShortcut != s.PopupShortcut {
		a.registerShortcut()
	}
	a.applyTray()
	if a.panel != nil {
		a.panel.SetSize(a.popup.Fit())
	}
	a.refresh()
}

// openPrefs shows the settings window, made the first time.
func (a *app) openPrefs() {
	if a.prefsWin == nil {
		a.prefs = &prefs.View{Host: a, S: &a.settings}
		a.prefsWin = mygo.NewWindow(mygo.WindowOptions{
			Title:   locales.PrefWindowTitle(),
			Width:   680,
			Height:  620,
			Hidden:  true,
			Content: ui.View(a.prefs.View),
		})
		a.prefsWin.OnClose(func(e *mygo.CloseEvent) {
			a.prefs.StopRecording()
			if !a.quitting {
				e.PreventDefault()
				a.prefsWin.Hide()
			}
		})
	}
	a.prefs.Load()
	if !a.prefsWin.IsVisible() {
		a.prefsWin.Center()
	}
	a.prefsWin.Show()
	mygo.App.Focus()
	a.prefsWin.Focus()
}

// The settings window's Host.

func (a *app) Recording(on bool) {
	if on {
		if a.shortcut != "" {
			mygo.GlobalShortcut.Unregister(a.shortcut)
			a.shortcut = ""
		}
		return
	}
	a.registerShortcut()
}

func (a *app) PinnedItems() []*history.Item {
	var out []*history.Item
	for _, it := range a.sortedItems() {
		if it.Pin != "" {
			out = append(out, it)
		}
	}
	return out
}

func (a *app) PinKeys(id int64) []string {
	var keys []string
	if it := a.engine.Item(id); it != nil && it.Pin != "" {
		keys = append(keys, it.Pin)
	}
	keys = append(keys, history.AvailablePins(a.engine.Items(), a.pinReserved()...)...)
	slices.Sort(keys)
	return keys
}

func (a *app) SetPin(id int64, key string) error     { return a.engine.SetPin(id, key) }
func (a *app) SetTitle(id int64, title string) error { return a.engine.SetTitle(id, title) }
func (a *app) SetText(id int64, text string) error   { return a.engine.SetText(id, text) }

func (a *app) Unpin(id int64) {
	if it := a.engine.Item(id); it != nil && it.Pin != "" {
		a.TogglePin(id)
	}
}

func (a *app) LaunchAtLogin() bool { return mygo.App.OpenAtLogin() }

func (a *app) SetLaunchAtLogin(on bool) error { return mygo.App.SetOpenAtLogin(on) }

func (a *app) AccessibilityTrusted() bool { return darwin.AccessibilityTrusted(false) }

func (a *app) RequestAccessibility() { darwin.AccessibilityTrusted(true) }

// StorageSize is the history database's size on disk.
func (a *app) StorageSize() string {
	fi, err := os.Stat(a.store.Path())
	if err != nil {
		return "?"
	}
	n := float64(fi.Size())
	for _, unit := range []string{"B", "KB", "MB", "GB"} {
		if n < 1000 || unit == "GB" {
			if unit == "B" {
				return fmt.Sprintf("%.0f %s", n, unit)
			}
			return fmt.Sprintf("%.1f %s", n, unit)
		}
		n /= 1000
	}
	return ""
}

// PickApp asks for an application and ignores its copies.
func (a *app) PickApp() {
	go func() {
		paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
			Parent:      a.prefsWin,
			DefaultPath: "/Applications",
			Filters:     []mygo.FileFilter{{Name: "Applications", Extensions: []string{"app"}}},
			Multiple:    true,
		})
		if err != nil || len(paths) == 0 {
			return
		}
		mygo.RunOnMain(func() {
			for _, p := range paths {
				if id := darwin.BundleID(p); id != "" && !slices.Contains(a.settings.IgnoredApps, id) {
					a.settings.IgnoredApps = append(slices.Clone(a.settings.IgnoredApps), id)
				}
			}
			a.Changed()
		})
	}()
}

func (a *app) Displays() []string {
	var out []string
	for i, d := range mygo.Screen.Displays() {
		name := d.Label
		if name == "" {
			name = fmt.Sprint(i + 1)
		}
		out = append(out, name)
	}
	return out
}

// resetSettings is the way out of capture-safe mode: it keeps the
// unreadable file aside and writes the defaults. Capture resumes only once
// both succeeded; otherwise the latch and the menu's reset item stay.
func (a *app) resetSettings() {
	aside, err := settings.Reset(a.settingsPath, time.Now())
	if err != nil {
		log.Print("myclip: reset settings: ", err)
		a.err = err
		a.refresh()
		return
	}
	if aside != "" {
		log.Print("myclip: unreadable settings kept as ", aside)
	}
	a.captureSafe = false
	a.err = nil
	a.settings = settings.Default()
	// A pause or ignore-next chosen while capture-safe must not swallow the
	// first copy after the reset.
	a.engine.SetPaused(false)
	a.engine.SetPrefs(a.settings.CapturePrefs())
	a.engine.SetCaptureSafe(false)
	if a.prefs != nil {
		a.prefs.Load() // its controls mirror the old values
	}
	a.Changed()
}
