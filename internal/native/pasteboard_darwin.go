//go:build darwin

// Package darwin binds the AppKit calls MyClip needs through purego (no cgo):
// NSPasteboard changeCount, raw type enumeration, data reads and writes, and
// the frontmost application's bundle ID.
package native

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

var (
	loadOnce sync.Once
	loadErr  error

	poolPush func() uintptr
	poolPop  func(uintptr)

	clsPasteboard, clsPasteboardItem, clsMutableArray, clsData, clsString, clsWorkspace objc.ID

	selGeneral, selChangeCount, selTypes, selItems, selCount, selObjectAtIndex,
	selDataForType, selUTF8String, selLength, selBytes, selClearContents,
	selWriteObjects, selSetDataForType, selAlloc, selInit, selRelease, selArray,
	selAddObject, selDataWithBytes, selStringWithUTF8, selShared, selFrontmost,
	selBundleID objc.SEL
)

func load() error {
	loadOnce.Do(func() {
		const flags = purego.RTLD_GLOBAL | purego.RTLD_NOW
		if _, err := purego.Dlopen("/System/Library/Frameworks/AppKit.framework/AppKit", flags); err != nil {
			loadErr = fmt.Errorf("native: load AppKit: %w", err)
			return
		}
		lib, err := purego.Dlopen("/usr/lib/libobjc.A.dylib", flags)
		if err != nil {
			loadErr = fmt.Errorf("native: load libobjc: %w", err)
			return
		}
		purego.RegisterLibFunc(&poolPush, lib, "objc_autoreleasePoolPush")
		purego.RegisterLibFunc(&poolPop, lib, "objc_autoreleasePoolPop")

		cls := func(name string) objc.ID {
			c := objc.ID(objc.GetClass(name))
			if c == 0 && loadErr == nil {
				loadErr = fmt.Errorf("native: class %s not found", name)
			}
			return c
		}
		clsPasteboard = cls("NSPasteboard")
		clsPasteboardItem = cls("NSPasteboardItem")
		clsMutableArray = cls("NSMutableArray")
		clsData = cls("NSData")
		clsString = cls("NSString")
		clsWorkspace = cls("NSWorkspace")

		s := objc.RegisterName
		selGeneral = s("generalPasteboard")
		selChangeCount = s("changeCount")
		selTypes = s("types")
		selItems = s("pasteboardItems")
		selCount = s("count")
		selObjectAtIndex = s("objectAtIndex:")
		selDataForType = s("dataForType:")
		selUTF8String = s("UTF8String")
		selLength = s("length")
		selBytes = s("bytes")
		selClearContents = s("clearContents")
		selWriteObjects = s("writeObjects:")
		selSetDataForType = s("setData:forType:")
		selAlloc = s("alloc")
		selInit = s("init")
		selRelease = s("release")
		selArray = s("array")
		selAddObject = s("addObject:")
		selDataWithBytes = s("dataWithBytes:length:")
		selStringWithUTF8 = s("stringWithUTF8String:")
		selShared = s("sharedWorkspace")
		selFrontmost = s("frontmostApplication")
		selBundleID = s("bundleIdentifier")
	})
	return loadErr
}

// withPool runs fn in an autorelease pool, pinned to one OS thread: a pool
// popped on another thread crashes.
func withPool(fn func()) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	p := poolPush()
	defer poolPop(p)
	fn()
}

func nsString(s string) objc.ID { return clsString.Send(selStringWithUTF8, s) }

func goString(str objc.ID) string {
	if str == 0 {
		return ""
	}
	p := str.Send(selUTF8String)
	if p == 0 {
		return ""
	}
	// Reinterpret the returned C pointer without a uintptr→Pointer
	// conversion, which vet rejects.
	b := *(**byte)(unsafe.Pointer(&p))
	n := 0
	for *(*byte)(unsafe.Add(unsafe.Pointer(b), n)) != 0 {
		n++
	}
	return string(unsafe.Slice(b, n))
}

func goBytes(data objc.ID) []byte {
	if data == 0 {
		return nil
	}
	n := objc.Send[uint](data, selLength)
	if n == 0 {
		return []byte{}
	}
	p := data.Send(selBytes)
	return append([]byte(nil), unsafe.Slice(*(**byte)(unsafe.Pointer(&p)), n)...)
}

