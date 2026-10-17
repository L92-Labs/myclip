package history

import "time"

// History is an in-memory history list + sessionLog for ModifiedType merges.
// Persistence is layered on later (SQLite store).
type History struct {
	items      []*Item
	sessionLog map[int]*Item // changeCount → item (Maccy sessionLog)
}

func New() *History {
	return &History{sessionLog: map[int]*Item{}}
}

func (h *History) Items() []*Item {
	out := make([]*Item, len(h.items))
	copy(out, h.items)
	return out
}

func (h *History) Len() int { return len(h.items) }

// FindSimilarItem ports History.findSimilarItem @ eb03eba:
// supersedes first, else isModified via sessionLog.
func (h *History) FindSimilarItem(item *Item) *Item {
	for _, existing := range h.items {
		if existing == item {
			continue
		}
		if existing.Supersedes(item) {
			return existing
		}
	}
	return h.isModified(item)
}

func (h *History) isModified(item *Item) *Item {
	mod, ok := item.Modified()
	if !ok {
		return nil
	}
	return h.sessionLog[mod]
}

// AddResult describes the merge outcome.
type AddResult struct {
	Merged   bool
	Item     *Item
	Replaced *Item // the stored item a merge replaced, else nil
}

// Add ports History.add merge matrix @ eb03eba (contents rule per source):
// - supersedes branch (isModified == nil): stored contents win
// - isModified branch: keep incoming (modifier) contents
// Own-write / foreign: full merge on match — bump count, preserve firstCopiedAt/pin/title.
func (h *History) Add(item *Item, changeCount int, fromOwnWrite bool, now time.Time) AddResult {
	item.LastCopiedAt = now
	if item.NumberOfCopies < 1 {
		item.NumberOfCopies = 1
	}

	existing := h.FindSimilarItem(item)
	if existing == nil {
		if item.FirstCopiedAt.IsZero() {
			item.FirstCopiedAt = now
		}
		h.items = append([]*Item{item}, h.items...)
		h.sessionLog[changeCount] = item
		return AddResult{Merged: false, Item: item}
	}

	// Merge into existing.
	if h.isModified(item) == nil {
		// supersedes branch: stored contents win
		item.Contents = existing.Contents
	}
	// else isModified branch: keep item.Contents (modifier bytes)

	item.ID = existing.ID // the merged item replaces the stored row
	item.FirstCopiedAt = existing.FirstCopiedAt
	item.NumberOfCopies += existing.NumberOfCopies
	item.Pin = existing.Pin
	item.Title = existing.Title
	if !fromOwnWrite {
		// Maccy: if !item.fromMaccy { item.application = existing.application }
		// For foreign writes we keep incoming application; for own-write we preserve existing.
		// Plan: own-write marker selects merge-vs-new attribution; application rules per .fromMaccy path.
		// fromMaccy true ⇒ do NOT overwrite application from existing... actually:
		// `if !item.fromMaccy { item.application = existingHistoryItem.application }`
		// so when fromOwnWrite (fromMaccy), keep item.Application; when foreign merge, take existing's app.
		item.Application = existing.Application
	}

	// Replace existing in list with merged item at front.
	out := make([]*Item, 0, len(h.items))
	for _, it := range h.items {
		if it == existing {
			continue
		}
		out = append(out, it)
	}
	h.items = append([]*Item{item}, out...)
	for k, v := range h.sessionLog {
		if v == existing {
			h.sessionLog[k] = item
		}
	}
	h.sessionLog[changeCount] = item
	return AddResult{Merged: true, Item: item, Replaced: existing}
}

// Load replaces the list with stored items, newest first. The session log
// starts empty: ModifiedType merges only match this session's writes.
func (h *History) Load(items []*Item) {
	h.items = append([]*Item(nil), items...)
	h.sessionLog = map[int]*Item{}
}

// Remove drops it from the list and the session log.
func (h *History) Remove(it *Item) {
	out := h.items[:0]
	for _, x := range h.items {
		if x != it {
			out = append(out, x)
		}
	}
	h.items = out
	for k, v := range h.sessionLog {
		if v == it {
			delete(h.sessionLog, k)
		}
	}
}

// Trim evicts the oldest unpinned items until at most maxItems unpinned
// items remain (Maccy's `size` counts unpinned only) and all contents fit in
// maxBytes. Zero disables a limit. It returns the evicted items.
func (h *History) Trim(maxItems int, maxBytes int64) []*Item {
	var evicted []*Item
	for {
		unpinned, total := 0, int64(0)
		var oldest *Item
		for _, it := range h.items {
			total += it.Size()
			if it.Pin == "" {
				unpinned++
				oldest = it // items are newest first
			}
		}
		over := (maxItems > 0 && unpinned > maxItems) || (maxBytes > 0 && total > maxBytes)
		if !over || oldest == nil {
			return evicted
		}
		h.Remove(oldest)
		evicted = append(evicted, oldest)
	}
}
