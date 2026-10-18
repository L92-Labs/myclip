//go:build linux

package native

import (
	"errors"
	"net/url"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ebitengine/purego"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/transfer"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// ErrUnsupported is returned where Linux has no equivalent.
var ErrUnsupported = errors.New("native: not supported on this platform")

// Linux clipboard targets (MIME types and X11 atoms) MyClip maps to the
// macOS type names the capture rules use. Types not listed keep their own
// name, so an ignore list can name them too.
var (
	textTargets = []string{"UTF8_STRING", "text/plain;charset=utf-8", "text/plain", "STRING", "TEXT", "COMPOUND_TEXT"}
	rtfTargets  = []string{"text/rtf", "application/rtf", "text/richtext"}

	// KDE's convention, which KeePassXC and others set on secrets: its mere
	// presence marks the copy concealed, before any data is read.
	passwordHint = "x-kde-passwordManagerHint"
	ownWrite     = "application/x-myclip-own-write"
	stackWrite   = "application/x-myclip-stack-" // + nonce
)

// uti maps a Linux target to the type MyClip's rules use.
func uti(target string) string {
	switch {
	case slices.Contains(textTargets, target):
		return pasteboard.TypeString
	case slices.Contains(rtfTargets, target):
		return pasteboard.TypeRTF
	case target == "text/html":
		return pasteboard.TypeHTML
	case target == "image/png":
		return pasteboard.TypePNG
	case target == "image/jpeg":
		return pasteboard.TypeJPEG
	case target == "image/tiff":
		return pasteboard.TypeTIFF
	case target == "text/uri-list":
		return pasteboard.TypeFileURL
	case strings.HasPrefix(target, "image/"):
		// Other image targets (BMP, HEIC…) are read converted to PNG.
		return pasteboard.TypePNG
	case target == passwordHint:
		return pasteboard.TypeConcealed
	case target == ownWrite:
		return pasteboard.TypeOwnWrite
	case strings.HasPrefix(target, stackWrite):
		return pasteboard.StackWritePrefix + strings.TrimPrefix(target, stackWrite)
	}
	return target
}

var (
	ownerOnce   sync.Once
	ownerErr    error
	changeCount atomic.Int64
	// ownerChanged must stay referenced: GTK calls it for the process's life.
	ownerChanged uintptr
)

// watchOwner counts CLIPBOARD owner changes (GTK's owner-change, via
// XFixes on X11): the Linux stand-in for NSPasteboard.changeCount. Call it
// on the main thread once GTK is up.
func watchOwner() error {
	ownerOnce.Do(func() {
		const flags = purego.RTLD_GLOBAL | purego.RTLD_NOW
		gtk, err := purego.Dlopen("libgtk-3.so.0", flags)
		if err != nil {
			ownerErr = err
			return
		}
		gobject, err := purego.Dlopen("libgobject-2.0.so.0", flags)
		if err != nil {
			ownerErr = err
			return
		}
		var (
			gdkAtomIntern    func(name string, onlyIfExists int32) uintptr
			gtkClipboardGet  func(atom uintptr) uintptr
			gSignalConnectDt func(instance uintptr, signal string, handler, data, destroy uintptr, flags int32) uint64
		)
		gdk, err := purego.Dlopen("libgdk-3.so.0", flags)
		if err != nil {
			ownerErr = err
			return
		}
		purego.RegisterLibFunc(&gdkAtomIntern, gdk, "gdk_atom_intern")
		purego.RegisterLibFunc(&gtkClipboardGet, gtk, "gtk_clipboard_get")
		purego.RegisterLibFunc(&gSignalConnectDt, gobject, "g_signal_connect_data")
		clip := gtkClipboardGet(gdkAtomIntern("CLIPBOARD", 0))
		if clip == 0 {
			ownerErr = errors.New("native: no GTK clipboard")
			return
		}
		ownerChanged = purego.NewCallback(func(clipboard, event, data uintptr) {
			changeCount.Add(1)
		})
		if gSignalConnectDt(clip, "owner-change", ownerChanged, 0, 0, 0) == 0 {
			ownerErr = errors.New("native: cannot watch the clipboard owner")
		}
	})
	return ownerErr
}

// General is the CLIPBOARD selection, read and written through MyGo.
type General struct{}

// NewGeneral starts watching the clipboard. Call it on the main thread.
func NewGeneral() (*General, error) {
	if err := watchOwner(); err != nil {
		return nil, err
	}
	return &General{}, nil
}

// ChangeCount counts clipboard owner changes since launch.
func (General) ChangeCount() int { return int(changeCount.Load()) }

// Snapshot lists the offered targets, raw and mapped, reading no data.
func (g General) Snapshot() (pasteboard.Snapshot, error) {
	s := pasteboard.Snapshot{ChangeCount: g.ChangeCount()}
	seen := map[string]bool{}
	for _, t := range mygo.Clipboard.AvailableFormats() {
		for _, name := range []string{t, uti(t)} {
			if !seen[name] {
				seen[name] = true
				s.Types = append(s.Types, name)
			}
		}
	}
	return s, nil
}

// ReadItems reads the representations MyClip keeps: text, rich text,
// HTML, images and files. Files become one item each, as on macOS.
func (g General) ReadItems() []pasteboard.Item {
	offered := mygo.Clipboard.AvailableFormats()
	has := func(ts ...string) bool {
		return slices.ContainsFunc(offered, func(o string) bool { return slices.Contains(ts, o) })
	}
	main := pasteboard.Item{Types: map[string][]byte{}}
	read := func(f transfer.Format, typ string) {
		if b, err := mygo.Clipboard.ReadFormat(f); err == nil {
			main.Types[typ] = b
		}
	}
	if has(textTargets...) {
		read(transfer.Text, pasteboard.TypeString)
	}
	if has("text/html") {
		read(transfer.HTML, pasteboard.TypeHTML)
	}
	for _, t := range rtfTargets {
		if has(t) {
			read(transfer.Format(t), pasteboard.TypeRTF)
			break
		}
	}
	switch {
	case has("image/png"):
		read(transfer.PNG, pasteboard.TypePNG)
	case has("image/jpeg"):
		read("image/jpeg", pasteboard.TypeJPEG)
	case slices.ContainsFunc(offered, func(o string) bool { return strings.HasPrefix(o, "image/") }):
		read(transfer.PNG, pasteboard.TypePNG) // MyGo converts other images to PNG
	}
	if has(ownWrite) {
		main.Types[pasteboard.TypeOwnWrite] = []byte{}
	}
	var items []pasteboard.Item
	if has("text/uri-list") {
		if b, err := mygo.Clipboard.ReadFormat(transfer.URIList); err == nil {
			for _, line := range strings.Split(string(b), "\n") {
				line = strings.TrimSpace(line)
				if u, err := url.Parse(line); err == nil && u.Scheme == "file" {
					items = append(items, pasteboard.Item{Types: map[string][]byte{pasteboard.TypeFileURL: []byte(line)}})
				}
			}
		}
	}
	if len(items) > 0 {
		// A file copy's text is its path list: the file items stand for it.
		if _, ok := main.Types[pasteboard.TypeOwnWrite]; ok {
			items[0].Types[pasteboard.TypeOwnWrite] = []byte{}
		}
		return items
	}
	if len(main.Types) == 0 {
		return nil
	}
	return []pasteboard.Item{main}
}

// Write replaces the clipboard. MyClip owns the selection afterwards and
// serves it until another app copies.
func (General) Write(items []pasteboard.Item) error {
	if len(items) == 0 {
		mygo.Clipboard.Clear()
		return nil
	}
	var reps []transfer.Representation
	var uris []string
	added := map[transfer.Format]bool{}
	add := func(f transfer.Format, b []byte) {
		if !added[f] {
			added[f] = true
			reps = append(reps, transfer.Bytes(f, b))
		}
	}
	for _, it := range items {
		for typ, v := range it.Types {
			switch typ {
			case pasteboard.TypeString:
				add(transfer.Text, v)
			case pasteboard.TypeHTML:
				add(transfer.HTML, v)
			case pasteboard.TypeRTF:
				add("text/rtf", v)
			case pasteboard.TypePNG:
				add(transfer.PNG, v)
			case pasteboard.TypeJPEG:
				add("image/jpeg", v)
			case pasteboard.TypeTIFF:
				add("image/tiff", v)
			case pasteboard.TypeFileURL:
				uris = append(uris, string(v))
			case pasteboard.TypeOwnWrite:
				add(transfer.Format(ownWrite), []byte{1})
			default:
				if nonce, ok := strings.CutPrefix(typ, pasteboard.StackWritePrefix); ok {
					add(transfer.Format(stackWrite+nonce), []byte{1})
				}
			}
		}
	}
	if len(uris) > 0 {
		add(transfer.URIList, []byte(strings.Join(uris, "\r\n")+"\r\n"))
		if !added[transfer.Text] {
			var paths []string
			for _, u := range uris {
				if p, err := url.Parse(u); err == nil {
					paths = append(paths, p.Path)
				}
			}
			add(transfer.Text, []byte(strings.Join(paths, "\n")))
		}
	}
	return mygo.Clipboard.Write(transfer.New(transfer.NewItem(reps...)))
}
