package transform_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpipe/internal/transform"
)

func TestNewRename_EmptyFrom(t *testing.T) {
	_, err := transform.NewRename("", "newKey")
	if err == nil {
		t.Fatal("expected error for empty 'from' field")
	}
}

func TestNewRename_EmptyTo(t *testing.T) {
	_, err := transform.NewRename("oldKey", "")
	if err == nil {
		t.Fatal("expected error for empty 'to' field")
	}
}

func TestRenameTransformer_JSON(t *testing.T) {
	tr, err := transform.NewRename("level", "severity")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := `{"level":"info","msg":"hello"}`
	out, keep := tr.Transform(input)
	if !keep {
		t.Fatal("expected line to be kept")
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, exists := obj["level"]; exists {
		t.Error("old field 'level' should not exist in output")
	}
	if v, exists := obj["severity"]; !exists {
		t.Error("new field 'severity' should exist in output")
	} else if v != "info" {
		t.Errorf("expected severity=info, got %v", v)
	}
}

func TestRenameTransformer_FieldMissing(t *testing.T) {
	tr, err := transform.NewRename("missing", "target")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := `{"msg":"hello"}`
	out, keep := tr.Transform(input)
	if !keep {
		t.Fatal("expected line to be kept")
	}
	if out != input {
		t.Errorf("expected unchanged output, got %q", out)
	}
}

func TestRenameTransformer_PlainText(t *testing.T) {
	tr, err := transform.NewRename("level", "severity")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := "plain text log line"
	out, keep := tr.Transform(input)
	if !keep {
		t.Fatal("expected line to be kept")
	}
	if out != input {
		t.Errorf("expected unchanged output for plain text, got %q", out)
	}
}
