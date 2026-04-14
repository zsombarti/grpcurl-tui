package grpc

import (
	"errors"
	"sync"
)

// BroadcasterPolicy holds configuration for the RequestBroadcaster.
type BroadcasterPolicy struct {
	MaxTargets int
}

// DefaultBroadcasterPolicy returns sensible defaults.
func DefaultBroadcasterPolicy() BroadcasterPolicy {
	return BroadcasterPolicy{
		MaxTargets: 10,
	}
}

// BroadcastResult holds the outcome for a single target.
type BroadcastResult struct {
	Target  string
	Payload map[string]any
	Err     error
}

// RequestBroadcaster fans out a payload to multiple named targets.
type RequestBroadcaster struct {
	mu      sync.RWMutex
	policy  BroadcasterPolicy
	targets []string
}

// NewRequestBroadcaster creates a RequestBroadcaster with the given policy.
// Invalid policies fall back to defaults.
func NewRequestBroadcaster(policy BroadcasterPolicy) *RequestBroadcaster {
	if policy.MaxTargets <= 0 {
		policy = DefaultBroadcasterPolicy()
	}
	return &RequestBroadcaster{policy: policy}
}

// AddTarget registers a named target. Returns an error if the target name is
// empty or the maximum number of targets has been reached.
func (b *RequestBroadcaster) AddTarget(name string) error {
	if name == "" {
		return errors.New("broadcaster: target name must not be empty")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.targets) >= b.policy.MaxTargets {
		return errors.New("broadcaster: maximum number of targets reached")
	}
	b.targets = append(b.targets, name)
	return nil
}

// Len returns the number of registered targets.
func (b *RequestBroadcaster) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.targets)
}

// Broadcast fans out the payload to all registered targets concurrently and
// returns one BroadcastResult per target.
func (b *RequestBroadcaster) Broadcast(payload map[string]any) []BroadcastResult {
	b.mu.RLock()
	targets := make([]string, len(b.targets))
	copy(targets, b.targets)
	b.mu.RUnlock()

	results := make([]BroadcastResult, len(targets))
	var wg sync.WaitGroup
	for i, t := range targets {
		wg.Add(1)
		go func(idx int, target string) {
			defer wg.Done()
			// Shallow-copy payload per target to avoid data races.
			copy := make(map[string]any, len(payload))
			for k, v := range payload {
				copy[k] = v
			}
			results[idx] = BroadcastResult{Target: target, Payload: copy}
		}(i, t)
	}
	wg.Wait()
	return results
}

// Clear removes all registered targets.
func (b *RequestBroadcaster) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.targets = b.targets[:0]
}
