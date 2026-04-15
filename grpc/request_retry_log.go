package grpc

import (
	"errors"
	"sync"
	"time"
)

const defaultRetryLogMaxSize = 200

// RetryLogEntry records a single retry attempt.
type RetryLogEntry struct {
	Method    string
	Attempt   int
	Err       error
	Timestamp time.Time
}

// RetryLog stores a bounded history of retry attempts.
type RetryLog struct {
	mu      sync.RWMutex
	entries []RetryLogEntry
	maxSize int
}

// NewRetryLog creates a RetryLog with the given capacity (0 → default).
func NewRetryLog(maxSize int) *RetryLog {
	if maxSize <= 0 {
		maxSize = defaultRetryLogMaxSize
	}
	return &RetryLog{maxSize: maxSize}
}

// Record appends a retry attempt entry.
func (r *RetryLog) Record(method string, attempt int, err error) error {
	if method == "" {
		return errors.New("retry log: method must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.entries) >= r.maxSize {
		r.entries = r.entries[1:]
	}
	r.entries = append(r.entries, RetryLogEntry{
		Method:    method,
		Attempt:   attempt,
		Err:       err,
		Timestamp: time.Now(),
	})
	return nil
}

// Len returns the number of recorded entries.
func (r *RetryLog) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}

// Entries returns a shallow copy of all entries.
func (r *RetryLog) Entries() []RetryLogEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]RetryLogEntry, len(r.entries))
	copy(out, r.entries)
	return out
}

// Clear removes all entries.
func (r *RetryLog) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = nil
}
