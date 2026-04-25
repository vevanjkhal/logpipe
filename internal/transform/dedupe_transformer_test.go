package transform_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/transform"
)

func TestDedupeTransformer_PassesUniqueLines(t *testing.T) {
	tr := transform.NewDedupe(10 * time.Second)
	got := tr.Transform("unique line")
	if got != "unique line" {
		t.Fatalf("want %q, got %q", "unique line", got)
	}
}

func TestDedupeTransformer_SuppressesDuplicate(t *testing.T) {
	tr := transform.NewDedupe(10 * time.Second)
	tr.Transform("dup")
	got := tr.Transform("dup")
	if got != "" {
		t.Fatalf("expected empty string for duplicate, got %q", got)
	}
}

func TestDedupeTransformer_AllowsDifferentLine(t *testing.T) {
	tr := transform.NewDedupe(10 * time.Second)
	tr.Transform("a")
	got := tr.Transform("b")
	if got != "b" {
		t.Fatalf("want %q, got %q", "b", got)
	}
}

func TestDedupeTransformer_ZeroWindowUsesDefault(t *testing.T) {
	// Should not panic and should function.
	tr := transform.NewDedupe(0)
	if tr.Transform("x") != "x" {
		t.Fatal("first line should pass through")
	}
	if tr.Transform("x") != "" {
		t.Fatal("duplicate should be suppressed")
	}
}
