package transform_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/transform"
)

func TestNewThrottle_NegativeRate(t *testing.T) {
	_, err := transform.NewThrottle(-1)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestThrottleTransformer_PassesLineWhenAllowed(t *testing.T) {
	tr, err := transform.NewThrottle(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, err := tr.Transform("hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello world" {
		t.Errorf("expected line to pass through, got %q", out)
	}
}

func TestThrottleTransformer_DropsLineWhenExceeded(t *testing.T) {
	// rate=1 means only 1 token in the bucket
	tr, err := transform.NewThrottle(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// First call consumes the single token.
	_, _ = tr.Transform("first")

	// Second call should be dropped.
	out, err := tr.Transform("second")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string for dropped line, got %q", out)
	}
}

func TestThrottleTransformer_DroppedCounter(t *testing.T) {
	tr, err := transform.NewThrottle(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tr.Transform("first")  // allowed
	tr.Transform("second") // dropped
	tr.Transform("third")  // dropped

	if got := tr.Dropped(); got != 2 {
		t.Errorf("expected 2 dropped, got %d", got)
	}
}

func TestThrottleTransformer_ResetClearsDropped(t *testing.T) {
	tr, err := transform.NewThrottle(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tr.Transform("first")  // allowed
	tr.Transform("second") // dropped

	tr.Reset()

	if got := tr.Dropped(); got != 0 {
		t.Errorf("expected 0 dropped after reset, got %d", got)
	}

	// After reset the bucket should be full again.
	out, _ := tr.Transform("after-reset")
	if out != "after-reset" {
		t.Errorf("expected line to pass after reset, got %q", out)
	}
}
