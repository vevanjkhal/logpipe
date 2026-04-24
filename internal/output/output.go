package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// LogEntry represents a structured log line with metadata.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
	Level     string    `json:"level,omitempty"`
	Message   string    `json:"message"`
}

// Writer defines the interface for log output destinations.
type Writer interface {
	Write(entry LogEntry) error
	Close() error
}

// TextWriter writes log entries as plain text lines.
type TextWriter struct {
	w io.Writer
}

// NewTextWriter creates a TextWriter writing to w.
func NewTextWriter(w io.Writer) *TextWriter {
	return &TextWriter{w: w}
}

// NewStdoutWriter creates a TextWriter that writes to os.Stdout.
func NewStdoutWriter() *TextWriter {
	return NewTextWriter(os.Stdout)
}

// Write formats the entry as plain text and writes it.
func (t *TextWriter) Write(entry LogEntry) error {
	_, err := fmt.Fprintf(t.w, "%s [%s] %s\n",
		entry.Timestamp.Format(time.RFC3339),
		entry.Source,
		entry.Message,
	)
	return err
}

// Close is a no-op for TextWriter.
func (t *TextWriter) Close() error { return nil }

// JSONWriter writes log entries as newline-delimited JSON.
type JSONWriter struct {
	enc *json.Encoder
	w   io.WriteCloser
}

// NewJSONWriter creates a JSONWriter writing to w.
func NewJSONWriter(w io.WriteCloser) *JSONWriter {
	return &JSONWriter{enc: json.NewEncoder(w), w: w}
}

// Write serialises the entry as JSON.
func (j *JSONWriter) Write(entry LogEntry) error {
	return j.enc.Encode(entry)
}

// Close closes the underlying writer.
func (j *JSONWriter) Close() error { return j.w.Close() }
