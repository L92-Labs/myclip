//go:build windows

package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// PrepareFile creates path 0600, or leaves an existing one. The sidecar
// checks of the Unix build have no Windows analogue here: access is gated
// by the %AppData% profile's ACLs, and Go's 0600 mode maps to an owner-only
// file on Windows.
//
// Deviation from the Unix recipe (PARITY D6 hardening): no O_NOFOLLOW and
// no same-user ownership check — Windows has no per-file stat uid, and
// symlink following is not a same-user threat model here.
func PrepareFile(path string) (created bool, err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, FileMode)
	if err == nil {
		err = f.Chmod(FileMode)
		return true, errors.Join(err, f.Close())
	}
	if !errors.Is(err, fs.ErrExist) {
		return false, fmt.Errorf("store: create %s: %w", path, err)
	}
	return false, nil
}
