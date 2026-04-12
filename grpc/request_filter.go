package grpc

import (
	"errors"
	"strings"
)

// FilterRule defines a single filter predicate on a request field.
type FilterRule struct {
	Field    string
	Operator string // "eq", "contains", "prefix"
	Value    string
}

// RequestFilter holds a set of rules and applies them to request payloads.
type RequestFilter struct {
	rules []FilterRule
}

// NewRequestFilter creates a new RequestFilter with no rules.
func NewRequestFilter() *RequestFilter {
	return &RequestFilter{
		rules: make([]FilterRule, 0),
	}
}

// AddRule appends a filter rule. Returns an error for invalid input.
func (f *RequestFilter) AddRule(field, operator, value string) error {
	field = strings.TrimSpace(field)
	if field == "" {
		return errors.New("request_filter: field must not be empty")
	}
	switch operator {
	case "eq", "contains", "prefix":
	default:
		return errors.New("request_filter: unsupported operator: " + operator)
	}
	f.rules = append(f.rules, FilterRule{Field: field, Operator: operator, Value: value})
	return nil
}

// Len returns the number of active filter rules.
func (f *RequestFilter) Len() int {
	return len(f.rules)
}

// Rules returns a copy of the current filter rules.
func (f *RequestFilter) Rules() []FilterRule {
	out := make([]FilterRule, len(f.rules))
	copy(out, f.rules)
	return out
}

// Match reports whether the given payload map satisfies all filter rules.
func (f *RequestFilter) Match(payload map[string]string) bool {
	for _, rule := range f.rules {
		v, ok := payload[rule.Field]
		if !ok {
			return false
		}
		switch rule.Operator {
		case "eq":
			if v != rule.Value {
				return false
			}
		case "contains":
			if !strings.Contains(v, rule.Value) {
				return false
			}
		case "prefix":
			if !strings.HasPrefix(v, rule.Value) {
				return false
			}
		}
	}
	return true
}

// Clear removes all filter rules.
func (f *RequestFilter) Clear() {
	f.rules = f.rules[:0]
}
