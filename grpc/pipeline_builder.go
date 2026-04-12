package grpc

import (
	"context"
	"fmt"
	"strings"
)

// PipelineBuilder provides a fluent API for constructing a RequestPipeline
// with common pre-built steps.
type PipelineBuilder struct {
	pipeline *RequestPipeline
}

// NewPipelineBuilder returns a PipelineBuilder backed by a fresh RequestPipeline.
func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{pipeline: NewRequestPipeline()}
}

// WithEnvSubstitution adds a step that replaces ${VAR} tokens in string values
// using the provided substitutor.
func (b *PipelineBuilder) WithEnvSubstitution(sub *EnvSubstitutor) *PipelineBuilder {
	b.pipeline.AddStep(PipelineStep{
		Name: "env-substitution",
		Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) {
			for k, v := range m {
				if s, ok := v.(string); ok {
					m[k] = sub.Substitute(s)
				}
			}
			return m, nil
		},
	})
	return b
}

// WithKeyNormalisation adds a step that lowercases all top-level keys.
func (b *PipelineBuilder) WithKeyNormalisation() *PipelineBuilder {
	b.pipeline.AddStep(PipelineStep{
		Name: "key-normalisation",
		Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) {
			norm := make(map[string]interface{}, len(m))
			for k, v := range m {
				norm[strings.ToLower(k)] = v
			}
			return norm, nil
		},
	})
	return b
}

// WithRequiredFields adds a step that returns an error if any of the specified
// field names are absent from the payload.
func (b *PipelineBuilder) WithRequiredFields(fields ...string) *PipelineBuilder {
	b.pipeline.AddStep(PipelineStep{
		Name: "required-fields",
		Handler: func(_ context.Context, m map[string]interface{}) (map[string]interface{}, error) {
			for _, f := range fields {
				if _, ok := m[f]; !ok {
					return nil, fmt.Errorf("required field %q is missing", f)
				}
			}
			return m, nil
		},
	})
	return b
}

// Build returns the constructed RequestPipeline.
func (b *PipelineBuilder) Build() *RequestPipeline {
	return b.pipeline
}
