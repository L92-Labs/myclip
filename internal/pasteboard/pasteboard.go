package pasteboard

// Item is one pasteboard item with its type→bytes map.
type Item struct {
	Types map[string][]byte // UTI → raw bytes (nil value allowed for presence-only types)
}

// Snapshot is a changeCount-scoped types-only view (N6 enumerate step).
type Snapshot struct {
	ChangeCount int
	// Types is NSPasteboard.types ∪ ⋃ pasteboardItems[i].types (N6 union).
	Types []string
}

// Reader is the pasteboard access seam. Darwin uses NSPasteboard; tests use Fake.
// Type enumeration must not imply data reads; ReadItems does.
type Reader interface {
	ChangeCount() int
	Snapshot() (Snapshot, error)
	ReadItems() []Item
}
