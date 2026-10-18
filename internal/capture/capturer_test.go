package capture

import (
	"strings"
	"testing"

	"github.com/L92-Labs/myclip/internal/pasteboard"
)

func text(s string) pasteboard.Item {
	return pasteboard.Item{Types: map[string][]byte{pasteboard.TypeString: []byte(s)}}
}

func concealed(s string) pasteboard.Item {
	return pasteboard.Item{Types: map[string][]byte{
		pasteboard.TypeString:    []byte(s),
		pasteboard.TypeConcealed: nil,
	}}
}

func newCap(t *testing.T) (*Capturer, *pasteboard.Fake) {
	t.Helper()
	f := pasteboard.NewFake()
	f.SetItems([]pasteboard.Item{text("before launch")}, nil)
	return NewCapturer(f, Prefs{}), f
}

func TestPreLaunchContentNotCaptured(t *testing.T) {
	c, f := newCap(t)
	if _, o := c.Check(); o != Unchanged {
		t.Fatalf("outcome %v", o)
	}
	if f.DataReads() != 0 {
		t.Fatalf("%d data reads", f.DataReads())
	}
}

func TestCapturesText(t *testing.T) {
	c, f := newCap(t)
	c.Frontmost = func() string { return "com.apple.TextEdit" }
	f.SetItems([]pasteboard.Item{text("hello")}, nil)
	r, o := c.Check()
	if o != Captured || len(r.Contents) != 1 || string(r.Contents[0].Value) != "hello" ||
		r.Application != "com.apple.TextEdit" || r.OwnWrite {
		t.Fatalf("%v %+v", o, r)
	}
	if _, o := c.Check(); o != Unchanged {
		t.Fatalf("second check %v", o)
	}
}

// N6 named test: the concealed marker on a later item still ignores the
// whole pasteboard before any data read.
func TestConcealedOnSecondItem(t *testing.T) {
	c, f := newCap(t)
	f.SetItems([]pasteboard.Item{text("plain"), concealed("hunter2")}, nil)
	f.ResetReads()
	r, o := c.Check()
	if o != Ignored || r != nil {
		t.Fatalf("%v %+v", o, r)
	}
	if f.DataReads() != 0 {
		t.Fatalf("%d data reads, want 0", f.DataReads())
	}
}

// N6 named test, c1 path: the pasteboard changes after the types-only
// snapshot ⇒ exit before any data read, change left for the next tick.
func TestConcealedMidReadC1(t *testing.T) {
	c, f := newCap(t)
	f.SetItems([]pasteboard.Item{text("plain")}, nil)
	f.AfterSnapshot = func() {
		f.AfterSnapshot = nil
		f.SetItems([]pasteboard.Item{concealed("hunter2")}, nil)
	}
	f.ResetReads()
	if _, o := c.Check(); o != Raced {
		t.Fatalf("outcome %v, want raced", o)
	}
	if f.DataReads() != 0 {
		t.Fatalf("%d data reads on the c1 path, want 0", f.DataReads())
	}
	// Next tick sees the concealed change and ignores it without reads.
	if _, o := c.Check(); o != Ignored {
		t.Fatalf("retry outcome %v, want ignored", o)
	}
	if f.DataReads() != 0 {
		t.Fatalf("%d data reads after retry", f.DataReads())
	}
}

// N6 named test, c2 path: the change lands during the read ⇒ the bytes read
// are discarded (no result), and the next tick ignores the concealed item.
func TestConcealedMidReadC2(t *testing.T) {
	c, f := newCap(t)
	f.SetItems([]pasteboard.Item{text("plain")}, nil)
	f.BeforeRead = func() {
		f.BeforeRead = nil
		f.SetItems([]pasteboard.Item{concealed("hunter2")}, nil)
	}
	r, o := c.Check()
	if o != Raced || r != nil {
		t.Fatalf("%v %+v, want raced with nothing kept", o, r)
	}
	f.ResetReads()
	if _, o := c.Check(); o != Ignored || f.DataReads() != 0 {
		t.Fatalf("retry %v with %d reads", o, f.DataReads())
	}
}

func TestIgnoreNextSkipsExactlyOne(t *testing.T) {
	c, f := newCap(t)
	c.UpdatePrefs(func(p *Prefs) { p.IgnoreEvents, p.IgnoreOnlyNextEvent = true, true })
	f.SetItems([]pasteboard.Item{text("skip me")}, nil)
	f.ResetReads()
	if _, o := c.Check(); o != Ignored || f.DataReads() != 0 {
		t.Fatalf("first %v, %d reads", o, f.DataReads())
	}
	if p := c.Prefs(); p.IgnoreEvents || p.IgnoreOnlyNextEvent {
		t.Fatalf("flags not cleared: %+v", p)
	}
	f.SetItems([]pasteboard.Item{text("keep me")}, nil)
	if r, o := c.Check(); o != Captured || string(r.Contents[0].Value) != "keep me" {
		t.Fatalf("second %v", o)
	}
}

