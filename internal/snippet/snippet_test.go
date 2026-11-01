package snippet

import (
	"strings"
	"testing"
	"time"
)

var now = time.Date(2026, 11, 1, 9, 7, 0, 0, time.UTC)

func TestExpand(t *testing.T) {
	clip := "from clipboard"
	tests := []struct{ in, want string }{
		{"{date}", "2026-11-01"},
		{"{time}", "09:07"},
		{"{clipboard}", "from clipboard"},
		{"on {date} at {time}", "on 2026-11-01 at 09:07"},
		{"id={uuid}", ""},
		{"{unknown} stays", "{unknown} stays"},
		{"{Date} is not a placeholder", "{Date} is not a placeholder"},
		{"no placeholders", "no placeholders"},
	}
	for _, tt := range tests {
		got := Expand(tt.in, now, func() string { return clip })
		if tt.want == "" && strings.HasPrefix(tt.in, "id=") {
			if !strings.HasPrefix(got, "id=") || len(got) != len("id=xxxxxxxx-xxxx-4xxx-xxxx-xxxxxxxxxxxx") {
				t.Fatalf("uuid: got %q", got)
			}
			continue
		}
		if got != tt.want {
			t.Fatalf("Expand(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
	// A nil clipboard reader expands to nothing, not a crash.
	if got := Expand("x{clipboard}y", now, nil); got != "xy" {
		t.Fatalf("nil clipboard: %q", got)
	}
}

func TestUUIDsDiffer(t *testing.T) {
	a := Expand("{uuid}", now, nil)
	b := Expand("{uuid}", now, nil)
	if a == b {
		t.Fatalf("two expands gave the same uuid %q", a)
	}
}

func TestHas(t *testing.T) {
	if !Has("see {date}") {
		t.Fatal("Has missed {date}")
	}
	if Has("plain text") {
		t.Fatal("plain text has no placeholder")
	}
}
