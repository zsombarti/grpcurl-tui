package grpc

import (
	"errors"
	"sync"
)

// MirrorTarget represents a destination for mirrored requests.
type MirrorTarget struct {
	Address string
	Weight  int // 1-100, percentage of requests to mirror
}

// RequestMirror duplicates outgoing requests to one or more shadow targets.
type RequestMirror struct {
	mu      sync.RWMutex
	targets []MirrorTarget
	maxSize int
}

// NewRequestMirror creates a new RequestMirror with the given max target size.
func NewRequestMirror(maxSize int) *RequestMirror {
	if maxSize <= 0 {
		maxSize = 8
	}
	return &RequestMirror{maxSize: maxSize}
}

// AddTarget registers a mirror target. Returns an error if the address is
// empty, weight is out of range, or the target list is full.
func (m *RequestMirror) AddTarget(t MirrorTarget) error {
	if t.Address == "" {
		return errors.New("mirror target address must not be empty")
	}
	if t.Weight < 1 || t.Weight > 100 {
		return errors.New("mirror target weight must be between 1 and 100")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.targets) >= m.maxSize {
		return errors.New("mirror target list is full")
	}
	m.targets = append(m.targets, t)
	return nil
}

// Targets returns a copy of the current mirror targets.
func (m *RequestMirror) Targets() []MirrorTarget {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]MirrorTarget, len(m.targets))
	copy(out, m.targets)
	return out
}

// Len returns the number of registered mirror targets.
func (m *RequestMirror) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.targets)
}

// Clear removes all registered mirror targets.
func (m *RequestMirror) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.targets = nil
}
