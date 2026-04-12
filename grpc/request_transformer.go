package grpc

import (
	"errors"
	"strings"
)

// TransformFn is a function that transforms a JSON payload map.
type TransformFn func(payload map[string]interface{}) (map[string]interface{}, error)

// RequestTransformer applies a chain of named transform functions to a payload.
type RequestTransformer struct {
	steps []transformStep
}

type transformStep struct {
	name string
	fn   TransformFn
}

// NewRequestTransformer creates a new RequestTransformer with no steps.
func NewRequestTransformer() *RequestTransformer {
	return &RequestTransformer{}
}

// AddStep appends a named transform step. Returns an error if name is empty or fn is nil.
func (t *RequestTransformer) AddStep(name string, fn TransformFn) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("transform step name must not be empty")
	}
	if fn == nil {
		return errors.New("transform function must not be nil")
	}
	t.steps = append(t.steps, transformStep{name: name, fn: fn})
	return nil
}

// Len returns the number of registered transform steps.
func (t *RequestTransformer) Len() int {
	return len(t.steps)
}

// Names returns the ordered list of step names.
func (t *RequestTransformer) Names() []string {
	names := make([]string, len(t.steps))
	for i, s := range t.steps {
		names[i] = s.name
	}
	return names
}

// Transform applies all registered steps in order to the given payload.
// It returns the transformed payload or the first error encountered.
func (t *RequestTransformer) Transform(payload map[string]interface{}) (map[string]interface{}, error) {
	if payload == nil {
		return nil, errors.New("payload must not be nil")
	}
	current := payload
	for _, step := range t.steps {
		result, err := step.fn(current)
		if err != nil {
			return nil, err
		}
		current = result
	}
	return current, nil
}

// Clear removes all registered transform steps.
func (t *RequestTransformer) Clear() {
	t.steps = nil
}
