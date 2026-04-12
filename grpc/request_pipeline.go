package grpc

import (
	"context"
	"fmt"
)

// PipelineStep represents a single processing step in a request pipeline.
type PipelineStep struct {
	Name    string
	Handler func(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error)
}

// RequestPipeline executes a sequence of steps on a request payload.
type RequestPipeline struct {
	steps []PipelineStep
}

// NewRequestPipeline creates an empty RequestPipeline.
func NewRequestPipeline() *RequestPipeline {
	return &RequestPipeline{
		steps: make([]PipelineStep, 0),
	}
}

// AddStep appends a step to the pipeline.
func (p *RequestPipeline) AddStep(step PipelineStep) {
	p.steps = append(p.steps, step)
}

// Len returns the number of steps in the pipeline.
func (p *RequestPipeline) Len() int {
	return len(p.steps)
}

// Run executes all pipeline steps in order, passing the output of each step
// as the input to the next. Returns the final transformed payload.
func (p *RequestPipeline) Run(ctx context.Context, initial map[string]interface{}) (map[string]interface{}, error) {
	payload := initial
	for _, step := range p.steps {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("pipeline cancelled before step %q: %w", step.Name, ctx.Err())
		default:
		}
		var err error
		payload, err = step.Handler(ctx, payload)
		if err != nil {
			return nil, fmt.Errorf("pipeline step %q failed: %w", step.Name, err)
		}
	}
	return payload, nil
}

// Clear removes all steps from the pipeline.
func (p *RequestPipeline) Clear() {
	p.steps = p.steps[:0]
}

// StepNames returns the names of all registered steps in order.
func (p *RequestPipeline) StepNames() []string {
	names := make([]string, len(p.steps))
	for i, s := range p.steps {
		names[i] = s.Name
	}
	return names
}
