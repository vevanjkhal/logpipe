package dedupe_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/dedupe"
)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestIsDuplicate_FirstLineNotDuplicate(t *testing.T) {
	d := dedupe.New()
	if d.IsDuplicate("hello") {
		t.Fatal("first line should not be a duplicate")
	}
}

func TestIsDuplicate_ConsecutiveDuplicate(t *testing.T) {
	now := time.Now()
	d := dedupe.New(dedupe.WithWindow(10 * time.Second))
	d.(*dedupe.Deduper) // ensure concrete type via unexported field test below

	// Use package-level New and cast via interface isn't needed — test directly.
	dd := newTestDeduper(now, 10*time.Second)
	dd.IsDuplicate("line")
	if !dd.IsDuplicate("line") {
		t.Fatal("second identical line should be duplicate")
	}
}

func TestIsDuplicate_DifferentLineNotDuplicate(t *testing.T) {
	now := time.Now()
	dd := newTestDeduper(now, 10*time.Second)
	dd.IsDuplicate("line1")
	if dd.IsDuplicate("line2") {
		t.Fatal("different line should not be duplicate")
	}
}

func TestIsDuplicate_WindowExpires(t *testing.T) {
	base := time.Now()
	dd := newTestDeduper(base, 1*time.Second)
	dd.IsDuplicate("line")

	// Advance time beyond window.
	dd.SetNow(base.Add(2 * time.Second))
	if dd.IsDuplicate("line") {
		t.Fatal("duplicate after window expiry should not be suppressed")
	}
}

func TestSuppressed_CountsCorrectly(t *testing.T) {
	now := time.Now()
	dd := newTestDeduper(now, 10*time.Second)
	dd.IsDuplicate("x")
	dd.IsDuplicate("x")
	dd.IsDuplicate("x")
	if got := dd.Suppressed(); got != 2 {
		t.Fatalf("want 2 suppressed, got %d", got)
	}
}

func TestReset_ClearsState(t *testing.T) {
	now := time.Now()
	dd := newTestDeduper(now, 10*time.Second)
	dd.IsDuplicate("x")
	dd.IsDuplicate("x")
	dd.Reset()
	if dd.Suppressed() != 0 {
		t.Fatal("after reset suppressed should be 0")
	}
	if dd.IsDuplicate("x") {
		t.Fatal("after reset same line should not be duplicate")
	}
}
