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

// NormalizeMany applies the normalization policy to each payload in the slice.
// Returns an error if any payload is nil, including the index of the offending entry.
func (n *RequestNormalizer) NormalizeMany(payloads []map[string]any) ([]map[string]any, error) {
	out := make([]map[string]any, 0, len(payloads))
	for i, p := range payloads {
		normalized, err := n.Normalize(p)
		if err != nil {
			return nil, fmt.Errorf("normalizer: error at index %d: %w", i, err)
		}
		out = append(out, normalized)
	}
	return out, nil
}

// Policy returns the current normalization policy.
func (n *RequestNormalizer) Policy() NormalizerPolicy {
	return n.policy
}
