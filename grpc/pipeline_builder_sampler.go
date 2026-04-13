package grpc

import (
	"errors"
	"fmt"
)

// SamplerStep wraps a RequestSampler as a pipeline step.
type SamplerStep struct {
	name    string
	sampler *RequestSampler
	method  string
}

// NewSamplerStep creates a pipeline step that samples payloads passing through.
func NewSamplerStep(name, method string, sampler *RequestSampler) (*SamplerStep, error) {
	if name == "" {
		return nil, errors.New("sampler step: name must not be empty")
	}
	if method == "" {
		return nil, errors.New("sampler step: method must not be empty")
	}
	if sampler == nil {
		return nil, errors.New("sampler step: sampler must not be nil")
	}
	return &SamplerStep{name: name, sampler: sampler, method: method}, nil
}

// Name returns the step name.
func (s *SamplerStep) Name() string { return s.name }

// Run samples the payload and passes it through unchanged.
func (s *SamplerStep) Run(payload map[string]interface{}) (map[string]interface{}, error) {
	if payload == nil {
		return nil, fmt.Errorf("sampler step %q: nil payload", s.name)
	}
	s.sampler.Sample(s.method, payload)
	return payload, nil
}
