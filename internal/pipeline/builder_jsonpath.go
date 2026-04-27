package pipeline

import (
	"fmt"

	"github.com/yourorg/logpipe/internal/transform"
)

// WithJSONPath adds a JSONPathTransformer to the pipeline. The transformer
// extracts the value of the named top-level JSON field and uses it as the new
// log line. When fallback is true, non-JSON lines or lines missing the field
// are passed through unchanged; otherwise they are dropped.
//
// Example:
//
//	builder.WithJSONPath("message", true)
func (b *Builder) WithJSONPath(field string, fallback bool) *Builder {
	tr, err := transform.NewJSONPath(field, fallback)
	if err != nil {
		panic(fmt.Sprintf("logpipe: WithJSONPath: %v", err))
	}
	b.transformers = append(b.transformers, tr)
	return b
}
