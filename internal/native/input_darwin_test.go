//go:build darwin

package native

import "testing"

func TestInputBindingsLoad(t *testing.T) {
	if err := loadInput(); err != nil {
		t.Fatal(err)
	}
	code, ok := layoutKeyCode('v')
	t.Logf("v key %d %v, mods %x, trusted %v, secure %v", code, ok, ModifierFlags(), AccessibilityTrusted(false), SecureInputEnabled())
	if p := AppPath("com.apple.finder"); p == "" {
		t.Error("no Finder path")
	}
	r, ok := FrontmostWindowFrame()
	t.Logf("frontmost frame %+v %v, Finder name %q", r, ok, AppName("com.apple.finder"))
}
