package grpc

import (
	"errors"
	"time"
)

// TimeoutPolicy defines configurable timeout settings for gRPC calls.
type TimeoutPolicy struct {
	DialTimeout    time.Duration
	RequestTimeout time.Duration
}

// DefaultTimeoutPolicy returns a TimeoutPolicy with sensible defaults.
func DefaultTimeoutPolicy() TimeoutPolicy {
	return TimeoutPolicy{
		DialTimeout:    5 * time.Second,
		RequestTimeout: 30 * time.Second,
	}
}

// TimeoutManager manages and validates timeout policies.
type TimeoutManager struct {
	policy TimeoutPolicy
}

// NewTimeoutManager creates a new TimeoutManager with the given policy.
func NewTimeoutManager(policy TimeoutPolicy) (*TimeoutManager, error) {
	if policy.DialTimeout <= 0 {
		return nil, errors.New("dial timeout must be positive")
	}
	if policy.RequestTimeout <= 0 {
		return nil, errors.New("request timeout must be positive")
	}
	return &TimeoutManager{policy: policy}, nil
}

// DialTimeout returns the configured dial timeout.
func (t *TimeoutManager) DialTimeout() time.Duration {
	return t.policy.DialTimeout
}

// RequestTimeout returns the configured request timeout.
func (t *TimeoutManager) RequestTimeout() time.Duration {
	return t.policy.RequestTimeout
}

// Policy returns the current TimeoutPolicy.
func (t *TimeoutManager) Policy() TimeoutPolicy {
	return t.policy
}
