package pipeline

import "github.com/yourorg/logpipe/internal/transform"

// WithConditional adds a transformer that applies inner only when the JSON
// field `field` equals `value`. For plain-text lines a substring match is
// used instead.
//
// Example:
//
//	builder.WithConditional("level", "error",
//	    transform.NewRedact("password"),
//	)
func (b *Builder) WithConditional(field, value string, inner transform.Transformer, opts ...transform.ConditionalOption) *Builder {
	b.transforms = append(b.transforms, transform.NewConditional(field, value, inner, opts...))
	return b
}

// WithConditionalInverted adds a transformer that applies inner only when the
// JSON field `field` does NOT equal `value`.
func (b *Builder) WithConditionalInverted(field, value string, inner transform.Transformer) *Builder {
	return b.WithConditional(field, value, inner, transform.WithInvert())
}
