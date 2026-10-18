// Package popup is MyClip's history popup, after Maccy 2.7.1's: a search
// field over the history, pinned items at the top or bottom, ⌘1–9 and pin
// keys, a preview beside the list, and a footer of commands. It draws in
// native UI and leaves the window, the pasteboard and storage to its Host.
package popup

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/egoist/mygo/ui"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/search"
	"github.com/L92-Labs/myclip/internal/settings"
	"github.com/L92-Labs/myclip/internal/shortcut"
	"github.com/L92-Labs/myclip/locales"
)

// Action is what selecting an item does after copying it.
type Action struct {
	Paste bool // paste into the frontmost app
	Plain bool // without formatting
}

// Resolve maps the modifiers of a selection to an action, as Maccy's
// HistoryItemAction: no modifiers follow the paste and formatting
// preferences; ⌘, ⌥ and ⌥⇧/⌘⇧ choose among copy, paste and paste plain.
func Resolve(mods ui.Modifiers, pasteByDefault, removeFormatting bool) (Action, bool) {
	pbd, rfd := pasteByDefault, removeFormatting
	copyOnly, paste, plain := Action{}, Action{Paste: true}, Action{Paste: true, Plain: true}
	switch mods {
	case 0:
		return Action{Paste: pbd, Plain: rfd}, true
	case ui.Cmd:
		switch {
		case !pbd:
			return copyOnly, true
		case !rfd:
			return paste, true
		default:
			return plain, true
		}
	case ui.Alt:
		switch {
		case pbd:
			return copyOnly, true
		case rfd:
			return plain, true
		default:
			return paste, true
		}
	case ui.Alt | ui.Shift:
		if pbd {
			return Action{}, false
		}
		if rfd {
			return paste, true
		}
		return plain, true
	case ui.Cmd | ui.Shift:
		if !pbd {
			return Action{}, false
		}
		if rfd {
			return paste, true
		}
		return plain, true
	}
	return Action{}, false
}

// Host is what the popup acts on. Its methods run on the main thread.
type Host interface {
	Select(id int64, a Action)
	// SelectMany copies (and pastes) several items together, in the order
	// shown: the paste stack.
	SelectMany(ids []int64, a Action)
	TogglePin(id int64)
	Delete(id int64)
	Clear(all bool)
	OpenSettings()
	About()
	Quit()
	Close()
	// Thumbnail returns the item's image scaled to at most maxHeight DIPs
	// high, nil while it loads or when it has none.
	Thumbnail(it *history.Item, maxHeight int) *ui.Bitmap
	// Preview returns the item's image for the preview, or nil.
	Preview(it *history.Item) *ui.Bitmap
	// PreviewSize is that image's own pixel size, before scaling, once
	// known.
	PreviewSize(it *history.Item) (w, h int, ok bool)
	AppIcon(bundleID string) *ui.Bitmap
	AppName(bundleID string) string
	// Resize asks for the window to take this size, after the frame.
	Resize(width, height int)
}

type row struct {
	it     *history.Item
	title  string
	ranges []search.Range
	key    ui.Key // ⌘-key selecting it: 1–9, or the pin's letter; 0 for none
}

// Popup is the popup's state. Settings points at the app's settings,
// which change on the main thread only.
type Popup struct {
	Host     Host
	Settings *settings.Settings

	items  []*history.Item
	rows   []row
	query  string
	chosen int
	// marked holds the rows ⇧-selected for a paste stack, by item, from
	// anchor to the chosen row.
	marked map[int64]bool
	anchor int

	list   ui.ListState
	search ui.Handle

	chosenID    int64
	chosenAt    time.Time
	previewOpen bool
	previewShut bool // closed by hand: no automatic preview until reopened

	hoverRow int
	hoverX   float32
	hoverY   float32

	confirm, confirmAll, dontAsk bool

	width, height int
}

// SetItems takes the history, sorted, and filters it again, keeping the
// item chosen.
func (p *Popup) SetItems(items []*history.Item) {
	p.items = items
	var id int64
	if p.chosen < len(p.rows) {
		id = p.rows[p.chosen].it.ID
	}
	p.filter()
	for i, r := range p.rows {
		if r.it.ID == id {
			p.choose(i)
			return
		}
	}
	p.chooseDefault()
}

