package settings

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestLoadMissingGivesDefaults(t *testing.T) {
	s, err := Load(filepath.Join(t.TempDir(), FileName))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s, Default()) {
		t.Fatalf("got %+v", s)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)
	s := Default()
	s.SearchMode = SearchFuzzy
	s.PasteByDefault = true
	s.IgnoredApps = []string{"com.example.app"}
	s.HistorySize = 50
	if err := Save(path, s); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("mode %v", fi.Mode().Perm())
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, s) {
		t.Fatalf("got %+v\nwant %+v", got, s)
	}
}

func TestCorruptGivesDefaultsAndError(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)
	os.WriteFile(path, []byte("{nope"), 0o600)
	s, err := Load(path)
	if err == nil {
		t.Fatal("no error")
	}
	if !reflect.DeepEqual(s, Default()) {
		t.Fatal("not defaults")
	}
}

func TestNormalize(t *testing.T) {
	s := Default()
	s.SearchMode = "bogus"
	s.HistorySize = 100000
	s.ClipboardCheckInterval = 5
	s.PinShortcut = ""
	s.WindowPositionX = 7
	s.Normalize()
	if s.SearchMode != SearchExact || s.HistorySize != MaxHistorySize ||
		s.ClipboardCheckInterval != MinCheckIntervalMs || s.PinShortcut != DefaultPinShortcut ||
		s.WindowPositionX != 0.5 {
		t.Fatalf("got %+v", s)
	}
}

func TestPartialFileKeepsDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)
	os.WriteFile(path, []byte(`{"pasteByDefault": true}`), 0o600)
	s, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !s.PasteByDefault || s.HistorySize != 200 || !s.ShowFooter {
		t.Fatalf("got %+v", s)
	}
}

func TestResetKeepsFileAsideAndWritesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, FileName)
	os.WriteFile(path, []byte("{"), 0o600)
	aside, err := Reset(path, time.Unix(1700000000, 0))
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(aside); string(b) != "{" {
		t.Fatalf("aside %q holds %q", aside, b)
	}
	if s, err := Load(path); err != nil || s.HistorySize != Default().HistorySize {
		t.Fatalf("after reset: %v %+v", err, s)
	}
}

// A reset that cannot write leaves the unreadable file where it was.
func TestResetFailureRestores(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, FileName)
	os.WriteFile(path, []byte("{"), 0o600)
	os.Chmod(dir, 0o500)
	defer os.Chmod(dir, 0o700)
	if _, err := Reset(path, time.Now()); err == nil {
		t.Fatal("reset succeeded in a read-only dir")
	}
	if b, _ := os.ReadFile(path); string(b) != "{" {
		t.Fatalf("original changed: %q", b)
	}
}
