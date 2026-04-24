package pipeline_test

import (
	"context"
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/output"
	"github.com/yourorg/logpipe/internal/pipeline"
	"github.com/yourorg/logpipe/internal/source"
)

func TestPipeline_RunNoFilter(t *testing.T) {
	src := source.NewReaderSource("test", strings.NewReader("line1\nline2\nline3\n"))
	var buf strings.Builder
	w := output.NewTextWriter(&buf)
	p := pipeline.New(src, nil, w)

	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	p.Close()

	got := buf.String()
	for _, want := range []string{"line1", "line2", "line3"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected output to contain %q, got: %s", want, got)
		}
	}
}

func TestPipeline_RunWithFilter(t *testing.T) {
	src := source.NewReaderSource("test", strings.NewReader(
		`{"level":"error","msg":"boom"}`+"\n"+
			`{"level":"info","msg":"ok"}`+"\n",
	))
	var buf strings.Builder
	w := output.NewTextWriter(&buf)

	f := filter.NewBuilder().Level("error").Build()
	p := pipeline.New(src, f, w)

	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	p.Close()

	got := buf.String()
	if !strings.Contains(got, "boom") {
		t.Errorf("expected 'boom' in output, got: %s", got)
	}
	if strings.Contains(got, "ok") {
		t.Errorf("unexpected 'ok' in output, got: %s", got)
	}
}

func TestPipeline_ContextCancel(t *testing.T) {
	pr, _ := strings.NewReader(""), (*strings.Reader)(nil)
	_ = pr
	src := source.NewReaderSource("test", strings.NewReader(""))
	var buf strings.Builder
	w := output.NewTextWriter(&buf)
	p := pipeline.New(src, nil, w)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := p.Run(ctx); err != nil {
		t.Fatalf("unexpected error on cancelled context: %v", err)
	}
}
