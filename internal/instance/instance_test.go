package instance

import (
	"errors"
	"net"
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

func TestToggleCommand(t *testing.T) {
	root := socketTestRoot(t)
	got := make(chan string, 2)
	c, err := Serve(root, func(cmd string) { got <- cmd })
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	info, _ := os.Stat(filepath.Join(root, "run", "myclip.sock"))
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("socket mode %o", info.Mode().Perm())
	}
	if err := Send(root, "clear"); err != nil { // unknown: dropped
		t.Fatal(err)
	}
	if err := Send(root, Toggle); err != nil {
		t.Fatal(err)
	}
	if cmd := <-got; cmd != Toggle {
		t.Fatalf("got %q", cmd)
	}
	select {
	case cmd := <-got:
		t.Fatalf("unknown command delivered: %q", cmd)
	default:
	}
}

func TestServeReplacesAStaleSocket(t *testing.T) {
	root := socketTestRoot(t)
	c1, err := Serve(root, func(string) {})
	if err != nil {
		t.Fatal(err)
	}
	c1.(*net.UnixListener).SetUnlinkOnClose(false)
	c1.Close() // a crash leaves the socket inode behind
	if info, err := os.Lstat(socketPath(root)); err != nil || info.Mode()&os.ModeSocket == 0 {
		t.Fatalf("stale socket missing: %v", err)
	}
	c2, err := Serve(root, func(string) {})
	if err != nil {
		t.Fatalf("relaunch: %v", err)
	}
	c2.Close()
}

func socketTestRoot(t *testing.T) string {
	t.Helper()
	// macOS limits Unix socket paths to 104 bytes. Its long TMPDIR plus
	// t.TempDir's test-name prefix can exceed that before the socket name.
	dir, err := os.MkdirTemp("/tmp", "myclip-test-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	root, err := DataRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	return root
}
