package grpc

import (
	"errors"
	"sync"
)

// RequestLabel associates a human-readable label with a gRPC method call.
type RequestLabel struct {
	Method string
	Label  string
	Color  string // optional color hint for UI rendering
}

// RequestLabeler stores and retrieves labels for gRPC requests.
type RequestLabeler struct {
	mu     sync.RWMutex
	labels map[string]RequestLabel // keyed by method
	maxLen int
}

// NewRequestLabeler creates a RequestLabeler with the given max capacity.
func NewRequestLabeler(maxLen int) *RequestLabeler {
	const defaultMax = 128
	if maxLen <= 0 {
		maxLen = defaultMax
	}
	return &RequestLabeler{
		labels: make(map[string]RequestLabel, maxLen),
		maxLen: maxLen,
	}
}

// Set assigns a label to a method, replacing any existing one.
func (r *RequestLabeler) Set(method, label, color string) error {
	if method == "" {
		return errors.New("request labeler: method must not be empty")
	}
	if label == "" {
		return errors.New("request labeler: label must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.labels[method]; !exists && len(r.labels) >= r.maxLen {
		// evict an arbitrary entry to stay within capacity
		for k := range r.labels {
			delete(r.labels, k)
			break
		}
	}
	r.labels[method] = RequestLabel{Method: method, Label: label, Color: color}
	return nil
}

// Get returns the label for a method and whether it was found.
func (r *RequestLabeler) Get(method string) (RequestLabel, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.labels[method]
	return l, ok
}

// Remove deletes the label for the given method.
func (r *RequestLabeler) Remove(method string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.labels, method)
}

// Len returns the number of stored labels.
func (r *RequestLabeler) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.labels)
}

// All returns a snapshot of all labels.
func (r *RequestLabeler) All() []RequestLabel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]RequestLabel, 0, len(r.labels))
	for _, l := range r.labels {
		out = append(out, l)
	}
	return out
}
