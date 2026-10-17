// Package settings holds MyClip's preferences: Maccy 2.7.1's Defaults keys
// (Defaults.Keys+Names.swift @ eb03eba) plus MyClip's own, kept as JSON in
// the data root (0600, written atomically).
package settings

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// FileName is the settings file inside the data root.
const FileName = "settings.json"

// Choices of the enum preferences, named as Maccy stores them.
const (
	SearchExact  = "exact"
	SearchFuzzy  = "fuzzy"
	SearchRegexp = "regexp"
	SearchMixed  = "mixed"

	SearchAlways       = "always"
	SearchDuringSearch = "duringSearch"

	HighlightBold      = "bold"
	HighlightColor     = "color"
	HighlightItalic    = "italic"
	HighlightUnderline = "underline"

	PinTop    = "top"
	PinBottom = "bottom"

	PopupCursor       = "cursor"
	PopupStatusItem   = "statusItem"
	PopupWindow       = "window"
	PopupCenter       = "center"
	PopupLastPosition = "lastPosition"

	SortLastCopiedAt   = "lastCopiedAt"
	SortFirstCopiedAt  = "firstCopiedAt"
	SortNumberOfCopies = "numberOfCopies"

	IconMyClip    = "maccy"
	IconClipboard = "clipboard"
	IconScissors  = "scissors"
	IconPaperclip = "paperclip"
)

var (
	SearchModes      = []string{SearchExact, SearchFuzzy, SearchRegexp, SearchMixed}
	SearchVisibility = []string{SearchAlways, SearchDuringSearch}
	HighlightStyles  = []string{HighlightBold, HighlightColor, HighlightItalic, HighlightUnderline}
	PinPositions     = []string{PinTop, PinBottom}
	PopupPositions   = []string{PopupCursor, PopupStatusItem, PopupWindow, PopupCenter, PopupLastPosition}
	SortOrders       = []string{SortLastCopiedAt, SortFirstCopiedAt, SortNumberOfCopies}
	MenuIcons        = []string{IconMyClip, IconClipboard, IconScissors, IconPaperclip}
)

// Shortcut defaults (Maccy KeyboardShortcuts.Name+Shortcuts.swift), in
// MyGo accelerator syntax.
const (
	DefaultPopupShortcut   = "Cmd+Shift+C"
	DefaultPinShortcut     = "Alt+P"
	DefaultDeleteShortcut  = "Alt+Backspace"
	DefaultPreviewShortcut = "Ctrl+Space"
)

// Limits of the numeric preferences.
const (
	MinHistorySize     = 1
	MaxHistorySize     = 999
	MinImageHeight     = 1
	MaxImageHeight     = 200
	MinPreviewDelay    = 200
	MaxPreviewDelay    = 100000
	MinCheckIntervalMs = 100 // D7
	MaxCheckIntervalMs = 3600_000
)

