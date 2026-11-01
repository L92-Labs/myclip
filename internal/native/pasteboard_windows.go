//go:build windows

package native

import (
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// ErrUnsupported is returned where Windows has no equivalent.
var ErrUnsupported = errors.New("native: not supported on this platform")

var (
	user32      = windows.NewLazySystemDLL("user32.dll")
	kernel32    = windows.NewLazySystemDLL("kernel32.dll")
	shell32     = windows.NewLazySystemDLL("shell32.dll")
	procSeq     = user32.NewProc("GetClipboardSequenceNumber")
	procOpen    = user32.NewProc("OpenClipboard")
	procClose   = user32.NewProc("CloseClipboard")
	procEmpty   = user32.NewProc("EmptyClipboard")
	procEnum    = user32.NewProc("EnumClipboardFormats")
	procGet     = user32.NewProc("GetClipboardData")
	procSet     = user32.NewProc("SetClipboardData")
	procName    = user32.NewProc("GetClipboardFormatNameW")
	procReg     = user32.NewProc("RegisterClipboardFormatW")
	procSize    = kernel32.NewProc("GlobalSize")
	procAlloc   = kernel32.NewProc("GlobalAlloc")
	procFree    = kernel32.NewProc("GlobalFree")
	procLock    = kernel32.NewProc("GlobalLock")
	procUnlck   = kernel32.NewProc("GlobalUnlock")
	procMove    = kernel32.NewProc("RtlMoveMemory")
	procClassEx = user32.NewProc("RegisterClassExW")
	procWindow  = user32.NewProc("CreateWindowExW")
	procDestroy = user32.NewProc("DestroyWindow")
	procDefProc = user32.NewProc("DefWindowProcW")
	procModule  = kernel32.NewProc("GetModuleHandleW")
)

// Standard clipboard formats; registered ones resolve through regFormats.
const (
	cfDIB         = 8
	cfHDrop       = 15
	cfUnicodeText = 13
)

const (
	gmemMoveable = 0x0002
	maxOpenTries = 10
)

// regFormats holds the registered formats MyClip uses, resolved once.
var regFormats struct {
	once sync.Once
	html uintptr // "HTML Format"
	rtf  uintptr // "Rich Text Format"
	png  uintptr // "PNG"
	own  uintptr // com.l92labs.myclip.own-write
}

func initFormats() {
	regFormats.html = registerFormat("HTML Format")
	regFormats.rtf = registerFormat("Rich Text Format")
	regFormats.png = registerFormat("PNG")
	regFormats.own = registerFormat(pasteboard.TypeOwnWrite)
}

func registerFormat(name string) uintptr {
	p, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return 0
	}
	v, _, _ := procReg.Call(uintptr(unsafe.Pointer(p)))
	return v
}

// uti maps a clipboard format to the type name MyClip's rules use. Types
// not listed keep their own name, so an ignore list can name them too.
func uti(format uintptr, name string) string {
	switch format {
	case cfUnicodeText:
		return pasteboard.TypeString
	case cfDIB, cfDIBV5:
		return pasteboard.TypePNG // read converted to PNG
	case cfHDrop:
		return pasteboard.TypeFileURL
	}
	switch name {
	case "HTML Format":
		return pasteboard.TypeHTML
	case "Rich Text Format":
		return pasteboard.TypeRTF
	case "PNG":
		return pasteboard.TypePNG
	case pasteboard.TypeOwnWrite:
		return pasteboard.TypeOwnWrite
	}
	if strings.HasPrefix(name, pasteboard.StackWritePrefix) {
		return name
	}
	return name
}

// General is the Windows clipboard.
type General struct{}

// NewGeneral returns the clipboard board.
func NewGeneral() (*General, error) {
	regFormats.once.Do(initFormats)
	return &General{}, nil
}

// ChangeCount is the system clipboard sequence number.
func (General) ChangeCount() int {
	v, _, _ := procSeq.Call()
	return int(v)
}

