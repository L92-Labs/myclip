// Package snippet expands clip placeholders at paste time: a clip marked
// as a snippet can carry {date}, {time}, {clipboard} and {uuid}, which
// become the real values when it goes back on the pasteboard.
package snippet

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

// Has reports whether s holds any known placeholder.
func Has(s string) bool {
	for _, p := range []string{"{date}", "{time}", "{clipboard}", "{uuid}"} {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}

// Expand replaces the known placeholders in s. Unknown {braced} text is
// left alone, so prose that merely looks like a placeholder survives.
func Expand(s string, now time.Time, clipboard func() string) string {
	if !Has(s) {
		return s
	}
	r := strings.NewReplacer(
		"{date}", now.Format("2006-01-02"),
		"{time}", now.Format("15:04"),
		"{clipboard}", func() string {
			if clipboard == nil {
				return ""
			}
			return clipboard()
		}(),
		"{uuid}", uuid(),
	)
	return r.Replace(s)
}

// uuid returns a random RFC 4122 version-4 UUID.
func uuid() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "" // no randomness: leave the placeholder gone rather than fake it
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
