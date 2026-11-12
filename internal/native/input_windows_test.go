//go:build windows

package native

import (
	"testing"
	"unsafe"
)

// SendInput takes sizeof(INPUT): the union is sized to MOUSEINPUT, larger
// than KEYBDINPUT. A smaller cbSize makes every call fail.
func TestKeyInputIsNativeInputSize(t *testing.T) {
	if got := unsafe.Sizeof(keyInput{}); got != 40 {
		t.Fatalf("keyInput is %d bytes, want the native INPUT's 40", got)
	}
}
