package grpc

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ScheduledRequest holds a request to be executed at a specific interval.
type ScheduledRequest struct {
	ID       string
	Address  string
	Method   string
	Payload  string
	Interval time.Duration
	CreatedAt time.Time
}

// SchedulerPolicy defines the configuration for the request scheduler.
type SchedulerPolicy struct {
	MaxJobs     int
	DefaultInterval time.Duration
}

// DefaultSchedulerPolicy returns sensible defaults.
func DefaultSchedulerPolicy() SchedulerPolicy {
	return SchedulerPolicy{
		MaxJobs:         10,
		DefaultInterval: 5 * time.Second,
	}
}

// RequestScheduler manages periodic gRPC request execution.
type RequestScheduler struct {
	mu      sync.Mutex
	policy  SchedulerPolicy
	jobs    map[string]context.CancelFunc
	requests map[string]ScheduledRequest
}

// NewRequestScheduler creates a new RequestScheduler with the given policy.
func NewRequestScheduler(policy SchedulerPolicy) *RequestScheduler {
	if policy.MaxJobs <= 0 {
		policy = DefaultSchedulerPolicy()
	}
	if policy.DefaultInterval <= 0 {
		policy.DefaultInterval = DefaultSchedulerPolicy().DefaultInterval
	}
	return &RequestScheduler{
		policy:   policy,
		jobs:     make(map[string]context.CancelFunc),
		requests: make(map[string]ScheduledRequest),
	}
}

// Add schedules a new request. Returns error if max jobs reached or ID is duplicate.
func (s *RequestScheduler) Add(req ScheduledRequest, fn func(ScheduledRequest)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req.ID == "" {
		return errors.New("scheduler: request ID must not be empty")
	}
	if _, exists := s.jobs[req.ID]; exists {
		return errors.New("scheduler: job with this ID already exists")
	}
	if len(s.jobs) >= s.policy.MaxJobs {
		return errors.New("scheduler: max jobs reached")
	}
	if req.Interval <= 0 {
		req.Interval = s.policy.DefaultInterval
	}
	req.CreatedAt = time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	s.jobs[req.ID] = cancel
	s.requests[req.ID] = req
	go func() {
		ticker := time.NewTicker(req.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fn(req)
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

// Remove stops and removes a scheduled job by ID.
func (s *RequestScheduler) Remove(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cancel, ok := s.jobs[id]
	if !ok {
		return false
	}
	cancel()
	delete(s.jobs, id)
	delete(s.requests, id)
	return true
}

// List returns a snapshot of all scheduled requests.
func (s *RequestScheduler) List() []ScheduledRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]ScheduledRequest, 0, len(s.requests))
	for _, r := range s.requests {
		out = append(out, r)
	}
	return out
}

// Len returns the number of active scheduled jobs.
func (s *RequestScheduler) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.jobs)
}

// StopAll cancels all running jobs.
func (s *RequestScheduler) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, cancel := range s.jobs {
		cancel()
		delete(s.jobs, id)
		delete(s.requests, id)
	}
}
