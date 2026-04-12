package grpc

import (
	"errors"
	"sync"
)

// EnricherFunc is a function that enriches a payload map with additional fields.
type EnricherFunc func(payload map[string]interface{}) (map[string]interface{}, error)

// RequestEnricher applies a chain of enrichment steps to a request payload.
type RequestEnricher struct {
	mu    sync.RWMutex
	steps []enricherStep
}

type enricherStep struct {
	name string
	fn   EnricherFunc
}

// NewRequestEnricher creates a new RequestEnricher with no steps.
func NewRequestEnricher() *RequestEnricher {
	return &RequestEnricher{}
}

// AddStep registers a named enrichment function.
func (e *RequestEnricher) AddStep(name string, fn EnricherFunc) error {
	if name == "" {
		return errors.New("enricher step name must not be empty")
	}
	if fn == nil {
		return errors.New("enricher step function must not be nil")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.steps = append(e.steps, enricherStep{name: name, fn: fn})
	return nil
}

// Len returns the number of registered enrichment steps.
func (e *RequestEnricher) Len() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.steps)
}

// Enrich applies all registered steps to the given payload in order.
// It returns the enriched payload or the first error encountered.
func (e *RequestEnricher) Enrich(payload map[string]interface{}) (map[string]interface{}, error) {
	if payload == nil {
		return nil, errors.New("payload must not be nil")
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	current := payload
	for _, step := range e.steps {
		result, err := step.fn(current)
		if err != nil {
			return nil, err
		}
		current = result
	}
	return current, nil
}

// StepNames returns the names of all registered steps in order.
func (e *RequestEnricher) StepNames() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	names := make([]string, len(e.steps))
	for i, s := range e.steps {
		names[i] = s.name
	}
	return names
}

// Clear removes all enrichment steps.
func (e *RequestEnricher) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.steps = nil
}
