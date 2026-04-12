package grpc

import (
	"strings"
	"unicode"
)

// SanitizerPolicy holds configuration for request sanitization.
type SanitizerPolicy struct {
	RedactedKeys   []string
	RedactedValue  string
	MaxValueLength int
}

// DefaultSanitizerPolicy returns a sensible default policy.
func DefaultSanitizerPolicy() SanitizerPolicy {
	return SanitizerPolicy{
		RedactedKeys:   []string{"password", "secret", "token", "authorization", "api_key"},
		RedactedValue:  "[REDACTED]",
		MaxValueLength: 1024,
	}
}

// RequestSanitizer scrubs sensitive fields from request payloads.
type RequestSanitizer struct {
	policy SanitizerPolicy
}

// NewRequestSanitizer creates a RequestSanitizer with the given policy.
// Falls back to DefaultSanitizerPolicy if MaxValueLength is zero.
func NewRequestSanitizer(policy SanitizerPolicy) *RequestSanitizer {
	if policy.MaxValueLength <= 0 {
		policy = DefaultSanitizerPolicy()
	}
	if policy.RedactedValue == "" {
		policy.RedactedValue = DefaultSanitizerPolicy().RedactedValue
	}
	return &RequestSanitizer{policy: policy}
}

// Sanitize returns a copy of the payload map with sensitive keys redacted
// and oversized string values truncated.
func (s *RequestSanitizer) Sanitize(payload map[string]interface{}) map[string]interface{} {
	if payload == nil {
		return nil
	}
	out := make(map[string]interface{}, len(payload))
	for k, v := range payload {
		if s.isRedacted(k) {
			out[k] = s.policy.RedactedValue
			continue
		}
		switch val := v.(type) {
		case string:
			out[k] = s.truncate(val)
		case map[string]interface{}:
			out[k] = s.Sanitize(val)
		default:
			out[k] = v
		}
	}
	return out
}

// Policy returns the active sanitizer policy.
func (s *RequestSanitizer) Policy() SanitizerPolicy { return s.policy }

func (s *RequestSanitizer) isRedacted(key string) bool {
	norm := strings.ToLower(strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return r
		}
		return '_'
	}, key))
	for _, rk := range s.policy.RedactedKeys {
		if strings.Contains(norm, strings.ToLower(rk)) {
			return true
		}
	}
	return false
}

func (s *RequestSanitizer) truncate(v string) string {
	if len(v) > s.policy.MaxValueLength {
		return v[:s.policy.MaxValueLength]
	}
	return v
}
