package source

import (
	"bufio"
	"context"
	"io"
	"os"
)

// StdinSource reads log lines from standard input.
type StdinSource struct {
	reader io.Reader
}

// NewStdinSource creates a new StdinSource reading from os.Stdin.
func NewStdinSource() *StdinSource {
	return &StdinSource{reader: os.Stdin}
}

// newStdinSourceFromReader creates a StdinSource from an arbitrary reader (for testing).
func newStdinSourceFromReader(r io.Reader) *StdinSource {
	return &StdinSource{reader: r}
}

// Name returns the name identifier for this source.
func (s *StdinSource) Name() string {
	return "stdin"
}

// Lines emits each line read from stdin onto the returned channel.
// The channel is closed when stdin reaches EOF or the context is cancelled.
func (s *StdinSource) Lines(ctx context.Context) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		scanner := bufio.NewScanner(s.reader)
		for scanner.Scan() {
			line := scanner.Text()
			select {
			case <-ctx.Done():
				return
			case ch <- line:
			}
		}
	}()

	return ch
}
