//go:build linux

package native

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

// ErrNoOCR says text recognition needs a tool this system lacks.
var ErrNoOCR = errors.New("native: install tesseract to search text in images")

// RecognizeText reads the text in an image with the tesseract command, on
// this computer, when it is installed. Call it off the main thread.
func RecognizeText(image []byte) (string, error) {
	if _, err := exec.LookPath("tesseract"); err != nil {
		return "", ErrNoOCR
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "tesseract", "stdin", "stdout")
	cmd.Stdin = bytes.NewReader(image)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	var lines []string
	for _, l := range strings.Split(string(out), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			lines = append(lines, l)
		}
	}
	return strings.Join(lines, "\n"), nil
}
