package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/pipeline"
	"github.com/yourorg/logpipe/internal/transform"
)

func TestBuilder_WithConditional_AppliesInner(t *testing.T) {
	var buf bytes.Buffer

	p := pipeline.NewBuilder().
		WithSource(strings.NewReader(`{"level":"error","msg":"boom"}`+"\n"+`{"level":"info","msg":"ok"}`+"\n")).
		WithConditional("level", "error", transform.NewAddField("flagged", "true")).
		WithWriter(&buf).
		Build()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := p.Run(ctx); err != nil && err != context.Canceled {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "flagged") {
		t.Errorf("expected error line to contain 'flagged', got:\n%s", output)
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}
	if strings.Contains(lines[1], "flagged") {
		t.Errorf("info line should NOT be flagged, got: %s", lines[1])
	}
}

func TestBuilder_WithConditionalInverted(t *testing.T) {
	var buf bytes.Buffer

	p := pipeline.NewBuilder().
		WithSource(strings.NewReader(`{"level":"info","msg":"ok"}`+"\n"+`{"level":"error","msg":"bad"}`+"\n")).
		WithConditionalInverted("level", "error", transform.NewAddField("non_error", "true")).
		WithWriter(&buf).
		Build()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := p.Run(ctx); err != nil && err != context.Canceled {
		t.Fatal(err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "non_error") {
		t.Errorf("info line should have non_error field, got: %s", lines[0])
	}
	if strings.Contains(lines[1], "non_error") {
		t.Errorf("error line should NOT have non_error field, got: %s", lines[1])
	}
}
