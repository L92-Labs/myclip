package main

import (
	"bytes"
	"errors"
	"image/png"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

func clip(id int64, s string) *history.Item {
	it := history.NewItem([]history.Content{{Type: pasteboard.TypeString, Value: []byte(s)}}, "com.apple.TextEdit", time.Now())
	it.ID = id
	return it
}

func labels(items []*mygo.MenuItem) []string {
	var out []string
	for _, it := range items {
		out = append(out, it.Label)
	}
	return out
}

func TestMenuEmpty(t *testing.T) {
	got := labels(menuItems(menuState{}, menuActions{}))
	if got[0] != "No clips yet — copy something" || got[len(got)-1] != "Quit MyClip" {
		t.Fatalf("labels %q", got)
	}
}

func TestMenuClipsCopyBack(t *testing.T) {
	var copied int64
	var paused bool
	s := menuState{Items: []*history.Item{clip(7, "  hello\n  world  "), clip(3, "second")}}
	items := menuItems(s, menuActions{
		Copy:      func(id int64) { copied = id },
		SetPaused: func(p bool) { paused = p },
	})
	if items[0].Label != "hello world" || items[0].Accelerator != "CmdOrCtrl+1" {
		t.Fatalf("first row %q %q", items[0].Label, items[0].Accelerator)
	}
	items[1].Click(items[1], nil)
	if copied != 3 {
		t.Fatalf("copied id %d, want 3", copied)
	}
	for _, it := range items {
		if it.Label == "Pause Capture" {
			it.Checked = true // MyGo toggles before Click
			it.Click(it, nil)
		}
	}
	if !paused {
		t.Fatal("pause not forwarded")
	}
	if !strings.Contains(strings.Join(labels(items), "|"), "2 clips · stored on this Mac only") {
		t.Fatalf("status missing: %q", labels(items))
	}
}

func TestMenuCapsAndWarnings(t *testing.T) {
	var many []*history.Item
	for i := range 30 {
		many = append(many, clip(int64(i+1), "x"))
	}
	items := menuItems(menuState{Items: many, Oversized: 2, Err: errors.New("disk full")}, menuActions{})
	clips := 0
	for _, it := range items {
		if it.Label == "x" {
			clips++
		}
	}
	joined := strings.Join(labels(items), "|")
	if clips != menuLimit || !strings.Contains(joined, "2 skipped (over 32 MB)") ||
		!strings.Contains(joined, "⚠ Not saved: disk full") {
		t.Fatalf("%d clips, labels %q", clips, joined)
	}
}

func TestTrayIconIsPNG(t *testing.T) {
	img, err := png.Decode(bytes.NewReader(trayIcon()))
	if err != nil || img.Bounds().Dx() != 32 {
		t.Fatalf("icon: %v", err)
	}
}
