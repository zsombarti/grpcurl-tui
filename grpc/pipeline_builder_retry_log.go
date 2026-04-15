package grpc

import (
	"context"
	"fmt"
)

// RetryLogStep is a pipeline step that records each pass through the
// pipeline as a retry attempt in a RetryLog.
type RetryLogStep struct {
	log     *RetryLog
	method  string
	attempt int
}

// NewRetryLogStep creates a pipeline step that logs every execution as
// attempt #N for the given method name.
func NewRetryLogStep(log *RetryLog, method string) (*RetryLogStep, error) {
	if log == nil {
		return nil, fmt.Errorf("retry log step: log must not be nil")
	}
	if method == "" {
		return nil, fmt.Errorf("retry log step: method must not be empty")
	}
	return &RetryLogStep{log: log, method: method}, nil
}

// Name satisfies the PipelineStep interface.
func (s *RetryLogStep) Name() string { return "retry-log:" + s.method }

// Run records the attempt and forwards the payload unchanged.
func (s *RetryLogStep) Run(_ context.Context, payload map[string]any) (map[string]any, error) {
	s.attempt++
	_ = s.log.Record(s.method, s.attempt, nil)
	return payload, nil
}
