//go:build darwin

package native

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// Paste and the Accessibility checks, after Maccy's Clipboard.paste and
// Accessibility.swift: a ⌘V posted with CGEvent, the V found in the
// current keyboard layout.

var (
	inputOnce sync.Once
	inputErr  error

	cgEventSourceCreate      func(state int32) uintptr
	cgEventSourceSetFilter   func(src uintptr, filter uint32, state int32)
	cgEventCreateKeyboard    func(src uintptr, key uint16, down bool) uintptr
	cgEventSetFlags          func(ev uintptr, flags uint64)
	cgEventPost              func(tap uint32, ev uintptr)
	cgWindowListCopyInfo     func(option, relativeTo uint32) uintptr
	cfRelease                func(ref uintptr)
	cfDataGetBytePtr         func(data uintptr) unsafe.Pointer
	axIsProcessTrusted       func(options uintptr) bool
	isSecureEventInput       func() bool
	tisCopyCurrentLayout     func() uintptr
	tisGetProperty           func(src, key uintptr) uintptr
	ucKeyTranslate           func(layout unsafe.Pointer, key, action uint16, mods, kbdType, options uint32, dead *uint32, maxLen uint, actualLen *uint, chars *uint16) int32
	lmGetKbdType             func() uint8
	keyLayoutData, localName uintptr // CFStringRef property keys

	clsEvent, clsDictionary, clsNumber objc.ID
	selModifierFlags, selDictWithObject, selNumberWithBool, selObjectForKey, selIntValue,
	selDoubleValue, selPID, selURLForBundleID, selPath objc.SEL
)

func loadInput() error {
	inputOnce.Do(func() {
		if inputErr = load(); inputErr != nil {
			return
		}
		const flags = purego.RTLD_GLOBAL | purego.RTLD_NOW
		open := func(path string) uintptr {
			lib, err := purego.Dlopen(path, flags)
			if err != nil && inputErr == nil {
				inputErr = fmt.Errorf("native: load %s: %w", filepath.Base(path), err)
			}
			return lib
		}
		cg := open("/System/Library/Frameworks/CoreGraphics.framework/CoreGraphics")
		cf := open("/System/Library/Frameworks/CoreFoundation.framework/CoreFoundation")
		as := open("/System/Library/Frameworks/ApplicationServices.framework/ApplicationServices")
		carbon := open("/System/Library/Frameworks/Carbon.framework/Carbon")
		if inputErr != nil {
			return
		}
		purego.RegisterLibFunc(&cgEventSourceCreate, cg, "CGEventSourceCreate")
		purego.RegisterLibFunc(&cgEventSourceSetFilter, cg, "CGEventSourceSetLocalEventsFilterDuringSuppressionState")
		purego.RegisterLibFunc(&cgEventCreateKeyboard, cg, "CGEventCreateKeyboardEvent")
		purego.RegisterLibFunc(&cgEventSetFlags, cg, "CGEventSetFlags")
		purego.RegisterLibFunc(&cgEventPost, cg, "CGEventPost")
		purego.RegisterLibFunc(&cgWindowListCopyInfo, cg, "CGWindowListCopyWindowInfo")
		purego.RegisterLibFunc(&cfRelease, cf, "CFRelease")
		purego.RegisterLibFunc(&cfDataGetBytePtr, cf, "CFDataGetBytePtr")
		purego.RegisterLibFunc(&axIsProcessTrusted, as, "AXIsProcessTrustedWithOptions")
		purego.RegisterLibFunc(&isSecureEventInput, carbon, "IsSecureEventInputEnabled")
		purego.RegisterLibFunc(&tisCopyCurrentLayout, carbon, "TISCopyCurrentKeyboardLayoutInputSource")
		purego.RegisterLibFunc(&tisGetProperty, carbon, "TISGetInputSourceProperty")
		purego.RegisterLibFunc(&ucKeyTranslate, carbon, "UCKeyTranslate")
		purego.RegisterLibFunc(&lmGetKbdType, carbon, "LMGetKbdType")
		global := func(name string) uintptr {
			p, err := purego.Dlsym(carbon, name)
			if err != nil {
				inputErr = fmt.Errorf("native: %s: %w", name, err)
				return 0
			}
			return **(**uintptr)(unsafe.Pointer(&p))
		}
		keyLayoutData = global("kTISPropertyUnicodeKeyLayoutData")
		localName = global("kTISPropertyLocalizedName")

		clsEvent = objc.ID(objc.GetClass("NSEvent"))
		clsDictionary = objc.ID(objc.GetClass("NSDictionary"))
		clsNumber = objc.ID(objc.GetClass("NSNumber"))
		s := objc.RegisterName
		selModifierFlags = s("modifierFlags")
		selDictWithObject = s("dictionaryWithObject:forKey:")
		selNumberWithBool = s("numberWithBool:")
		selObjectForKey = s("objectForKey:")
		selIntValue = s("intValue")
		selDoubleValue = s("doubleValue")
		selPID = s("processIdentifier")
		selURLForBundleID = s("URLForApplicationWithBundleIdentifier:")
		selPath = s("path")
	})
	return inputErr
}

