// Package routing provides log-line routing based on field matching.
// Lines can be dispatched to different outputs depending on their content.
package routing

import (
	"encoding/json"
	"strings"
)

// Rule defines a single routing rule: if the log line matches the
// condition, it is forwarded to the named destination.
type Rule struct {
	// Field is the JSON key to inspect. If empty, the whole line is matched.
	Field string
	// Contains is the substring that must appear in the field value (or line).
	Contains string
	// Destination is an opaque label used by the caller to select an output.
	Destination string
}

// Router evaluates an ordered list of Rules against each log line and
// returns the destination label for the first matching rule.
// If no rule matches, DefaultDestination is returned.
type Router struct {
	rules              []Rule
	DefaultDestination string
}

// New creates a Router with the provided rules and a default destination.
func New(defaultDest string, rules ...Rule) *Router {
	return &Router{
		rules:              rules,
		DefaultDestination: defaultDest,
	}
}

// Route returns the destination label for line.
func (r *Router) Route(line string) string {
	for _, rule := range r.rules {
		if r.matches(line, rule) {
			return rule.Destination
		}
	}
	return r.DefaultDestination
}

// matches checks whether line satisfies rule.
func (r *Router) matches(line string, rule Rule) bool {
	if rule.Field == "" {
		return strings.Contains(line, rule.Contains)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return false
	}
	val, ok := obj[rule.Field]
	if !ok {
		return false
	}
	str, ok := val.(string)
	if !ok {
		return false
	}
	return strings.Contains(str, rule.Contains)
}
