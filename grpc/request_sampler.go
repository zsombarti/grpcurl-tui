package grpc

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// DefaultSamplerPolicy returns a sensible default sampling policy.
func DefaultSamplerPolicy() SamplerPolicy {
	return SamplerPolicy{
		Rate:    1.0, // 100% sampling by default
		MaxSize: 256,
	}
}

// SamplerPolicy controls how requests are sampled.
type SamplerPolicy struct {
	Rate    float64 // 0.0–1.0
	MaxSize int
}

// SampledRequest holds a sampled request record.
type SampledRequest struct {
	Method    string
	Payload   map[string]interface{}
	Timestamp time.Time
}

// RequestSampler probabilistically samples requests for inspection.
type RequestSampler struct {
	mu      sync.Mutex
	policy  SamplerPolicy
	samples []SampledRequest
	rng     *rand.Rand
}

// NewRequestSampler creates a new RequestSampler with the given policy.
func NewRequestSampler(policy SamplerPolicy) (*RequestSampler, error) {
	if policy.Rate < 0 || policy.Rate > 1.0 {
		policy = DefaultSamplerPolicy()
	}
	if policy.MaxSize <= 0 {
		policy = DefaultSamplerPolicy()
	}
	return &RequestSampler{
		policy:  policy,
		samples: make([]SampledRequest, 0, policy.MaxSize),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Sample probabilistically records the request based on the configured rate.
func (s *RequestSampler) Sample(method string, payload map[string]interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.rng.Float64() > s.policy.Rate {
		return false
	}
	rec := SampledRequest{
		Method:    method,
		Payload:   payload,
		Timestamp: time.Now(),
	}
	if len(s.samples) >= s.policy.MaxSize {
		s.samples = s.samples[1:]
	}
	s.samples = append(s.samples, rec)
	return true
}

// Len returns the number of stored samples.
func (s *RequestSampler) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.samples)
}

// All returns a copy of all sampled requests.
func (s *RequestSampler) All() []SampledRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]SampledRequest, len(s.samples))
	copy(out, s.samples)
	return out
}

// Clear removes all stored samples.
func (s *RequestSampler) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) == 0 {
		return errors.New("sampler: nothing to clear")
	}
	s.samples = s.samples[:0]
	return nil
}

// Policy returns the active sampler policy.
func (s *RequestSampler) Policy() SamplerPolicy {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.policy
}
