//go:build !darwin && !linux

package native

import (
	"errors"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// ErrUnsupported is returned off macOS; other platforms come later.
var ErrUnsupported = errors.New("native: the system pasteboard needs macOS")

// General is unavailable off macOS.
type General struct{ pasteboard.Board }

// NewGeneral fails off macOS.
func NewGeneral() (*General, error) { return nil, ErrUnsupported }

// FrontmostBundleID is empty off macOS.
func FrontmostBundleID() string { return "" }

// Modifier flags (NSEventModifierFlags).
const (
	ModShift   = 1 << 17
	ModControl = 1 << 18
	ModOption  = 1 << 19
	ModCommand = 1 << 20
)

// Rect is a screen rectangle in points.
type Rect struct{ X, Y, Width, Height float64 }

func ModifierFlags() uint                   { return 0 }
func AccessibilityTrusted(prompt bool) bool { return false }
func SecureInputEnabled() bool              { return false }
func Paste() error                          { return ErrUnsupported }
func FrontmostWindowFrame() (Rect, bool)    { return Rect{}, false }
func AppPath(bundleID string) string        { return "" }
func AppName(bundleID string) string        { return bundleID }
func BundleID(path string) string           { return "" }
func ResolveFileURL(s string) string        { return s }

// Prepare runs before the app starts; nothing to do here.
func Prepare() {}
