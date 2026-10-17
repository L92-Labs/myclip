package instance

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDataRootModes(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "com.l92labs.myclip")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := DataRoot(dir); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{dir, filepath.Join(dir, "run")} {
		info, err := os.Stat(d)
		if err != nil || info.Mode().Perm() != 0o700 {
			t.Fatalf("%s: %v %v", d, info.Mode().Perm(), err)
		}
	}
	if _, err := DataRoot("relative/dir"); err == nil {
		t.Fatal("relative data root accepted")
	}
}

func TestSecondLockFailsThenRelaunchReusesFile(t *testing.T) {
	root, err := DataRoot(filepath.Join(t.TempDir(), "root"))
	if err != nil {
		t.Fatal(err)
	}
	first, err := Lock(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Lock(root); !errors.Is(err, ErrRunning) {
		t.Fatalf("second lock: %v, want ErrRunning", err)
	}
	first.Close() // the holder exits; the file stays
	if !Exists(filepath.Join(root, "run", "myclip.lock")) {
		t.Fatal("lock file was removed")
	}
	again, err := Lock(root)
	if err != nil {
		t.Fatalf("relaunch: %v", err)
	}
	again.Close()
}

func TestLockRefusesSymlink(t *testing.T) {
	root, _ := DataRoot(filepath.Join(t.TempDir(), "root"))
	target := filepath.Join(t.TempDir(), "elsewhere")
	os.WriteFile(target, nil, 0o644)
	if err := os.Symlink(target, filepath.Join(root, "run", "myclip.lock")); err != nil {
		t.Fatal(err)
	}
	if f, err := Lock(root); err == nil {
		f.Close()
		t.Fatal("lock followed a symlink")
	}
}
