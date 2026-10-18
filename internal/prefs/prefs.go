// Package prefs is MyClip's settings window, after Maccy 2.7.1's panes:
// General, Storage, Appearance, Pins, Ignore and Advanced. It edits the
// app's Settings in place and tells its Host after each change.
package prefs

import (
	"fmt"
	"regexp"
	"runtime"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
	"github.com/L92-Labs/myclip/internal/settings"
	"github.com/L92-Labs/myclip/internal/shortcut"
	"github.com/L92-Labs/myclip/locales"
)

// Host applies and backs the settings. Its methods run on the main thread.
type Host interface {
	// Changed applies and saves the settings after an edit.
	Changed()
	// Recording pauses the global shortcut while one is being recorded.
	Recording(on bool)
	PinnedItems() []*history.Item
	// PinKeys are the keys the item may take, its own included.
	PinKeys(id int64) []string
	SetPin(id int64, key string) error
	SetTitle(id int64, title string) error
	SetText(id int64, text string) error
	Unpin(id int64)
	LaunchAtLogin() bool
	SetLaunchAtLogin(on bool) error
	AccessibilityTrusted() bool
	RequestAccessibility()
	StorageSize() string
	// PickApp asks for an application to ignore and adds it.
	PickApp()
	AppName(bundleID string) string
	AppIcon(bundleID string) *ui.Bitmap
	Displays() []string
}

// Content type groups of the Storage pane (StorageType @ eb03eba).
var (
	filesTypes  = []string{pasteboard.TypeFileURL}
	imagesTypes = []string{pasteboard.TypePNG, pasteboard.TypeTIFF, pasteboard.TypeJPEG, pasteboard.TypeHEIC}
	textTypes   = []string{pasteboard.TypeHTML, pasteboard.TypeRTF, pasteboard.TypeString}
)

// View is the settings window's state.
type View struct {
	Host Host
	S    *settings.Settings

	tab, ignoreTab int
	recording      string // the shortcut being recorded

	// Mirrors of settings for widgets that edit their own types.
	historySize, imageHeight, previewDelay, interval, width, height float64
	files, images, text, loginItem                                  bool
	pinTo                                                           int

	newType, newRegexp, regexpErr string
	newApp                        string // Linux: a WM_CLASS to ignore
	titles, texts                 map[int64]*string
}

// Load reads the settings into the widgets, as the window opens.
func (v *View) Load() {
	s := v.S
	v.historySize = float64(s.HistorySize)
	v.imageHeight = float64(s.ImageMaxHeight)
	v.previewDelay = float64(s.PreviewDelay)
	v.interval = float64(s.ClipboardCheckInterval)
	v.width, v.height = float64(s.WindowWidth), float64(s.WindowHeight)
	has := func(types []string) bool {
		return slices.ContainsFunc(types, func(t string) bool { return slices.Contains(s.EnabledPasteboardTypes, t) })
	}
	v.files, v.images, v.text = has(filesTypes), has(imagesTypes), has(textTypes)
	v.pinTo = 0
	if s.PinTo == settings.PinBottom {
		v.pinTo = 1
	}
	v.loginItem = v.Host.LaunchAtLogin()
	v.titles, v.texts = map[int64]*string{}, map[int64]*string{}
	v.recording = ""
}

func (v *View) changed() { v.Host.Changed() }

// View builds the window.
func (v *View) View(c *ui.Context) {
	t := c.Theme()
	c.Root().Background(t.Background)
	c.OnShortcut(ui.Cmd, ui.KeyW, func() { v.stopRecording() })
	ui.Column(c).Fill().Children(func() {
		ui.Row(c).Padding(12, 16, 0, 16).Justify(ui.Center).Children(func() {
			if ui.Tabs(c, &v.tab, locales.PrefGeneral(), locales.PrefStorage(), locales.PrefAppearance(),
				locales.PrefPins(), locales.PrefIgnore(), locales.PrefAdvanced()).Changed() {
				v.stopRecording()
			}
		})
		ui.Scroll(c.Key(v.tab)).Grow(1).Children(func() {
			ui.Column(c).Padding(16, 24).Gap(12).Children(func() {
				switch v.tab {
				case 0:
					v.general(c)
				case 1:
					v.storage(c)
				case 2:
					v.appearance(c)
				case 3:
					v.pins(c)
				case 4:
					v.ignore(c)
				case 5:
					v.advanced(c)
				}
			})
		})
	})
}

