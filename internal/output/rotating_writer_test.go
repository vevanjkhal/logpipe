package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRotatingWriter_Write(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	rw, err := NewRotatingWriter(path, 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rw.Close()

	data := []byte("hello log\n")
	n, err := rw.Write(data)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}
	if n != len(data) {
		t.Errorf("expected %d bytes written, got %d", len(data), n)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if string(content) != string(data) {
		t.Errorf("expected %q, got %q", data, content)
	}
}

func TestRotatingWriter_Rotates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	// maxBytes set to 20 to force rotation quickly
	rw, err := NewRotatingWriter(path, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rw.Close()

	for i := 0; i < 5; i++ {
		_, err := rw.Write([]byte("0123456789\n")) // 11 bytes each
		if err != nil {
			t.Fatalf("write %d error: %v", i, err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("readdir error: %v", err)
	}
	if len(entries) < 2 {
		t.Errorf("expected at least 2 files after rotation, got %d", len(entries))
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".log") {
			t.Errorf("unexpected file extension: %s", e.Name())
		}
	}
}

func TestRotatingWriter_InvalidPath(t *testing.T) {
	_, err := NewRotatingWriter("/nonexistent/dir/test.log", 1024)
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestRotatingWriter_Path(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	rw, err := NewRotatingWriter(path, 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rw.Close()
	if rw.Path() != path {
		t.Errorf("expected path %q, got %q", path, rw.Path())
	}
}
