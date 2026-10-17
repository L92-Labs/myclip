package locales

import (
	"strings"
	"testing"
)

func TestTranslations(t *testing.T) {
	if got := In("fr").Clear(); got != "Supprimer" {
		t.Errorf("fr Clear = %q", got)
	}
	if got := In("es").TrayClips(1); !strings.HasPrefix(got, "1 elemento") {
		t.Errorf("es TrayClips(1) = %q", got)
	}
	if got := In("fr").TrayClips(3); !strings.HasPrefix(got, "3 éléments") {
		t.Errorf("fr TrayClips(3) = %q", got)
	}
	if got := In("es").PreviewPinKey("⌥P"); got != "Presiona ⌥P para anclar." {
		t.Errorf("es PreviewPinKey = %q", got)
	}
}
