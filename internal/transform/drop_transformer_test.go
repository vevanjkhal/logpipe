package transform_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/transform"
)

func TestNewDrop_EmptySlice_PassesAll(t *testing.T) {
	tr := transform.NewDrop(nil)
	got := tr.Transform("hello world")
	if got != "hello world" {
		t.Fatalf("expected line to pass through, got %q", got)
	}
}

func TestDropTransformer_MatchingSubstring_DropsLine(t *testing.T) {
	tr := transform.NewDrop([]string{"ERROR"})
	got := tr.Transform("2024/01/01 ERROR something broke")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestDropTransformer_NoMatch_PassesThrough(t *testing.T) {
	tr := transform.NewDrop([]string{"ERROR"})
	got := tr.Transform("2024/01/01 INFO all good")
	if got != "2024/01/01 INFO all good" {
		t.Fatalf("expected line to pass through, got %q", got)
	}
}

func TestDropTransformer_MultipleSubstrings_DropOnAnyMatch(t *testing.T) {
	tr := transform.NewDrop([]string{"ERROR", "WARN"})

	if got := tr.Transform("ERROR: disk full"); got != "" {
		t.Fatalf("expected drop on ERROR, got %q", got)
	}
	if got := tr.Transform("WARN: low memory"); got != "" {
		t.Fatalf("expected drop on WARN, got %q", got)
	}
	if got := tr.Transform("INFO: running"); got == "" {
		t.Fatal("expected INFO line to pass through")
	}
}

func TestDropTransformer_EmptySubstringIgnored(t *testing.T) {
	// An empty string in the list should not cause every line to be dropped.
	tr := transform.NewDrop([]string{"", "ERROR"})
	got := tr.Transform("INFO: ok")
	if got != "INFO: ok" {
		t.Fatalf("expected line to pass through, got %q", got)
	}
}

func TestDropTransformer_CaseSensitive(t *testing.T) {
	tr := transform.NewDrop([]string{"error"})
	// uppercase ERROR should not be dropped
	got := tr.Transform("ERROR: something")
	if got == "" {
		t.Fatal("expected case-sensitive match to pass through")
	}
}
