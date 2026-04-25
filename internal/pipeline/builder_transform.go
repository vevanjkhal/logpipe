package pipeline

import "github.com/logpipe/logpipe/internal/transform"

// WithRedact adds a RedactTransformer to the pipeline builder.
// The provided keys will have their values masked in JSON log lines.
func (b *Builder) WithRedact(keys []string, maskVal string) *Builder {
	b.transformers = append(b.transformers, transform.NewRedact(keys, maskVal))
	return b
}

// WithAddField adds an AddFieldTransformer to the pipeline builder.
// The provided key/value pair will be injected into JSON log lines.
func (b *Builder) WithAddField(key, value string) *Builder {
	b.transformers = append(b.transformers, transform.NewAddField(key, value))
	return b
}

// WithPrefix adds a PrefixTransformer to the pipeline builder.
// Every log line will be prepended with the given prefix string.
func (b *Builder) WithPrefix(prefix string) *Builder {
	b.transformers = append(b.transformers, transform.NewPrefix(prefix))
	return b
}
