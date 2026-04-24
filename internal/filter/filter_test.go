package filter_test

import (
	"regexp"
	"testing"

	"github.com/yourorg/logpipe/internal/filter"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		input string
		want  filter.Level
		ok    bool
	}{
		{"debug", filter.LevelDebug, true},
		{"INFO", filter.LevelInfo, true},
		{"Warn", filter.LevelWarn, true},
		{"error", filter.LevelError, true},
		{"trace", filter.LevelDebug, false},
	}
	for _, c := range cases {
		got, ok := filter.ParseLevel(c.input)
		if ok != c.ok {
			t.Errorf("ParseLevel(%q) ok=%v, want %v", c.input, ok, c.ok)
		}
		if ok && got != c.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestFilter(t *testing.T) {
	entry := filter.Entry{
		Level:   filter.LevelInfo,
		Message: "user logged in",
		Fields:  map[string]string{"service": "auth", "user": "alice"},
		Raw:     `{"level":"info","msg":"user logged in"}`,
	}

	tests := []struct {
		name string
		opts filter.Options
		want bool
	}{
		{"pass all", filter.Options{MinLevel: filter.LevelDebug}, true},
		{"level too low", filter.Options{MinLevel: filter.LevelError}, false},
		{"pattern match", filter.Options{MinLevel: filter.LevelDebug, Pattern: regexp.MustCompile("logged")}, true},
		{"pattern no match", filter.Options{MinLevel: filter.LevelDebug, Pattern: regexp.MustCompile("logout")}, false},
		{"field match", filter.Options{MinLevel: filter.LevelDebug, Fields: map[string]string{"service": "auth"}}, true},
		{"field mismatch", filter.Options{MinLevel: filter.LevelDebug, Fields: map[string]string{"service": "api"}}, false},
		{"field missing", filter.Options{MinLevel: filter.LevelDebug, Fields: map[string]string{"region": "us"}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter.Filter(entry, tt.opts); got != tt.want {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
