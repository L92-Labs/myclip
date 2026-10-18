// MyClip is a local-only clipboard history for the macOS menu bar.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/L92-Labs/myclip/internal/engine"
	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/instance"
	"github.com/L92-Labs/myclip/internal/native"
	"github.com/L92-Labs/myclip/internal/popup"
	"github.com/L92-Labs/myclip/internal/prefs"
	"github.com/L92-Labs/myclip/internal/settings"
	"github.com/L92-Labs/myclip/internal/store"
	"github.com/L92-Labs/myclip/locales"
)

const identifier = "com.l92labs.myclip"

// dataDir is where MyClip keeps everything. It is built from the home dir,
// never a literal "~", and replaces MyGo's default userData dir.
func dataDir() (string, error) {
	if runtime.GOOS == "darwin" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", identifier), nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, identifier), nil
}

func main() {
	native.Prepare() // before MyGo starts the toolkit (Linux: X11 under Wayland)
	// ① data root 0700 (+ run/), ② pinned as userData before anything else.
	dir, err := dataDir()
	if err == nil {
		dir, err = instance.DataRoot(dir)
	}
	if err != nil {
		log.Fatal("myclip: data dir: ", err)
	}
	mygo.App.SetPath(mygo.PathUserData, dir)

	// ②b single instance: a second launch exits before any capture work.
	// A second launch, or `myclip --toggle` from a desktop shortcut, shows
	// or hides the running MyClip's popup.
	lock, err := instance.Lock(dir)
	if errors.Is(err, instance.ErrRunning) {
		if err := instance.Send(dir, instance.Toggle); err != nil {
			fmt.Fprintln(os.Stderr, "myclip: already running")
		}
		os.Exit(0)
	}
	if err != nil {
		log.Fatal("myclip: ", err)
	}
	defer lock.Close()

	// ③ packaged apps load only the bundled SQLite, never an env override.
	if mygo.App.IsPackaged() {
		os.Unsetenv("MYGO_SQLITE_LIBRARY")
	}

	// ③b pre-create the history DB 0600 before SQLite can create it 0644.
	dbPath := filepath.Join(dir, "history.db")
	if _, err := store.PrepareFile(dbPath); err != nil {
		log.Fatal("myclip: ", err)
	}

	mygo.App.SetName("MyClip")
	mygo.App.SetActivationPolicy(mygo.ActivationPolicyAccessory)
	mygo.App.OnWindowAllClosed(func() {}) // lives in the menu bar
	a := &app{dir: dir, settingsPath: filepath.Join(dir, settings.FileName),
		showAtLaunch: slices.Contains(os.Args[1:], "--toggle")}
	mygo.App.WhenReady(func() {
		if err := a.start(dbPath); err != nil {
			log.Print("myclip: ", err)
			mygo.App.Quit()
		}
	})
	// Windows hide rather than close, except as the app quits.
	mygo.App.OnBeforeQuit(func(*mygo.QuitEvent) { a.quitting = true })
	mygo.App.OnWillQuit(func(*mygo.QuitEvent) { a.stop() })
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

// app is the running menu bar app. Everything but the engine's poll loop
// runs on the main thread.
type app struct {
	dir          string
	showAtLaunch bool      // started by `myclip --toggle`
	commands     io.Closer // the toggle socket
	quitting     bool
	settingsPath string
	settings     settings.Settings
	applied      settings.Settings // what the app last applied
	// captureSafe: the settings file could not be read, so the user's
	// ignore rules are unknown. Capture stays off and the file is never
	// overwritten until the user resets settings from the menu.
	captureSafe bool
	resetErr    error // why the last settings reset failed, shown in the menu

	engine *engine.Engine
	store  *store.Store
	board  *native.General
	tray   *mygo.Tray
	ctx    context.Context
	cancel context.CancelFunc
	poll   context.CancelFunc // the poll loop, restarted when its interval changes
	err    error              // last store error, shown in the menu

	panel    *mygo.Window
	offBlur  func()
	popup    *popup.Popup
	shortcut string // the registered popup shortcut

	prefsWin *mygo.Window
	prefs    *prefs.View

	images *imageCache
	cycle  cycler
	ocr    ocrQueue
}

func (a *app) start(dbPath string) error {
	var err error
	if a.settings, err = settings.Load(a.settingsPath); err != nil {
		log.Print("myclip: settings unreadable, capture off until reset: ", err)
		a.captureSafe = true
	}
	if a.board, err = native.NewGeneral(); err != nil {
		return err
	}
	a.ctx, a.cancel = context.WithCancel(context.Background())
	if a.store, err = store.Open(a.ctx, dbPath); err != nil {
		return err
	}
	a.engine, err = engine.New(a.ctx, engine.Config{
		Board:     a.board,
		Frontmost: native.FrontmostBundleID,
		Store:     a.store,
		Prefs:     a.settings.CapturePrefs(),
		Limits:    a.limits(),
		RunOnMain: mygo.RunOnMain,
		OnChange:  func() { mygo.RunOnMain(a.refresh) },
		// Every copy made in another app sounds, merged into an old clip or
		// not; MyClip's own selections sound in afterCopy.
		OnCopy: func(*history.Item) { mygo.RunOnMain(a.playCopySound) },
		OnNew: func(it *history.Item) {
			a.ocr.add(it)
			mygo.RunOnMain(func() { a.notify(it) })
		},
		OnError: func(err error) {
			log.Print("myclip: store: ", err)
			mygo.RunOnMain(func() { a.err = err; a.refresh() })
		},
	})
	if err != nil {
		return err
	}
	a.engine.SetCaptureSafe(a.captureSafe)
	a.cycle.held = native.ModifierFlags
	a.ocr.save = a.engine.SetRecognizedText
	a.ocr.item = a.engine.Item
	for _, it := range a.engine.Items() {
		a.ocr.add(it) // images kept from before: read once
	}
	history.ResolveFileURL = native.ResolveFileURL
	a.images = newImageCache(func() {
		if a.panel != nil {
			a.panel.Invalidate()
		}
	})
	a.popup = &popup.Popup{Host: a, Settings: &a.settings}
	a.panel = mygo.NewWindow(mygo.WindowOptions{
		Panel:         true,
		Hidden:        true,
		Frameless:     true,
		Vibrancy:      mygo.VibrancyPopover,
		Width:         a.settings.WindowWidth,
		Height:        400,
		DisableResize: true,
		Content:       ui.View(a.popup.View),
	})
	a.panel.OnClose(func(e *mygo.CloseEvent) {
		if !a.quitting {
			e.PreventDefault()
			a.Close()
		}
	})
	a.offBlur = a.panel.OnBlur(a.Close)

	a.applied = a.settings
	a.applyTray()
	a.registerShortcut()
	a.restartPoll()
	a.refresh()
	a.showOnLaunch(os.Getenv("MYCLIP_SHOW"))
	if c, err := instance.Serve(a.dir, func(string) { mygo.RunOnMain(a.toggleFromCommand) }); err != nil {
		log.Print("myclip: toggle command: ", err)
	} else {
		a.commands = c
	}
	if a.showAtLaunch {
		a.showPopup()
	}
	return nil
}

// showOnLaunch opens a window at launch, for screenshots and manual
// checks: MYCLIP_SHOW=popup, popup-preview or prefs[:tab].
func (a *app) showOnLaunch(what string) {
	switch {
	case what == "popup" || what == "popup-preview":
		a.showPopup()
		a.offBlur() // stays open without the keyboard
		if what == "popup-preview" {
			a.popup.TogglePreview()
		}
		// On the primary display, where screenshots look.
		area := mygo.Screen.PrimaryDisplay().WorkArea
		w, h := a.popup.Fit()
		a.panel.SetPosition(area.X+(area.Width-w)/2, area.Y+(area.Height-h)/2)
		a.logBounds(a.panel)
	case strings.HasPrefix(what, "prefs"):
		a.openPrefs()
		var tab int
		fmt.Sscan(strings.TrimPrefix(what, "prefs:"), &tab)
		a.prefs.SetTab(tab)
		a.logBounds(a.prefsWin)
	}
}

// logBounds prints a window's bounds once it has settled, for screenshots.
func (a *app) logBounds(w *mygo.Window) {
	time.AfterFunc(time.Second, func() {
		mygo.RunOnMain(func() {
			b := w.Bounds()
			fmt.Fprintf(os.Stderr, "bounds %d,%d,%d,%d\n", b.X, b.Y, b.Width, b.Height)
		})
	})
}

func (a *app) stop() {
	if a.engine != nil && a.settings.ClearOnQuit {
		// Maccy's clearOnQuit: unpinned items go, as History.clear().
		a.clearHistory(false)
	}
	if a.commands != nil {
		a.commands.Close()
	}
	if a.cancel != nil {
		a.cancel()
	}
	if a.tray != nil {
		a.tray.Destroy()
	}
	if a.store != nil {
		a.store.Close()
	}
}

func (a *app) limits() engine.Limits {
	l := engine.DefaultLimits
	l.MaxItems = a.settings.HistorySize
	return l
}

// restartPoll (re)starts the pasteboard poll at the settings' interval.
func (a *app) restartPoll() {
	if a.poll != nil {
		a.poll()
	}
	ctx, cancel := context.WithCancel(a.ctx)
	a.poll = cancel
	go a.engine.Run(ctx, a.settings.CheckInterval())
}

// sortedItems is the history in the order the popup shows it.
func (a *app) sortedItems() []*history.Item {
	return history.Sort(a.engine.Items(), a.settings.SortBy, a.settings.PinTo)
}

// refresh redraws what shows the history: the tray and an open popup.
// Main thread only.
func (a *app) refresh() {
	a.images.forget(a.engine.Items())
	a.updateTray()
	if a.panel != nil && a.panel.IsVisible() {
		a.popup.SetItems(a.sortedItems())
		a.panel.Invalidate()
	}
	if a.prefsWin != nil {
		a.prefsWin.Invalidate()
	}
}

// notify shows Maccy's copy notification, when the user turned it on.
func (a *app) notify(it *history.Item) {
	if !a.settings.NotifyOnCopy || !mygo.NotificationsSupported() {
		return
	}
	body := it.DisplayTitle()
	if _, _, ok := it.Image(); ok && it.Title == "" {
		body = locales.Image()
	}
	go func() {
		n := mygo.NewNotification(mygo.NotificationOptions{ID: "copy", Title: "MyClip", Body: body})
		if err := n.Show(); err != nil {
			log.Print("myclip: notification: ", err)
		}
	}()
}

// clearHistory clears the history, and the system clipboard with Maccy's
// clearSystemClipboard.
func (a *app) clearHistory(all bool) {
	if err := a.engine.Clear(all); err != nil {
		log.Print("myclip: clear: ", err)
	}
	if a.settings.ClearSystemClipboard {
		mygo.Clipboard.Clear()
	}
	mygo.ClearNotifications()
}
