package source

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestStdinSource_Name(t *testing.T) {
	s := NewStdinSource()
	if got := s.Name(); got != "stdin" {
		t.Errorf("expected name %q, got %q", "stdin", got)
	}
}

func TestStdinSource_Lines(t *testing.T) {
	input := "line one\nline two\nline three\n"
	s := newStdinSourceFromReader(strings.NewReader(input))

	ctx := context.Background()
	ch := s.Lines(ctx)

	expected := []string{"line one", "line two", "line three"}
	var got []string
	for line := range ch {
		got = append(got, line)
	}

	if len(got) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(got))
	}
	for i, want := range expected {
		if got[i] != want {
			t.Errorf("line %d: expected %q, got %q", i, want, got[i])
		}
	}
}

func TestStdinSource_ContextCancel(t *testing.T) {
	// Use a pipe so the reader blocks indefinitely until we cancel.
	pr, pw := newBlockingPipe()
	defer pw.Close()

	s := newStdinSourceFromReader(pr)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	ch := s.Lines(ctx)

	// Drain channel; it should close after context deadline.
	for range ch {
	}

	if ctx.Err() == nil {
		t.Error("expected context to be cancelled or timed out")
	}
}
