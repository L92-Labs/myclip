//go:build windows

package native

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ErrNoOCR says text recognition is unavailable: this Windows has no OCR
// language installed for the user's languages.
var ErrNoOCR = errors.New("native: no Windows OCR language is installed")

// RecognizeText reads the text in an image with the Windows OCR engine
// (through PowerShell, see ocr_script.go). Call it off the main thread.
func RecognizeText(image []byte) (string, error) {
	// The script finds the image by path: a temp file named after the
	// bytes' hash, so concurrent calls and retries reuse one file instead
	// of piling up copies, deleted after the run.
	sum := sha256.Sum256(image)
	path := filepath.Join(os.TempDir(), "myclip-ocr-"+hex.EncodeToString(sum[:8])+".png")
	if err := os.WriteFile(path, image, 0o600); err != nil {
		return "", err
	}
	defer os.Remove(path)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive",
		"-NoLogo", "-ExecutionPolicy", "Bypass", "-EncodedCommand", powershellCommand())
	cmd.Env = append(os.Environ(), "MYCLIP_OCR_INPUT="+path)
	out, err := cmd.Output()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 3 {
			return "", ErrNoOCR
		}
		return "", err
	}
	var lines []string
	for _, l := range strings.Split(string(bytes.TrimSpace(out)), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			lines = append(lines, l)
		}
	}
	return strings.Join(lines, "\n"), nil
}
