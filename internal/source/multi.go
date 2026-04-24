package source

import (
	"context"
	"sync"
)

// Entry is a log line paired with the name of its originating source.
type Entry struct {
	Source string
	Line   string
}

// MultiSource fans in lines from multiple Source implementations into a
// single channel of Entry values.
type MultiSource struct {
	sources []Source
}

// NewMultiSource creates a MultiSource from the provided sources.
func NewMultiSource(sources ...Source) *MultiSource {
	return &MultiSource{sources: sources}
}

// Add appends a source to the MultiSource.
func (m *MultiSource) Add(s Source) {
	m.sources = append(m.sources, s)
}

// Lines starts all underlying sources and merges their output into one
// channel. The channel is closed once all sources have finished or the
// context is cancelled.
func (m *MultiSource) Lines(ctx context.Context) (<-chan Entry, error) {
	out := make(chan Entry)
	var wg sync.WaitGroup

	for _, src := range m.sources {
		ch, err := src.Lines(ctx)
		if err != nil {
			return nil, err
		}

		wg.Add(1)
		go func(name string, lines <-chan string) {
			defer wg.Done()
			for line := range lines {
				select {
				case <-ctx.Done():
					return
				case out <- Entry{Source: name, Line: line}:
				}
			}
		}(src.Name(), ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out, nil
}
