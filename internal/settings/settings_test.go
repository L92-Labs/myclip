package settings

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
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
