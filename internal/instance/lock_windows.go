//go:build windows

package instance

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// Lock takes the fail-closed lock on root\run\myclip.lock with LockFileEx.
// The file is never unlinked (a crash-left file must not block relaunch);
// Windows releases the lock when the process dies. Keep the returned file
// open for the process lifetime.
func Lock(root string) (*os.File, error) {
	path := filepath.Join(root, "run", "myclip.lock")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, fmt.Errorf("instance: open lock: %w", err)
	}
	// Lock the first byte exclusively, non-blocking.
	var ov windows.Overlapped
	err = windows.LockFileEx(windows.Handle(f.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &ov)
	if err != nil {
		f.Close()
		if errors.Is(err, windows.ERROR_LOCK_VIOLATION) || errors.Is(err, windows.ERROR_IO_PENDING) {
			return nil, ErrRunning
		}
		return nil, fmt.Errorf("instance: lock: %w", err)
	}
	return f, nil
}
