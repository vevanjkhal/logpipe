package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/ratelimit"
)

func TestNew_ZeroRateDefaultsToOne(t *testing.T) {
	l := ratelimit.New(0)
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
	// Should allow at least one line immediately.
	if !l.Allow() {
		t.Error("expected first Allow() to return true")
	}
}

func TestAllow_BurstUpToRate(t *testing.T) {
	const rate = 5.0
	l := ratelimit.New(rate)

	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}
	// Burst capacity equals rate, so exactly 5 should be allowed.
	if allowed != rate {
		t.Errorf("expected %d allowed, got %d", int(rate), allowed)
	}
}

func TestDropped_CountsThrottledLines(t *testing.T) {
	l := ratelimit.New(2)

	for i := 0; i < 5; i++ {
		l.Allow()
	}

	if l.Dropped() != 3 {
		t.Errorf("expected 3 dropped, got %d", l.Dropped())
	}
}

func TestReset_RefillsTokensAndClearsDropped(t *testing.T) {
	l := ratelimit.New(3)

	// Exhaust the bucket.
	for i := 0; i < 6; i++ {
		l.Allow()
	}
	if l.Dropped() == 0 {
		t.Fatal("expected some dropped lines before reset")
	}

	l.Reset()

	if l.Dropped() != 0 {
		t.Errorf("expected 0 dropped after reset, got %d", l.Dropped())
	}
	// After reset the bucket should be full again.
	if !l.Allow() {
		t.Error("expected Allow() to return true after reset")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	l := ratelimit.New(1000) // 1000 lines/sec

	// Drain the bucket.
	for i := 0; i < 1000; i++ {
		l.Allow()
	}

	// Wait long enough for at least one token to refill.
	time.Sleep(5 * time.Millisecond)

	if !l.Allow() {
		t.Error("expected Allow() to return true after waiting for token refill")
	}
}
