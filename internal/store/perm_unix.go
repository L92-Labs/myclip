//go:build !windows

package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"
)

// PrepareFile creates path O_EXCL 0600, or tightens an existing regular file
// owned by this user to 0600. Existing sidecars get the same check. It
// reports whether the file was created now. It never follows a symlink.
func PrepareFile(path string) (created bool, err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, FileMode)
	if err == nil {
		// umask cannot widen 0600, but a pre-set ACL-free fchmod pins it.
		err = f.Chmod(FileMode)
		return true, errors.Join(err, f.Close())
	}
	if !errors.Is(err, fs.ErrExist) {
		return false, fmt.Errorf("store: create %s: %w", path, err)
	}
	if err := tighten(path); err != nil {
		return false, err
	}
	for _, s := range Sidecars {
		if err := tighten(path + s); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return false, err
		}
	}
	return false, nil
}

// tighten opens path without following links, requires a regular file owned
// by this user, and sets it to 0600.
func tighten(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return fmt.Errorf("store: refuse %s: %w", path, err)
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("store: stat %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("store: refuse %s: not a regular file", path)
	}
	if st, ok := info.Sys().(*syscall.Stat_t); ok && int(st.Uid) != os.Getuid() {
		return fmt.Errorf("store: refuse %s: owned by uid %d", path, st.Uid)
	}
	if info.Mode().Perm() != FileMode {
		if err := f.Chmod(FileMode); err != nil {
			return fmt.Errorf("store: chmod %s: %w", path, err)
		}
	}
	return nil
}
