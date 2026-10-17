// MyClip is a local-only clipboard history for the macOS menu bar.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/egoist/mygo"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/darwin"
	"github.com/L92-Labs/myclip/internal/engine"
	"github.com/L92-Labs/myclip/internal/instance"
	"github.com/L92-Labs/myclip/internal/store"
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
	lock, err := instance.Lock(dir)
	if errors.Is(err, instance.ErrRunning) {
		fmt.Fprintln(os.Stderr, "myclip: already running")
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
	a := &app{}
	mygo.App.WhenReady(func() {
		if err := a.start(dbPath); err != nil {
			log.Print("myclip: ", err)
			mygo.App.Quit()
		}
	})
	mygo.App.OnWillQuit(func(*mygo.QuitEvent) { a.stop() })
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

// app is the running menu bar app. Its fields are set once in start, on
// the main thread.
type app struct {
	engine *engine.Engine
	store  *store.Store
	tray   *mygo.Tray
	cancel context.CancelFunc
	err    error // last store error, shown in the menu
}

func (a *app) start(dbPath string) error {
	board, err := darwin.NewGeneral()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	if a.store, err = store.Open(ctx, dbPath); err != nil {
		return err
	}
	a.engine, err = engine.New(ctx, engine.Config{
		Board:     board,
		Frontmost: darwin.FrontmostBundleID,
		Store:     a.store,
		RunOnMain: mygo.RunOnMain,
		OnChange:  func() { mygo.RunOnMain(a.refresh) },
		OnError: func(err error) {
			log.Print("myclip: store: ", err)
			mygo.RunOnMain(func() { a.err = err })
		},
	})
	if err != nil {
		return err
	}
	a.tray, err = mygo.NewTray(mygo.TrayOptions{
		Icon:           trayIcon(),
		IconIsTemplate: true,
		ToolTip:        "MyClip",
	})
	if err != nil {
		return err
	}
	a.refresh()
	go a.engine.Run(ctx, capture.DefaultInterval)
	return nil
}

func (a *app) stop() {
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

// refresh rebuilds the tray menu from the history. Main thread only.
func (a *app) refresh() {
	if a.tray == nil {
		return
	}
	a.tray.SetMenu(mygo.NewMenu(menuItems(menuState{
		Items:     a.engine.Items(),
		Paused:    a.engine.Paused(),
		Oversized: a.engine.Oversized(),
		Err:       a.err,
	}, menuActions{
		Copy: func(id int64) {
			if err := a.engine.CopyBack(id); err != nil {
				log.Print("myclip: copy: ", err)
			}
		},
		SetPaused:  a.engine.SetPaused,
		IgnoreNext: func() { a.engine.IgnoreNext(); a.refresh() },
	})))
}
