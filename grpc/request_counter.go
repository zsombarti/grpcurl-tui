package grpc

import (
	"sync"
	"time"
)

// RequestCounter tracks per-method call counts and error rates.
type RequestCounter struct {
	mu      sync.RWMutex
	counts  map[string]int
	errors  map[string]int
	reset   time.Time
}

// NewRequestCounter creates a new RequestCounter.
func NewRequestCounter() *RequestCounter {
	return &RequestCounter{
		counts: make(map[string]int),
		errors: make(map[string]int),
		reset:  time.Now(),
	}
}

// Increment records a call for the given method.
func (c *RequestCounter) Increment(method string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[method]++
}

// IncrementError records an error for the given method.
func (c *RequestCounter) IncrementError(method string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errors[method]++
}

// Count returns the total call count for a method.
func (c *RequestCounter) Count(method string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.counts[method]
}

// ErrorCount returns the total error count for a method.
func (c *RequestCounter) ErrorCount(method string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.errors[method]
}

// Methods returns all tracked method names.
func (c *RequestCounter) Methods() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.counts))
	for k := range c.counts {
		keys = append(keys, k)
	}
	return keys
}

// Reset clears all counters and records the reset time.
func (c *RequestCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts = make(map[string]int)
	c.errors = make(map[string]int)
	c.reset = time.Now()
}

// Since returns the time elapsed since the last reset.
func (c *RequestCounter) Since() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Since(c.reset)
}
