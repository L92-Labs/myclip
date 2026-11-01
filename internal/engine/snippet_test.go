package engine

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/L92-Labs/myclip/internal/pasteboard"
	"github.com/L92-Labs/myclip/internal/snippet"
)

func TestSnippetPersistsThroughStore(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "see {date}")
	id := r.e.Items()[0].ID
	if err := r.e.SetSnippet(id, true); err != nil {
		t.Fatal(err)
	}
	got, err := r.s.Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !got[0].Snippet {
		t.Fatalf("snippet flag not stored: %+v", got)
	}
	if err := r.e.SetSnippet(id, false); err != nil {
		t.Fatal(err)
	}
	if got, _ = r.s.Load(context.Background()); got[0].Snippet {
		t.Fatal("snippet flag not cleared")
	}
}

func TestCopyBackExpandsSnippet(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "on {date} at {time}: {clipboard} {uuid}")
	// What sits on the pasteboard now (someone else's copy) is what
	// {clipboard} expands to.
	r.f.SetItems([]pasteboard.Item{{Types: map[string][]byte{
		pasteboard.TypeString: []byte("the clipboard text"),
	}}}, nil)
	id := r.e.Items()[0].ID
	if err := r.e.SetSnippet(id, true); err != nil {
		t.Fatal(err)
	}
	if err := r.e.CopyBack(id, false); err != nil {
		t.Fatal(err)
	}
	written := string(r.f.ReadItems()[0].Types[pasteboard.TypeString])
	if !strings.Contains(written, "the clipboard text") {
		t.Fatalf("{clipboard} not expanded: %q", written)
	}
	if strings.Contains(written, "{date}") || strings.Contains(written, "{uuid}") {
		t.Fatalf("placeholders remain: %q", written)
	}
	if _, err := time.Parse("2006-01-02", snippet.Expand("{date}", time.Now(), nil)); err != nil {
		t.Fatal(err)
	}
	// The stored clip keeps its raw text: expansion is paste-time only.
	if got := r.e.Items()[0].ContentText(); !strings.Contains(got, "{date}") {
		t.Fatalf("stored clip was rewritten: %q", got)
	}
}

func TestCopyBackPlainSnippetLeavesOthersUntouched(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "plain clip")
	if err := r.e.CopyBack(r.e.Items()[0].ID, false); err != nil {
		t.Fatal(err)
	}
	if got := string(r.f.ReadItems()[0].Types[pasteboard.TypeString]); got != "plain clip" {
		t.Fatalf("non-snippet rewritten: %q", got)
	}
	if r.f.ReadItems()[0].Types[pasteboard.TypeOwnWrite] == nil {
		t.Fatal("own-write marker lost")
	}
}