// Modifier flags (NSEventModifierFlags).
const (
	ModShift   = 1 << 17
	ModControl = 1 << 18
	ModOption  = 1 << 19
	ModCommand = 1 << 20
)

// ModifierFlags returns the modifier keys held now.
func ModifierFlags() uint {
	if loadInput() != nil {
		return 0
	}
	return objc.Send[uint](clsEvent, selModifierFlags) & (ModShift | ModControl | ModOption | ModCommand)
}

// AccessibilityTrusted reports whether MyClip may post key events; with
// prompt, macOS asks the user to allow it in System Settings.
func AccessibilityTrusted(prompt bool) bool {
	if loadInput() != nil {
		return false
	}
	if !prompt {
		return axIsProcessTrusted(0)
	}
	var ok bool
	withPool(func() {
		opts := clsDictionary.Send(selDictWithObject,
			clsNumber.Send(selNumberWithBool, true), nsString("AXTrustedCheckOptionPrompt"))
		ok = axIsProcessTrusted(uintptr(opts))
	})
	return ok
}

// SecureInputEnabled reports whether a password field (or another app)
// holds secure keyboard input, which drops posted key events.
func SecureInputEnabled() bool {
	return loadInput() == nil && isSecureEventInput()
}

// qwertyV is the virtual key code of V on a US keyboard (kVK_ANSI_V).
const qwertyV = 9

// Paste posts ⌘V to the frontmost app. It needs Accessibility access.
func Paste() error {
	if err := loadInput(); err != nil {
		return err
	}
	key := uint16(qwertyV)
	if code, ok := layoutKeyCode('v'); ok {
		key = code
	}
	// Command plus the left-command device bit (Flycut #18).
	const flags = ModCommand | 0x000008
	src := cgEventSourceCreate(0) // kCGEventSourceStateCombinedSessionState
	if src != 0 {
		defer cfRelease(src)
		// Hold back local keyboard events while pasting.
		cgEventSourceSetFilter(src, 1|4, 0) // permit local mouse + system defined
	}
	for _, down := range []bool{true, false} {
		ev := cgEventCreateKeyboard(src, key, down)
		if ev == 0 {
			return fmt.Errorf("native: CGEventCreateKeyboardEvent failed")
		}
		cgEventSetFlags(ev, flags)
		cgEventPost(1, ev) // kCGSessionEventTap
		cfRelease(ev)
	}
	return nil
}

