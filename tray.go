package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/L92-Labs/myclip/internal/darwin"
	"github.com/L92-Labs/myclip/internal/settings"
	"github.com/L92-Labs/myclip/locales"
)

// recentTitleRunes bounds the recent copy shown beside the menu icon.
const recentTitleRunes = 30

// applyTray shows or hides the menu icon (showInStatusBar) and draws it
// (menuIcon), dimmed while capture is paused.
func (a *app) applyTray() {
	// In capture-safe mode the icon stays: its menu holds the way out.
	if !a.settings.ShowInStatusBar && !a.captureSafe {
		if a.tray != nil {
			a.tray.Destroy()
			a.tray = nil
		}
		return
	}
	icon := trayIcon(a.settings.MenuIcon, a.engine.Paused() || a.captureSafe)
	if a.tray != nil {
		a.tray.SetIcon(icon, true)
		a.updateTray()
		return
	}
	var err error
	a.tray, err = mygo.NewTray(mygo.TrayOptions{Icon: icon, IconIsTemplate: true, ToolTip: "MyClip"})
	if err != nil {
		log.Print("myclip: tray: ", err)
		return
	}
	a.tray.OnClick(a.trayClick)
	a.tray.OnRightClick(func() {
		a.tray.SetMenu(mygo.NewMenu(a.trayMenu()))
		a.tray.PopUpMenu()
		a.tray.SetMenu(nil)
	})
	a.updateTray()
}

// trayClick is Maccy's performStatusItemClick: ⌥-click pauses or resumes
// capture, ⌥⇧-click ignores only the next copy, a click opens the popup
// below the icon.
func (a *app) trayClick() {
	mods := darwin.ModifierFlags()
	if mods&darwin.ModOption != 0 && a.captureSafe {
		return // capture is off until the settings are reset
	}
	if mods&darwin.ModOption != 0 {
		switch {
		case a.engine.Paused():
			a.setPaused(false)
		case mods&darwin.ModShift != 0:
			a.engine.IgnoreNext()
			a.applyTray()
		default:
			a.setPaused(true)
		}
		return
	}
	if a.panel.IsVisible() {
		a.Close()
		return
	}
	saved := a.settings.PopupPosition
	a.settings.PopupPosition = settings.PopupStatusItem
	a.showPopup()
	a.settings.PopupPosition = saved
}

// setPaused pauses or resumes capture (ignoreEvents, kept in the settings).
func (a *app) setPaused(paused bool) {
	a.engine.SetPaused(paused)
	a.settings.IgnoreEvents = paused
	a.applied.IgnoreEvents = paused
	a.saveSettings()
	a.applyTray()
}

// updateTray shows the recent copy beside the icon when asked
// (showRecentCopyInMenuBar).
func (a *app) updateTray() {
	if a.tray == nil {
		return
	}
	title := ""
	if a.settings.ShowRecentCopyInMenuBar {
		if items := a.engine.Items(); len(items) > 0 {
			title = recentTitle(items[0].DisplayTitle())
		}
	}
	a.tray.SetTitle(title)
}

// recentTitle shortens a title for the menu bar.
func recentTitle(s string) string {
	if r := []rune(s); len(r) > recentTitleRunes {
		return string(r[:recentTitleRunes-1]) + "…"
	}
	return s
}

