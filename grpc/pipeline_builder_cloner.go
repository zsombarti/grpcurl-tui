package grpc

import (
	"context"
	"fmt"
)

// ClonerStep is a PipelineStep that deep-clones the payload before passing
// it downstream, preventing mutation side-effects in multi-step pipelines.
type ClonerStep struct {
	name   string
	cloner *RequestCloner
}

// NewClonerStep returns a PipelineStep that clones payloads using the
// supplied RequestCloner. If cloner is nil a default one is used.
func NewClonerStep(cloner *RequestCloner) *ClonerStep {
	if cloner == nil {
		cloner = NewRequestCloner(DefaultClonerPolicy())
	}
	return &ClonerStep{name: "clone", cloner: cloner}
}

// Name returns the step identifier used in pipeline logs.
func (s *ClonerStep) Name() string { return s.name }

// Run clones the payload and returns the copy.
func (s *ClonerStep) Run(ctx context.Context, payload map[string]any) (map[string]any, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	cloned, err := s.cloner.Clone(payload)
	if err != nil {
		return nil, fmt.Errorf("cloner_step %q: %w", s.name, err)
	}
	return cloned, nil
}

// WithName returns a new ClonerStep with the given name.
func (s *ClonerStep) WithName(name string) *ClonerStep {
	return &ClonerStep{name: name, cloner: s.cloner}
}
