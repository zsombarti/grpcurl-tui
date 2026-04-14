package grpc

import (
	"errors"
	"sync"
)

// FencePolicy defines configuration for the request fence.
type FencePolicy struct {
	MaxConcurrent int
}

// DefaultFencePolicy returns a FencePolicy with sensible defaults.
func DefaultFencePolicy() FencePolicy {
	return FencePolicy{
		MaxConcurrent: 10,
	}
}

// RequestFence limits the number of concurrent in-flight requests.
type RequestFence struct {
	mu     sync.Mutex
	policy FencePolicy
	active int
}

// NewRequestFence creates a new RequestFence with the given policy.
// Invalid policies fall back to defaults.
func NewRequestFence(policy FencePolicy) *RequestFence {
	def := DefaultFencePolicy()
	if policy.MaxConcurrent <= 0 {
		policy.MaxConcurrent = def.MaxConcurrent
	}
	return &RequestFence{policy: policy}
}

// Acquire attempts to acquire a slot. Returns an error if the fence is full.
func (f *RequestFence) Acquire() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.active >= f.policy.MaxConcurrent {
		return errors.New("fence: max concurrent requests reached")
	}
	f.active++
	return nil
}

// Release releases a previously acquired slot.
func (f *RequestFence) Release() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.active > 0 {
		f.active--
	}
}

// Active returns the current number of in-flight requests.
func (f *RequestFence) Active() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.active
}

// Policy returns the current fence policy.
func (f *RequestFence) Policy() FencePolicy {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.policy
}
