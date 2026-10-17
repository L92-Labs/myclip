package pasteboard

import "sync"

// Fake is an in-memory pasteboard for Linux unit tests.
// Type enumeration does not count as a data read; ReadItems does.
type Fake struct {
	mu sync.Mutex

	changeCount int
	items       []Item
	boardTypes  []string // NSPasteboard.types surface (may be a superset of item types)

	dataReads int
}

func NewFake() *Fake {
	return &Fake{}
}

func (f *Fake) ChangeCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.changeCount
}

func (f *Fake) DataReads() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.dataReads
}

func (f *Fake) ResetReads() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.dataReads = 0
}

// SetItems replaces the pasteboard contents and bumps changeCount.
// boardTypes, if non-nil, sets the NSPasteboard.types surface; otherwise it is
// the union of all item types.
func (f *Fake) SetItems(items []Item, boardTypes []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.items = cloneItems(items)
	if boardTypes != nil {
		f.boardTypes = append([]string(nil), boardTypes...)
	} else {
		f.boardTypes = unionItemTypes(f.items)
	}
	f.changeCount++
}

// Snapshot returns changeCount + N6 type union without reading item bytes.
func (f *Fake) Snapshot() (Snapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return Snapshot{
		ChangeCount: f.changeCount,
		Types:       unionStrings(f.boardTypes, unionItemTypes(f.items)),
	}, nil
}

// ReadItems copies item bytes and increments DataReads once per type value read.
func (f *Fake) ReadItems() []Item {
	f.mu.Lock()
	defer f.mu.Unlock()
	items := make([]Item, len(f.items))
	for i, it := range f.items {
		items[i] = Item{Types: make(map[string][]byte, len(it.Types))}
		for k, v := range it.Types {
			f.dataReads++
			if v == nil {
				items[i].Types[k] = nil
				continue
			}
			cp := make([]byte, len(v))
			copy(cp, v)
			items[i].Types[k] = cp
		}
	}
	return items
}

func cloneItems(in []Item) []Item {
	out := make([]Item, len(in))
	for i, it := range in {
		out[i] = Item{Types: make(map[string][]byte, len(it.Types))}
		for k, v := range it.Types {
			if v == nil {
				out[i].Types[k] = nil
				continue
			}
			cp := make([]byte, len(v))
			copy(cp, v)
			out[i].Types[k] = cp
		}
	}
	return out
}

func unionItemTypes(items []Item) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, it := range items {
		for t := range it.Types {
			if _, ok := seen[t]; ok {
				continue
			}
			seen[t] = struct{}{}
			out = append(out, t)
		}
	}
	return out
}

func unionStrings(a, b []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, xs := range [][]string{a, b} {
		for _, t := range xs {
			if _, ok := seen[t]; ok {
				continue
			}
			seen[t] = struct{}{}
			out = append(out, t)
		}
	}
	return out
}