// SetTab shows a pane: 0 General … 5 Advanced.
func (v *View) SetTab(i int) { v.tab = max(0, min(i, 5)) }

// StopRecording ends a shortcut recording, as the window hides.
func (v *View) StopRecording() { v.stopRecording() }

func (v *View) stopRecording() {
	if v.recording != "" {
		v.recording = ""
		v.Host.Recording(false)
	}
}

func muted(c *ui.Context, s string) {
	ui.Text(c, s).FontSize(11).TextColor(c.Theme().TextMuted).Wrap()
}

// check is a check box writing to a bool setting.
func (v *View) check(c *ui.Context, on *bool, label string) ui.Element {
	cb := ui.Checkbox(c, on, label)
	if cb.Changed() {
		v.changed()
	}
	return cb
}

// pick is a drop-down choosing one of values, shown by labels.
func (v *View) pick(c *ui.Context, value *string, values, labels []string) {
	i := max(0, slices.Index(values, *value))
	label := labels[i]
	if ui.Select(c, &label, labels).Changed() {
		if j := slices.Index(labels, label); j >= 0 {
			*value = values[j]
			v.changed()
		}
	}
}

// number is a number input writing an int setting.
func (v *View) number(c *ui.Context, mirror *float64, target *int, lo, hi, step float64) {
	if ui.NumberInput(c, mirror, lo, hi, step).Width(150).Changed() {
		*target = int(*mirror)
		v.changed()
	}
}

func (v *View) general(c *ui.Context) {
	s := v.S
	ui.Form(c, func() {
		ui.Field(c, locales.PrefStartup(), func() {
			if ui.Checkbox(c, &v.loginItem, locales.PrefLaunchAtLogin()).Changed() {
				if err := v.Host.SetLaunchAtLogin(v.loginItem); err != nil {
					v.loginItem = !v.loginItem
					c.Toast(err.Error())
				}
			}
		})
		ui.Field(c, locales.PrefOpen(), func() { v.recorder(c, "open", &s.PopupShortcut, true) })
		ui.Field(c, locales.PrefPin(), func() { v.recorder(c, "pin", &s.PinShortcut, false) })
		ui.Field(c, locales.PrefDelete(), func() { v.recorder(c, "delete", &s.DeleteShortcut, false) })
		ui.Field(c, locales.PrefPreview(), func() { v.recorder(c, "preview", &s.PreviewShortcut, false) }).
			Description(locales.PrefShortcutsDescription())
		ui.Field(c, locales.PrefSearch(), func() {
			v.pick(c, &s.SearchMode, settings.SearchModes, []string{
				locales.PrefSearchExact(), locales.PrefSearchFuzzy(), locales.PrefSearchRegex(), locales.PrefSearchMixed()})
		})
		ui.Field(c, locales.PrefBehavior(), func() {
			ui.Column(c).Gap(6).Children(func() {
				v.check(c, &s.PasteByDefault, locales.PrefPasteAutomatically())
				v.check(c, &s.RemoveFormattingByDefault, locales.PrefPasteWithoutFormatting())
			})
		}).Description(v.modifiersText())
		ui.Field(c, "", func() {
			ui.Column(c).Gap(6).Children(func() {
				if v.Host.AccessibilityTrusted() {
					muted(c, locales.PrefPasteGranted())
				} else {
					muted(c, locales.PrefPasteAccess())
					if locales.CanRequestPasteAccess() && ui.Button(c, locales.PrefAccessibilityOpen()).AlignSelf(ui.Start).Clicked() {
						v.Host.RequestAccessibility()
					}
				}
			})
		})
		ui.Field(c, locales.PrefNotifications(), func() {
			v.check(c, &s.NotifyOnCopy, locales.PrefNotifyOnCopy())
		})
	})
}

