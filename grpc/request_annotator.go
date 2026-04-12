package grpc

import (
	"errors"
	"sync"
	"time"
)

const defaultAnnotatorMaxSize = 256

// Annotation holds a key-value note attached to a request.
type Annotation struct {
	Method    string
	Key       string
	Value     string
	CreatedAt time.Time
}

// RequestAnnotator stores free-form annotations keyed by method + key.
type RequestAnnotator struct {
	mu      sync.RWMutex
	entries []Annotation
	maxSize int
}

// NewRequestAnnotator returns a new RequestAnnotator with the given capacity.
// If maxSize <= 0 the default is used.
func NewRequestAnnotator(maxSize int) *RequestAnnotator {
	if maxSize <= 0 {
		maxSize = defaultAnnotatorMaxSize
	}
	return &RequestAnnotator{maxSize: maxSize}
}

// Annotate adds an annotation for the given method and key.
func (a *RequestAnnotator) Annotate(method, key, value string) error {
	if method == "" {
		return errors.New("method must not be empty")
	}
	if key == "" {
		return errors.New("key must not be empty")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	// Overwrite existing entry with same method+key.
	for i, e := range a.entries {
		if e.Method == method && e.Key == key {
			a.entries[i].Value = value
			a.entries[i].CreatedAt = time.Now()
			return nil
		}
	}
	if len(a.entries) >= a.maxSize {
		a.entries = a.entries[1:]
	}
	a.entries = append(a.entries, Annotation{
		Method:    method,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	})
	return nil
}

// Get returns the annotation value for the given method and key.
func (a *RequestAnnotator) Get(method, key string) (string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, e := range a.entries {
		if e.Method == method && e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}

// All returns a copy of all annotations.
func (a *RequestAnnotator) All() []Annotation {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]Annotation, len(a.entries))
	copy(out, a.entries)
	return out
}

// Len returns the number of stored annotations.
func (a *RequestAnnotator) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.entries)
}

// Clear removes all annotations.
func (a *RequestAnnotator) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = nil
}