// Reset readies the popup to show: no query, the first unpinned item
// chosen, no preview, the search field focused.
func (p *Popup) Reset() {
	p.query = ""
	p.previewOpen, p.previewShut = false, false
	p.confirm = false
	p.marked = nil
	p.hoverRow = -2 // the pointer resting where the popup opens chooses nothing
	p.filter()
	p.chooseDefault()
	p.search.Focus()
}

// TogglePreview shows or hides the preview, as its shortcut does.
func (p *Popup) TogglePreview() {
	p.previewOpen = !p.previewOpen
	p.previewShut = !p.previewOpen
}

// Query is the search text.
func (p *Popup) Query() string { return p.query }

// filter applies the query to the items and numbers the unpinned rows.
func (p *Popup) filter() {
	s := p.Settings
	titles := make([]string, len(p.items))
	texts := make([]string, len(p.items)) // what the query searches
	for i, it := range p.items {
		titles[i] = it.ListTitle(s.ShowSpecialSymbols)
		texts[i] = titles[i]
		if texts[i] == "" { // an image row: search the text read in it, else its own text
			if texts[i] = it.RecognizedText(); texts[i] == "" {
				texts[i] = it.PreviewText()
			}
		}
	}
	results := search.Search(s.SearchMode, p.query, texts)
	p.rows = p.rows[:0]
	n := 0
	for _, r := range results {
		it := p.items[r.Index]
		rw := row{it: it, title: titles[r.Index], ranges: r.Ranges}
		if rw.title == "" {
			rw.ranges = nil // the ranges are in the image's text, not shown in the row
		}
		if it.Pin != "" {
			rw.key, _ = shortcut.Letter([]rune(it.Pin)[0])
		} else if n < 9 {
			n++
			rw.key = ui.Key0 + ui.Key(n)
		}
		p.rows = append(p.rows, rw)
	}
}

// chooseDefault chooses the first unpinned row with no query (Maccy
// selects unpinnedItems.first), else the first row.
func (p *Popup) chooseDefault() {
	if p.query == "" {
		for i, r := range p.rows {
			if r.it.Pin == "" {
				p.choose(i)
				return
			}
		}
	}
	p.choose(0)
}

func (p *Popup) choose(i int) {
	p.chosen = max(0, min(i, len(p.rows)-1))
	if len(p.rows) > 0 {
		p.list.ScrollIntoView(p.chosen)
	}
}

func (p *Popup) current() *history.Item {
	if p.chosen < len(p.rows) {
		return p.rows[p.chosen].it
	}
	return nil
}

func (p *Popup) setQuery(q string) {
	p.query = q
	p.marked = nil
	p.filter()
	p.chooseDefault()
}

func (p *Popup) selectRow(i int, mods ui.Modifiers) {
	if i >= len(p.rows) {
		return
	}
	a, ok := Resolve(mods, p.Settings.PasteByDefault, p.Settings.RemoveFormattingByDefault)
	if !ok {
		return
	}
	if ids := p.Marked(); len(ids) > 1 && p.marked[p.rows[i].it.ID] {
		p.Host.SelectMany(ids, a)
		return
	}
	p.Host.Select(p.rows[i].it.ID, a)
}

// ChooseNext moves the choice down, round to the top (hold-to-cycle).
func (p *Popup) ChooseNext() {
	if len(p.rows) == 0 {
		return
	}
	p.marked = nil
	p.choose((p.chosen + 1) % len(p.rows))
}

// SelectChosen selects the chosen row as Enter does (hold-to-cycle's
// release).
func (p *Popup) SelectChosen() { p.selectRow(p.chosen, 0) }

// Marked returns the ⇧-selected items' IDs in the order shown.
func (p *Popup) Marked() []int64 {
	var ids []int64
	for _, r := range p.rows {
		if p.marked[r.it.ID] {
			ids = append(ids, r.it.ID)
		}
	}
	return ids
}