// modifiersText is Maccy's "Modifiers" note for the current preferences.
func (v *View) modifiersText() string {
	mods := func(a, b string) string {
		if !v.S.PasteByDefault {
			return a
		}
		return b
	}
	copyKey := mods("⌘", "⌥")
	pasteKey, plainKey := mods("⌥", "⌘"), mods("⌥⇧", "⌘⇧")
	if v.S.RemoveFormattingByDefault {
		pasteKey, plainKey = plainKey, pasteKey
	}
	return locales.PrefModifiers(copyKey, pasteKey, plainKey)
}

// recorder shows a shortcut, and records a new one after a click: a key
// with ⌘, ⌃ or ⌥; Escape cancels and, where clearable, Backspace clears.
func (v *View) recorder(c *ui.Context, key string, value *string, clearable bool) {
	t := c.Theme()
	recording := v.recording == key
	label := shortcut.Symbols(*value)
	if label == "" {
		label = locales.PrefNone()
	}
	if recording {
		label = locales.PrefRecording()
	}
	box := ui.Row(c.Key("rec-"+key)).Focusable().Padding(4, 10).Radius(6).MinWidth(150).
		Border(1, t.Border).Background(t.Surface).Label(locales.PrefRecord())
	if recording {
		box.Border(2, t.Accent)
	}
	if box.Clicked() && !recording {
		v.stopRecording()
		v.recording = key
		v.Host.Recording(true)
	}
	if recording && !box.Focused() && !box.FocusWithin() {
		box.Focus()
	}
	box.HandleInput(func(ev ui.InputEvent) bool {
		if v.recording != key || ev.Kind != ui.InputKeyDown {
			return v.recording == key && ev.Kind == ui.InputText
		}
		switch {
		case ev.Key == ui.KeyEscape && ev.Mods == 0:
			v.stopRecording()
		case (ev.Key == ui.KeyBackspace || ev.Key == ui.KeyDelete) && ev.Mods == 0:
			if clearable {
				*value = ""
				v.stopRecording()
				v.changed()
			}
		case ev.Mods&(ui.Super|ui.Ctrl|ui.Alt) != 0: // Super is ⌘ (ui.Cmd is Ctrl off macOS)
			ch := shortcut.Chord{Mods: ev.Mods, Key: ev.Key}
			if _, ok := shortcut.Parse(ch.String()); ok {
				*value = ch.String()
				v.stopRecording()
				v.changed()
			}
		}
		return true
	})
	box.Children(func() {
		ui.Text(c, label).FontSize(13).Grow(1)
		if clearable && *value != "" && !recording {
			x := ui.Text(c.Key("clear"), "✕").FontSize(11).TextColor(t.TextMuted).Label(locales.PrefRemove())
			if x.Clicked() {
				*value = ""
				v.changed()
			}
		}
	})
}

func (v *View) storage(c *ui.Context) {
	s := v.S
	ui.Form(c, func() {
		ui.Field(c, locales.PrefSave(), func() {
			ui.Column(c).Gap(6).Children(func() {
				for _, g := range []struct {
					on    *bool
					label string
				}{{&v.files, locales.PrefFiles()}, {&v.images, locales.PrefImages()}, {&v.text, locales.PrefText()}} {
					if ui.Checkbox(c, g.on, g.label).Changed() {
						v.saveTypes()
					}
				}
			})
		}).Description(locales.PrefSaveDescription())
		ui.Field(c, locales.PrefSize(), func() {
			ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
				v.number(c, &v.historySize, &s.HistorySize, settings.MinHistorySize, settings.MaxHistorySize, 1)
				muted(c, locales.PrefSizeOnDisk(v.Host.StorageSize()))
			})
		}).Description(locales.PrefSizeDescription())
		ui.Field(c, locales.PrefSortBy(), func() {
			v.pick(c, &s.SortBy, settings.SortOrders, []string{
				locales.PrefSortLast(), locales.PrefSortFirst(), locales.PrefSortCopies()})
		})
	})
}

// saveTypes writes the Save check boxes to enabledPasteboardTypes,
// keeping the order of Maccy's default set.
func (v *View) saveTypes() {
	var types []string
	for _, t := range pasteboard.DefaultEnabledPasteboardTypes {
		if v.files && slices.Contains(filesTypes, t) || v.images && slices.Contains(imagesTypes, t) ||
			v.text && slices.Contains(textTypes, t) {
			types = append(types, t)
		}
	}
	v.S.EnabledPasteboardTypes = append([]string{}, types...)
	v.changed()
}

