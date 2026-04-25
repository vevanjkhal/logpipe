package sampling_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/sampling"
)

func TestNew_InvalidRate(t *testing.T) {
	cases := []float64{0, -0.1, 1.1, 2.0}
	for _, r := range cases {
		_, err := sampling.New(r)
		if err == nil {
			t.Errorf("expected error for rate %v, got nil", r)
		}
	}
}

func TestNew_ValidRate(t *testing.T) {
	s, err := sampling.New(0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestAllow_RateOne_PassesAll(t *testing.T) {
	s, _ := sampling.New(1.0, sampling.WithSeed(42))
	for i := 0; i < 100; i++ {
		if !s.Allow("line") {
			t.Fatal("rate=1.0 should pass all lines")
		}
	}
	passed, dropped := s.Snapshot()
	if passed != 100 {
		t.Errorf("expected passed=100, got %d", passed)
	}
	if dropped != 0 {
		t.Errorf("expected dropped=0, got %d", dropped)
	}
}

func TestAllow_RateZeroPoint5_ApproximatelyHalf(t *testing.T) {
	s, _ := sampling.New(0.5, sampling.WithSeed(1234))
	const total = 10000
	var passed int
	for i := 0; i < total; i++ {
		if s.Allow("some log line") {
			passed++
		}
	}
	ratio := float64(passed) / float64(total)
	// Allow 10% tolerance around 0.5
	if ratio < 0.40 || ratio > 0.60 {
		t.Errorf("expected ratio near 0.5, got %.3f", ratio)
	}
}

func TestSnapshot_CountsMatch(t *testing.T) {
	s, _ := sampling.New(1.0, sampling.WithSeed(0))
	s.Allow("a")
	s.Allow("b")
	s.Allow("c")
	passed, dropped := s.Snapshot()
	if passed != 3 {
		t.Errorf("expected passed=3, got %d", passed)
	}
	if dropped != 0 {
		t.Errorf("expected dropped=0, got %d", dropped)
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	s, _ := sampling.New(1.0, sampling.WithSeed(0))
	s.Allow("a")
	s.Allow("b")
	s.Reset()
	passed, dropped := s.Snapshot()
	if passed != 0 || dropped != 0 {
		t.Errorf("expected zeroed counters after Reset, got passed=%d dropped=%d", passed, dropped)
	}
}
