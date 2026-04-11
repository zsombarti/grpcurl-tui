package grpc

import (
	"errors"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker.
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreakerPolicy holds configuration for the circuit breaker.
type CircuitBreakerPolicy struct {
	MaxFailures  int
	OpenDuration time.Duration
}

// DefaultCircuitBreakerPolicy returns sensible defaults.
func DefaultCircuitBreakerPolicy() CircuitBreakerPolicy {
	return CircuitBreakerPolicy{
		MaxFailures:  5,
		OpenDuration: 10 * time.Second,
	}
}

// CircuitBreaker guards calls and trips open after too many failures.
type CircuitBreaker struct {
	mu       sync.Mutex
	policy   CircuitBreakerPolicy
	failures int
	state    CircuitState
	openedAt time.Time
}

// NewCircuitBreaker creates a CircuitBreaker with the given policy.
// Falls back to defaults if MaxFailures <= 0 or OpenDuration <= 0.
func NewCircuitBreaker(p CircuitBreakerPolicy) *CircuitBreaker {
	def := DefaultCircuitBreakerPolicy()
	if p.MaxFailures <= 0 {
		p.MaxFailures = def.MaxFailures
	}
	if p.OpenDuration <= 0 {
		p.OpenDuration = def.OpenDuration
	}
	return &CircuitBreaker{policy: p, state: StateClosed}
}

// Allow returns an error if the circuit is open and not yet ready to probe.
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case StateOpen:
		if time.Since(cb.openedAt) >= cb.policy.OpenDuration {
			cb.state = StateHalfOpen
			return nil
		}
		return errors.New("circuit breaker is open")
	default:
		return nil
	}
}

// RecordSuccess resets failure count and closes the circuit.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = StateClosed
}

// RecordFailure increments failure count and may trip the circuit open.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.failures >= cb.policy.MaxFailures {
		cb.state = StateOpen
		cb.openedAt = time.Now()
	}
}

// State returns the current circuit state.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Failures returns the current failure count.
func (cb *CircuitBreaker) Failures() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.failures
}
