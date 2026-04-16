package grpc

import (
	"errors"
	"sync"
	"time"
)

// DefaultQuotaPolicy returns sensible quota defaults.
func DefaultQuotaPolicy() QuotaPolicy {
	return QuotaPolicy{
		MaxRequests: 100,
		WindowSize:  time.Minute,
	}
}

// QuotaPolicy defines the quota window and limit.
type QuotaPolicy struct {
	MaxRequests int
	WindowSize  time.Duration
}

// QuotaEntry records a method's usage within the current window.
type QuotaEntry struct {
	Method    string
	Count     int
	WindowEnd time.Time
}

// RequestQuota tracks per-method request quotas.
type RequestQuota struct {
	mu     sync.Mutex
	policy QuotaPolicy
	usage  map[string]*QuotaEntry
}

// NewRequestQuota creates a RequestQuota with the given policy.
func NewRequestQuota(p QuotaPolicy) *RequestQuota {
	if p.MaxRequests <= 0 || p.WindowSize <= 0 {
		p = DefaultQuotaPolicy()
	}
	return &RequestQuota{
		policy: p,
		usage:  make(map[string]*QuotaEntry),
	}
}

// Allow checks whether the method is within quota and increments the counter.
func (q *RequestQuota) Allow(method string) error {
	if method == "" {
		return errors.New("quota: method must not be empty")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	now := time.Now()
	entry, ok := q.usage[method]
	if !ok || now.After(entry.WindowEnd) {
		q.usage[method] = &QuotaEntry{
			Method:    method,
			Count:     1,
			WindowEnd: now.Add(q.policy.WindowSize),
		}
		return nil
	}
	if entry.Count >= q.policy.MaxRequests {
		return errors.New("quota: limit exceeded for method " + method)
	}
	entry.Count++
	return nil
}

// Usage returns a snapshot of current quota usage.
func (q *RequestQuota) Usage() []QuotaEntry {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]QuotaEntry, 0, len(q.usage))
	for _, e := range q.usage {
		out = append(out, *e)
	}
	return out
}

// Reset clears all quota counters.
func (q *RequestQuota) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.usage = make(map[string]*QuotaEntry)
}

// Len returns the number of tracked methods.
func (q *RequestQuota) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.usage)
}
