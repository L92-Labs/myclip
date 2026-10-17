package store

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/pasteboard"
)

func open(t *testing.T) (*Store, string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.Chmod(dir, DirMode); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "history.db")
	s, err := Open(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s, path
}

func item(text string, at time.Time) *history.Item {
	return history.NewItem([]history.Content{
		{Type: pasteboard.TypeString, Value: []byte(text)},
		{Type: pasteboard.TypeHTML, Value: []byte("<b>" + text + "</b>")},
	}, "com.apple.TextEdit", at)
}

func mode(t *testing.T, path string) os.FileMode {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	return info.Mode().Perm()
}

func TestCreatesFile0600(t *testing.T) {
	old := setUmask(0)
	defer setUmask(old)
	s, path := open(t)
	if m := mode(t, path); m != FileMode {
		t.Fatalf("db mode %o, want 600", m)
	}
	if err := s.Save(context.Background(), item("hello", time.Now())); err != nil {
		t.Fatal(err)
	}
	// DELETE-mode journals are transient; any sidecar left must be 0600.
	for _, sc := range Sidecars {
		if info, err := os.Lstat(path + sc); err == nil && info.Mode().Perm() != FileMode {
			t.Errorf("%s mode %o, want 600", sc, info.Mode().Perm())
		}
	}
}

func TestJournalModeDeleteAndSecureDelete(t *testing.T) {
	s, path := open(t)
	ctx := context.Background()
	if m, err := s.JournalMode(ctx); err != nil || m != "delete" {
		t.Fatalf("journal_mode %q, %v", m, err)
	}
	if _, err := os.Stat(path + "-wal"); err == nil {
		t.Fatal("a -wal file exists in DELETE mode")
	}
}

func TestPrepareTightensExisting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.db")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path+"-journal", nil, 0o644); err != nil {
		t.Fatal(err)
	}
	created, err := PrepareFile(path)
	if err != nil || created {
		t.Fatalf("PrepareFile = %v, %v; want existing file accepted", created, err)
	}
	if m := mode(t, path); m != FileMode {
		t.Errorf("db mode %o after tighten", m)
	}
	if m := mode(t, path+"-journal"); m != FileMode {
		t.Errorf("journal mode %o after tighten", m)
	}
}

func TestPrepareRefusesSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "elsewhere")
	if err := os.WriteFile(target, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "history.db")
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
	if _, err := PrepareFile(path); err == nil {
		t.Fatal("PrepareFile followed a symlink")
	}
	if m := mode(t, target); m != 0o644 {
		t.Errorf("symlink target changed to %o", m)
	}
}

func TestRoundTripUpdateDelete(t *testing.T) {
	s, path := open(t)
	ctx := context.Background()
	t0 := time.Unix(1_700_000_000, 0)
	a, b := item("first", t0), item("second", t0.Add(time.Second))
	for _, it := range []*history.Item{a, b} {
		if err := s.Save(ctx, it); err != nil {
			t.Fatal(err)
		}
	}
	if a.ID == 0 || b.ID == 0 || a.ID == b.ID {
		t.Fatalf("ids %d %d", a.ID, b.ID)
	}

	a.LastCopiedAt = t0.Add(time.Minute)
	a.NumberOfCopies = 3
	a.Pin = "b"
	a.Contents = a.Contents[:1]
	if err := s.Save(ctx, a); err != nil {
		t.Fatal(err)
	}

	s.Close()
	s2, err := Open(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()
	got, err := s2.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ID != a.ID || got[1].ID != b.ID {
		t.Fatalf("order/ids: %+v", got)
	}
	g := got[0]
	if g.NumberOfCopies != 3 || g.Pin != "b" || !g.LastCopiedAt.Equal(a.LastCopiedAt) ||
		g.Application != "com.apple.TextEdit" || len(g.Contents) != 1 || g.Text() != "first" {
		t.Fatalf("updated item: %+v", g)
	}
	if len(got[1].Contents) != 2 || got[1].Contents[1].Type != pasteboard.TypeHTML {
		t.Fatalf("contents order: %+v", got[1].Contents)
	}

	if err := s2.Delete(ctx, a.ID); err != nil {
		t.Fatal(err)
	}
	got, _ = s2.Load(ctx)
	if len(got) != 1 || got[0].ID != b.ID {
		t.Fatalf("after delete: %+v", got)
	}
}

// Forensics: a deleted clip's bytes must not remain anywhere in the file.
func TestDeletedContentLeavesNoBytes(t *testing.T) {
	s, path := open(t)
	ctx := context.Background()
	secret := "sentinel-7f3a9c-do-not-keep"
	it := item(secret, time.Now())
	if err := s.Save(ctx, it); err != nil {
		t.Fatal(err)
	}
	// Control: the committed clip is in the file, so the check below is real.
	if raw, _ := os.ReadFile(path); !bytes.Contains(raw, []byte(secret)) {
		t.Fatal("control: saved content not found in the database file")
	}
	if err := s.Delete(ctx, it.ID); err != nil {
		t.Fatal(err)
	}
	s.Close()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte(secret)) {
		t.Fatal("deleted content found in the database file")
	}
}
