package transform_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/transform"
)

func TestBufferTransformer_PassesLineThrough(t *testing.T) {
	tx := transform.NewBuffer(10)
	out, ok := tx.Transform("hello world")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if out != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", out)
	}
}

func TestBufferTransformer_RecordsLines(t *testing.T) {
	tx := transform.NewBuffer(5)
	lines := []string{"a", "b", "c"}
	for _, l := range lines {
		tx.Transform(l)
	}
	snap := tx.Buffer().Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(snap))
	}
	for i, want := range lines {
		if snap[i] != want {
			t.Fatalf("pos %d: want %q got %q", i, want, snap[i])
		}
	}
}

func TestBufferTransformer_RespectsCapacity(t *testing.T) {
	tx := transform.NewBuffer(2)
	tx.Transform("first")
	tx.Transform("second")
	tx.Transform("third")
	snap := tx.Buffer().Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(snap))
	}
	if snap[0] != "second" || snap[1] != "third" {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
}

func TestBufferTransformer_EmptySnapshot(t *testing.T) {
	tx := transform.NewBuffer(4)
	if snap := tx.Buffer().Snapshot(); snap != nil {
		t.Fatalf("expected nil snapshot, got %v", snap)
	}
}
