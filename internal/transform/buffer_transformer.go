package transform

import "github.com/logpipe/logpipe/internal/buffer"

// BufferTransformer passes every line through unchanged but records it in an
// in-memory RingBuffer so callers can retrieve recent history at any time.
type BufferTransformer struct {
	buf *buffer.RingBuffer
}

// NewBuffer returns a Transformer that stores the last capacity lines.
// The underlying RingBuffer is accessible via Buffer().
func NewBuffer(capacity int) *BufferTransformer {
	return &BufferTransformer{buf: buffer.New(capacity)}
}

// Transform records the line and returns it unmodified.
func (t *BufferTransformer) Transform(line string) (string, bool) {
	t.buf.Write(line)
	return line, true
}

// Buffer returns the underlying RingBuffer for snapshot access.
func (t *BufferTransformer) Buffer() *buffer.RingBuffer {
	return t.buf
}
