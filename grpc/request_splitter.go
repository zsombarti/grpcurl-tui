package grpc

import (
	"errors"
	"sync"
)

// SplitterPolicy controls how requests are fanned out.
type SplitterPolicy struct {
	MaxTargets int
}

// DefaultSplitterPolicy returns sensible defaults.
func DefaultSplitterPolicy() SplitterPolicy {
	return SplitterPolicy{MaxTargets: 8}
}

// SplitResult holds the target address and any error for a single fan-out leg.
type SplitResult struct {
	Target string
	Payload map[string]interface{}
	Err    error
}

// RequestSplitter fans a single request payload out to multiple target addresses.
type RequestSplitter struct {
	mu      sync.RWMutex
	policy  SplitterPolicy
	targets []string
}

// NewRequestSplitter creates a RequestSplitter with the given policy.
// Falls back to defaults when the policy is invalid.
func NewRequestSplitter(policy SplitterPolicy) *RequestSplitter {
	if policy.MaxTargets <= 0 {
		policy = DefaultSplitterPolicy()
	}
	return &RequestSplitter{policy: policy}
}

// AddTarget registers a target address for fan-out.
// Returns an error when the address is empty or the cap is reached.
func (s *RequestSplitter) AddTarget(address string) error {
	if address == "" {
		return errors.New("splitter: target address must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.targets) >= s.policy.MaxTargets {
		return errors.New("splitter: max targets reached")
	}
	s.targets = append(s.targets, address)
	return nil
}

// Len returns the number of registered targets.
func (s *RequestSplitter) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.targets)
}

// Split copies payload to every registered target and returns one SplitResult per target.
// A nil payload is treated as an empty map.
func (s *RequestSplitter) Split(payload map[string]interface{}) []SplitResult {
	if payload == nil {
		payload = map[string]interface{}{}
	}
	s.mu.RLock()
	targets := make([]string, len(s.targets))
	copy(targets, s.targets)
	s.mu.RUnlock()

	results := make([]SplitResult, 0, len(targets))
	for _, t := range targets {
		copy := make(map[string]interface{}, len(payload))
		for k, v := range payload {
			copy[k] = v
		}
		results = append(results, SplitResult{Target: t, Payload: copy})
	}
	return results
}

// ClearTargets removes all registered targets.
func (s *RequestSplitter) ClearTargets() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.targets = nil
}
