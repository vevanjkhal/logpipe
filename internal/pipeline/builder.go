package pipeline

import (
	"io"
	"os"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/output"
	"github.com/yourorg/logpipe/internal/source"
)

// Builder provides a fluent API for constructing a Pipeline.
type Builder struct {
	src    source.Source
	fb     *filter.Builder
	writer output.Writer
}

// NewBuilder returns a new pipeline Builder.
func NewBuilder() *Builder {
	return &Builder{fb: filter.NewBuilder()}
}

// Source sets the log source.
func (b *Builder) Source(s source.Source) *Builder {
	b.src = s
	return b
}

// Level restricts output to lines at or above the given log level.
func (b *Builder) Level(level string) *Builder {
	b.fb.Level(level)
	return b
}

// Keyword restricts output to lines containing the given keyword.
func (b *Builder) Keyword(kw string) *Builder {
	b.fb.Keyword(kw)
	return b
}

// OutputWriter sets a custom writer for pipeline output.
func (b *Builder) OutputWriter(w output.Writer) *Builder {
	b.writer = w
	return b
}

// Stdout configures the pipeline to write plain-text output to stdout.
func (b *Builder) Stdout() *Builder {
	b.writer = output.NewStdoutWriter()
	return b
}

// JSON configures the pipeline to write JSON output to the given writer.
func (b *Builder) JSON(w io.Writer) *Builder {
	b.writer = output.NewJSONWriter(w)
	return b
}

// Build assembles and returns the Pipeline. Panics if no source is set.
func (b *Builder) Build() *Pipeline {
	if b.src == nil {
		panic("pipeline: source must be set before Build()")
	}
	if b.writer == nil {
		b.writer = output.NewStdoutWriter()
	}
	_ = os.Stdout // ensure os import used
	return New(b.src, b.fb.Build(), b.writer)
}
