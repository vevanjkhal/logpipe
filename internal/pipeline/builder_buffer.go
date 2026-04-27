package pipeline

import "github.com/logpipe/logpipe/internal/transform"

// WithBuffer adds a ring-buffer transformer to the pipeline that retains the
// last capacity log lines in memory. The returned *transform.BufferTransformer
// can be used outside the pipeline to inspect recent history via Snapshot().
//
// Example:
//
//	buf := transform.NewBuffer(200)
//	pipeline.NewBuilder().
//	    WithSource(src).
//	    WithBufferTransformer(buf).
//	    Build()
//	// later: buf.Buffer().Snapshot()
func (b *Builder) WithBufferTransformer(bt *transform.BufferTransformer) *Builder {
	b.transformers = append(b.transformers, bt)
	return b
}

// WithBuffer is a convenience method that creates a BufferTransformer with the
// given capacity, registers it, and returns the Builder for chaining.
// Use WithBufferTransformer when you need a handle to the transformer itself.
func (b *Builder) WithBuffer(capacity int) *Builder {
	return b.WithBufferTransformer(transform.NewBuffer(capacity))
}
