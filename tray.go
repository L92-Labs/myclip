package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/egoist/mygo"

	"github.com/L92-Labs/myclip/internal/history"
)

// menuLimit is how many clips the M1 tray lists; the popup (M2) shows all.
const menuLimit = 20

type menuState struct {
	Items     []*history.Item
	Paused    bool
	Oversized int
	Err       error
}

type menuActions struct {
	Copy       func(id int64)
	SetPaused  func(bool)
	IgnoreNext func()
}

// menuItems lays out the tray menu: recent clips (click copies one back),
// capture controls, a local-only status line, and Quit.
func menuItems(s menuState, do menuActions) []*mygo.MenuItem {
	var items []*mygo.MenuItem
	if len(s.Items) == 0 {
		items = append(items, &mygo.MenuItem{Label: "No clips yet — copy something", Disabled: true})
	}
	for i, it := range s.Items {
		if i == menuLimit {
			break
		}
		id := it.ID
		label := it.DisplayTitle()
		if it.Pin != "" {
			label = "📌 " + label
		}
		mi := &mygo.MenuItem{
			Label:   label,
			ToolTip: it.Application,
			Click:   func(*mygo.MenuItem, *mygo.Window) { do.Copy(id) },
		}
		if i < 9 {
			mi.Accelerator = fmt.Sprintf("CmdOrCtrl+%d", i+1)
		}
		items = append(items, mi)
	}

	status := fmt.Sprintf("%d clips · stored on this Mac only", len(s.Items))
	if len(s.Items) == 1 {
		status = "1 clip · stored on this Mac only"
	}
	items = append(items,
		mygo.Separator(),
		&mygo.MenuItem{
			Label:   "Pause Capture",
			Type:    mygo.MenuItemCheckbox,
			Checked: s.Paused,
			Click:   func(mi *mygo.MenuItem, _ *mygo.Window) { do.SetPaused(mi.Checked) },
		},
		&mygo.MenuItem{
			Label:    "Ignore Next Copy",
			Disabled: s.Paused,
			Click:    func(*mygo.MenuItem, *mygo.Window) { do.IgnoreNext() },
		},
		mygo.Separator(),
		&mygo.MenuItem{Label: status, Disabled: true},
	)
	if s.Oversized > 0 {
		items = append(items, &mygo.MenuItem{
			Label:    fmt.Sprintf("%d skipped (over 32 MB)", s.Oversized),
			Disabled: true,
		})
	}
	if s.Err != nil {
		items = append(items, &mygo.MenuItem{Label: "⚠ Not saved: " + s.Err.Error(), Disabled: true})
	}
	return append(items, mygo.Separator(), &mygo.MenuItem{Label: "Quit MyClip", Role: mygo.RoleQuit})
}

// trayIcon draws a clipboard, black on transparent for macOS to tint.
func trayIcon() []byte {
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
	var b bytes.Buffer
	_ = png.Encode(&b, img)
	return b.Bytes()
}
