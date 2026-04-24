package pipeline_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/output"
	"github.com/yourorg/logpipe/internal/pipeline"
	"github.com/yourorg/logpipe/internal/source"
)

func TestBuilder_DefaultsToStdout(t *testing.T) {
	src := source.NewReaderSource("x", strings.NewReader(""))
	p := pipeline.NewBuilder().Source(src).Build()
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
	p.Close()
}

func TestBuilder_WithJSONWriter(t *testing.T) {
	var buf strings.Builder
	src := source.NewReaderSource("x", strings.NewReader(`{"level":"info","msg":"hello"}`+"\n"))
	p := pipeline.NewBuilder().
		Source(src).
		JSON(&buf).
		Build()

	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	p.Close()

	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected JSON output to contain 'hello', got: %s", buf.String())
	}
}

func TestBuilder_WithKeyword(t *testing.T) {
	var buf strings.Builder
	src := source.NewReaderSource("x", strings.NewReader("keep this\nskip that\n"))
	w := output.NewTextWriter(&buf)
	p := pipeline.NewBuilder().
		Source(src).
		Keyword("keep").
		OutputWriter(w).
		Build()

	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	p.Close()

	got := buf.String()
	if !strings.Contains(got, "keep this") {
		t.Errorf("expected 'keep this' in output, got: %s", got)
	}
	if strings.Contains(got, "skip that") {
		t.Errorf("unexpected 'skip that' in output, got: %s", got)
	}
}

func TestBuilder_PanicsWithoutSource(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when source is not set")
		}
	}()
	pipeline.NewBuilder().Build()
}
