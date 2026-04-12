package grpc

import (
	"errors"
	"time"
)

// SchedulerStep wraps a RequestScheduler as a pipeline step that registers
// the current payload as a new scheduled job.
type SchedulerStep struct {
	scheduler *RequestScheduler
	jobID     string
	address   string
	method    string
	interval  time.Duration
	handler   func(ScheduledRequest)
}

// NewSchedulerStep creates a pipeline step that schedules the payload for
// periodic replay via the provided scheduler.
func NewSchedulerStep(
	scheduler *RequestScheduler,
	jobID, address, method string,
	interval time.Duration,
	handler func(ScheduledRequest),
) (*SchedulerStep, error) {
	if scheduler == nil {
		return nil, errors.New("scheduler step: scheduler must not be nil")
	}
	if jobID == "" {
		return nil, errors.New("scheduler step: jobID must not be empty")
	}
	if handler == nil {
		return nil, errors.New("scheduler step: handler must not be nil")
	}
	return &SchedulerStep{
		scheduler: scheduler,
		jobID:     jobID,
		address:   address,
		method:    method,
		interval:  interval,
		handler:   handler,
	}, nil
}

// Run registers the payload as a scheduled job and passes it through unchanged.
func (s *SchedulerStep) Run(payload string) (string, error) {
	req := ScheduledRequest{
		ID:      s.jobID,
		Address: s.address,
		Method:  s.method,
		Payload: payload,
		Interval: s.interval,
	}
	if err := s.scheduler.Add(req, s.handler); err != nil {
		return payload, err
	}
	return payload, nil
}

// Name returns the identifier for this pipeline step.
func (s *SchedulerStep) Name() string {
	return "scheduler:" + s.jobID
}
