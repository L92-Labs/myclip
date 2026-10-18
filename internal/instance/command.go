package instance

import (
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// The running MyClip takes one command from later launches, over a socket
// in the 0700 run dir: "toggle" shows or hides the popup. It reads nothing
// and changes nothing else, so a desktop shortcut can run `myclip --toggle`
// (Wayland compositors keep global keys for themselves).

// Toggle is the only command.
const Toggle = "toggle"

func socketPath(root string) string { return filepath.Join(root, "run", "myclip.sock") }

// Serve listens for commands until the returned closer runs; fn gets each
// known command. Call it while holding the lock: a socket left by a crash
// is removed first, which only the lock holder may do.
func Serve(root string, fn func(cmd string)) (io.Closer, error) {
	path := socketPath(root)
	if info, err := os.Lstat(path); err == nil {
		if info.Mode()&os.ModeSocket == 0 {
			return nil, errors.New("instance: " + path + " is not a socket")
		}
		os.Remove(path)
	}
	l, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(path, 0o600); err != nil {
		l.Close()
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