// Settings are the preferences. JSON names follow Maccy's Defaults keys.
type Settings struct {
	// General
	PopupShortcut             string `json:"popupShortcut"`
	PinShortcut               string `json:"pinShortcut"`
	DeleteShortcut            string `json:"deleteShortcut"`
	PreviewShortcut           string `json:"togglePreviewShortcut"`
	SearchMode                string `json:"searchMode"`
	PasteByDefault            bool   `json:"pasteByDefault"`
	RemoveFormattingByDefault bool   `json:"removeFormattingByDefault"`
	NotifyOnCopy              bool   `json:"notifyOnCopy"` // MyClip: Maccy notifies once allowed

	// Storage
	EnabledPasteboardTypes []string `json:"enabledPasteboardTypes"`
	HistorySize            int      `json:"historySize"`
	SortBy                 string   `json:"sortBy"`

	// Appearance
	PopupPosition            string  `json:"popupPosition"`
	PopupScreen              int     `json:"popupScreen"` // 0 = active screen, else 1-based display
	PinTo                    string  `json:"pinTo"`
	ImageMaxHeight           int     `json:"imageMaxHeight"`
	PreviewDelay             int     `json:"previewDelay"` // ms
	OpenPreviewAutomatically bool    `json:"openPreviewAutomatically"`
	HighlightMatch           string  `json:"highlightMatch"`
	ShowSpecialSymbols       bool    `json:"showSpecialSymbols"`
	ShowInStatusBar          bool    `json:"showInStatusBar"`
	MenuIcon                 string  `json:"menuIcon"`
	ShowRecentCopyInMenuBar  bool    `json:"showRecentCopyInMenuBar"`
	ShowSearch               bool    `json:"showSearch"`
	SearchVisibility         string  `json:"searchVisibility"`
	ShowTitle                bool    `json:"showTitle"`
	ShowApplicationIcons     bool    `json:"showApplicationIcons"`
	ShowHexColorSwatch       bool    `json:"showHexColorSwatch"`
	ShowFooter               bool    `json:"showFooter"`
	WindowWidth              int     `json:"windowWidth"`
	WindowHeight             int     `json:"windowHeight"`
	WindowPositionX          float64 `json:"windowPositionX"` // fraction of the screen, top middle anchor
	WindowPositionY          float64 `json:"windowPositionY"`
	PreviewWidth             int     `json:"previewWidth"`

	// Ignore
	IgnoredApps               []string `json:"ignoredApps"`
	IgnoreAllAppsExceptListed bool     `json:"ignoreAllAppsExceptListed"`
	IgnoredPasteboardTypes    []string `json:"ignoredPasteboardTypes"`
	IgnoreRegexp              []string `json:"ignoreRegexp"`

	// Advanced
	IgnoreEvents           bool `json:"ignoreEvents"`
	ClearOnQuit            bool `json:"clearOnQuit"`
	ClearSystemClipboard   bool `json:"clearSystemClipboard"`
	SuppressClearAlert     bool `json:"suppressClearAlert"`
	ClipboardCheckInterval int  `json:"clipboardCheckIntervalMs"`
}

// Default returns Maccy's defaults.
func Default() Settings {
	return Settings{
		PopupShortcut:   DefaultPopupShortcut,
		PinShortcut:     DefaultPinShortcut,
		DeleteShortcut:  DefaultDeleteShortcut,
		PreviewShortcut: DefaultPreviewShortcut,
		SearchMode:      SearchExact,

		EnabledPasteboardTypes: slices.Clone(pasteboard.DefaultEnabledPasteboardTypes),
		HistorySize:            200,
		SortBy:                 SortLastCopiedAt,

		PopupPosition:            PopupCursor,
		PinTo:                    PinTop,
		ImageMaxHeight:           40,
		PreviewDelay:             1500,
		OpenPreviewAutomatically: true,
		HighlightMatch:           HighlightBold,
		ShowSpecialSymbols:       true,
		ShowInStatusBar:          true,
		MenuIcon:                 IconMyClip,
		ShowSearch:               true,
		SearchVisibility:         SearchAlways,
		ShowTitle:                true,
		ShowHexColorSwatch:       true,
		ShowFooter:               true,
		WindowWidth:              450,
		WindowHeight:             800,
		WindowPositionX:          0.5,
		WindowPositionY:          0.8,
		PreviewWidth:             400,

		IgnoredApps:            []string{},
		IgnoredPasteboardTypes: slices.Clone(pasteboard.DefaultIgnoredPasteboardTypes),
		IgnoreRegexp:           []string{},

		ClipboardCheckInterval: int(capture.DefaultInterval / time.Millisecond),
	}
}