// openClipboard retries briefly: the clipboard is a shared resource and
// another process may hold it. hwnd is the owner: writes pass a real
// window, because EmptyClipboard with no owner leaves the clipboard
// unable to hold data (SetClipboardData then fails).
func openClipboard(hwnd uintptr) bool {
	for i := 0; i < maxOpenTries; i++ {
		if v, _, _ := procOpen.Call(hwnd); v != 0 {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

// The clipboard owner: a message-only window. Writes must open the
// clipboard with one, so the emptied clipboard keeps an owner and the
// data published into it sticks.
var (
	ownerOnce   sync.Once
	ownerWindow uintptr
)

// clipboardOwnerWindow makes the message-only window once and returns it
// (0 when creation failed, which Write reports rather than publishing
// ownerless).
func clipboardOwnerWindow() uintptr {
	ownerOnce.Do(func() {
		name, err := windows.UTF16PtrFromString("MyClipClipboard")
		if err != nil {
			return
		}
		wc := wndClassEx{
			Size:      uint32(unsafe.Sizeof(wndClassEx{})),
			ClassName: name,
			WndProc:   windows.NewCallback(clipboardWndProc),
		}
		wc.Instance, _, _ = procModule.Call(0) // the executable itself
		atom, _, _ := procClassEx.Call(uintptr(unsafe.Pointer(&wc)))
		if atom == 0 {
			return
		}
		// (exStyle, class, name, style, x, y, w, h, parent, menu, instance,
		// param) — parent HWND_MESSAGE: invisible, no taskbar presence.
		ownerWindow, _, _ = procWindow.Call(0, atom, 0, 0, 0, 0, 0, 0,
			^uintptr(2), 0, wc.Instance, 0)
	})
	return ownerWindow
}

type wndClassEx struct {
	Size       uint32
	Styles     uint32
	WndProc    uintptr
	Extra      int32
	Window     int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	SmallIcon  uintptr
}

// clipboardWndProc is the owner's window procedure: everything defaults.
func clipboardWndProc(hwnd, msg, w, l uintptr) uintptr {
	v, _, _ := procDefProc.Call(hwnd, msg, w, l)
	return v
}

func (g General) Snapshot() (pasteboard.Snapshot, error) {
	s := pasteboard.Snapshot{ChangeCount: g.ChangeCount()}
	if !openClipboard(0) {
		return s, errors.New("native: the clipboard is busy")
	}
	defer procClose.Call()
	seen := map[string]bool{}
	for f := enumFormat(0); f != 0; f = enumFormat(f) {
		name := formatName(f)
		for _, n := range []string{name, uti(f, name)} {
			if !seen[n] {
				seen[n] = true
				s.Types = append(s.Types, n)
			}
		}
	}
	return s, nil
}

func enumFormat(prev uintptr) uintptr {
	v, _, _ := procEnum.Call(prev)
	return v
}

// cfDIBV5 is the V5 bitmap header format.
const cfDIBV5 = 17

func formatName(f uintptr) string {
	switch f {
	case cfUnicodeText:
		return "CF_UNICODETEXT"
	case cfDIB:
		return "CF_DIB"
	case cfDIBV5:
		return "CF_DIBV5"
	case cfHDrop:
		return "CF_HDROP"
	}
	buf := make([]uint16, 128)
	n, _, _ := procName.Call(f, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if n == 0 {
		return ""
	}
	return windows.UTF16ToString(buf[:n])
}

// readData returns the bytes of one clipboard format, or nil. The locked
// block's address stays a uintptr and RtlMoveMemory copies out of it, so
// no pointer conversion vet would question is needed.
func readData(f uintptr) []byte {
	h, _, _ := procGet.Call(f)
	if h == 0 {
		return nil
	}
	p, _, _ := procLock.Call(h)
	if p == 0 {
		return nil
	}
	defer procUnlck.Call(h)
	n, _, _ := procSize.Call(h)
	if n == 0 || n > 1<<31 {
		return nil
	}
	b := make([]byte, n)
	procMove.Call(uintptr(unsafe.Pointer(unsafe.SliceData(b))), p, n)
	return b
}

func (General) ReadItems() []pasteboard.Item {
	if !openClipboard(0) {
		return nil
	}
	defer procClose.Call()
	main := pasteboard.Item{Types: map[string][]byte{}}
	var fileItems []pasteboard.Item
	for f := enumFormat(0); f != 0; f = enumFormat(f) {
		name := formatName(f)
		if name == "" {
			continue
		}
		if f == cfDIB || f == cfDIBV5 {
			// Bitmaps offered only as a DIB are converted after the loop.
			continue
		}
		switch uti(f, name) {
		case pasteboard.TypeString:
			if _, ok := main.Types[pasteboard.TypeString]; !ok {
				if b := readData(f); b != nil {
					main.Types[pasteboard.TypeString] = []byte(utf16BytesToString(b))
				}
			}
		case pasteboard.TypeHTML:
			if b := readData(f); b != nil {
				main.Types[pasteboard.TypeHTML] = cfhtmlParse(b)
			}
		case pasteboard.TypeRTF:
			if b := readData(f); b != nil {
				main.Types[pasteboard.TypeRTF] = b
			}
		case pasteboard.TypePNG:
			if _, ok := main.Types[pasteboard.TypePNG]; !ok {
				if b := readData(f); b != nil {
					main.Types[pasteboard.TypePNG] = b
				}
			}
		case pasteboard.TypeFileURL:
			if len(fileItems) == 0 {
				if b := readData(f); b != nil {
					for _, path := range dropFiles(b) {
						fileItems = append(fileItems, pasteboard.Item{Types: map[string][]byte{
							pasteboard.TypeFileURL: []byte(fileURL(path)),
						}})
					}
				}
			}
		default:
			if name == pasteboard.TypeOwnWrite || strings.HasPrefix(name, pasteboard.StackWritePrefix) {
				main.Types[name] = []byte{}
			}
		}
	}
	// A bitmap the copy offered only as a DIB is stored converted to PNG.
	if _, ok := main.Types[pasteboard.TypePNG]; !ok {
		for _, f := range []uintptr{cfDIB, cfDIBV5} {
			if b := readData(f); b != nil {
				if pngData, ok := dibToPNG(b); ok {
					main.Types[pasteboard.TypePNG] = pngData
					break
				}
			}
		}
	}
	if len(fileItems) > 0 {
		// A file copy's text is its path list: the file items stand for it.
		for typ := range main.Types {
			if typ == pasteboard.TypeOwnWrite || strings.HasPrefix(typ, pasteboard.StackWritePrefix) {
				fileItems[0].Types[typ] = []byte{}
			}
		}
		return fileItems
	}
	if len(main.Types) == 0 {
		return nil
	}
	return []pasteboard.Item{main}
}

// Write replaces the clipboard contents.
func (General) Write(items []pasteboard.Item) error {
	// Writes need a real owner window: an emptied clipboard without one
	// cannot hold data (SetClipboardData fails silently).
	hwnd := clipboardOwnerWindow()
	if hwnd == 0 {
		return errors.New("native: no clipboard owner window")
	}
	if len(items) == 0 {
		if !openClipboard(hwnd) {
			return errors.New("native: the clipboard is busy")
		}
		defer procClose.Call()
		procEmpty.Call()
		return nil
	}
	if !openClipboard(hwnd) {
		return errors.New("native: the clipboard is busy")
	}
	defer procClose.Call()
	procEmpty.Call()

	var files []string
	var errs []error
	wrote := map[string]bool{}
	// Types the clipboard keeps without data: MyClip's own-write and
	// paste-stack markers, from any item.
	for _, it := range items {
		for typ, v := range it.Types {
			if len(v) == 0 && (typ == pasteboard.TypeOwnWrite || strings.HasPrefix(typ, pasteboard.StackWritePrefix)) {
				wrote[typ] = true
				errs = append(errs, setClipboard(registerFormat(typ), nil))
			}
		}
	}
	// One clipboard item: the first data of each type, plus the files.
	for _, it := range items {
		for typ, v := range it.Types {
			if len(v) == 0 {
				continue
			}
			if typ == pasteboard.TypeFileURL {
				// Every item's file URL lands in the one CF_HDROP; never
				// deduplicated away, or a multi-file copy keeps one file.
				if p, err := url.Parse(string(v)); err == nil && p.Path != "" {
					files = append(files, strings.TrimPrefix(p.Path, "/"))
				}
				continue
			}
			if wrote[typ] {
				continue
			}
			wrote[typ] = true
			switch typ {
			case pasteboard.TypeString:
				errs = append(errs, setClipboard(cfUnicodeText, utf16StringToBytes(string(v))))
			case pasteboard.TypeHTML:
				errs = append(errs, setClipboard(regFormats.html, cfhtmlBuild(v)))
			case pasteboard.TypeRTF:
				errs = append(errs, setClipboard(regFormats.rtf, v))
			case pasteboard.TypePNG:
				errs = append(errs, setClipboard(regFormats.png, v))
				if dib, ok := pngToDIB(v); ok {
					errs = append(errs, setClipboard(cfDIB, dib))
				}
			}
		}
	}
	if len(files) > 0 {
		errs = append(errs, setClipboard(cfHDrop, dropFilesBuild(files)))
	}
	return errors.Join(errs...)
}

// setClipboard puts data on the clipboard, transferring ownership of the
// global allocation; empty data publishes a zero-length block, which is
// how the own-write and paste-stack markers ride along. (NULL would ask
// for delayed rendering, whose success is indistinguishable from failure.)
// A failure is reported: the caller clears the clipboard first, so a
// silent failure would leave it empty while claiming success.
func setClipboard(f uintptr, data []byte) error {
	if f == 0 {
		return errors.New("native: clipboard format did not register")
	}
	h, _, _ := procAlloc.Call(gmemMoveable, uintptr(len(data)))
	if h == 0 {
		return errors.New("native: clipboard allocation failed")
	}
	if len(data) > 0 {
		p, _, _ := procLock.Call(h)
		if p == 0 {
			procFree.Call(h)
			return errors.New("native: clipboard lock failed")
		}
		procMove.Call(p, uintptr(unsafe.Pointer(unsafe.SliceData(data))), uintptr(len(data)))
		procUnlck.Call(h)
	}
	if v, _, _ := procSet.Call(f, h); v == 0 { // the clipboard owns h now
		procFree.Call(h)
		return errors.New("native: SetClipboardData failed")
	}
	return nil
}

// utf16BytesToString decodes UTF-16LE clipboard bytes (dropping the NUL).
// utf16StringToBytes encodes a string as NUL-terminated UTF-16LE.
// Both live in utf16.go, pure and tested everywhere.

// fileURL is a Windows path as a file URL, as storage keeps them.
func fileURL(path string) string {
	u := url.URL{Scheme: "file", Path: "/" + strings.ReplaceAll(path, "\\", "/")}
	return u.String()
}
