package transform_test

import (
	"encoding/json"
	"testing"

	"github.com/logpipe/logpipe/internal/transform"
)

func TestRedactTransformer_JSON(t *testing.T) {
	r := transform.NewRedact([]string{"password", "token"}, "")
	input := `{"user":"alice","password":"secret","token":"abc123"}`
	out := r.Transform(input)

	var obj map[string]any
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["password"] != "***" {
		t.Errorf("expected password to be redacted, got %v", obj["password"])
	}
	if obj["token"] != "***" {
		t.Errorf("expected token to be redacted, got %v", obj["token"])
	}
	if obj["user"] != "alice" {
		t.Errorf("expected user to be unchanged, got %v", obj["user"])
	}
}

func TestRedactTransformer_NonJSON(t *testing.T) {
	r := transform.NewRedact([]string{"password"}, "")
	input := "plain text log line"
	if got := r.Transform(input); got != input {
		t.Errorf("expected non-JSON line to be unchanged, got %q", got)
	}
}

func TestAddFieldTransformer_JSON(t *testing.T) {
	a := transform.NewAddField("env", "production")
	input := `{"level":"info","msg":"started"}`
	out := a.Transform(input)

	var obj map[string]any
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["env"] != "production" {
		t.Errorf("expected env=production, got %v", obj["env"])
	}
}

func TestAddFieldTransformer_NonJSON(t *testing.T) {
	a := transform.NewAddField("env", "production")
	input := "not json"
	if got := a.Transform(input); got != input {
		t.Errorf("expected non-JSON line to be unchanged, got %q", got)
	}
}

func TestPrefixTransformer(t *testing.T) {
	p := transform.NewPrefix("[APP]")
	input := "hello world"
	expected := "[APP] hello world"
	if got := p.Transform(input); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPrefixTransformer_Empty(t *testing.T) {
	p := transform.NewPrefix("")
	input := "hello world"
	if got := p.Transform(input); got != input {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestChain(t *testing.T) {
	chain := transform.Chain{
		transform.NewAddField("env", "staging"),
		transform.NewRedact([]string{"secret"}, "REDACTED"),
	}
	input := `{"msg":"ok","secret":"xyz"}`
	out := chain.Transform(input)

	var obj map[string]any
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["env"] != "staging" {
		t.Errorf("expected env=staging, got %v", obj["env"])
	}
	if obj["secret"] != "REDACTED" {
		t.Errorf("expected secret=REDACTED, got %v", obj["secret"])
	}
}