// Normalize replaces out-of-range or unknown values with their defaults.
func (s *Settings) Normalize() {
	d := Default()
	oneOf := func(v *string, choices []string, def string) {
		if !slices.Contains(choices, *v) {
			*v = def
		}
	}
	oneOf(&s.SearchMode, SearchModes, d.SearchMode)
	oneOf(&s.SearchVisibility, SearchVisibility, d.SearchVisibility)
	oneOf(&s.HighlightMatch, HighlightStyles, d.HighlightMatch)
	oneOf(&s.PinTo, PinPositions, d.PinTo)
	oneOf(&s.PopupPosition, PopupPositions, d.PopupPosition)
	oneOf(&s.SortBy, SortOrders, d.SortBy)
	oneOf(&s.MenuIcon, MenuIcons, d.MenuIcon)
	clamp := func(v *int, lo, hi, def int) {
		if *v == 0 {
			*v = def
		}
		*v = max(lo, min(*v, hi))
	}
	clamp(&s.HistorySize, MinHistorySize, MaxHistorySize, d.HistorySize)
	clamp(&s.ImageMaxHeight, MinImageHeight, MaxImageHeight, d.ImageMaxHeight)
	clamp(&s.PreviewDelay, MinPreviewDelay, MaxPreviewDelay, d.PreviewDelay)
	clamp(&s.ClipboardCheckInterval, MinCheckIntervalMs, MaxCheckIntervalMs, d.ClipboardCheckInterval)
	clamp(&s.WindowWidth, 300, 2000, d.WindowWidth)
	clamp(&s.WindowHeight, 200, 2000, d.WindowHeight)
	clamp(&s.PreviewWidth, 200, 1000, d.PreviewWidth)
	if s.PopupScreen < 0 {
		s.PopupScreen = 0
	}
	if s.WindowPositionX < 0 || s.WindowPositionX > 1 || s.WindowPositionY < 0 || s.WindowPositionY > 1 {
		s.WindowPositionX, s.WindowPositionY = d.WindowPositionX, d.WindowPositionY
	}
	// The popup shortcut may be empty (none); in-popup ones fall back.
	orDefault := func(v *string, def string) {
		if *v == "" {
			*v = def
		}
	}
	orDefault(&s.PinShortcut, d.PinShortcut)
	orDefault(&s.DeleteShortcut, d.DeleteShortcut)
	orDefault(&s.PreviewShortcut, d.PreviewShortcut)
	if s.EnabledPasteboardTypes == nil {
		s.EnabledPasteboardTypes = d.EnabledPasteboardTypes
	}
	if s.IgnoredPasteboardTypes == nil {
		s.IgnoredPasteboardTypes = d.IgnoredPasteboardTypes
	}
	if s.IgnoredApps == nil {
		s.IgnoredApps = []string{}
	}
	if s.IgnoreRegexp == nil {
		s.IgnoreRegexp = []string{}
	}
}

// CapturePrefs are the settings the capturer reads. Pause state stays with
// the capturer (ignore-next is never persisted).
func (s Settings) CapturePrefs() capture.Prefs {
	return capture.Prefs{
		IgnoreEvents:              s.IgnoreEvents,
		EnabledPasteboardTypes:    slices.Clone(s.EnabledPasteboardTypes),
		IgnoredPasteboardTypes:    slices.Clone(s.IgnoredPasteboardTypes),
		IgnoredApps:               slices.Clone(s.IgnoredApps),
		IgnoreAllAppsExceptListed: s.IgnoreAllAppsExceptListed,
		IgnoreRegexp:              slices.Clone(s.IgnoreRegexp),
	}
}

// CheckInterval is the pasteboard poll cadence.
func (s Settings) CheckInterval() time.Duration {
	return time.Duration(s.ClipboardCheckInterval) * time.Millisecond
}

// Load reads path. A missing file gives the defaults; a corrupt one gives
// the defaults and the error, so the caller can report it and keep going.
func Load(path string) (Settings, error) {
	s := Default()
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return s, err
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return Default(), err
	}
	s.Normalize()
	return s, nil
}

// Save writes s to path atomically, mode 0600.
func Save(path string, s Settings) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".settings-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(append(b, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
