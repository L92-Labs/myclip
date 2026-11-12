package engine

import (
	"context"
	"testing"
	"time"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
	"github.com/L92-Labs/myclip/internal/settings"
)

func importItem(s string) *history.Item {
	return &history.Item{
		FirstCopiedAt: time.Unix(1_000, 0), LastCopiedAt: time.Unix(2_000, 0),
		NumberOfCopies: 3,
		Contents:       []history.Content{{Type: pasteboard.TypeString, Value: []byte(s)}},
	}
}

func TestImportAddsAndSkips(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "already here")
	added, skipped, err := r.e.ImportItems([]*history.Item{
		importItem("new clip"),
		importItem("already here"), // same contents: skipped
	})
	if err != nil {
		t.Fatal(err)
	}
	if added != 1 || skipped != 1 {
		t.Fatalf("added %d skipped %d", added, skipped)
	}
	it := r.e.find(2)
	if it == nil || string(it.Contents[0].Value) != "new clip" {
		t.Fatalf("imported item missing: %+v", it)
	}
	// Timestamps are preserved, not re-set to now.
	if it.FirstCopiedAt != time.Unix(1_000, 0) || it.NumberOfCopies != 3 {
		t.Fatalf("re-timed: %+v", it)
	}
	// The import survives a reload.
	got, err := r.s.Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("reloaded %d items", len(got))
	}
}

func TestImportPinCollisionMovesKey(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "existing")
	if err := r.e.SetPin(r.e.Items()[0].ID, "b"); err != nil {
		t.Fatal(err)
	}
	pinned := importItem("imported pinned")
	pinned.Pin = "b" // taken
	added, _, err := r.e.ImportItems([]*history.Item{pinned})
	if err != nil || added != 1 {
		t.Fatalf("added %d err %v", added, err)
	}
	it := r.e.Items()
	if it[0].Pin != "b" {
		t.Fatalf("existing pin moved: %+v", it[0])
	}
	if it[1].Pin != "c" { // the next free key
		t.Fatalf("imported pin: %+v", it[1])
	}
}

func TestImportPinDroppedWhenNoKeysFree(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "existing")
	if err := r.e.SetPin(r.e.Items()[0].ID, "b"); err != nil {
		t.Fatal(err)
	}
	// No free keys: every supported key is reserved.
	r.e.SetRules(nil, history.SupportedPins...)
	pinned := importItem("pinned")
	pinned.Pin = "b" // taken, and no other key may be used
	added, _, err := r.e.ImportItems([]*history.Item{pinned})
	if err != nil || added != 1 {
		t.Fatalf("added %d err %v", added, err)
	}
	if got := r.e.Items()[0].Pin; got != "b" || r.e.Items()[1].Pin != "" {
		t.Fatalf("pins: %q %q", r.e.Items()[0].Pin, r.e.Items()[1].Pin)
	}
}

// A pin an in-popup shortcut reserves is a conflict too, even on an empty
// history: the popup checks pinned rows before its own shortcuts.
func TestImportReservedPinMovesKey(t *testing.T) {
	r := newRig(t, Limits{})
	r.e.SetRules(nil, "p") // Alt+P: the pin key p is reserved
	imp := importItem("imported pinned")
	imp.Pin = "p"
	added, _, err := r.e.ImportItems([]*history.Item{imp})
	if err != nil || added != 1 {
		t.Fatalf("added %d err %v", added, err)
	}
	if got := r.e.Items()[0].Pin; got == "p" || got == "" {
		t.Fatalf("reserved pin kept: %q", got)
	}
}

// Imported clips get the live rules' effects, as a fresh capture would.
func TestImportAppliesExpiryRules(t *testing.T) {
	r := newRig(t, Limits{})
	r.e.SetRules([]Rule{mustRule(t, `one-time`, settings.RuleExpire, 30)})
	imp := importItem("one-time code 42")
	added, _, err := r.e.ImportItems([]*history.Item{imp})
	if err != nil || added != 1 {
		t.Fatalf("added %d err %v", added, err)
	}
	it := r.e.Items()[0]
	at, ok := r.e.expires[it.ID]
	if !ok {
		t.Fatal("no expiry after import")
	}
	if !at.Equal(imp.LastCopiedAt.Add(30 * time.Second)) {
		t.Fatalf("expiry %v, want last copy + 30s", at)
	}
	// And the sweep really deletes it when the deadline has run out.
	r.clock = imp.LastCopiedAt.Add(31 * time.Second)
	r.e.Tick(context.Background())
	if items := r.e.Items(); len(items) != 0 {
		t.Fatalf("imported clip outlived its rule: %+v", items)
	}
}
