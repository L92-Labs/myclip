//go:build linux

package native

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Modifier flags, as on macOS; Linux reports none globally.
const (
	ModShift   = 1 << 17
	ModControl = 1 << 18
	ModOption  = 1 << 19
	ModCommand = 1 << 20
)

// Rect is a screen rectangle in points.
type Rect struct{ X, Y, Width, Height float64 }

// Prepare runs before the app starts. Under Wayland, an app may read the
// clipboard only while focused; MyClip asks GTK for X11 (XWayland), whose
// clipboard the compositor keeps in sync, so it sees copies made anywhere.
func Prepare() {
	if os.Getenv("WAYLAND_DISPLAY") != "" && os.Getenv("DISPLAY") != "" && os.Getenv("GDK_BACKEND") == "" {
		os.Setenv("GDK_BACKEND", "x11")
	}
}

// ModifierFlags is 0: Linux has no global modifier state to read.
func ModifierFlags() uint { return 0 }

// paster is the command that types Ctrl+V into the focused window.
func paster() []string {
	look := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
	wtype := []string{"wtype", "-M", "ctrl", "v", "-m", "ctrl"}
	ydotool := []string{"ydotool", "key", "29:1", "47:1", "47:0", "29:0"}
	xdotool := []string{"xdotool", "key", "--clearmodifiers", "ctrl+v"}
	// In a Wayland session (even with MyClip on XWayland) the focused app
	// may be a native Wayland one, which xdotool's XTEST cannot reach.
	order := [][]string{xdotool, ydotool}
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		order = [][]string{wtype, ydotool, xdotool}
	}
	for _, cmd := range order {
		if cmd[0] == "xdotool" && os.Getenv("DISPLAY") == "" {
			continue
		}
		if look(cmd[0]) {
			return cmd
		}
	}
	return nil
}

// AccessibilityTrusted reports whether MyClip can paste: a key-typing tool
// (xdotool on X11, wtype or ydotool on Wayland) is installed.
func AccessibilityTrusted(prompt bool) bool { return paster() != nil }

// SecureInputEnabled is false: Linux has no such switch.
func SecureInputEnabled() bool { return false }

// ErrNoPaster says why a paste did not happen.
var ErrNoPaster = errors.New("install xdotool (X11) or wtype / ydotool (Wayland) to paste; the item was copied")

// Paste types Ctrl+V into the focused window.
func Paste() error {
	cmd := paster()
	if cmd == nil {
		return ErrNoPaster
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, cmd[0], cmd[1:]...).Run()
}

var (
	activeWindow = regexp.MustCompile(`window id # (0x[0-9a-fA-F]+)`)
	wmClass      = regexp.MustCompile(`WM_CLASS = "([^"]*)", "([^"]*)"`)
)

func xprop(args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	out, err := exec.CommandContext(ctx, "xprop", args...).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

// FrontmostBundleID is the focused window's WM_CLASS (e.g. "KeePassXC") on
// X11, the name the Ignore › Applications list takes on Linux; "" when it
// cannot be read (xprop missing, native Wayland windows).
func FrontmostBundleID() string {
	if os.Getenv("DISPLAY") == "" {
		return ""
	}
	m := activeWindow.FindStringSubmatch(xprop("-root", "-notype", "_NET_ACTIVE_WINDOW"))
	if m == nil || m[1] == "0x0" {
		return ""
	}
	c := wmClass.FindStringSubmatch(xprop("-id", m[1], "-notype", "WM_CLASS"))
	if c == nil {
		return ""
	}
	return c[2]
}

func FrontmostWindowFrame() (Rect, bool) { return Rect{}, false }
func AppPath(bundleID string) string     { return "" }

// AppName is the WM_CLASS itself.
func AppName(bundleID string) string { return strings.TrimSpace(bundleID) }

func BundleID(path string) string    { return "" }
func ResolveFileURL(s string) string { return s }
