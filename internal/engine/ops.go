package engine

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
	"github.com/L92-Labs/myclip/internal/snippet"
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
	e.persist.Lock()
	e.mu.Lock()
	e.cfg.Limits = l
	evicted := e.hist.Trim(l.MaxItems, l.MaxBytes)
	for _, x := range evicted {
		delete(e.expires, x.ID)
	}
	e.mu.Unlock()
	err := e.deleteStored(evicted)
	e.persist.Unlock()
	e.report(err)
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
	return e.updateIf(id, nil, fn)
}

// updateIf is update when ok (nil: always) holds for the stored item, both
// checked and applied under the lock.
func (e *Engine) updateIf(id int64, ok func(*history.Item) bool, fn func(*history.Item)) error {
	e.persist.Lock()
	err := e.updateHeld(id, ok, fn)
	e.persist.Unlock()
	if err == ErrNoItem {
		return err
	}
	e.report(err)
	e.changed()
	return err
}

// updateHeld changes and saves the item; the caller holds persist.
func (e *Engine) updateHeld(id int64, ok func(*history.Item) bool, fn func(*history.Item)) error {
	e.mu.Lock()
	old := e.find(id)
	if old == nil || ok != nil && !ok(old) {
		e.mu.Unlock()
		return ErrNoItem
	}
	it := *old
	it.Contents = slices.Clone(old.Contents)
	fn(&it)
	e.hist.Replace(old, &it)
	e.mu.Unlock()
	if beforeSave != nil {
		beforeSave()
	}
	return e.cfg.Store.Save(context.Background(), &it)
}

// beforeSave is a test hook between an update's change and its write.
var beforeSave func()

// ImportItems merges imported items into the history and the store,
// preserving their timestamps. An item matching an existing clip
// (same contents) is skipped; a pin the history already uses, or a pin
// key an in-popup shortcut reserves, moves to the next free key, or is
// dropped when none is free. It reports how many were added and skipped.
func (e *Engine) ImportItems(items []*history.Item) (added, skipped int, err error) {
	e.persist.Lock()
	for _, cand := range items {
		e.mu.Lock()
		if e.hist.FindSimilarItem(cand) != nil {
			e.mu.Unlock()
			skipped++
			continue
		}
		taken := false
		if cand.Pin != "" {
			for _, x := range e.hist.Items() {
				if x.Pin == cand.Pin {
					taken = true
					break
				}
			}
			// A reserved key is a conflict even on empty history: the
			// popup checks pinned rows before its own shortcuts.
			if !taken && slices.Contains(e.pinReserved, cand.Pin) {
				taken = true
			}
		}
		if taken {
			if free := history.AvailablePins(e.hist.Items(), e.pinReserved...); len(free) > 0 {
				cand.Pin = free[0]
			} else {
				cand.Pin = ""
			}
		}
		e.mu.Unlock()
		if err := e.cfg.Store.Save(context.Background(), cand); err != nil {
			e.report(err)
			e.persist.Unlock()
			return added, skipped, err
		}
		e.mu.Lock()
		e.hist.Import(cand)
		e.mu.Unlock()
		added++
	}
	e.persist.Unlock()
	// The import may carry clips a live expiry rule covers (or pin rule
	// matches): give them the rules' effects, as a fresh capture would.
	e.applyRules()
	e.changed()
	return added, skipped, nil
}

// Delete removes an item.
func (e *Engine) Delete(id int64) error {
	e.persist.Lock()
	e.mu.Lock()
	it := e.find(id)
	if it != nil {
		e.hist.Remove(it)
		delete(e.expires, id) // its deadline dies with it (the id may be reused)
	}
	e.mu.Unlock()
	if it == nil {
		e.persist.Unlock()
		return ErrNoItem
	}
	err := e.cfg.Store.Delete(context.Background(), id)
	e.persist.Unlock()
	e.report(err)
	e.changed()
	return err
}