// extend moves the choice to row j and marks the rows from the anchor to
// it, as ⇧↑/⇧↓ and ⇧-click do in lists.
func (p *Popup) extend(j int) {
	if len(p.rows) == 0 {
		return
	}
	if len(p.marked) == 0 {
		p.anchor = p.chosen
	}
	p.choose(j)
	p.marked = map[int64]bool{}
	lo, hi := min(p.anchor, p.chosen), max(p.anchor, p.chosen)
	for k := lo; k <= hi && k < len(p.rows); k++ {
		p.marked[p.rows[k].it.ID] = true
	}
}

func (p *Popup) requestClear(all bool) {
	if p.Settings.SuppressClearAlert {
		p.Host.Clear(all)
		return
	}
	p.confirm, p.confirmAll, p.dontAsk = true, all, false
}

func (p *Popup) searchShown() bool {
	s := p.Settings
	return s.ShowSearch && (s.SearchVisibility == settings.SearchAlways || p.query != "")
}

// key handles a key pressed in the popup (KeyChord @ eb03eba) and reports
// whether it took it. The search field passes it every key first, so
// chords win over its editing keys.
func (p *Popup) key(mods ui.Modifiers, k ui.Key) bool {
	is := func(m ui.Modifiers, key ui.Key) bool { return mods == m && k == key }
	if k == ui.KeyEnter {
		p.selectRow(p.chosen, mods)
		return true
	}
	if mods != 0 && mods != ui.Shift {
		for i, r := range p.rows {
			if r.key == k {
				if _, ok := Resolve(mods, p.Settings.PasteByDefault, p.Settings.RemoveFormattingByDefault); ok {
					p.selectRow(i, mods)
					return true
				}
			}
		}
	}
	chord := func(accel string) bool {
		ch, ok := shortcut.Parse(accel)
		return ok && is(ch.Mods, ch.Key)
	}
	switch {
	case is(ui.Shift, ui.KeyDown):
		p.extend(p.chosen + 1)
	case is(ui.Shift, ui.KeyUp):
		p.extend(p.chosen - 1)
	case is(0, ui.KeyDown), is(ui.Ctrl, ui.KeyN), is(ui.Ctrl, ui.KeyJ):
		p.marked = nil
		p.choose(p.chosen + 1)
	case is(0, ui.KeyUp), is(ui.Ctrl, ui.KeyP), p.chosen > 0 && is(ui.Ctrl, ui.KeyK):
		p.marked = nil
		p.choose(p.chosen - 1)
	case is(ui.Cmd, ui.KeyDown), is(ui.Alt, ui.KeyDown), is(ui.Ctrl|ui.Alt, ui.KeyN), is(0, ui.KeyPageDown):
		p.choose(len(p.rows) - 1)
	case is(ui.Cmd, ui.KeyUp), is(ui.Alt, ui.KeyUp), is(ui.Ctrl|ui.Alt, ui.KeyP), is(0, ui.KeyPageUp):
		p.choose(0)
	case is(ui.Cmd|ui.Alt, ui.KeyBackspace):
		p.requestClear(false)
	case is(ui.Cmd|ui.Alt|ui.Shift, ui.KeyBackspace):
		p.requestClear(true)
	case is(ui.Ctrl, ui.KeyU):
		p.setQuery("")
	case is(ui.Ctrl, ui.KeyW):
		q := strings.TrimRightFunc(p.query, unicode.IsSpace)
		p.setQuery(q[:strings.LastIndexFunc(q, unicode.IsSpace)+1])
	case is(ui.Ctrl, ui.KeyH):
		if r := []rune(p.query); len(r) > 0 {
			p.setQuery(string(r[:len(r)-1]))
		}
	case chord(p.Settings.PinShortcut):
		if it := p.current(); it != nil {
			p.Host.TogglePin(it.ID)
		}
	case chord(p.Settings.DeleteShortcut):
		if it := p.current(); it != nil {
			p.Host.Delete(it.ID)
		}
	case chord(p.Settings.PreviewShortcut):
		p.TogglePreview()
	case is(ui.Cmd, ui.KeyComma):
		p.Host.OpenSettings()
	case is(ui.Cmd, ui.KeyQ):
		p.Host.Quit()
	case is(0, ui.KeyEscape):
		p.Host.Close()
	default:
		return false
	}
	return true
}

