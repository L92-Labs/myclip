//go:build !windows

package instance

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// Lock takes the fail-closed flock on root/run/myclip.lock. The file is
// never unlinked (a crash-left file must not block relaunch); the kernel
// drops the lock when the process dies. Keep the returned file open for the
// process lifetime.
func Lock(root string) (*os.File, error) {
	path := filepath.Join(root, "run", "myclip.lock")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0o600)
	if err != nil {
		return nil, fmt.Errorf("instance: open lock: %w", err)
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if !info.Mode().IsRegular() {
		f.Close()
		return nil, fmt.Errorf("instance: %s is not a regular file", path)
	}
	if st, ok := info.Sys().(*syscall.Stat_t); ok && int(st.Uid) != os.Getuid() {
		f.Close()
		return nil, fmt.Errorf("instance: %s is owned by uid %d", path, st.Uid)
	}
	if info.Mode().Perm() != 0o600 {
		if err := f.Chmod(0o600); err != nil {
			f.Close()
			return nil, err
		}
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, ErrRunning
		}
		return nil, fmt.Errorf("instance: flock: %w", err)
	}
	return f, nil
}
