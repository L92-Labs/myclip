// Smart rules: regular expressions over a clip's text that pin it or
// expire it, applied as clips are captured and when the rules change.
package engine

import (
	"context"
	"regexp"
	"time"

	"github.com/L92-Labs/myclip/internal/history"
	"github.com/L92-Labs/myclip/internal/settings"
)

// Rule is a compiled smart rule. The first rule whose pattern matches a
// clip's text wins.
type Rule struct {
	Pattern *regexp.Regexp
	Action  string // settings.RulePin or settings.RuleExpire
	Seconds int    // the expire delay
}

// RulesFromSettings compiles the settings' rules, keeping only the ones
// that compile (Normalize has usually dropped the rest already).
func RulesFromSettings(in []settings.Rule) []Rule {
	var out []Rule
	for _, r := range in {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			continue
		}
		action := r.Action
		if action != settings.RulePin && action != settings.RuleExpire {
			continue
		}
		out = append(out, Rule{Pattern: re, Action: action, Seconds: r.Seconds})
	}
	return out
}

// SetRules installs rules and applies them to the history at once: an
// expire rule sets each match's deadline from its last copy, a pin rule
// pins unpinned matches while a key is free. reserved are the pin keys
// the in-popup shortcuts use.
func (e *Engine) SetRules(rules []Rule, reserved ...string) {
	e.mu.Lock()
	e.rules, e.pinReserved = rules, reserved
	e.mu.Unlock()
	e.applyRules()
}

// applyRules recomputes expiries and pins the unpinned matches.
func (e *Engine) applyRules() {
	e.mu.Lock()
	rules, reserved := e.rules, e.pinReserved
	expires := map[int64]time.Time{}
	for _, it := range e.hist.Items() {
		if r := match(rules, it); r != nil && r.Action == settings.RuleExpire {
			delay := time.Duration(max(settings.MinRuleSeconds, r.Seconds)) * time.Second
			expires[it.ID] = it.LastCopiedAt.Add(delay)
		}
	}
	e.expires = expires
	e.mu.Unlock()

	for {
		// Find the target, check its pin is still free and allocate the key
		// in one locked section: two concurrent appliers, or a capture,
		// must not take the same key (it is assigned to the history here,
		// before the persistence lock, so another goroutine cannot claim it
		// in between).
		e.mu.Lock()
		var target *history.Item
		for _, it := range e.hist.Items() {
			if it.Pin == "" {
				if r := match(rules, it); r != nil && r.Action == settings.RulePin {
					target = it
					break
				}
			}
		}
		if target == nil {
			e.mu.Unlock()
			return
		}
		free := history.AvailablePins(e.hist.Items(), reserved...)
		if len(free) == 0 {
			e.mu.Unlock()
			return
		}
		pinned := *target
		pinned.Pin = free[0]
		e.hist.Replace(target, &pinned)
		e.mu.Unlock()

		e.persist.Lock()
		err := e.cfg.Store.Save(context.Background(), &pinned)
		e.persist.Unlock()
		e.report(err)
		e.changed()
	}
}

// sweepExpired deletes the clips whose expiry has run out; the poll loop
// calls it each tick.
func (e *Engine) sweepExpired() {
	now := e.cfg.Now()
	e.mu.Lock()
	var due []int64
	for id, at := range e.expires {
		if !now.Before(at) {
			due = append(due, id)
		}
	}
	e.mu.Unlock()
	for _, id := range due {
		e.Delete(id) // already gone (ErrNoItem) counts as swept
		e.mu.Lock()
		delete(e.expires, id)
		e.mu.Unlock()
	}
}

// match returns the first rule matching the clip's text: what it carries
// (plain, RTF, HTML, file paths), its title, and the text read in an
// image.
func match(rules []Rule, it *history.Item) *Rule {
	if len(rules) == 0 {
		return nil
	}
	text := it.ContentText()
	if it.Title != "" {
		text += "\n" + it.Title
	}
	if ocr := it.RecognizedText(); ocr != "" {
		text += "\n" + ocr
	}
	for i := range rules {
		if rules[i].Pattern.MatchString(text) {
			return &rules[i]
		}
	}
	return nil
}
