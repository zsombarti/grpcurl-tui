package grpc

import (
	"errors"
	"sync"
	"time"
)

const defaultAggregatorMaxSize = 256

// AggregatedResult holds the merged result of multiple request payloads.
type AggregatedResult struct {
	Method    string
	Payloads  []map[string]interface{}
	Merged    map[string]interface{}
	CreatedAt time.Time
}

// RequestAggregator merges multiple request payloads into a single result.
type RequestAggregator struct {
	mu      sync.Mutex
	results []AggregatedResult
	maxSize int
}

// NewRequestAggregator creates a new RequestAggregator with the given max size.
// If maxSize <= 0, the default is used.
func NewRequestAggregator(maxSize int) *RequestAggregator {
	if maxSize <= 0 {
		maxSize = defaultAggregatorMaxSize
	}
	return &RequestAggregator{maxSize: maxSize}
}

// Aggregate merges the provided payloads for the given method and stores the result.
// Keys from later payloads overwrite earlier ones on conflict.
func (a *RequestAggregator) Aggregate(method string, payloads []map[string]interface{}) (AggregatedResult, error) {
	if method == "" {
		return AggregatedResult{}, errors.New("method must not be empty")
	}
	if len(payloads) == 0 {
		return AggregatedResult{}, errors.New("at least one payload required")
	}
	merged := make(map[string]interface{})
	for _, p := range payloads {
		for k, v := range p {
			merged[k] = v
		}
	}
	result := AggregatedResult{
		Method:    method,
		Payloads:  payloads,
		Merged:    merged,
		CreatedAt: time.Now(),
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.results) >= a.maxSize {
		a.results = a.results[1:]
	}
	a.results = append(a.results, result)
	return result, nil
}

// Len returns the number of stored aggregated results.
func (a *RequestAggregator) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.results)
}

// All returns a copy of all stored aggregated results.
func (a *RequestAggregator) All() []AggregatedResult {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]AggregatedResult, len(a.results))
	copy(out, a.results)
	return out
}

// Clear removes all stored results.
func (a *RequestAggregator) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.results = nil
}
