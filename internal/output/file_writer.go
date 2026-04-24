package output

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// FileWriter writes log entries to a file on disk.
type FileWriter struct {
	mu   sync.Mutex
	file *os.File
	path string
	json bool
}

// NewFileWriter creates a new FileWriter that appends to the given path.
// If json is true, entries are formatted as JSON; otherwise plain text.
func NewFileWriter(path string, json bool) (*FileWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("output: open file %q: %w", path, err)
	}
	return &FileWriter{file: f, path: path, json: json}, nil
}

// Write formats and appends a log entry to the file.
func (fw *FileWriter) Write(source, line string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	var out string
	if fw.json {
		out = formatJSON(source, line, time.Now())
	} else {
		out = formatText(source, line, time.Now())
	}

	_, err := fmt.Fprintln(fw.file, out)
	if err != nil {
		return fmt.Errorf("output: write to file %q: %w", fw.path, err)
	}
	return nil
}

// Close closes the underlying file.
func (fw *FileWriter) Close() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.file.Close()
}

// Path returns the file path this writer is writing to.
func (fw *FileWriter) Path() string {
	return fw.path
}
