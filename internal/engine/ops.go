package engine

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

// ErrNoItem is returned for an ID the history no longer has.
var ErrNoItem = errors.New("engine: no such item")

// Item returns the item with id, or nil.
func (e *Engine) Item(id int64) *history.Item {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.find(id)
}

func (e *Engine) find(id int64) *history.Item {
	for _, x := range e.hist.Items() {
		if x.ID == id {
			return x
		}
	}
	return nil
}

// SetPrefs applies capture settings, keeping a pending ignore-next.
func (e *Engine) SetPrefs(p capture.Prefs) {
	e.cap.UpdatePrefs(func(cur *capture.Prefs) {
		next := cur.IgnoreOnlyNextEvent && cur.IgnoreEvents
		*cur = p
		if next {
			cur.IgnoreEvents, cur.IgnoreOnlyNextEvent = true, true
		}
	})
	e.changed()
}

// SetLimits changes the history caps and trims to them now.
func (e *Engine) SetLimits(l Limits) {
	e.mu.Lock()
	e.cfg.Limits = l
	evicted := e.hist.Trim(l.MaxItems, l.MaxBytes)
	e.mu.Unlock()
	e.report(e.deleteStored(evicted))
	e.changed()
}

func (e *Engine) deleteStored(items []*history.Item) error {
	var ids []int64
	for _, x := range items {
		if x.ID != 0 {
			ids = append(ids, x.ID)
		}
	}
	return e.cfg.Store.Delete(context.Background(), ids...)
}

// update replaces the item with id by a modified copy, so items already
// handed out stay unchanged, and saves it.
func (e *Engine) update(id int64, fn func(*history.Item)) error {
	e.mu.Lock()
	old := e.find(id)
	if old == nil {
		e.mu.Unlock()
		return ErrNoItem
	}
	it := *old
	it.Contents = slices.Clone(old.Contents)
	fn(&it)
	e.hist.Replace(old, &it)
	e.mu.Unlock()
	err := e.cfg.Store.Save(context.Background(), &it)
	e.report(err)
	e.changed()
	return err
}

// Delete removes an item.
func (e *Engine) Delete(id int64) error {
	e.mu.Lock()
	it := e.find(id)
	if it != nil {
		e.hist.Remove(it)
	}
	e.mu.Unlock()
	if it == nil {
		return ErrNoItem
	}
	err := e.cfg.Store.Delete(context.Background(), id)
	e.report(err)
	e.changed()
	return err
}

// Clear removes the unpinned items, or every item with all
// (History.clear / clearAll).
func (e *Engine) Clear(all bool) error {
	e.mu.Lock()
	var gone []*history.Item
	for _, it := range e.hist.Items() {
		if all || it.Pin == "" {
			e.hist.Remove(it)
			gone = append(gone, it)
		}
	}
	e.mu.Unlock()
	err := e.deleteStored(gone)
	e.report(err)
	e.changed()
	return err
}

// TogglePin unpins a pinned item, or pins it to the first free key not in
// reserved. It returns the new pin ("" when unpinned or no key is free).
func (e *Engine) TogglePin(id int64, reserved ...string) (string, error) {
	var pin string
	err := e.update(id, func(it *history.Item) {
		if it.Pin != "" {
			it.Pin = ""
			return
		}
		if free := history.AvailablePins(e.hist.Items(), reserved...); len(free) > 0 {
			pin = free[0]
			it.Pin = pin
		}
	})
	return pin, err
}

// SetPin gives a pinned item another key; a key another item has is
// refused.
func (e *Engine) SetPin(id int64, key string) error {
	e.mu.Lock()
	for _, x := range e.hist.Items() {
		if x.ID != id && x.Pin == key {
			e.mu.Unlock()
			return fmt.Errorf("engine: key %q is taken", key)
		}
	}
	e.mu.Unlock()
	return e.update(id, func(it *history.Item) { it.Pin = key })
}

// SetTitle renames an item; "" goes back to the title from its contents.
func (e *Engine) SetTitle(id int64, title string) error {
	return e.update(id, func(it *history.Item) { it.Title = title })
}

// SetText replaces a text item's contents with plain text (Maccy's pin
// content edit, which drops formatting).
func (e *Engine) SetText(id int64, text string) error {
	return e.update(id, func(it *history.Item) {
		it.Contents = []history.Content{{Type: pasteboard.TypeString, Value: []byte(text)}}
	})
}

// CopyBack puts a history item back on the pasteboard, marked as MyClip's
// own write; the next tick merges it into the item (copy-back bump). With
// plain, formatting is removed as Clipboard.clearFormatting does: only the
// string, plus file URLs, when there is a string.
func (e *Engine) CopyBack(id int64, plain bool) error {
	e.mu.Lock()
	it := e.find(id)
	var items []pasteboard.Item
	if it != nil {
		items = boardItems(it.Contents, plain)
	}
	e.mu.Unlock()
	if it == nil {
		return ErrNoItem
	}
	var err error
	e.cfg.RunOnMain(func() { err = e.cfg.Board.Write(items) })
	return err
}

// boardItems lays contents out as pasteboard items: one item with every
// representation, and a further item per extra file URL, so several copied
// files paste as several (Clipboard.copy writes file URLs as objects).
func boardItems(contents []history.Content, plain bool) []pasteboard.Item {
	if plain && slices.ContainsFunc(contents, func(c history.Content) bool { return c.Type == pasteboard.TypeString }) {
		contents = slices.DeleteFunc(slices.Clone(contents), func(c history.Content) bool {
			return c.Type != pasteboard.TypeString && c.Type != pasteboard.TypeFileURL
		})
	}
	first := pasteboard.Item{Types: map[string][]byte{pasteboard.TypeOwnWrite: {}}}
	items := []pasteboard.Item{first}
	for _, c := range contents {
		if _, dup := first.Types[c.Type]; dup && c.Type == pasteboard.TypeFileURL {
			items = append(items, pasteboard.Item{Types: map[string][]byte{c.Type: c.Value}})
			continue
		}
		first.Types[c.Type] = c.Value
	}
	return items
}
