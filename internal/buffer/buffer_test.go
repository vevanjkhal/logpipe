package buffer_test

import (
	"fmt"
	"testing"

	"github.com/logpipe/logpipe/internal/buffer"
)

func TestNew_CapBelowOne_DefaultsToOne(t *testing.T) {
	b := buffer.New(0)
	if b.Cap() != 1 {
		t.Fatalf("expected cap 1, got %d", b.Cap())
	}
}

func TestWrite_SingleLine(t *testing.T) {
	b := buffer.New(5)
	b.Write("hello")
	snap := b.Snapshot()
	if len(snap) != 1 || snap[0] != "hello" {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
}

func TestSnapshot_OrderedOldestFirst(t *testing.T) {
	b := buffer.New(3)
	for i := 1; i <= 3; i++ {
		b.Write(fmt.Sprintf("line%d", i))
	}
	snap := b.Snapshot()
	expected := []string{"line1", "line2", "line3"}
	for i, v := range expected {
		if snap[i] != v {
			t.Fatalf("pos %d: want %q got %q", i, v, snap[i])
		}
	}
}

func TestWrite_EvictsOldestWhenFull(t *testing.T) {
	b := buffer.New(3)
	for i := 1; i <= 5; i++ {
		b.Write(fmt.Sprintf("line%d", i))
	}
	snap := b.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(snap))
	}
	expected := []string{"line3", "line4", "line5"}
	for i, v := range expected {
		if snap[i] != v {
			t.Fatalf("pos %d: want %q got %q", i, v, snap[i])
		}
	}
}

func TestSnapshot_EmptyBuffer_ReturnsNil(t *testing.T) {
	b := buffer.New(4)
	if snap := b.Snapshot(); snap != nil {
		t.Fatalf("expected nil, got %v", snap)
	}
}

func TestLen_TracksCount(t *testing.T) {
	b := buffer.New(10)
	for i := 0; i < 7; i++ {
		b.Write("x")
	}
	if b.Len() != 7 {
		t.Fatalf("expected 7, got %d", b.Len())
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	b := buffer.New(5)
	b.Write("a")
	b.Write("b")
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", b.Len())
	}
	if snap := b.Snapshot(); snap != nil {
		t.Fatalf("expected nil after reset, got %v", snap)
	}
}