func nsData(b []byte) objc.ID {
	var p unsafe.Pointer
	if len(b) > 0 {
		p = unsafe.Pointer(&b[0])
	}
	d := clsData.Send(selDataWithBytes, p, uint(len(b)))
	runtime.KeepAlive(b)
	return d
}

func arrayStrings(arr objc.ID) []string {
	if arr == 0 {
		return nil
	}
	n := objc.Send[uint](arr, selCount)
	out := make([]string, 0, n)
	for i := uint(0); i < n; i++ {
		out = append(out, goString(arr.Send(selObjectAtIndex, i)))
	}
	return out
}

// General is the system pasteboard. Its methods may be called from any
// goroutine; MyClip calls them on the main thread, as AppKit prefers.
type General struct{}

// NewGeneral loads AppKit and returns the general pasteboard.
func NewGeneral() (*General, error) {
	if err := load(); err != nil {
		return nil, err
	}
	return &General{}, nil
}

func (General) board() objc.ID { return clsPasteboard.Send(selGeneral) }

// ChangeCount reads NSPasteboard.changeCount; it reads no contents.
func (g General) ChangeCount() int {
	var n int
	withPool(func() { n = objc.Send[int](g.board(), selChangeCount) })
	return n
}

// Snapshot enumerates NSPasteboard.types ∪ every item's types (N6 union)
// without reading any data.
func (g General) Snapshot() (pasteboard.Snapshot, error) {
	var s pasteboard.Snapshot
	withPool(func() {
		pb := g.board()
		s.ChangeCount = objc.Send[int](pb, selChangeCount)
		seen := map[string]bool{}
		add := func(ts []string) {
			for _, t := range ts {
				if !seen[t] {
					seen[t] = true
					s.Types = append(s.Types, t)
				}
			}
		}
		add(arrayStrings(pb.Send(selTypes)))
		items := pb.Send(selItems)
		if items == 0 {
			return
		}
		for i := uint(0); i < objc.Send[uint](items, selCount); i++ {
			add(arrayStrings(items.Send(selObjectAtIndex, i).Send(selTypes)))
		}
	})
	return s, nil
}

// ReadItems copies every type's data of every pasteboard item.
func (g General) ReadItems() []pasteboard.Item {
	var out []pasteboard.Item
	withPool(func() {
		items := g.board().Send(selItems)
		if items == 0 {
			return
		}
		n := objc.Send[uint](items, selCount)
		for i := uint(0); i < n; i++ {
			pi := items.Send(selObjectAtIndex, i)
			types := pi.Send(selTypes)
			it := pasteboard.Item{Types: map[string][]byte{}}
			for j := uint(0); j < objc.Send[uint](types, selCount); j++ {
				t := types.Send(selObjectAtIndex, j)
				it.Types[goString(t)] = goBytes(pi.Send(selDataForType, t))
			}
			out = append(out, it)
		}
	})
	return out
}

// Write replaces the pasteboard with items, one NSPasteboardItem each.
func (g General) Write(items []pasteboard.Item) error {
	var ok bool
	withPool(func() {
		pb := g.board()
		objc.Send[int](pb, selClearContents)
		arr := clsMutableArray.Send(selArray)
		for _, it := range items {
			pi := clsPasteboardItem.Send(selAlloc).Send(selInit)
			for t, v := range it.Types {
				objc.Send[bool](pi, selSetDataForType, nsData(v), nsString(t))
			}
			arr.Send(selAddObject, pi)
			pi.Send(selRelease) // the array holds it
		}
		ok = objc.Send[bool](pb, selWriteObjects, arr)
	})
	if !ok {
		return fmt.Errorf("native: pasteboard rejected the write")
	}
	return nil
}

// FrontmostBundleID is the bundle ID of the active application, or "".
func FrontmostBundleID() string {
	if load() != nil {
		return ""
	}
	var id string
	withPool(func() {
		app := clsWorkspace.Send(selShared).Send(selFrontmost)
		if app != 0 {
			id = goString(app.Send(selBundleID))
		}
	})
	return id
}

// Prepare runs before the app starts; nothing to do on macOS.
func Prepare() {}
