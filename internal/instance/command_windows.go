//go:build windows

package instance

import (
	"io"
	"net"
	"path/filepath"
	"strings"
	"time"
)

// Serve listens on \\.\pipe\<name> for one command at a time; fn gets each
// known command. The pipe is per data root, so two MyClips with different
// data roots do not cross wires.
func Serve(root string, fn func(cmd string)) (io.Closer, error) {
	l, err := net.Listen("pipe", pipePath(root))
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
	conn, err := net.DialTimeout("pipe", pipePath(root), time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(cmd + "\n"))
	return err
}

func pipePath(root string) string {
	// A short, stable name: pipe names are flat and length-bound.
	return `myclip-` + filepath.Base(filepath.Clean(root))
}