func (v *View) appearance(c *ui.Context) {
	s := v.S
	ui.Form(c, func() {
		ui.Field(c, locales.PrefPopupAt(), func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				v.pick(c, &s.PopupPosition, settings.PopupPositions, []string{
					locales.PrefPopupCursor(), locales.PrefPopupMenuIcon(), locales.PrefPopupWindow(),
					locales.PrefPopupScreenCenter(), locales.PrefPopupLast()})
				if s.PopupPosition == settings.PopupLastPosition && ui.Button(c, locales.PrefResetPosition()).Clicked() {
					d := settings.Default()
					s.WindowPositionX, s.WindowPositionY = d.WindowPositionX, d.WindowPositionY
					v.changed()
				}
			})
		})
		if s.PopupPosition == settings.PopupCenter || s.PopupPosition == settings.PopupLastPosition {
			ui.Field(c, locales.PrefScreen(), func() {
				labels := append([]string{locales.PrefActiveScreen()}, v.Host.Displays()...)
				values := make([]string, len(labels))
				for i := range values {
					values[i] = fmt.Sprint(i)
				}
				cur := fmt.Sprint(min(s.PopupScreen, len(labels)-1))
				v.pick(c, &cur, values, labels)
				fmt.Sscan(cur, &s.PopupScreen)
			})
		}
		ui.Field(c, locales.PrefWindowSize(), func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				v.number(c, &v.width, &s.WindowWidth, 300, 2000, 10)
				ui.Text(c, "×")
				v.number(c, &v.height, &s.WindowHeight, 200, 2000, 10)
			})
		})
		ui.Field(c, locales.PrefPinTo(), func() {
			if ui.Segmented(c, &v.pinTo, locales.PrefPinTop(), locales.PrefPinBottom()).Changed() {
				s.PinTo = settings.PinPositions[v.pinTo]
				v.changed()
			}
		})
		ui.Field(c, locales.PrefImageHeight(), func() {
			v.number(c, &v.imageHeight, &s.ImageMaxHeight, settings.MinImageHeight, settings.MaxImageHeight, 1)
		}).Description(locales.PrefImageHeightDescription())
		ui.Field(c, locales.PrefPreviewDelay(), func() {
			ui.Column(c).Gap(6).Children(func() {
				v.number(c, &v.previewDelay, &s.PreviewDelay, settings.MinPreviewDelay, settings.MaxPreviewDelay, 100)
				v.check(c, &s.OpenPreviewAutomatically, locales.PrefOpenPreviewAutomatically())
			})
		}).Description(locales.PrefPreviewDelayDescription())
		ui.Field(c, locales.PrefHighlight(), func() {
			v.pick(c, &s.HighlightMatch, settings.HighlightStyles, []string{
				locales.PrefHighlightBold(), locales.PrefHighlightColor(), locales.PrefHighlightItalic(), locales.PrefHighlightUnderline()})
		})
		ui.Field(c, locales.PrefMenuIconStyle(), func() {
			ui.Column(c).Gap(6).Children(func() {
				v.check(c, &s.ShowInStatusBar, locales.PrefShowMenuIcon())
				if s.ShowInStatusBar {
					v.pick(c, &s.MenuIcon, settings.MenuIcons, []string{
						locales.PrefIconMyClip(), locales.PrefIconClipboard(), locales.PrefIconScissors(), locales.PrefIconPaperclip()})
					v.check(c, &s.ShowRecentCopyInMenuBar, locales.PrefShowRecentCopy())
				} else {
					muted(c, locales.MenuIconHiddenHere())
				}
			})
		})
		ui.Field(c, locales.PrefSearchVisibility(), func() {
			ui.Column(c).Gap(6).Children(func() {
				v.check(c, &s.ShowSearch, locales.PrefShowSearchField())
				if s.ShowSearch {
					v.pick(c, &s.SearchVisibility, settings.SearchVisibility, []string{
						locales.PrefSearchAlways(), locales.PrefSearchDuringSearch()})
					v.check(c, &s.ShowTitle, locales.PrefShowTitle())
				}
			})
		})
		ui.Field(c, locales.PrefShow(), func() {
			ui.Column(c).Gap(6).Children(func() {
				v.check(c, &s.ShowSpecialSymbols, locales.PrefShowSpecialSymbols())
				v.check(c, &s.ShowApplicationIcons, locales.PrefShowAppIcons())
				v.check(c, &s.ShowHexColorSwatch, locales.PrefShowHexSwatch())
				v.check(c, &s.ShowFooter, locales.PrefShowFooter())
				if !s.ShowFooter {
					muted(c, locales.FooterWarningHere())
				}
			})
		})
	})
}

