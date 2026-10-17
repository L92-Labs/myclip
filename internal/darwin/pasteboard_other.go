//go:build !darwin

package darwin

import (
	"errors"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// ErrUnsupported is returned off macOS; other platforms come later.
var ErrUnsupported = errors.New("darwin: the system pasteboard needs macOS")

// General is unavailable off macOS.
type General struct{ pasteboard.Board }

// NewGeneral fails off macOS.
func NewGeneral() (*General, error) { return nil, ErrUnsupported }

// FrontmostBundleID is empty off macOS.
func FrontmostBundleID() string { return "" }
