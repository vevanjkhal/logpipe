package output

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// RotatingWriter writes log entries to a file, rotating it when it exceeds maxBytes.
type RotatingWriter struct {
	mu       sync.Mutex
	path     string
	maxBytes int64
	file     *os.File
	written  int64
}

// NewRotatingWriter creates a new RotatingWriter. maxBytes is the max file size before rotation.
func NewRotatingWriter(path string, maxBytes int64) (*RotatingWriter, error) {
	if maxBytes <= 0 {
		maxBytes = 10 * 1024 * 1024 // 10 MB default
	}
	rw := &RotatingWriter{path: path, maxBytes: maxBytes}
	if err := rw.openFile(); err != nil {
		return nil, err
	}
	return rw, nil
}

func (rw *RotatingWriter) openFile() error {
	f, err := os.OpenFile(rw.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("rotating_writer: open %s: %w", rw.path, err)
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return fmt.Errorf("rotating_writer: stat %s: %w", rw.path, err)
	}
	rw.file = f
	rw.written = info.Size()
	return nil
}

func (rw *RotatingWriter) rotate() error {
	if err := rw.file.Close(); err != nil {
		return err
	}
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	ext := filepath.Ext(rw.path)
	base := rw.path[:len(rw.path)-len(ext)]
	rotated := fmt.Sprintf("%s.%s%s", base, timestamp, ext)
	if err := os.Rename(rw.path, rotated); err != nil {
		return fmt.Errorf("rotating_writer: rename: %w", err)
	}
	return rw.openFile()
}

// Write writes data to the current log file, rotating if necessary.
func (rw *RotatingWriter) Write(p []byte) (int, error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	if rw.written+int64(len(p)) > rw.maxBytes {
		if err := rw.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := rw.file.Write(p)
	rw.written += int64(n)
	return n, err
}

// Close closes the underlying file.
func (rw *RotatingWriter) Close() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	if rw.file != nil {
		return rw.file.Close()
	}
	return nil
}

// Path returns the current log file path.
func (rw *RotatingWriter) Path() string {
	return rw.path
}
