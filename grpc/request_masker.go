package grpc

import (
	"errors"
	"strings"
)

// MaskerPolicy controls how sensitive fields are masked.
type MaskerPolicy struct {
	MaskChar  string
	KeepFirst int
	KeepLast  int
}

// DefaultMaskerPolicy returns sensible masking defaults.
func DefaultMaskerPolicy() MaskerPolicy {
	return MaskerPolicy{
		MaskChar:  "*",
		KeepFirst: 0,
		KeepLast:  0,
	}
}

// RequestMasker masks specified fields in a payload map.
type RequestMasker struct {
	policy MaskerPolicy
	fields map[string]struct{}
}

// NewRequestMasker creates a RequestMasker with the given policy.
func NewRequestMasker(policy MaskerPolicy) *RequestMasker {
	if policy.MaskChar == "" {
		policy = DefaultMaskerPolicy()
	}
	return &RequestMasker{
		policy: policy,
		fields: make(map[string]struct{}),
	}
}

// AddField registers a field name to be masked.
func (m *RequestMasker) AddField(field string) error {
	field = strings.TrimSpace(field)
	if field == "" {
		return errors.New("request masker: field name must not be empty")
	}
	m.fields[strings.ToLower(field)] = struct{}{}
	return nil
}

// Len returns the number of registered mask fields.
func (m *RequestMasker) Len() int {
	return len(m.fields)
}

// Mask returns a copy of payload with sensitive fields masked.
func (m *RequestMasker) Mask(payload map[string]any) map[string]any {
	if payload == nil {
		return nil
	}
	out := make(map[string]any, len(payload))
	for k, v := range payload {
		if _, ok := m.fields[strings.ToLower(k)]; ok {
			out[k] = m.maskValue(v)
		} else {
			out[k] = v
		}
	}
	return out
}

func (m *RequestMasker) maskValue(v any) string {
	s, ok := v.(string)
	if !ok {
		s = "****"
		return s
	}
	n := len(s)
	kf := m.policy.KeepFirst
	kl := m.policy.KeepLast
	if kf+kl >= n {
		return strings.Repeat(m.policy.MaskChar, n)
	}
	prefix := s[:kf]
	suffix := ""
	if kl > 0 {
		suffix = s[n-kl:]
	}
	midLen := n - kf - kl
	return prefix + strings.Repeat(m.policy.MaskChar, midLen) + suffix
}
