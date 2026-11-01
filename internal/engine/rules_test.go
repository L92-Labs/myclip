package engine

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/settings"
)

func mustRule(t *testing.T, pattern, action string, seconds int) Rule {
	t.Helper()
	re, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatal(err)
	}
	return Rule{Pattern: re, Action: action, Seconds: seconds}
}

func TestRulePinsCapturedClips(t *testing.T) {
	r := newRig(t, Limits{})
	r.e.SetRules([]Rule{mustRule(t, `^secret`, settings.RulePin, 0)}, "p", "q")
	r.copy(t, "secret stuff")
	// The first pins go to the first free keys; a and p, q are reserved.
	items := r.e.Items()
	if len(items) != 1 || items[0].Pin != "b" {
		t.Fatalf("items %+v", items[0])
	}
	stored, _ := r.s.Load(context.Background())
	if len(stored) != 1 || stored[0].Pin != "b" {
		t.Fatalf("stored %+v", stored)
	}
	r.copy(t, "secret again")
	if items := r.e.Items(); items[0].Pin != "c" {
		t.Fatalf("second pin %q", items[0].Pin)
	}
	r.copy(t, "secret three")
	if items := r.e.Items(); items[0].Pin != "d" {
		t.Fatalf("third pin %q", items[0].Pin)
	}
}

// With every pin key reserved, a matching clip stays unpinned.
func TestRulePinWithNoFreeKeys(t *testing.T) {
	r := newRig(t, Limits{})
	r.e.SetRules([]Rule{mustRule(t, `^secret`, settings.RulePin, 0)}, history.SupportedPins...)
	r.copy(t, "secret stuff")
	if items := r.e.Items(); len(items) != 1 || items[0].Pin != "" {
		t.Fatalf("items %+v", items)
	}
}

func TestRuleFirstMatchWins(t *testing.T) {
	r := newRig(t, Limits{})
	r.e.SetRules([]Rule{
		mustRule(t, `otp`, settings.RuleExpire, 5),
		mustRule(t, `otp|code`, settings.RulePin, 0),
	})
	r.copy(t, "your otp is 1234")
	items := r.e.Items()
	if len(items) != 1 || items[0].Pin != "" {
		t.Fatalf("items %+v", items)
	}
	if got, ok := r.e.expires[items[0].ID]; !ok || got.IsZero() {
		t.Fatalf("no expiry: %v", r.e.expires)
	}
}

func TestRuleExpiresClips(t *testing.T) {
	r := newRig(t, Limits{})
	r.e.SetRules([]Rule{mustRule(t, `one-time`, settings.RuleExpire, 10)})
	r.copy(t, "one-time code 9")
	items := r.e.Items()
	if len(items) != 1 {
		t.Fatalf("items %+v", items)
	}
	id := items[0].ID
	if at := r.e.expires[id]; at.IsZero() {
		t.Fatal("no expiry set")
	}
	r.clock = r.clock.Add(11 * time.Second)
	r.e.Tick(context.Background())
	if items := r.e.Items(); len(items) != 0 {
		t.Fatalf("expired clip kept: %+v", items)
	}
	if _, ok := r.e.expires[id]; ok {
		t.Fatal("expiry kept for a gone clip")
	}
	stored, _ := r.s.Load(context.Background())
	if len(stored) != 0 {
		t.Fatalf("expired clip stored: %+v", stored)
	}
}

// SetRules applies to the history at once, with deadlines from each
// clip's last copy.
func TestSetRulesAppliesToHistory(t *testing.T) {
	r := newRig(t, Limits{})
	r.copy(t, "temporary thing")
	it := r.e.Items()[0]
	r.e.SetRules([]Rule{mustRule(t, `temporary`, settings.RuleExpire, 60)})
	at := r.e.expires[it.ID]
	if at.IsZero() || !at.Equal(it.LastCopiedAt.Add(60*time.Second)) {
		t.Fatalf("expiry %v, last copied %v", at, it.LastCopiedAt)
	}
	// Dropping the rule drops the deadline.
	r.e.SetRules(nil)
	if _, ok := r.e.expires[it.ID]; ok {
		t.Fatal("deadline kept without the rule")
	}
	// A pin rule pins the stored clip.
	r.e.SetRules([]Rule{mustRule(t, `temporary`, settings.RulePin, 0)})
	if items := r.e.Items(); items[0].Pin == "" {
		t.Fatalf("not pinned: %+v", items[0])
	}
}

func TestRuleMatchesTitleAndOCR(t *testing.T) {
	r := newRig(t, Limits{})
	r.e.SetRules([]Rule{mustRule(t, `invoice`, settings.RuleExpire, 5)})
	r.copy(t, "receipt scan")
	it := r.e.Items()[0]
	if _, ok := r.e.expires[it.ID]; ok {
		t.Fatal("matched without the word")
	}
	r.e.SetTitle(it.ID, "invoice march")
	r.e.SetRules([]Rule{mustRule(t, `invoice`, settings.RuleExpire, 5)})
	if _, ok := r.e.expires[it.ID]; !ok {
		t.Fatal("title not matched")
	}
}

func TestRulesFromSettingsDropsBadPatterns(t *testing.T) {
	in := []settings.Rule{
		{Pattern: "(", Action: settings.RulePin},
		{Pattern: "ok", Action: "nope"},
		{Pattern: "good", Action: settings.RuleExpire, Seconds: 5},
	}
	got := RulesFromSettings(in)
	if len(got) != 1 || got[0].Action != settings.RuleExpire || got[0].Seconds != 5 {
		t.Fatalf("got %+v", got)
	}
}

// OCR arriving later can be what a rule matches: the clip is pinned or
// given its expiry then.
func TestRuleAppliesAfterOCR(t *testing.T) {
	r := newRig(t, Limits{})
	r.e.SetRules([]Rule{mustRule(t, `receipt`, settings.RuleExpire, 30)})
	r.copy(t, "png")
	it := r.e.Items()[0]
	if _, ok := r.e.expires[it.ID]; ok {
		t.Fatal("matched before the reading")
	}
	digest := ImageDigest(it)
	if err := r.e.SetRecognizedText(it.ID, digest, "Acme receipt total 4"); err != nil {
		t.Fatal(err)
	}
	at, ok := r.e.expires[it.ID]
	if !ok {
		t.Fatal("expiry not set after OCR")
	}
	if !at.Equal(it.LastCopiedAt.Add(30 * time.Second)) {
		t.Fatalf("expiry %v", at)
	}
}
