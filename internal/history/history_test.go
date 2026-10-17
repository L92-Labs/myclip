package history_test

import (
	"testing"
	"time"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

func textItem(s string) *history.Item {
	return history.NewItem([]history.Content{{
		Type:  pasteboard.TypeString,
		Value: []byte(s),
	}}, "com.example.app", time.Unix(100, 0))
}

func TestSupersedesMergePreservesStoredContents(t *testing.T) {
	h := history.New()
	a := textItem("hello")
	h.Add(a, 1, false, time.Unix(100, 0))

	// Same text again → supersedes merge; stored contents win (identical here).
	b := textItem("hello")
	res := h.Add(b, 2, false, time.Unix(200, 0))
	if !res.Merged {
		t.Fatal("expected merge")
	}
	if h.Len() != 1 {
		t.Fatalf("len=%d want 1", h.Len())
	}
	got := h.Items()[0]
	if got.NumberOfCopies != 2 {
		t.Fatalf("copies=%d want 2", got.NumberOfCopies)
	}
	if got.FirstCopiedAt != time.Unix(100, 0) {
		t.Fatalf("firstCopiedAt=%v", got.FirstCopiedAt)
	}
	if string(got.Contents[0].Value) != "hello" {
		t.Fatalf("contents=%q", got.Contents[0].Value)
	}
}

func TestModifiedMergeKeepsModifierContents(t *testing.T) {
	h := history.New()
	orig := textItem("old")
	h.Add(orig, 10, false, time.Unix(100, 0))

	// Modifier: same supersede set elsewhere? isModified path: item carries
	// ModifiedType whose value is prior changeCount 10, with new text bytes.
	mod := history.NewItem([]history.Content{
		{Type: pasteboard.TypeString, Value: []byte("new")},
		{Type: pasteboard.TypeModified, Value: []byte("10")},
	}, "com.example.app", time.Unix(200, 0))
	// Ensure supersedes does NOT match (different string bytes) so findSimilar
	// falls through to isModified.
	res := h.Add(mod, 11, false, time.Unix(200, 0))
	if !res.Merged {
		t.Fatal("expected isModified merge")
	}
	if h.Len() != 1 {
		t.Fatalf("len=%d want 1 (no dup)", h.Len())
	}
	got := h.Items()[0]
	if string(got.Contents[0].Value) != "new" {
		t.Fatalf("modifier contents lost: %q", got.Contents[0].Value)
	}
	if got.NumberOfCopies != 2 {
		t.Fatalf("copies=%d want 2", got.NumberOfCopies)
	}
}

func TestOwnWriteMergePreservesPinAndTitle(t *testing.T) {
	h := history.New()
	a := textItem("pinme")
	a.Pin = "a"
	a.Title = "titled"
	h.Add(a, 1, false, time.Unix(100, 0))

	// Own-write copy-back with marker (stripped from stored reps).
	b := history.NewItem([]history.Content{
		{Type: pasteboard.TypeString, Value: []byte("pinme")},
		{Type: pasteboard.TypeOwnWrite, Value: []byte("")},
	}, "com.l92labs.myclip", time.Unix(200, 0))
	if !history.HasOwnWriteMarker([]history.Content{
		{Type: pasteboard.TypeString, Value: []byte("pinme")},
		{Type: pasteboard.TypeOwnWrite, Value: []byte("")},
	}) {
		t.Fatal("marker not detected")
	}
	res := h.Add(b, 2, true, time.Unix(200, 0))
	if !res.Merged {
		t.Fatal("expected merge")
	}
	got := h.Items()[0]
	if got.Pin != "a" || got.Title != "titled" {
		t.Fatalf("pin/title not preserved: pin=%q title=%q", got.Pin, got.Title)
	}
	for _, c := range got.Contents {
		if c.Type == pasteboard.TypeOwnWrite {
			t.Fatal("own-write marker must not be stored as a representation")
		}
	}
	if got.NumberOfCopies != 2 {
		t.Fatalf("copies=%d", got.NumberOfCopies)
	}
}

func TestNoDupOnDistinctItems(t *testing.T) {
	h := history.New()
	h.Add(textItem("a"), 1, false, time.Unix(1, 0))
	h.Add(textItem("b"), 2, false, time.Unix(2, 0))
	if h.Len() != 2 {
		t.Fatalf("len=%d want 2", h.Len())
	}
}
