package grpc

import (
	"errors"
	"sync"
	"time"
)

const defaultCheckpointMaxSize = 100

// Checkpoint represents a named save-point for a request payload.
type Checkpoint struct {
	Name      string
	Method    string
	Payload   map[string]interface{}
	CreatedAt time.Time
}

// RequestCheckpointer stores named checkpoints for request payloads.
type RequestCheckpointer struct {
	mu      sync.RWMutex
	items   []Checkpoint
	maxSize int
}

// NewRequestCheckpointer creates a new RequestCheckpointer with the given max size.
// If maxSize <= 0 the default is used.
func NewRequestCheckpointer(maxSize int) *RequestCheckpointer {
	if maxSize <= 0 {
		maxSize = defaultCheckpointMaxSize
	}
	return &RequestCheckpointer{maxSize: maxSize}
}

// Save stores a checkpoint. Returns an error if name or method is empty.
func (c *RequestCheckpointer) Save(name, method string, payload map[string]interface{}) error {
	if name == "" {
		return errors.New("checkpoint name must not be empty")
	}
	if method == "" {
		return errors.New("method must not be empty")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	cp := Checkpoint{
		Name:      name,
		Method:    method,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
	// Replace existing checkpoint with same name.
	for i, existing := range c.items {
		if existing.Name == name {
			c.items[i] = cp
			return nil
		}
	}
	if len(c.items) >= c.maxSize {
		c.items = c.items[1:]
	}
	c.items = append(c.items, cp)
	return nil
}

// Get retrieves a checkpoint by name. Returns false if not found.
func (c *RequestCheckpointer) Get(name string) (Checkpoint, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, cp := range c.items {
		if cp.Name == name {
			return cp, true
		}
	}
	return Checkpoint{}, false
}

// Delete removes a checkpoint by name. Returns false if not found.
func (c *RequestCheckpointer) Delete(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, cp := range c.items {
		if cp.Name == name {
			c.items = append(c.items[:i], c.items[i+1:]...)
			return true
		}
	}
	return false
}

// Len returns the number of stored checkpoints.
func (c *RequestCheckpointer) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// All returns a copy of all checkpoints.
func (c *RequestCheckpointer) All() []Checkpoint {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Checkpoint, len(c.items))
	copy(out, c.items)
	return out
}