// input routes the search field's input: keys to key, and typed text to
// the field only while there is a search.
func (p *Popup) input(ev ui.InputEvent) bool {
	switch ev.Kind {
	case ui.InputKeyDown:
		return p.key(ev.Mods, ev.Key)
	case ui.InputText:
		return !p.Settings.ShowSearch
	}
	return false
}

// View builds the popup.
func (p *Popup) View(c *ui.Context) {
	s := p.Settings
	t := c.Theme()
	root := c.Root()
	if c.Vibrancy() {
		root.Background(ui.Transparent)
	} else {
		root.Background(t.Background)
	}
	if !p.confirm && !p.search.Focused(c) {
		p.search.Focus(c)
	}
	p.autoPreview(c)

	listWidth := float32(s.WindowWidth)
	ui.Row(c).Fill().Children(func() {
		ui.Column(c).Width(listWidth).FillHeight().Children(func() {
			p.header(c)
			p.listView(c)
			if s.ShowFooter {
				p.footer(c)
			}
		})
		if p.previewOpen && p.current() != nil {
			ui.Divider(c).Vertical()
			p.preview(c, p.current())
		}
	})
	p.confirmDialog(c)
	p.fit()
}

func (p *Popup) header(c *ui.Context) {
	s := p.Settings
	shown := p.searchShown()
	t := c.Theme()
	hdr := ui.Row(c).Padding(10, 12, 8, 12).Gap(8).AlignItems(ui.Center)
	if !shown {
		// Kept, unseen, to hold the focus and take the keys, and what is
		// typed with searchVisibility duringSearch: it shows with the
		// first letter.
		hdr.Height(0).Padding(0).Clip().Opacity(0)
	}
	hdr.Children(func() {
		if s.ShowTitle && shown {
			ui.Text(c, "MyClip").Bold().FontSize(13)
		}
		field := ui.Row(c).Grow(1).Gap(6).Padding(5, 8).Radius(6).Background(t.Surface).AlignItems(ui.Center)
		field.Children(func() {
			ui.Icon(c, searchIcon).FontSize(13).TextColor(t.TextMuted)
			if ui.TextInputBase(c.Key("search"), &p.query).Bind(&p.search).Grow(1).FontSize(13).
				Placeholder(locales.SearchPlaceholder()).Label(locales.SearchPlaceholder()).
				HandleInput(p.input).Changed() {
				p.marked = nil // a new query: the marked rows may be gone
				p.filter()
				p.chooseDefault()
			}
		})
	})
	if shown {
		ui.Divider(c)
	}
}

func (p *Popup) listView(c *ui.Context) {
	t := c.Theme()
	p.list.Key = func(i int) any { return p.rows[i].it.ID }
	p.list.Label = func(i int) string { return p.rows[i].title }
	ui.List(c.Key("history"), &p.list, len(p.rows), func(i int) { p.row(c, i) }).
		Grow(1).Padding(4, 5).Children(func() {
		if len(p.rows) == 0 {
			msg := locales.NoClips()
			if p.query != "" {
				msg = locales.NoMatches()
			}
			ui.Text(c, msg).TextColor(t.TextMuted).Padding(10, 10).FontSize(13)
		}
	})
}

var hexColor = regexp.MustCompile(`^#?([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$`)

