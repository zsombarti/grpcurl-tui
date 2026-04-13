package grpc

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// CorrelationEntry holds a correlation ID mapped to a request method and timestamp.
type CorrelationEntry struct {
	ID        string
	Method    string
	Timestamp time.Time
}

// RequestCorrelator assigns and tracks correlation IDs for outgoing requests.
type RequestCorrelator struct {
	mu      sync.RWMutex
	entries []CorrelationEntry
	maxSize int
}

// NewRequestCorrelator creates a new RequestCorrelator with the given max size.
// If maxSize <= 0 it defaults to 200.
func NewRequestCorrelator(maxSize int) *RequestCorrelator {
	if maxSize <= 0 {
		maxSize = 200
	}
	return &RequestCorrelator{maxSize: maxSize}
}

// Assign generates a new correlation ID for the given method, stores it, and
// returns the ID.
func (c *RequestCorrelator) Assign(method string) (string, error) {
	if method == "" {
		return "", errors.New("request_correlator: method must not be empty")
	}
	id := uuid.NewString()
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= c.maxSize {
		c.entries = c.entries[1:]
	}
	c.entries = append(c.entries, CorrelationEntry{
		ID:        id,
		Method:    method,
		Timestamp: time.Now(),
	})
	return id, nil
}

// Lookup returns the CorrelationEntry for the given ID, or an error if not found.
func (c *RequestCorrelator) Lookup(id string) (CorrelationEntry, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, e := range c.entries {
		if e.ID == id {
			return e, nil
		}
	}
	return CorrelationEntry{}, errors.New("request_correlator: id not found")
}

// Len returns the number of tracked correlation entries.
func (c *RequestCorrelator) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// Clear removes all stored correlation entries.
func (c *RequestCorrelator) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = nil
}
