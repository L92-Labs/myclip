package locales

import "runtime"

// The messages that name the platform: macOS says Mac and Accessibility,
// Linux says this computer and names its paste tools.

var linux = runtime.GOOS == "linux"

func pick(mac, other func() string) string {
	if linux {
		return other()
	}
	return mac()
}

// ClipsHere is the tray's clip count and where they are stored.
func ClipsHere(n int) string {
	if linux {
		return TrayClipsLinux(n)
	}
	return TrayClips(n)
}

func AboutDetailHere() string   { return pick(AboutDetail, AboutDetailLinux) }
func PasteNeedsAccess() string  { return pick(PasteNeedsAccessibility, PasteNeedsTool) }
func PrefPasteAccess() string   { return pick(PrefAccessibility, PrefPasteTool) }
func PrefPasteGranted() string  { return pick(PrefAccessibilityGranted, PrefPasteToolFound) }
func FooterWarningHere() string { return pick(PrefFooterWarning, PrefFooterWarningLinux) }
func MenuIconHiddenHere() string {
	return pick(PrefMenuIconHidden, PrefMenuIconHiddenLinux)
}
func LocalOnlyHere() string { return pick(PrefLocalOnly, PrefLocalOnlyLinux) }

// CanRequestPasteAccess reports whether the OS has a permission to ask for
// (macOS Accessibility); on Linux the user installs a tool instead.
func CanRequestPasteAccess() bool { return !linux }
