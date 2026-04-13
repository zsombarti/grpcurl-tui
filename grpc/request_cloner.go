package grpc

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ClonerPolicy configures deep-copy behaviour.
type ClonerPolicy struct {
	MaxDepth int
}

// DefaultClonerPolicy returns sensible defaults.
func DefaultClonerPolicy() ClonerPolicy {
	return ClonerPolicy{MaxDepth: 64}
}

// RequestCloner deep-copies JSON payloads.
type RequestCloner struct {
	policy ClonerPolicy
}

// NewRequestCloner creates a RequestCloner with the given policy.
// Invalid policies fall back to defaults.
func NewRequestCloner(p ClonerPolicy) *RequestCloner {
	if p.MaxDepth <= 0 {
		p = DefaultClonerPolicy()
	}
	return &RequestCloner{policy: p}
}

// Clone returns a deep copy of the supplied payload map.
// Returns an error if payload is nil or exceeds MaxDepth.
func (c *RequestCloner) Clone(payload map[string]any) (map[string]any, error) {
	if payload == nil {
		return nil, errors.New("request_cloner: nil payload")
	}
	cloned, err := deepCopyMap(payload, 0, c.policy.MaxDepth)
	if err != nil {
		return nil, fmt.Errorf("request_cloner: %w", err)
	}
	return cloned, nil
}

// CloneJSON round-trips through JSON to produce a deep copy.
func (c *RequestCloner) CloneJSON(payload map[string]any) (map[string]any, error) {
	if payload == nil {
		return nil, errors.New("request_cloner: nil payload")
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("request_cloner: marshal: %w", err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("request_cloner: unmarshal: %w", err)
	}
	return out, nil
}

// Policy returns the active policy.
func (c *RequestCloner) Policy() ClonerPolicy { return c.policy }

func deepCopyMap(m map[string]any, depth, max int) (map[string]any, error) {
	if depth > max {
		return nil, errors.New("max depth exceeded")
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		copied, err := deepCopyValue(v, depth+1, max)
		if err != nil {
			return nil, err
		}
		out[k] = copied
	}
	return out, nil
}

func deepCopyValue(v any, depth, max int) (any, error) {
	switch val := v.(type) {
	case map[string]any:
		return deepCopyMap(val, depth, max)
	case []any:
		slice := make([]any, len(val))
		for i, elem := range val {
			copied, err := deepCopyValue(elem, depth+1, max)
			if err != nil {
				return nil, err
			}
			slice[i] = copied
		}
		return slice, nil
	default:
		return val, nil
	}
}
