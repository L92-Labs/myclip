//go:build windows

package instance

import (
	"os"
	"testing"
)

func socketTestRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "myclip-test-")
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

// The Windows tests run without the stale-socket fixture, which names a
// unix socket; the pipe server needs no crash cleanup.

func TestServeSendWindows(t *testing.T) {
	root := socketTestRoot(t)
	c, err := Serve(root, func(string) {})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if err := Send(root, Toggle); err != nil {
		t.Fatalf("send: %v", err)
	}
}
