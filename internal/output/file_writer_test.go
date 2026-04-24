package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileWriter_WritePlainText(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	fw, err := NewFileWriter(path, false)
	if err != nil {
		t.Fatalf("NewFileWriter: %v", err)
	}
	defer fw.Close()

	if err := fw.Write("myapp", "hello world"); err != nil {
		t.Fatalf("Write: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "myapp") {
		t.Errorf("expected source %q in output, got: %s", "myapp", content)
	}
	if !strings.Contains(content, "hello world") {
		t.Errorf("expected line %q in output, got: %s", "hello world", content)
	}
}

func TestFileWriter_WriteJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")

	fw, err := NewFileWriter(path, true)
	if err != nil {
		t.Fatalf("NewFileWriter: %v", err)
	}
	defer fw.Close()

	if err := fw.Write("svc", "structured entry"); err != nil {
		t.Fatalf("Write: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, `"source"`) {
		t.Errorf("expected JSON key 'source' in output, got: %s", content)
	}
	if !strings.Contains(content, "structured entry") {
		t.Errorf("expected line in JSON output, got: %s", content)
	}
}

func TestFileWriter_InvalidPath(t *testing.T) {
	_, err := NewFileWriter("/nonexistent/dir/log.txt", false)
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestFileWriter_Path(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	fw, err := NewFileWriter(path, false)
	if err != nil {
		t.Fatalf("NewFileWriter: %v", err)
	}
	defer fw.Close()

	if fw.Path() != path {
		t.Errorf("Path() = %q, want %q", fw.Path(), path)
	}
}
