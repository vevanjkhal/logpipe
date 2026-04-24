package source

import (
	"bufio"
	"context"
	"io"
	"os"
)

// Source represents a log input source.
type Source interface {
	// Name returns the identifier for this source.
	Name() string
	// Lines returns a channel that emits log lines.
	Lines(ctx context.Context) (<-chan string, error)
}

// FileSource reads log lines from a file (tail-like behavior).
type FileSource struct {
	name string
	path string
}

// NewFileSource creates a new FileSource for the given file path.
func NewFileSource(name, path string) *FileSource {
	return &FileSource{name: name, path: path}
}

// Name returns the source name.
func (f *FileSource) Name() string {
	return f.name
}

// Lines opens the file and streams each line over the returned channel.
// The channel is closed when the context is cancelled or EOF is reached.
func (f *FileSource) Lines(ctx context.Context) (<-chan string, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- scanner.Text():
			}
		}
	}()

	return ch, nil
}

// ReaderSource reads log lines from any io.Reader.
type ReaderSource struct {
	name   string
	reader io.Reader
}

// NewReaderSource creates a new ReaderSource wrapping the given reader.
func NewReaderSource(name string, r io.Reader) *ReaderSource {
	return &ReaderSource{name: name, reader: r}
}

// Name returns the source name.
func (r *ReaderSource) Name() string {
	return r.name
}

// Lines streams each line from the underlying reader over the returned channel.
func (r *ReaderSource) Lines(ctx context.Context) (<-chan string, error) {
	ch := make(chan string)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(r.reader)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- scanner.Text():
			}
		}
	}()
	return ch, nil
}
