//go:build windows

package instance

import (
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Windows 10 1803+ has AF_UNIX sockets, which Go's net package speaks.
// The pipe-shaped path would need hand-rolled win32 code; the socket sits
// in the data root, where DataRoot has made everything 0700.
func socketPath(root string) string { return filepath.Join(root, "run", "myclip.sock") }

// Serve listens for commands until the returned closer runs; fn gets each
// known command. Call it while holding the lock: a socket left by a crash
// is removed first, which only the lock holder may do.
func Serve(root string, fn func(cmd string)) (io.Closer, error) {
	path := socketPath(root)
	// A socket left by a crash is removed; only the lock holder runs this.
	// (Windows' Lstat does not flag socket files, so anything stale goes.)
	if _, err := os.Lstat(path); err == nil {
		os.Remove(path)
	}
	l, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.SetReadDeadline(time.Now().Add(time.Second))
			buf := make([]byte, 16)
			n, _ := io.ReadFull(conn, buf)
			conn.Close()
			if cmd := strings.TrimSpace(string(buf[:n])); cmd == Toggle {
				fn(cmd)
			}
		}
	}()
	return l, nil
}

// Send gives the running MyClip a command.
func Send(root, cmd string) error {
	conn, err := net.DialTimeout("unix", socketPath(root), time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(cmd + "\n"))
	return err
}
