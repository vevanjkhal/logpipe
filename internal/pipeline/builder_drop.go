package pipeline

import "github.com/yourorg/logpipe/internal/transform"

// WithDrop registers a DropTransformer on the builder that discards any line
// containing at least one of the provided substrings. Multiple calls append
// additional substrings to a single transformer.
//
// Example:
//
//	pipeline.NewBuilder().
//	    WithSource(src).
//	    WithDrop([]string{"healthcheck", "debug"}).
//	    Build()
func (b *Builder) WithDrop(substrings []string) *Builder {
	if len(substrings) == 0 {
		return b
	}
	b.transformers = append(b.transformers, transform.NewDrop(substrings))
	return b
}
