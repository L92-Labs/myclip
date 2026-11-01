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
	user32    = windows.NewLazySystemDLL("user32.dll")
	kernel32  = windows.NewLazySystemDLL("kernel32.dll")
	shell32   = windows.NewLazySystemDLL("shell32.dll")
	procSeq   = user32.NewProc("GetClipboardSequenceNumber")
	procOpen  = user32.NewProc("OpenClipboard")
	procClose = user32.NewProc("CloseClipboard")
	procEmpty = user32.NewProc("EmptyClipboard")
	procEnum  = user32.NewProc("EnumClipboardFormats")
	procGet   = user32.NewProc("GetClipboardData")
	procSet   = user32.NewProc("SetClipboardData")
	procName  = user32.NewProc("GetClipboardFormatNameW")
	procReg   = user32.NewProc("RegisterClipboardFormatW")
	procSize  = kernel32.NewProc("GlobalSize")
	procAlloc = kernel32.NewProc("GlobalAlloc")
	procFree  = kernel32.NewProc("GlobalFree")
	procLock  = kernel32.NewProc("GlobalLock")
	procUnlck = kernel32.NewProc("GlobalUnlock")
	procMove  = kernel32.NewProc("RtlMoveMemory")
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
// another process may hold it.
func openClipboard() bool {
	for i := 0; i < maxOpenTries; i++ {
		if v, _, _ := procOpen.Call(0); v != 0 {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

func (g General) Snapshot() (pasteboard.Snapshot, error) {
	s := pasteboard.Snapshot{ChangeCount: g.ChangeCount()}
	if !openClipboard() {
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
	if !openClipboard() {
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
	if len(items) == 0 {
		if !openClipboard() {
			return errors.New("native: the clipboard is busy")
		}
		defer procClose.Call()
		procEmpty.Call()
		return nil
	}
	if !openClipboard() {
		return errors.New("native: the clipboard is busy")
	}
	defer procClose.Call()
	procEmpty.Call()

	var files []string
	wrote := map[string]bool{}
	// Types the clipboard keeps without data: MyClip's own-write and
	// paste-stack markers, from any item.
	for _, it := range items {
		for typ, v := range it.Types {
			if len(v) == 0 && (typ == pasteboard.TypeOwnWrite || strings.HasPrefix(typ, pasteboard.StackWritePrefix)) {
				wrote[typ] = true
				setClipboard(registerFormat(typ), nil)
			}
		}
	}
	// One clipboard item: the first data of each type, plus the files.
	for _, it := range items {
		for typ, v := range it.Types {
			if len(v) == 0 || wrote[typ] {
				continue
			}
			wrote[typ] = true
			switch typ {
			case pasteboard.TypeString:
				setClipboard(cfUnicodeText, utf16StringToBytes(string(v)))
			case pasteboard.TypeHTML:
				setClipboard(regFormats.html, cfhtmlBuild(v))
			case pasteboard.TypeRTF:
				setClipboard(regFormats.rtf, v)
			case pasteboard.TypePNG:
				setClipboard(regFormats.png, v)
				if dib, ok := pngToDIB(v); ok {
					setClipboard(cfDIB, dib)
				}
			case pasteboard.TypeFileURL:
				if p, err := url.Parse(string(v)); err == nil && p.Path != "" {
					files = append(files, strings.TrimPrefix(p.Path, "/"))
				}
			}
		}
	}
	if len(files) > 0 {
		setClipboard(cfHDrop, dropFilesBuild(files))
	}
	return nil
}

// setClipboard puts data on the clipboard, transferring ownership of the
// global allocation; nil data registers the format with no data.
func setClipboard(f uintptr, data []byte) {
	if f == 0 {
		return
	}
	if data == nil {
		procSet.Call(f, 0)
		return
	}
	h, _, _ := procAlloc.Call(gmemMoveable, uintptr(len(data)))
	if h == 0 {
		return
	}
	p, _, _ := procLock.Call(h)
	if p == 0 {
		procFree.Call(h)
		return
	}
	procMove.Call(p, uintptr(unsafe.Pointer(unsafe.SliceData(data))), uintptr(len(data)))
	procUnlck.Call(h)
	procSet.Call(f, h) // the clipboard owns h now
}

// utf16BytesToString decodes UTF-16LE clipboard bytes (dropping the NUL).
// utf16StringToBytes encodes a string as NUL-terminated UTF-16LE.
// Both live in utf16.go, pure and tested everywhere.

// fileURL is a Windows path as a file URL, as storage keeps them.
func fileURL(path string) string {
	u := url.URL{Scheme: "file", Path: "/" + strings.ReplaceAll(path, "\\", "/")}
	return u.String()
}