// trayMenu is the right-click menu: open, capture controls, preferences,
// a local-only status line, and Quit.
func (a *app) trayMenu() []*mygo.MenuItem {
	n := len(a.engine.Items())
	var items []*mygo.MenuItem
	if a.captureSafe {
		items = append(items,
			&mygo.MenuItem{Label: locales.TrayCaptureSafe(), Disabled: true},
			&mygo.MenuItem{Label: locales.TrayResetSettings(), Click: func(*mygo.MenuItem, *mygo.Window) { a.resetSettings() }},
			mygo.Separator(),
		)
	}
	items = append(items, []*mygo.MenuItem{
		{Label: locales.TrayOpen(), Click: func(*mygo.MenuItem, *mygo.Window) { a.showPopup() }},
		mygo.Separator(),
		{
			Label:    locales.TrayPause(),
			Type:     mygo.MenuItemCheckbox,
			Checked:  a.engine.Paused() && !a.engine.IgnoringNext(),
			Disabled: a.captureSafe, // capture is already off; reset is the way out
			Click:    func(mi *mygo.MenuItem, _ *mygo.Window) { a.setPaused(mi.Checked) },
		},
		{
			Label:    locales.TrayIgnoreNext(),
			Type:     mygo.MenuItemCheckbox,
			Checked:  a.engine.IgnoringNext(),
			Disabled: a.captureSafe,
			Click: func(mi *mygo.MenuItem, _ *mygo.Window) {
				if mi.Checked {
					a.engine.IgnoreNext()
				} else {
					a.engine.SetPaused(false)
				}
				a.applyTray()
			},
		},
		mygo.Separator(),
		{Label: locales.Preferences(), Accelerator: "Cmd+,", Click: func(*mygo.MenuItem, *mygo.Window) { a.openPrefs() }},
		{Label: locales.About(), Click: func(*mygo.MenuItem, *mygo.Window) { a.About() }},
		mygo.Separator(),
		{Label: locales.TrayClips(n), Disabled: true},
	}...)
	if o := a.engine.Oversized(); o > 0 {
		items = append(items, &mygo.MenuItem{Label: locales.TraySkipped(o), Disabled: true})
	}
	if a.err != nil {
		items = append(items, &mygo.MenuItem{Label: locales.TrayNotSaved(a.err.Error()), Disabled: true})
	}
	return append(items, mygo.Separator(), &mygo.MenuItem{Label: locales.QuitMyClip(), Role: mygo.RoleQuit})
}

var trayIcons = map[string]*ui.SVG{
	settings.IconClipboard: ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="8" height="4" x="8" y="2" rx="1"/><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/></svg>`)),
	settings.IconScissors:  ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="6" cy="6" r="3"/><path d="M8.12 8.12 12 12"/><path d="M20 4 8.12 15.88"/><circle cx="6" cy="18" r="3"/><path d="M14.8 14.8 20 20"/></svg>`)),
	settings.IconPaperclip: ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48"/></svg>`)),
}

// trayIcon draws the menu icon, black on transparent for macOS to tint;
// paused, it is faint (Maccy's appearsDisabled).
func trayIcon(name string, paused bool) []byte {
	var img *image.NRGBA
	if svg, ok := trayIcons[name]; ok {
		rgba := ui.Render(func(c *ui.Context) {
			c.Root().Background(ui.Transparent)
			ui.Icon(c, svg).FontSize(16).TextColor(ui.RGB(0, 0, 0))
		}, 16, 16, 2)
		img = image.NewNRGBA(rgba.Bounds())
		draw.Draw(img, img.Bounds(), rgba, image.Point{}, draw.Src)
	} else {
		img = clipboardIcon()
	}
	if paused {
		for i := 3; i < len(img.Pix); i += 4 {
			img.Pix[i] = uint8(int(img.Pix[i]) * 2 / 5)
		}
	}
	var b bytes.Buffer
	_ = png.Encode(&b, img)
	return b.Bytes()
}

// clipboardIcon is MyClip's own icon: a clipboard with lines of text.
func clipboardIcon() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	ink := image.NewUniform(color.NRGBA{A: 255})
	for _, r := range []image.Rectangle{
		// board outline
		image.Rect(6, 5, 26, 7), image.Rect(6, 28, 26, 30),
		image.Rect(6, 7, 8, 28), image.Rect(24, 7, 26, 28),
		// clip
		image.Rect(11, 2, 21, 4), image.Rect(11, 4, 13, 9), image.Rect(19, 4, 21, 9),
		image.Rect(11, 8, 21, 10),
		// lines of text
		image.Rect(11, 14, 21, 16), image.Rect(11, 19, 21, 21), image.Rect(11, 24, 17, 26),
	} {
		draw.Draw(img, r, ink, image.Point{}, draw.Src)
	}
	return img
}