// row draws one history item: its app's icon, its title or image, a color
// swatch for a hex color, and its shortcut.
func (p *Popup) row(c *ui.Context, i int) {
	s := p.Settings
	t := c.Theme()
	r := p.rows[i]
	chosen := i == p.chosen
	text, muted := t.Text, t.TextMuted
	el := ui.Row(c).Padding(4, 8).Gap(8).Radius(5).AlignItems(ui.Center).MinHeight(24)
	if chosen {
		el.Background(t.Accent)
		text, muted = t.AccentText, t.AccentText.Alpha(0.75)
	} else if p.marked[r.it.ID] {
		el.Background(t.Accent.Alpha(0.3)) // in the paste stack
	}
	if x, y, ok := el.PointerPosition(); ok && el.Hovered() {
		if p.hoverRow == -2 {
			p.hoverRow, p.hoverX, p.hoverY = i, x, y // where the popup opened
		} else if (i != p.hoverRow || x != p.hoverX || y != p.hoverY) && c.Modifiers()&ui.Shift == 0 {
			// (⇧ held: the pointer is heading for a ⇧-click, which extends
			// from the row chosen now.)
			p.hoverRow, p.hoverX, p.hoverY = i, x, y
			if i != p.chosen {
				p.choose(i)
			}
		}
	}
	if el.Clicked() {
		if el.ClickModifiers() == ui.Shift {
			p.extend(i)
		} else {
			p.selectRow(i, el.ClickModifiers())
		}
	}
	el.ContextMenu(func(m *ui.Menu) { p.rowMenu(m, r.it) })
	el.Children(func() {
		if s.ShowApplicationIcons {
			if icon := p.Host.AppIcon(r.it.Application); icon != nil {
				ui.Image(c, icon).Size(16, 16)
			} else {
				ui.Box(c).Size(16, 16)
			}
		}
		if r.title == "" {
			if bmp := p.Host.Thumbnail(r.it, s.ImageMaxHeight); bmp != nil {
				w, h := bmp.Size()
				dh := float32(min(h, s.ImageMaxHeight))
				dw := float32(w) * dh / float32(max(h, 1))
				ui.Image(c, bmp).Size(min(dw, 340), dh).Fit(ui.Contain).Radius(3)
			} else {
				ui.Text(c, locales.Image()).TextColor(muted).FontSize(13)
			}
			ui.Spacer(c)
		} else {
			if s.ShowHexColorSwatch && hexColor.MatchString(strings.TrimSpace(r.title)) {
				ui.Box(c).Size(12, 12).Radius(3).Background(ui.Hex(strings.TrimSpace(r.title))).Border(1, t.Border)
			}
			p.title(c, r, text).Grow(1).Shrink(1)
		}
		if r.key != 0 {
			ui.Text(c, p.keyLabel(c, r)).FontSize(12).TextColor(muted)
		}
	})
}

// keyLabel shows the row's chord for the modifiers held: ⌘1 at rest, ⌥1
// while Option is held, as Maccy's KeyShortcut.isVisible.
func (p *Popup) keyLabel(c *ui.Context, r row) string {
	mods := ui.Cmd
	switch held := c.Modifiers(); held {
	case ui.Alt, ui.Alt | ui.Shift, ui.Cmd | ui.Shift:
		if _, ok := Resolve(held, p.Settings.PasteByDefault, p.Settings.RemoveFormattingByDefault); ok {
			mods = held
		}
	}
	return shortcut.Chord{Mods: mods, Key: r.key}.Symbols()
}

// title draws the row's title with the search matches highlighted.
func (p *Popup) title(c *ui.Context, r row, color ui.Color) ui.Element {
	if len(r.ranges) == 0 {
		return ui.Text(c, r.title).SingleLine().FontSize(13).TextColor(color)
	}
	var spans []ui.Span
	at := 0
	for _, rg := range r.ranges {
		if rg.Start < at || rg.End > len(r.title) {
			continue
		}
		if rg.Start > at {
			spans = append(spans, ui.Span{Text: r.title[at:rg.Start], Color: color})
		}
		m := ui.Span{Text: r.title[rg.Start:rg.End], Color: color}
		switch p.Settings.HighlightMatch {
		case settings.HighlightItalic:
			m.Italic = true
		case settings.HighlightUnderline:
			m.Underline = true
		case settings.HighlightColor:
			m.Background, m.Color = ui.RGB(255, 230, 0), ui.RGB(0, 0, 0)
		default:
			m.Weight = 700
		}
		spans = append(spans, m)
		at = rg.End
	}
	if at < len(r.title) {
		spans = append(spans, ui.Span{Text: r.title[at:], Color: color})
	}
	return ui.RichText(c, spans...).SingleLine().FontSize(13)
}

