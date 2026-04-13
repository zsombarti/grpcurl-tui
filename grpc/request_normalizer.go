package grpc

import (
	"errors"
	"strings"
)

// NormalizerPolicy defines how request payloads are normalized.
type NormalizerPolicy struct {
	TrimStrings    bool
	LowercaseKeys  bool
	RemoveNullKeys bool
}

// DefaultNormalizerPolicy returns a sensible default normalization policy.
func DefaultNormalizerPolicy() NormalizerPolicy {
	return NormalizerPolicy{
		TrimStrings:    true,
		LowercaseKeys:  false,
		RemoveNullKeys: true,
	}
}

// RequestNormalizer normalizes request payloads according to a policy.
type RequestNormalizer struct {
	policy NormalizerPolicy
}

// NewRequestNormalizer creates a new RequestNormalizer with the given policy.
// If the policy is zero-valued, the default policy is used.
func NewRequestNormalizer(policy NormalizerPolicy) *RequestNormalizer {
	default_ := DefaultNormalizerPolicy()
	if !policy.TrimStrings && !policy.LowercaseKeys && !policy.RemoveNullKeys {
		policy = default_
	}
	return &RequestNormalizer{policy: policy}
}

// Normalize applies the normalization policy to the given payload map.
// Returns an error if payload is nil.
func (n *RequestNormalizer) Normalize(payload map[string]any) (map[string]any, error) {
	if payload == nil {
		return nil, errors.New("normalizer: payload must not be nil")
	}
	out := make(map[string]any, len(payload))
	for k, v := range payload {
		key := k
		if n.policy.LowercaseKeys {
			key = strings.ToLower(k)
		}
		if n.policy.RemoveNullKeys && v == nil {
			continue
		}
		if n.policy.TrimStrings {
			if s, ok := v.(string); ok {
				v = strings.TrimSpace(s)
			}
		}
		out[key] = v
	}
	return out, nil
}

// Policy returns the current normalization policy.
func (n *RequestNormalizer) Policy() NormalizerPolicy {
	return n.policy
}
