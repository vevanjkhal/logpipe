package filter

import (
	"fmt"
	"regexp"
)

// Builder provides a fluent API for constructing filter Options.
type Builder struct {
	opts Options
	errs []error
}

// NewBuilder returns a Builder with default options (LevelDebug, no restrictions).
func NewBuilder() *Builder {
	return &Builder{
		opts: Options{
			MinLevel: LevelDebug,
			Fields:   make(map[string]string),
		},
	}
}

// WithMinLevel sets the minimum log level by name (e.g. "warn").
func (b *Builder) WithMinLevel(level string) *Builder {
	l, ok := ParseLevel(level)
	if !ok {
		b.errs = append(b.errs, fmt.Errorf("unknown log level: %q", level))
		return b
	}
	b.opts.MinLevel = l
	return b
}

// WithPattern sets a regex pattern that log messages must match.
func (b *Builder) WithPattern(pattern string) *Builder {
	if pattern == "" {
		return b
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		b.errs = append(b.errs, fmt.Errorf("invalid pattern %q: %w", pattern, err))
		return b
	}
	b.opts.Pattern = re
	return b
}

// WithField adds a required key=value field constraint.
func (b *Builder) WithField(key, value string) *Builder {
	if key == "" {
		b.errs = append(b.errs, fmt.Errorf("field key must not be empty"))
		return b
	}
	b.opts.Fields[key] = value
	return b
}

// Build returns the constructed Options or an error if any configuration was invalid.
func (b *Builder) Build() (Options, error) {
	if len(b.errs) > 0 {
		return Options{}, fmt.Errorf("filter build errors: %v", b.errs)
	}
	return b.opts, nil
}
