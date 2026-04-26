package transform_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/transform"
)

func TestNewRegexFilter_InvalidPattern(t *testing.T) {
	_, err := transform.NewRegexFilter("[", false)
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestRegexFilter_MatchingLine(t *testing.T) {
	f, err := transform.NewRegexFilter(`error`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, err := f.Transform("level=error msg=something")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected line to pass through, got empty string")
	}
}

func TestRegexFilter_NonMatchingLine(t *testing.T) {
	f, err := transform.NewRegexFilter(`error`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, err := f.Transform("level=info msg=ok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected line to be suppressed, got %q", out)
	}
}

func TestRegexFilter_Inverted_MatchSuppressed(t *testing.T) {
	f, err := transform.NewRegexFilter(`debug`, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// line matches pattern — with invert=true it should be suppressed
	out, err := f.Transform("level=debug msg=verbose")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected line to be suppressed, got %q", out)
	}
}

func TestRegexFilter_Inverted_NonMatchPasses(t *testing.T) {
	f, err := transform.NewRegexFilter(`debug`, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// line does not match pattern — with invert=true it should pass
	out, err := f.Transform("level=info msg=hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected line to pass through, got empty string")
	}
}

func TestRegexFilter_EmptyLine(t *testing.T) {
	f, err := transform.NewRegexFilter(`error`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, err := f.Transform("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty line to be suppressed, got %q", out)
	}
}