func (v *View) pins(c *ui.Context) {
	t := c.Theme()
	muted(c, locales.PrefPinsDescription())
	items := v.Host.PinnedItems()
	if len(items) == 0 {
		ui.Text(c, locales.PrefNoPins(shortcut.Symbols(v.S.PinShortcut))).TextColor(t.TextMuted).Padding(12, 0)
		return
	}
	ui.Grid(c).ColumnTracks(ui.Fixed(70), ui.Fr(1), ui.Fr(2), ui.FitContent()).GapX(12).GapY(8).AlignItems(ui.Center).Children(func() {
		for _, h := range []string{locales.PrefKey(), locales.PrefAlias(), locales.PrefContent(), ""} {
			ui.Text(c, h).FontSize(11).Bold().TextColor(t.TextMuted)
		}
		for _, it := range items {
			id := it.ID
			key := it.Pin
			keys := v.Host.PinKeys(id)
			if ui.Select(c.Key(fmt.Sprint("key", id)), &key, keys).Changed() {
				if err := v.Host.SetPin(id, key); err != nil {
					c.Toast(err.Error())
				}
			}
			title := v.titles[id]
			if title == nil {
				s := it.Title
				title = &s
				v.titles[id] = title
			}
			if ui.EditableText(c.Key(fmt.Sprint("title", id)), title).Changed() {
				v.Host.SetTitle(id, *title)
			}
			if textOnly(it) {
				text := v.texts[id]
				if text == nil {
					s := it.Text()
					text = &s
					v.texts[id] = text
				}
				if ui.EditableText(c.Key(fmt.Sprint("text", id)), text).Changed() {
					v.Host.SetText(id, *text)
				}
			} else {
				ui.Text(c, locales.PrefNotText()).FontSize(12).TextColor(t.TextMuted)
			}
			if ui.Button(c.Key(fmt.Sprint("unpin", id)), locales.PrefUnpin()).Clicked() {
				v.Host.Unpin(id)
			}
		}
	})
}

// textOnly reports whether the item's content is text Maccy lets one edit.
func textOnly(it *history.Item) bool {
	if _, _, img := it.Image(); img || len(it.FilePaths()) > 0 {
		return false
	}
	_, ok := it.Value(pasteboard.TypeString)
	return ok
}