func (p *Popup) rowMenu(m *ui.Menu, it *history.Item) {
	id := it.ID
	if m.Item(locales.MenuCopy()).Chosen() {
		p.Host.Select(id, Action{})
	}
	if m.Item(locales.MenuPaste()).Chosen() {
		p.Host.Select(id, Action{Paste: true})
	}
	if m.Item(locales.MenuPastePlain()).Chosen() {
		p.Host.Select(id, Action{Paste: true, Plain: true})
	}
	m.Separator()
	pin := locales.MenuPin()
	if it.Pin != "" {
		pin = locales.MenuUnpin()
	}
	if m.Item(pin).Chosen() {
		p.Host.TogglePin(id)
	}
	if m.Item(locales.MenuDelete()).Chosen() {
		p.Host.Delete(id)
	}
}

func (p *Popup) footer(c *ui.Context) {
	t := c.Theme()
	ui.Divider(c)
	all := c.Modifiers()&ui.Shift != 0
	ui.Column(c).Padding(4, 5).Children(func() {
		clear, clearKey := locales.Clear(), "⌥⌘⌫"
		if all {
			clear, clearKey = locales.ClearAll(), "⇧⌥⌘⌫"
		}
		p.footerItem(c, t, clear, clearKey, func() { p.requestClear(all) })
		p.footerItem(c, t, locales.Preferences(), "⌘,", p.Host.OpenSettings)
		p.footerItem(c, t, locales.About(), "", p.Host.About)
		p.footerItem(c, t, locales.Quit(), "⌘Q", p.Host.Quit)
	})
}

func (p *Popup) footerItem(c *ui.Context, t *ui.Theme, label, key string, do func()) {
	el := ui.Row(c.Key(label)).Padding(4, 8).Radius(5).AlignItems(ui.Center).MinHeight(22)
	text, muted := t.Text, t.TextMuted
	if el.Hovered() {
		el.Background(t.Accent)
		text, muted = t.AccentText, t.AccentText.Alpha(0.75)
	}
	el.OnClick(do).Children(func() {
		ui.Text(c, label).FontSize(13).TextColor(text).Grow(1)
		if key != "" {
			ui.Text(c, key).FontSize(12).TextColor(muted)
		}
	})
}

// autoPreview opens the preview once an item has been chosen for the
// preview delay (openPreviewAutomatically).
func (p *Popup) autoPreview(c *ui.Context) {
	it := p.current()
	if it == nil {
		return
	}
	now := c.Now()
	if it.ID != p.chosenID {
		p.chosenID, p.chosenAt = it.ID, now
	}
	if p.previewOpen || p.previewShut || !p.Settings.OpenPreviewAutomatically {
		return
	}
	delay := time.Duration(p.Settings.PreviewDelay) * time.Millisecond
	if wait := p.chosenAt.Add(delay).Sub(now); wait > 0 {
		c.After(wait)
		return
	}
	p.previewOpen = true
}

