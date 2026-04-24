package source

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestReaderSource_Name(t *testing.T) {
	rs := NewReaderSource("stdin", strings.NewReader(""))
	if rs.Name() != "stdin" {
		t.Errorf("expected name 'stdin', got %q", rs.Name())
	}
}

func TestReaderSource_Lines(t *testing.T) {
	input := "line one\nline two\nline three"
	rs := NewReaderSource("test", strings.NewReader(input))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch, err := rs.Lines(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got []string
	for line := range ch {
		got = append(got, line)
	}

	expected := []string{"line one", "line two", "line three"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(got))
	}
	for i, l := range expected {
		if got[i] != l {
			t.Errorf("line %d: expected %q, got %q", i, l, got[i])
		}
	}
}

func TestReaderSource_ContextCancel(t *testing.T) {
	// Large input to ensure we can cancel mid-stream.
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("log line\n")
	}
	rs := NewReaderSource("cancel-test", strings.NewReader(sb.String()))

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	ch, err := rs.Lines(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Drain; should stop quickly due to cancellation.
	for range ch {
	}
}

func TestFileSource_Name(t *testing.T) {
	fs := NewFileSource("myfile", "/tmp/test.log")
	if fs.Name() != "myfile" {
		t.Errorf("expected 'myfile', got %q", fs.Name())
	}
}

func TestFileSource_Lines(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "logpipe-*.log")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = tmp.WriteString("alpha\nbeta\ngamma\n")
	tmp.Close()

	fs := NewFileSource("tmpfile", tmp.Name())
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch, err := fs.Lines(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got []string
	for line := range ch {
		got = append(got, line)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(got), got)
	}
}
