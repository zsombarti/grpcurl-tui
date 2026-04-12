package grpc

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ThrottlePolicy defines the configuration for request throttling.
type ThrottlePolicy struct {
	MaxConcurrent int
	QueueTimeout  time.Duration
}

// DefaultThrottlePolicy returns sensible defaults.
func DefaultThrottlePolicy() ThrottlePolicy {
	return ThrottlePolicy{
		MaxConcurrent: 10,
		QueueTimeout:  5 * time.Second,
	}
}

// RequestThrottle limits the number of concurrent in-flight gRPC requests.
type RequestThrottle struct {
	mu     sync.Mutex
	policy ThrottlePolicy
	sem    chan struct{}
	active int
}

// NewRequestThrottle creates a RequestThrottle with the given policy.
// Falls back to defaults if the policy is invalid.
func NewRequestThrottle(p ThrottlePolicy) *RequestThrottle {
	def := DefaultThrottlePolicy()
	if p.MaxConcurrent <= 0 {
		p.MaxConcurrent = def.MaxConcurrent
	}
	if p.QueueTimeout <= 0 {
		p.QueueTimeout = def.QueueTimeout
	}
	return &RequestThrottle{
		policy: p,
		sem:    make(chan struct{}, p.MaxConcurrent),
	}
}

// Acquire attempts to acquire a slot within the queue timeout.
// Returns an error if the context is cancelled or the queue timeout elapses.
func (t *RequestThrottle) Acquire(ctx context.Context) error {
	deadline := time.After(t.policy.QueueTimeout)
	select {
	case t.sem <- struct{}{}:
		t.mu.Lock()
		t.active++
		t.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-deadline:
		return errors.New("request throttle: queue timeout exceeded")
	}
}

// Release frees a previously acquired slot.
func (t *RequestThrottle) Release() {
	select {
	case <-t.sem:
		t.mu.Lock()
		if t.active > 0 {
			t.active--
		}
		t.mu.Unlock()
	default:
	}
}

// Active returns the number of currently held slots.
func (t *RequestThrottle) Active() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.active
}

// Policy returns the current throttle policy.
func (t *RequestThrottle) Policy() ThrottlePolicy {
	return t.policy
}
