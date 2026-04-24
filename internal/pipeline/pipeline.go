package pipeline

import (
	"context"
	"sync"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/output"
	"github.com/yourorg/logpipe/internal/source"
)

// Pipeline connects a source to a writer through an optional filter.
type Pipeline struct {
	src    source.Source
	filter *filter.Filter
	writer output.Writer
}

// New creates a new Pipeline.
func New(src source.Source, f *filter.Filter, w output.Writer) *Pipeline {
	return &Pipeline{
		src:    src,
		filter: f,
		writer: w,
	}
}

// Run starts the pipeline, reading lines from the source, applying the filter,
// and writing matching lines to the writer. It blocks until the context is
// cancelled or the source is exhausted.
func (p *Pipeline) Run(ctx context.Context) error {
	lines, err := p.src.Lines(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for line := range lines {
			if p.filter != nil && !p.filter.Match(line) {
				continue
			}
			_ = p.writer.Write(line)
		}
	}()

	wg.Wait()
	return nil
}

// Close closes the underlying writer.
func (p *Pipeline) Close() error {
	return p.writer.Close()
}