// Clear removes the unpinned items, or every item with all
// (History.clear / clearAll).
func (e *Engine) Clear(all bool) error {
	e.persist.Lock()
	e.mu.Lock()
	var gone []*history.Item
	for _, it := range e.hist.Items() {
		if all || it.Pin == "" {
			e.hist.Remove(it)
			delete(e.expires, it.ID)
			gone = append(gone, it)
		}
	}
	e.mu.Unlock()
	err := e.deleteStored(gone)
	e.persist.Unlock()
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

// MaxRecognizedRunes bounds the text kept from an image.
const MaxRecognizedRunes = 20000

// ImageDigest is the SHA-256 of an item's image (zero without one).
func ImageDigest(it *history.Item) [32]byte {
	data, _, ok := it.Image()
	if !ok {
		return [32]byte{}
	}
	return sha256.Sum256(data)
}

// SetRecognizedText stores the text read in an item's image, replacing an
// earlier reading. An empty text is stored too: it marks the image read.
// The storage limits apply again, since the item grew.
// OCR text uses only the bytes left under the capture-time per-item cap;
// an empty reading still marks a full item as processed.
//
// The reading lands only on the image read: the item stored under id must
// still hold an image with digest image (SHA-256), checked under the same
// lock as the update, as an ID may be taken by a newer clip meanwhile.
func (e *Engine) SetRecognizedText(id int64, image [32]byte, text string) error {
	if r := []rune(text); len(r) > MaxRecognizedRunes {
		text = string(r[:MaxRecognizedRunes])
	}
	same := func(it *history.Item) bool { return ImageDigest(it) == image }
	limit := e.cap.MaxItemBytes
	if limit == 0 {
		limit = capture.DefaultMaxItemBytes
	}
	e.persist.Lock()
	err := e.updateHeld(id, same, func(it *history.Item) {
		it.Contents = slices.DeleteFunc(it.Contents, func(c history.Content) bool {
			return c.Type == pasteboard.TypeRecognizedText
		})
		// Reclaim an earlier reading's bytes before budgeting its replacement.
		if remaining := max(limit-it.Size(), 0); int64(len(text)) > remaining {
			end := int(remaining)
			for end > 0 && !utf8.RuneStart(text[end]) {
				end--
			}
			text = text[:end]
		}
		it.Contents = append(it.Contents, history.Content{Type: pasteboard.TypeRecognizedText, Value: []byte(text)})
	})
	if err == nil {
		e.mu.Lock()
		evicted := e.hist.Trim(e.cfg.Limits.MaxItems, e.cfg.Limits.MaxBytes)
		for _, x := range evicted {
			delete(e.expires, x.ID)
		}
		e.mu.Unlock()
		err = e.deleteStored(evicted)
	}
	e.persist.Unlock()
	if err == ErrNoItem {
		return err
	}
	e.report(err)
	e.changed()
	// The reading may be what a rule matches: pin or expire the item now,
	// as the rules would on their next change.
	e.applyRules()
	return err
}

// SetText replaces a text item's contents with plain text (Maccy's pin
// content edit, which drops formatting).
func (e *Engine) SetText(id int64, text string) error {
	return e.update(id, func(it *history.Item) {
		it.Contents = []history.Content{{Type: pasteboard.TypeString, Value: []byte(text)}}
	})
}

// SetSnippet marks or unmarks a clip as a snippet: its text carries
// placeholders ({date}, {time}, {clipboard}, {uuid}) that expand at
// paste time.
func (e *Engine) SetSnippet(id int64, on bool) error {
	return e.update(id, func(it *history.Item) { it.Snippet = on })
}

// CopyBack puts a history item back on the pasteboard, marked as MyClip's
// own write; the next tick merges it into the item (copy-back bump). With
// plain, formatting is removed as Clipboard.clearFormatting does: only the
// string, plus file URLs, when there is a string. A snippet's text has its
// placeholders expanded first; the {clipboard} value is the clipboard's
// text before this write.
func (e *Engine) CopyBack(id int64, plain bool) error {
	e.mu.Lock()
	it := e.find(id)
	var contents []history.Content
	if it != nil {
		contents = it.Contents
	}
	e.mu.Unlock()
	if it == nil {
		return ErrNoItem
	}
	var err error
	e.cfg.RunOnMain(func() {
		if it.Snippet {
			var clip string
			if read := e.cfg.Board.ReadItems(); len(read) > 0 {
				clip = string(read[0].Types[pasteboard.TypeString])
			}
			contents = expandContents(contents, time.Now(), clip)
		}
		items := boardItems(contents, plain)
		if it.Snippet {
			// The expanded text is no longer the stored template, so the
			// ordinary own-write merge cannot match it: capture would file
			// every paste as a new clip and evict the template. Skip the
			// write exactly once, like a paste-stack copy.
			var nonce [8]byte
			rand.Read(nonce[:])
			marker := pasteboard.StackWritePrefix + hex.EncodeToString(nonce[:])
			items[0].Types[marker] = []byte{}
			e.cap.ExpectStack(marker)
			defer func() {
				if err != nil {
					e.cap.ExpectStack("")
				}
			}()
		}
		err = e.cfg.Board.Write(items)
	})
	return err
}

// expandContents replaces the placeholders in a snippet's text content.
// The rich formats are left with the raw placeholders: snippets are
// text-first, and rewriting HTML or RTF reliably is not worth it.
func expandContents(contents []history.Content, now time.Time, clip string) []history.Content {
	out := make([]history.Content, 0, len(contents))
	for _, c := range contents {
		if c.Type == pasteboard.TypeString {
			v := snippet.Expand(string(c.Value), now, func() string { return clip })
			c.Value = []byte(v)
		}
		out = append(out, c)
	}
	return out
}

// CopyMany puts several items on the pasteboard as one copy, in the order
// given (the paste stack): their files when every item is files, else their
// texts joined by new lines; images without text are left out. The combined
// copy is not recorded as a new clip.
func (e *Engine) CopyMany(ids []int64) error {
	e.mu.Lock()
	var texts []string
	var files []history.Content
	allFiles := true
	for _, id := range ids {
		it := e.find(id)
		if it == nil {
			continue
		}
		if len(it.FilePaths()) == 0 {
			allFiles = false
		}
		for _, c := range it.Contents {
			if c.Type == pasteboard.TypeFileURL {
				files = append(files, c)
			}
		}
		if s := strings.TrimRight(it.ContentText(), "\n"); s != "" { // never a title
			texts = append(texts, s)
		}
	}
	e.mu.Unlock()
	var items []pasteboard.Item
	switch {
	case allFiles && len(files) > 0:
		items = boardItems(files, false)
		delete(items[0].Types, pasteboard.TypeOwnWrite)
	case len(texts) > 0:
		items = []pasteboard.Item{{Types: map[string][]byte{pasteboard.TypeString: []byte(strings.Join(texts, "\n"))}}}
	default:
		return ErrNoItem
	}
	// The combination is not a clip of its own: a one-shot marker makes
	// capture skip exactly this change — not a copy the user makes right
	// after, nor a later replay of the marker by another app.
	var nonce [8]byte
	rand.Read(nonce[:])
	marker := pasteboard.StackWritePrefix + hex.EncodeToString(nonce[:])
	items[0].Types[marker] = []byte{}
	var err error
	e.cfg.RunOnMain(func() {
		e.cap.ExpectStack(marker)
		if err = e.cfg.Board.Write(items); err != nil {
			e.cap.ExpectStack("")
		}
	})
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
		if c.Type == pasteboard.TypeRecognizedText {
			continue // MyClip's own reading, not part of the copy
		}
		if _, dup := first.Types[c.Type]; dup && c.Type == pasteboard.TypeFileURL {
			items = append(items, pasteboard.Item{Types: map[string][]byte{c.Type: c.Value}})
			continue
		}
		first.Types[c.Type] = c.Value
	}
	return items
}
