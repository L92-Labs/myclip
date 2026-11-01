//go:build windows

package native

import "errors"

// ErrNoOCR says text recognition is not available on this build: the
// Windows OCR API needs WinRT interop, which is not wired up yet.
var ErrNoOCR = errors.New("native: text recognition is not available on Windows yet")

// RecognizeText is unavailable on this build.
func RecognizeText([]byte) (string, error) { return "", ErrNoOCR }