func (p *Popup) preview(c *ui.Context, it *history.Item) {
	s := p.Settings
	t := c.Theme()
	ui.Column(c).Width(float32(s.PreviewWidth)).FillHeight().Padding(12).Gap(10).Children(func() {
		if bmp := p.Host.Preview(it); bmp != nil {
			ui.Image(c, bmp).Fit(ui.ScaleDown).MaxHeight(float32(s.WindowHeight) * 0.6).MaxWidth(float32(s.PreviewWidth - 24))
			if text := it.RecognizedText(); text != "" {
				// The text read in the image, which search finds.
				if r := []rune(text); len(r) > 10000 {
					text = string(r[:10000]) + "…"
				}
				ui.Scroll(c).Grow(1).Shrink(1).Children(func() {
					ui.Text(c, text).Wrap().FontSize(11).TextColor(t.TextMuted).Selectable()
				})
			}
		} else {
			text := it.PreviewText()
			if r := []rune(text); len(r) > 10000 {
				text = string(r[:10000]) + "…"
			}
			ui.Scroll(c).Grow(1).Shrink(1).Children(func() {
				ui.Text(c, text).Wrap().FontSize(12).Selectable()
			})
		}
		ui.Divider(c)
		detail := func(label, value string) {
			ui.Row(c).Gap(6).Children(func() {
				ui.Text(c, label).FontSize(11).TextColor(t.TextMuted)
				ui.Text(c, value).FontSize(11).Grow(1).Shrink(1)
			})
		}
		if bmp := p.Host.Preview(it); bmp != nil {
			// The bitmap is scaled down for display; the header has the
			// copied image's own size.
			w, h, ok := p.Host.PreviewSize(it)
			if !ok {
				w, h = bmp.Size()
			}
			detail(locales.PreviewDimensions(), fmt.Sprintf("%d×%d", w, h))
		}
		if it.Application != "" {
			detail(locales.PreviewApplication(), p.Host.AppName(it.Application))
		}
		const layout = "Jan 2, 2006 at 15:04:05"
		detail(locales.PreviewFirstCopyTime(), it.FirstCopiedAt.Format(layout))
		detail(locales.PreviewLastCopyTime(), it.LastCopiedAt.Format(layout))
		detail(locales.PreviewNumberOfCopies(), fmt.Sprint(it.NumberOfCopies))
		ui.Divider(c)
		hint := func(s string) { ui.Text(c, s).FontSize(11).TextColor(t.TextMuted) }
		if it.Pin == "" {
			hint(locales.PreviewPinKey(shortcut.Symbols(s.PinShortcut)))
		} else {
			hint(locales.PreviewUnpinKey(shortcut.Symbols(s.PinShortcut)))
		}
		hint(locales.PreviewDeleteKey(shortcut.Symbols(s.DeleteShortcut)))
		hint(locales.PreviewToggleKey(shortcut.Symbols(s.PreviewShortcut)))
	})
}

func (p *Popup) confirmDialog(c *ui.Context) {
	t := c.Theme()
	ui.DialogBase(c, &p.confirm, func(backdrop, panel ui.Element) {
		backdrop.Background(ui.RGBA(0, 0, 0, 0.3))
		panel.Width(300).Padding(16).Gap(10).Radius(10).Background(t.Background).Border(1, t.Border)
		ui.Text(c, locales.ClearAlertMessage()).Bold().Wrap()
		ui.Text(c, locales.ClearAlertComment()).FontSize(12).TextColor(t.TextMuted).Wrap()
		ui.Checkbox(c, &p.dontAsk, locales.DontAskAgain())
		ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
			if ui.Button(c, locales.ClearAlertCancel()).Clicked() {
				p.confirm = false
			}
			if ui.PrimaryButton(c, locales.ClearAlertConfirm()).AutoFocus().Clicked() {
				p.confirm = false
				if p.dontAsk {
					p.Settings.SuppressClearAlert = true
				}
				p.Host.Clear(p.confirmAll)
			}
		})
	})
}

// Fit is the window size for the rows, as Maccy's popup sizes itself to
// its content, up to the window size of the settings.
func (p *Popup) Fit() (width, height int) {
	s := p.Settings
	h := 10
	if p.searchShown() {
		h += 46
	}
	for _, r := range p.rows {
		if r.title == "" {
			h += min(s.ImageMaxHeight, 200) + 10
		} else {
			h += 30
		}
		if h > s.WindowHeight {
			break
		}
	}
	if len(p.rows) == 0 {
		h += 36
	}
	if s.ShowFooter {
		h += 4*30 + 10
	}
	if p.confirm || p.previewOpen {
		h = max(h, 360)
	}
	h = max(120, min(h, s.WindowHeight))
	w := s.WindowWidth
	if p.previewOpen && p.current() != nil {
		w += s.PreviewWidth + 1
	}
	return w, h
}

// fit asks the host for the window size when it changes.
func (p *Popup) fit() {
	if w, h := p.Fit(); w != p.width || h != p.height {
		p.width, p.height = w, h
		p.Host.Resize(w, h)
	}
}

var searchIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>`))
