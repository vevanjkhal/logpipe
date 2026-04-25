package transform

import (
	"encoding/json"
	"regexp"
	"strings"
)

// MaskTransformer replaces substrings matching a regex pattern with a mask string.
type MaskTransformer struct {
	pattern *regexp.Regexp
	mask    string
}

// NewMask creates a MaskTransformer that replaces all occurrences of pattern
// with mask in each log line. If mask is empty, "***" is used.
func NewMask(pattern string, mask string) (*MaskTransformer, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if mask == "" {
		mask = "***"
	}
	return &MaskTransformer{pattern: re, mask: mask}, nil
}

// Transform applies the mask to the log line. If the line is valid JSON,
// the mask is applied to string values only; otherwise it is applied to the
// raw line.
func (m *MaskTransformer) Transform(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err == nil {
		for k, v := range obj {
			if s, ok := v.(string); ok {
				obj[k] = m.pattern.ReplaceAllString(s, m.mask)
			}
		}
		b, err := json.Marshal(obj)
		if err == nil {
			return string(b)
		}
	}
	return m.pattern.ReplaceAllString(line, m.mask)
}

// TruncateTransformer shortens log lines that exceed a maximum byte length.
type TruncateTransformer struct {
	maxLen int
	suffix string
}

// NewTruncate creates a TruncateTransformer. Lines longer than maxLen bytes
// are cut and suffixed with suffix (default: "...").
func NewTruncate(maxLen int, suffix string) *TruncateTransformer {
	if suffix == "" {
		suffix = "..."
	}
	if maxLen <= 0 {
		maxLen = 256
	}
	return &TruncateTransformer{maxLen: maxLen, suffix: suffix}
}

// Transform truncates line to maxLen bytes, appending the suffix when trimmed.
func (t *TruncateTransformer) Transform(line string) string {
	if len(line) <= t.maxLen {
		return line
	}
	return strings.TrimRight(line[:t.maxLen], " ") + t.suffix
}
