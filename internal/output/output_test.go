package output

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"
)

var testEntry = LogEntry{
	Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	Source:    "app",
	Level:     "info",
	Message:   "server started",
}

func TestTextWriter_Write(t *testing.T) {
	var buf bytes.Buffer
	w := NewTextWriter(&buf)

	if err := w.Write(testEntry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "[app]") {
		t.Errorf("expected source in output, got: %q", got)
	}
	if !strings.Contains(got, "server started") {
		t.Errorf("expected message in output, got: %q", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("expected newline at end, got: %q", got)
	}
}

func TestTextWriter_Close(t *testing.T) {
	w := NewTextWriter(io.Discard)
	if err := w.Close(); err != nil {
		t.Errorf("Close() should return nil, got: %v", err)
	}
}

func TestJSONWriter_Write(t *testing.T) {
	pr, pw := io.Pipe()
	jw := NewJSONWriter(pw)

	done := make(chan LogEntry, 1)
	go func() {
		var entry LogEntry
		if err := json.NewDecoder(pr).Decode(&entry); err == nil {
			done <- entry
		}
	}()

	if err := jw.Write(testEntry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := <-done
	if got.Source != testEntry.Source {
		t.Errorf("expected source %q, got %q", testEntry.Source, got.Source)
	}
	if got.Message != testEntry.Message {
		t.Errorf("expected message %q, got %q", testEntry.Message, got.Message)
	}

	_ = jw.Close()
}
