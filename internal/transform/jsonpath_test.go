package transform_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/transform"
)

func TestNewJSONPath_EmptyField(t *testing.T) {
	_, err := transform.NewJSONPath("", false)
	if err == nil {
		t.Fatal("expected error for empty field, got nil")
	}
}

func TestJSONPathTransformer_ExtractsStringField(t *testing.T) {
	tr, _ := transform.NewJSONPath("msg", false)
	got := tr.Transform(`{"level":"info","msg":"hello world"}`)
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestJSONPathTransformer_ExtractsNumericField(t *testing.T) {
	tr, _ := transform.NewJSONPath("code", false)
	got := tr.Transform(`{"code":404}`)
	if got != "404" {
		t.Errorf("expected '404', got %q", got)
	}
}

func TestJSONPathTransformer_MissingField_NoFallback(t *testing.T) {
	tr, _ := transform.NewJSONPath("missing", false)
	got := tr.Transform(`{"level":"info"}`)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestJSONPathTransformer_MissingField_WithFallback(t *testing.T) {
	line := `{"level":"info"}`
	tr, _ := transform.NewJSONPath("missing", true)
	got := tr.Transform(line)
	if got != line {
		t.Errorf("expected original line, got %q", got)
	}
}

func TestJSONPathTransformer_NonJSON_NoFallback(t *testing.T) {
	tr, _ := transform.NewJSONPath("msg", false)
	got := tr.Transform("plain text line")
	if got != "" {
		t.Errorf("expected empty string for non-JSON, got %q", got)
	}
}

func TestJSONPathTransformer_NonJSON_WithFallback(t *testing.T) {
	tr, _ := transform.NewJSONPath("msg", true)
	got := tr.Transform("plain text line")
	if got != "plain text line" {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestJSONPathTransformer_NullValue_WithFallback(t *testing.T) {
	line := `{"msg":null}`
	tr, _ := transform.NewJSONPath("msg", true)
	got := tr.Transform(line)
	if got != line {
		t.Errorf("expected original line for null value with fallback, got %q", got)
	}
}
