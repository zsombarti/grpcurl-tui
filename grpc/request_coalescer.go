package grpc

import (
	"errors"
	"sync"
	"time"
)

// DefaultCoalescerPolicy returns a sensible default coalescing policy.
func DefaultCoalescerPolicy() CoalescerPolicy {
	return CoalescerPolicy{
		Window:  50 * time.Millisecond,
		MaxKeys: 256,
	}
}

// CoalescerPolicy controls how requests are coalesced.
type CoalescerPolicy struct {
	Window  time.Duration
	MaxKeys int
}

// coalesceEntry holds in-flight coalesced calls for a single key.
type coalesceEntry struct {
	result map[string]any
	err    error
	done   chan struct{}
}

// RequestCoalescer merges concurrent identical requests into a single in-flight call.
type RequestCoalescer struct {
	mu     sync.Mutex
	policy CoalescerPolicy
	inflight map[string]*coalesceEntry
}

// NewRequestCoalescer creates a new RequestCoalescer with the given policy.
// Falls back to the default policy when the provided policy is invalid.
func NewRequestCoalescer(policy CoalescerPolicy) *RequestCoalescer {
	if policy.Window <= 0 || policy.MaxKeys <= 0 {
		policy = DefaultCoalescerPolicy()
	}
	return &RequestCoalescer{
		policy:   policy,
		inflight: make(map[string]*coalesceEntry),
	}
}

// Policy returns the active coalescer policy.
func (c *RequestCoalescer) Policy() CoalescerPolicy {
	return c.policy
}

// Len returns the number of currently in-flight coalesced keys.
func (c *RequestCoalescer) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.inflight)
}

// Do executes fn for the given key, or waits for an in-flight call with the
// same key to complete and returns its result. Concurrent callers sharing a
// key receive the same result.
func (c *RequestCoalescer) Do(key string, fn func() (map[string]any, error)) (map[string]any, error) {
	if key == "" {
		return nil, errors.New("coalescer: key must not be empty")
	}

	c.mu.Lock()
	if len(c.inflight) >= c.policy.MaxKeys {
		c.mu.Unlock()
		return nil, errors.New("coalescer: max in-flight keys reached")
	}
	if entry, ok := c.inflight[key]; ok {
		c.mu.Unlock()
		<-entry.done
		return entry.result, entry.err
	}
	entry := &coalesceEntry{done: make(chan struct{})}
	c.inflight[key] = entry
	c.mu.Unlock()

	entry.result, entry.err = fn()
	close(entry.done)

	time.AfterFunc(c.policy.Window, func() {
		c.mu.Lock()
		delete(c.inflight, key)
		c.mu.Unlock()
	})

	return entry.result, entry.err
}
