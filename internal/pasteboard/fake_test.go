package pasteboard_test

import (
	"testing"

	"github.com/L92-Labs/myclip/internal/capture"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

func TestConcealedOnSecondItemZeroReads(t *testing.T) {
	f := pasteboard.NewFake()
	f.SetItems([]pasteboard.Item{
		{Types: map[string][]byte{pasteboard.TypeString: []byte("visible")}},
		{Types: map[string][]byte{pasteboard.TypeConcealed: []byte("")}},
	}, nil)
	f.ResetReads()

	snap, err := f.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if f.DataReads() != 0 {
		t.Fatalf("type enumeration must not read data; reads=%d", f.DataReads())
	}
	prefs := capture.Prefs{
		EnabledPasteboardTypes: pasteboard.DefaultEnabledPasteboardTypes,
		IgnoredPasteboardTypes: pasteboard.DefaultIgnoredPasteboardTypes,
	}
	if !capture.ShouldIgnoreTypes(snap.Types, prefs) {
		t.Fatal("concealed on any item must ignore capture")
	}
	// Ignore path: never ReadItems.
	if f.DataReads() != 0 {
		t.Fatalf("reads=%d after ignore decision", f.DataReads())
	}
}

func TestN6UnionIncludesBoardTypesSuperset(t *testing.T) {
	f := pasteboard.NewFake()
	f.SetItems(
		[]pasteboard.Item{{Types: map[string][]byte{pasteboard.TypeString: []byte("x")}}},
		[]string{pasteboard.TypeString, pasteboard.TypeConcealed}, // board-level concealed
	)
	snap, _ := f.Snapshot()
	found := false
	for _, typ := range snap.Types {
		if typ == pasteboard.TypeConcealed {
			found = true
		}
	}
	if !found {
		t.Fatal("board types must appear in N6 union")
	}
}
