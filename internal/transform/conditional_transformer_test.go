package transform_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/transform"
)

// stubTransformer appends a marker to any line it receives.
type stubTransformer struct{}

func (s *stubTransformer) Transform(line string) (string, error) {
	return line + "[transformed]", nil
}

func TestConditional_JSONFieldMatches_AppliesInner(t *testing.T) {
	inner := &stubTransformer{}
	ct := transform.NewConditional("level", "error", inner)

	out, err := ct.Transform(`{"level":"error","msg":"boom"}`)
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"level":"error","msg":"boom"}[transformed]` {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestConditional_JSONFieldNoMatch_PassesThrough(t *testing.T) {
	inner := &stubTransformer{}
	ct := transform.NewConditional("level", "error", inner)

	line := `{"level":"info","msg":"ok"}`
	out, err := ct.Transform(line)
	if err != nil {
		t.Fatal(err)
	}
	if out != line {
		t.Fatalf("expected passthrough, got: %s", out)
	}
}

func TestConditional_JSONFieldMissing_PassesThrough(t *testing.T) {
	inner := &stubTransformer{}
	ct := transform.NewConditional("level", "error", inner)

	line := `{"msg":"no level field"}`
	out, err := ct.Transform(line)
	if err != nil {
		t.Fatal(err)
	}
	if out != line {
		t.Fatalf("expected passthrough, got: %s", out)
	}
}

func TestConditional_PlainText_SubstringMatch(t *testing.T) {
	inner := &stubTransformer{}
	ct := transform.NewConditional("", "ERROR", inner)

	out, err := ct.Transform("2024/01/01 ERROR something failed")
	if err != nil {
		t.Fatal(err)
	}
	if out != "2024/01/01 ERROR something failed[transformed]" {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestConditional_PlainText_NoMatch_PassesThrough(t *testing.T) {
	inner := &stubTransformer{}
	ct := transform.NewConditional("", "ERROR", inner)

	line := "2024/01/01 INFO all good"
	out, err := ct.Transform(line)
	if err != nil {
		t.Fatal(err)
	}
	if out != line {
		t.Fatalf("expected passthrough, got: %s", out)
	}
}

func TestConditional_Inverted_AppliesWhenNoMatch(t *testing.T) {
	inner := &stubTransformer{}
	ct := transform.NewConditional("level", "error", inner, transform.WithInvert())

	// field is "info" — does NOT match "error" → inverted → inner applied
	out, err := ct.Transform(`{"level":"info","msg":"ok"}`)
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"level":"info","msg":"ok"}[transformed]` {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestConditional_Inverted_SkipsWhenMatch(t *testing.T) {
	inner := &stubTransformer{}
	ct := transform.NewConditional("level", "error", inner, transform.WithInvert())

	line := `{"level":"error","msg":"boom"}`
	out, err := ct.Transform(line)
	if err != nil {
		t.Fatal(err)
	}
	if out != line {
		t.Fatalf("expected passthrough, got: %s", out)
	}
}