// layoutKeyCode finds the key typing ch in the current keyboard layout.
// Layouts that switch to QWERTY while ⌘ is held ("Dvorak - QWERTY ⌘")
// paste with the QWERTY key (Maccy #482).
func layoutKeyCode(ch rune) (uint16, bool) {
	src := tisCopyCurrentLayout()
	if src == 0 {
		return 0, false
	}
	defer cfRelease(src)
	if name := tisGetProperty(src, localName); name != 0 && strings.HasSuffix(goString(objc.ID(name)), "⌘") {
		return qwertyV, true
	}
	data := tisGetProperty(src, keyLayoutData)
	if data == 0 {
		return 0, false
	}
	layout := cfDataGetBytePtr(data)
	kbd := uint32(lmGetKbdType())
	for code := uint16(0); code < 128; code++ {
		var dead uint32
		var n uint
		var buf [4]uint16
		// kUCKeyActionDisplay, no modifiers, kUCKeyTranslateNoDeadKeysBit.
		if ucKeyTranslate(layout, code, 3, 0, kbd, 1, &dead, uint(len(buf)), &n, &buf[0]) != 0 {
			continue
		}
		if n == 1 && rune(buf[0]) == ch {
			return code, true
		}
	}
	return 0, false
}

// Rect is a screen rectangle in points, origin at the top left of the
// primary display (as MyGo's coordinates).
type Rect struct{ X, Y, Width, Height float64 }

// FrontmostWindowFrame is the frontmost app's top window, for the popup's
// "window center" position.
func FrontmostWindowFrame() (Rect, bool) {
	if loadInput() != nil {
		return Rect{}, false
	}
	var r Rect
	var ok bool
	withPool(func() {
		app := clsWorkspace.Send(selShared).Send(selFrontmost)
		if app == 0 {
			return
		}
		pid := objc.Send[int32](app, selPID)
		list := cgWindowListCopyInfo(1|16, 0) // on screen only, no desktop elements
		if list == 0 {
			return
		}
		defer cfRelease(list)
		arr := objc.ID(list)
		num := func(d objc.ID, key string) objc.ID { return d.Send(selObjectForKey, nsString(key)) }
		for i := uint(0); i < objc.Send[uint](arr, selCount); i++ {
			w := arr.Send(selObjectAtIndex, i)
			if objc.Send[int32](num(w, "kCGWindowOwnerPID"), selIntValue) != pid ||
				objc.Send[int32](num(w, "kCGWindowLayer"), selIntValue) != 0 {
				continue
			}
			b := num(w, "kCGWindowBounds")
			if b == 0 {
				continue
			}
			f := func(k string) float64 { return objc.Send[float64](num(b, k), selDoubleValue) }
			r, ok = Rect{f("X"), f("Y"), f("Width"), f("Height")}, true
			return
		}
	})
	return r, ok
}

// AppPath is where the app with bundleID is installed, or "".
func AppPath(bundleID string) string {
	if bundleID == "" || loadInput() != nil {
		return ""
	}
	var p string
	withPool(func() {
		u := clsWorkspace.Send(selShared).Send(selURLForBundleID, nsString(bundleID))
		if u != 0 {
			p = goString(u.Send(selPath))
		}
	})
	return p
}

// AppName is the app's name as Finder shows its bundle, or the bundle ID.
func AppName(bundleID string) string {
	if p := AppPath(bundleID); p != "" {
		return strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
	}
	return bundleID
}

// BundleID reads the bundle identifier of the app at path, or "".
func BundleID(path string) string {
	if loadInput() != nil {
		return ""
	}
	var id string
	withPool(func() {
		b := objc.ID(objc.GetClass("NSBundle")).Send(objc.RegisterName("bundleWithPath:"), nsString(path))
		if b != 0 {
			id = goString(b.Send(selBundleID))
		}
	})
	return id
}

// ResolveFileURL turns a file reference URL into a file path URL; it
// returns s unchanged when it cannot.
func ResolveFileURL(s string) string {
	if loadInput() != nil {
		return s
	}
	out := s
	withPool(func() {
		u := objc.ID(objc.GetClass("NSURL")).Send(objc.RegisterName("URLWithString:"), nsString(s))
		if u == 0 {
			return
		}
		if p := u.Send(objc.RegisterName("filePathURL")); p != 0 {
			out = goString(p.Send(objc.RegisterName("absoluteString")))
		}
	})
	return out
}
