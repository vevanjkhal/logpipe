package filter

import (
	"regexp"
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
}

// ParseLevel converts a string to a Level, returning LevelDebug and false if unknown.
func ParseLevel(s string) (Level, bool) {
	l, ok := levelNames[strings.ToLower(s)]
	return l, ok
}

// Options holds the filtering criteria for log entries.
type Options struct {
	// MinLevel filters out entries below this severity.
	MinLevel Level
	// Pattern, if non-nil, requires log messages to match.
	Pattern *regexp.Regexp
	// Fields requires all specified key=value pairs to be present.
	Fields map[string]string
}

// Entry represents a single parsed log entry.
type Entry struct {
	Level   Level
	Message string
	Fields  map[string]string
	Raw     string
}

// Filter evaluates whether an Entry satisfies the given Options.
func Filter(e Entry, opts Options) bool {
	if e.Level < opts.MinLevel {
		return false
	}
	if opts.Pattern != nil && !opts.Pattern.MatchString(e.Message) {
		return false
	}
	for k, v := range opts.Fields {
		if got, ok := e.Fields[k]; !ok || got != v {
			return false
		}
	}
	return true
}
