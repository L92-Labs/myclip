// Package instance holds MyClip's data root and its single-instance lock.
package instance

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ErrRunning means another MyClip holds the lock for this data root.
var ErrRunning = errors.New("instance: MyClip is already running")

// DataRoot creates dir and dir/run 0700 (tightening existing ones) and
// returns dir. It refuses a symlinked dir.
func DataRoot(dir string) (string, error) {
	if !filepath.IsAbs(dir) {
		return "", fmt.Errorf("instance: data root %q is not absolute", dir)
	}
	for _, d := range []string{dir, filepath.Join(dir, "run")} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			return "", err
		}
		info, err := os.Lstat(d)
		if err != nil {
			return "", err
		}
		if !info.IsDir() {
			return "", fmt.Errorf("instance: %s is not a directory", d)
		}
		if info.Mode().Perm() != 0o700 {
			if err := os.Chmod(d, 0o700); err != nil {
				return "", err
			}
		}
	}
	return dir, nil
}

// Exists reports whether path exists (used to tell a fresh install apart).
func Exists(path string) bool {
	_, err := os.Lstat(path)
	return !errors.Is(err, fs.ErrNotExist)
}
