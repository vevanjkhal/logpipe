package transform

import (
	"fmt"
	"regexp"
)

// RegexFilter suppresses lines that match (or do not match) a regular expression.
// When invert is false, only lines matching the pattern are passed through.
// When invert is true, only lines that do NOT match are passed through.
type RegexFilter struct {
	pattern *regexp.Regexp
	invert  bool
}

// NewRegexFilter creates a transformer that keeps lines matching pattern.
// Set invert to true to keep lines that do NOT match.
func NewRegexFilter(pattern string, invert bool) (*RegexFilter, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("regexfilter: invalid pattern %q: %w", pattern, err)
	}
	return &RegexFilter{pattern: re, invert: invert}, nil
}

// Transform returns the line unchanged if it passes the filter, or an empty
// string to signal suppression.
func (r *RegexFilter) Transform(line string) (string, error) {
	matched := r.pattern.MatchString(line)
	if r.invert {
		matched = !matched
	}
	if !matched {
		return "", nil
	}
	return line, nil
}