func (v *View) ignore(c *ui.Context) {
	s := v.S
	t := c.Theme()
	ui.Segmented(c, &v.ignoreTab, locales.PrefApplicationsTab(), locales.PrefPasteboardTypesTab(), locales.PrefRegexpTab())
	switch v.ignoreTab {
	case 0:
		muted(c, locales.PrefIgnoredAppsDescription())
		v.check(c, &s.IgnoreAllAppsExceptListed, locales.PrefIgnoreAllExceptListed())
		v.list(c, &s.IgnoredApps, func(id string) {
			if icon := v.Host.AppIcon(id); icon != nil {
				ui.Image(c, icon).Size(16, 16)
			}
			ui.Text(c, v.Host.AppName(id)).FontSize(13)
			ui.Text(c, id).FontSize(11).TextColor(t.TextMuted).Grow(1).SingleLine()
		})
		if runtime.GOOS == "linux" {
			// Linux apps are named by their window class (WM_CLASS), typed in.
			ui.Row(c).Gap(8).Children(func() {
				add := func() {
					if name := strings.TrimSpace(v.newApp); name != "" && !slices.Contains(s.IgnoredApps, name) {
						s.IgnoredApps = append(s.IgnoredApps, name)
						v.newApp = ""
						v.changed()
					}
				}
				ui.TextInput(c, &v.newApp).Placeholder(locales.PrefNewApp()).Grow(1).OnSubmit(add)
				ui.Button(c, locales.PrefAdd()).Disabled(strings.TrimSpace(v.newApp) == "").OnClick(add)
			})
		} else if ui.Button(c, locales.PrefAddApp()).AlignSelf(ui.Start).Clicked() {
			v.Host.PickApp()
		}
	case 1:
		muted(c, locales.PrefIgnoredTypesDescription())
		v.list(c, &s.IgnoredPasteboardTypes, func(typ string) { ui.Text(c, typ).FontSize(13).Grow(1) })
		ui.Row(c).Gap(8).Children(func() {
			add := func() {
				if v.newType != "" && !slices.Contains(s.IgnoredPasteboardTypes, v.newType) {
					s.IgnoredPasteboardTypes = append(s.IgnoredPasteboardTypes, v.newType)
					v.newType = ""
					v.changed()
				}
			}
			ui.TextInput(c, &v.newType).Placeholder(locales.PrefNewType()).Grow(1).OnSubmit(add)
			ui.Button(c, locales.PrefAdd()).Disabled(v.newType == "").OnClick(add)
			if ui.Button(c, locales.PrefReset()).Clicked() {
				s.IgnoredPasteboardTypes = slices.Clone(pasteboard.DefaultIgnoredPasteboardTypes)
				v.changed()
			}
		})
	case 2:
		muted(c, locales.PrefIgnoredRegexpsDescription())
		v.list(c, &s.IgnoreRegexp, func(re string) { ui.Text(c, re).FontSize(13).Font("Menlo").Grow(1) })
		ui.Row(c).Gap(8).Children(func() {
			add := func() {
				if _, err := regexp.Compile(v.newRegexp); err != nil {
					v.regexpErr = locales.PrefInvalidRegexp()
					return
				}
				if v.newRegexp != "" && !slices.Contains(s.IgnoreRegexp, v.newRegexp) {
					s.IgnoreRegexp = append(s.IgnoreRegexp, v.newRegexp)
					v.newRegexp, v.regexpErr = "", ""
					v.changed()
				}
			}
			ui.Field(c, "", func() {
				if ui.TextInput(c, &v.newRegexp).Placeholder(locales.PrefNewRegexp()).OnSubmit(add).Changed() {
					v.regexpErr = ""
				}
			}).Error(v.regexpErr).Grow(1)
			ui.Button(c, locales.PrefAdd()).Disabled(v.newRegexp == "").OnClick(add)
		})
	}
}

// list shows strings with a Remove button each.
func (v *View) list(c *ui.Context, items *[]string, show func(string)) {
	t := c.Theme()
	ui.Column(c).Border(1, t.Border).Radius(6).Padding(4).MinHeight(60).Children(func() {
		remove := -1
		for i, s := range *items {
			ui.Row(c.Key(s)).Gap(8).Padding(4, 6).AlignItems(ui.Center).Children(func() {
				show(s)
				if ui.Button(c, locales.PrefRemove()).Clicked() {
					remove = i
				}
			})
		}
		if remove >= 0 {
			*items = slices.Delete(slices.Clone(*items), remove, remove+1)
			v.changed()
		}
	})
}

func (v *View) advanced(c *ui.Context) {
	s := v.S
	ui.Form(c, func() {
		ui.Field(c, locales.PrefTurnOff(), func() {
			v.check(c, &s.IgnoreEvents, locales.PrefTurnOff())
		}).Description(locales.PrefTurnOffDescription() + "\n" + locales.PrefTurnOffViaMenuIcon())
		ui.Field(c, locales.PrefHistory(), func() {
			ui.Column(c).Gap(6).Children(func() {
				v.check(c, &s.ClearOnQuit, locales.PrefClearOnQuit()).Tooltip(locales.PrefClearOnQuitDescription())
				v.check(c, &s.ClearSystemClipboard, locales.PrefClearSystemClipboard()).Tooltip(locales.PrefClearSystemClipboardDescription())
			})
		})
		ui.Field(c, locales.PrefCheckInterval(), func() {
			v.number(c, &v.interval, &s.ClipboardCheckInterval, settings.MinCheckIntervalMs, 60000, 100)
		}).Description(locales.PrefCheckIntervalDescription())
	})
	muted(c, locales.LocalOnlyHere())
}
