//go:build windows

package native

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Modifier flags, as on macOS (NSEventModifierFlags), which is what the
// app's hotkey and modifier tracking compare against.
const (
	ModShift   = 1 << 17
	ModControl = 1 << 18
	ModOption  = 1 << 19
	ModCommand = 1 << 20
)

// Rect is a screen rectangle in points.
type Rect struct{ X, Y, Width, Height float64 }

var (
	procAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	procSendInput     = user32.NewProc("SendInput")
	procForeground    = user32.NewProc("GetForegroundWindow")
	procWindowThread  = user32.NewProc("GetWindowThreadProcessId")
	procWindowRect    = user32.NewProc("GetWindowRect")
	procDpiForWindow  = user32.NewProc("GetDpiForWindow")
	procOpenProcess   = kernel32.NewProc("OpenProcess")
	procCloseHandle   = kernel32.NewProc("CloseHandle")
	procProcessImage  = kernel32.NewProc("QueryFullProcessImageNameW")
)

// Virtual-key codes for the modifiers MyClip tracks.
const (
	vkShift   = 0x10
	vkControl = 0x11
	vkMenu    = 0x12
	vkLWin    = 0x5B
	vkRWin    = 0x5C
)

// Prepare runs before the app starts; nothing to do here.
func Prepare() {}

func keyDown(vk int) bool {
	v, _, _ := procAsyncKeyState.Call(uintptr(vk))
	return v&0x8000 != 0
}

// ModifierFlags reports the keys held right now: Shift, Control, Alt
// (as Option), and the Windows keys (as Command).
func ModifierFlags() uint {
	var f uint
	if keyDown(vkShift) {
		f |= ModShift
	}
	if keyDown(vkControl) {
		f |= ModControl
	}
	if keyDown(vkMenu) {
		f |= ModOption
	}
	if keyDown(vkLWin) || keyDown(vkRWin) {
		f |= ModCommand
	}
	return f
}

// AccessibilityTrusted is always true: SendInput needs no permission grant.
func AccessibilityTrusted(bool) bool { return true }

// SecureInputEnabled is false: Windows has no such switch.
func SecureInputEnabled() bool { return false }

// keyInput is one INPUT structure for a keyboard event: KEYBDINPUT in the
// union's tail, padded out to the native INPUT size — SendInput requires
// sizeof(INPUT), which the union sizes to MOUSEINPUT (40 bytes on amd64).
type keyInput struct {
	kind      uint32
	reserved  uint32
	vk        uint16
	scan      uint16
	flags     uint32
	time      uint32
	extraInfo uintptr
	_pad      [8]byte // the union's tail: MOUSEINPUT is larger than KEYBDINPUT
}

const (
	inputKeyboard = 1
	keyKeyUp      = 2 // KEYEVENTF_KEYUP
)

func tapKey(vk uint16) error {
	down := keyInput{kind: inputKeyboard, vk: vk}
	up := down
	up.flags = keyKeyUp
	in := []keyInput{down, up}
	n, _, _ := procSendInput.Call(uintptr(len(in)), uintptr(unsafe.Pointer(&in[0])), unsafe.Sizeof(keyInput{}))
	if n == 0 {
		return errors.New("native: SendInput failed")
	}
	return nil
}

// Paste types Ctrl+V into the focused window.
func Paste() error {
	down := keyInput{kind: inputKeyboard, vk: vkControl}
	up := down
	up.flags = keyKeyUp
	v := keyInput{kind: inputKeyboard, vk: 'V'}
	vUp := v
	vUp.flags = keyKeyUp
	in := []keyInput{down, v, vUp, up}
	n, _, _ := procSendInput.Call(uintptr(len(in)), uintptr(unsafe.Pointer(&in[0])), unsafe.Sizeof(keyInput{}))
	if n == 0 {
		return errors.New("native: SendInput failed")
	}
	return nil
}

type rect struct{ Left, Top, Right, Bottom int32 }

// processImage returns the executable path of a window's process.
func processImage(hwnd uintptr) string {
	var pid uint32
	procWindowThread.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return ""
	}
	const (
		processQueryLimitedInformation = 0x1000
	)
	h, _, _ := procOpenProcess.Call(processQueryLimitedInformation, 0, uintptr(pid))
	if h == 0 {
		return ""
	}
	defer procCloseHandle.Call(h)
	buf := make([]uint16, windows.MAX_PATH)
	n := uint32(len(buf))
	// (hProcess, dwFlags, lpExeName, lpdwSize)
	ok, _, _ := procProcessImage.Call(h, 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n)))
	if ok == 0 {
		return ""
	}
	return windows.UTF16ToString(buf[:n])
}

// FrontmostBundleID is the focused window's executable path — the identity
// the Ignore › Applications list takes on Windows; "" when it cannot be
// read (a system window, an elevated process).
func FrontmostBundleID() string {
	hwnd, _, _ := procForeground.Call()
	return processImage(hwnd)
}

// windowDPI scales physical pixels to the 96-DPI points the app's
// geometry uses.
func windowDPI(hwnd uintptr) float64 {
	if v, _, _ := procDpiForWindow.Call(hwnd); v != 0 {
		return float64(v) / 96
	}
	return 1
}

func FrontmostWindowFrame() (Rect, bool) {
	hwnd, _, _ := procForeground.Call()
	if hwnd == 0 {
		return Rect{}, false
	}
	var r rect
	if v, _, _ := procWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r))); v == 0 {
		return Rect{}, false
	}
	d := windowDPI(hwnd)
	return Rect{
		X:      float64(r.Left) / d,
		Y:      float64(r.Top) / d,
		Width:  float64(r.Right-r.Left) / d,
		Height: float64(r.Bottom-r.Top) / d,
	}, true
}

// AppPath is the identity itself: on Windows FrontmostBundleID already
// returns the executable path.
func AppPath(bundleID string) string {
	if bundleID == "" {
		return ""
	}
	if _, err := windows.UTF16FromString(bundleID); err != nil {
		return ""
	}
	return bundleID
}

// AppName is the executable's base name without ".exe".
func AppName(bundleID string) string {
	name := filepath.Base(strings.TrimSpace(bundleID))
	return strings.TrimSuffix(name, ".exe")
}

// BundleID is the cleaned path itself.
func BundleID(path string) string { return filepath.Clean(strings.TrimSpace(path)) }

// ResolveFileURL turns a file URL into a local Windows path, leaving other
// strings alone.
func ResolveFileURL(s string) string {
	u, err := url.Parse(s)
	if err != nil || u.Scheme != "file" || u.Path == "" {
		return s
	}
	p := u.Path // file:///C:/dir/file → /C:/dir/file
	if len(p) >= 3 && p[0] == '/' && p[2] == ':' {
		p = p[1:]
	}
	return strings.ReplaceAll(p, "/", `\`)
}
