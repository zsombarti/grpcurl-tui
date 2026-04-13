package grpc

import (
	"errors"
	"strings"
)

// RedactorPolicy holds configuration for the request redactor.
type RedactorPolicy struct {
	Placeholder string
}

// DefaultRedactorPolicy returns a sensible default policy.
func DefaultRedactorPolicy() RedactorPolicy {
	return RedactorPolicy{
		Placeholder: "[REDACTED]",
	}
}

// RequestRedactor removes or replaces sensitive fields from a payload map.
type RequestRedactor struct {
	policy RedactorPolicy
	fields map[string]struct{}
}

// NewRequestRedactor creates a new RequestRedactor with the given policy.
// Falls back to the default policy if the placeholder is empty.
func NewRequestRedactor(policy RedactorPolicy) *RequestRedactor {
	if policy.Placeholder == "" {
		policy = DefaultRedactorPolicy()
	}
	return &RequestRedactor{
		policy: policy,
		fields: make(map[string]struct{}),
	}
}

// AddField registers a field name (case-insensitive) to be redacted.
func (r *RequestRedactor) AddField(field string) error {
	if strings.TrimSpace(field) == "" {
		return errors.New("redactor: field name must not be empty")
	}
	r.fields[strings.ToLower(field)] = struct{}{}
	return nil
}

// Len returns the number of registered redaction fields.
func (r *RequestRedactor) Len() int {
	return len(r.fields)
}

// Redact returns a copy of the payload with sensitive fields replaced.
func (r *RequestRedactor) Redact(payload map[string]any) map[string]any {
	if payload == nil {
		return nil
	}
	out := make(map[string]any, len(payload))
	for k, v := range payload {
		if _, ok := r.fields[strings.ToLower(k)]; ok {
			out[k] = r.policy.Placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// Policy returns the active redactor policy.
func (r *RequestRedactor) Policy() RedactorPolicy {
	return r.policy
}