func TestPausedIgnoresUntilResumed(t *testing.T) {
	c, f := newCap(t)
	c.UpdatePrefs(func(p *Prefs) { p.IgnoreEvents = true })
	for i := range 3 {
		f.SetItems([]pasteboard.Item{text(strings.Repeat("x", i+1))}, nil)
		if _, o := c.Check(); o != Ignored {
			t.Fatalf("paused check %d: %v", i, o)
		}
	}
	if f.DataReads() != 0 {
		t.Fatalf("%d reads while paused", f.DataReads())
	}
	c.UpdatePrefs(func(p *Prefs) { p.IgnoreEvents = false })
	f.SetItems([]pasteboard.Item{text("after")}, nil)
	if _, o := c.Check(); o != Captured {
		t.Fatalf("resumed %v", o)
	}
}

func TestIgnoredAppReadsNothing(t *testing.T) {
	c, f := newCap(t)
	c.Frontmost = func() string { return "com.1password.1password" }
	c.UpdatePrefs(func(p *Prefs) { p.IgnoredApps = []string{"com.1password.1password"} })
	f.SetItems([]pasteboard.Item{text("secret")}, nil)
	f.ResetReads()
	if _, o := c.Check(); o != Ignored || f.DataReads() != 0 {
		t.Fatalf("%v with %d reads", o, f.DataReads())
	}
}

func TestRegexpDropsMatchingItemOnly(t *testing.T) {
	c, f := newCap(t)
	c.UpdatePrefs(func(p *Prefs) { p.IgnoreRegexp = []string{`^\d{4} \d{4} \d{4} \d{4}$`} })
	f.SetItems([]pasteboard.Item{text("4111 1111 1111 1111"), text("keep")}, nil)
	r, o := c.Check()
	if o != Captured || len(r.Contents) != 1 || string(r.Contents[0].Value) != "keep" {
		t.Fatalf("%v %+v", o, r)
	}
	f.SetItems([]pasteboard.Item{text("4111 1111 1111 1111")}, nil)
	if _, o := c.Check(); o != Empty {
		t.Fatalf("all-matching outcome %v", o)
	}
}

func TestOversizedSkipped(t *testing.T) {
	c, f := newCap(t)
	c.MaxItemBytes = 8
	f.SetItems([]pasteboard.Item{text("way more than eight bytes")}, nil)
	r, o := c.Check()
	if o != Oversized || len(r.Contents) != 0 || r.Size == 0 {
		t.Fatalf("%v %+v", o, r)
	}
}

func TestOwnWriteFlagged(t *testing.T) {
	c, f := newCap(t)
	it := text("copied back")
	it.Types[pasteboard.TypeOwnWrite] = []byte{}
	f.SetItems([]pasteboard.Item{it}, nil)
	if r, o := c.Check(); o != Captured || !r.OwnWrite {
		t.Fatalf("%v %+v", o, r)
	}
}

func TestClampInterval(t *testing.T) {
	if ClampInterval(0) != DefaultInterval || ClampInterval(MinInterval/2) != MinInterval {
		t.Fatal("D7 floor not applied")
	}
}

// The capture-safe latch survives unpausing (tray resume, ⌥-click).
func TestCaptureSafeSurvivesResume(t *testing.T) {
	c, f := newCap(t)
	c.SetCaptureSafe(true)
	c.UpdatePrefs(func(p *Prefs) { p.IgnoreEvents = false })
	f.SetItems([]pasteboard.Item{text("secret")}, nil)
	f.ResetReads()
	if _, o := c.Check(); o != Ignored || f.DataReads() != 0 {
		t.Fatalf("%v with %d reads", o, f.DataReads())
	}
	// A copy made after the last safe-mode check but before the reset is
	// never read.
	f.SetItems([]pasteboard.Item{text("pending secret")}, nil)
	f.ResetReads()
	c.SetCaptureSafe(false)
	if _, o := c.Check(); o != Unchanged || f.DataReads() != 0 {
		t.Fatalf("pending copy: %v with %d reads", o, f.DataReads())
	}
	f.SetItems([]pasteboard.Item{text("after reset")}, nil)
	if _, o := c.Check(); o != Captured {
		t.Fatalf("after reset %v", o)
	}
}
