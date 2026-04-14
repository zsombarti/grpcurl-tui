package grpc

import (
	"fmt"
	"sync"
	"time"
)

// CircuitLogEntry records a single circuit breaker state transition.
type CircuitLogEntry struct {
	Method    string
	FromState string
	ToState   string
	Reason    string
	Timestamp time.Time
}

// CircuitLog maintains an ordered, bounded log of circuit breaker transitions
// for display and audit purposes.
type CircuitLog struct {
	mu      sync.RWMutex
	entries []CircuitLogEntry
	maxSize int
}

// NewCircuitLog creates a CircuitLog with the given maximum number of entries.
// If maxSize is less than 1 it defaults to 100.
func NewCircuitLog(maxSize int) *CircuitLog {
	if maxSize < 1 {
		maxSize = 100
	}
	return &CircuitLog{
		entries: make([]CircuitLogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Record appends a state transition entry. If the log is full the oldest
// entry is evicted to make room.
func (cl *CircuitLog) Record(method, from, to, reason string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	entry := CircuitLogEntry{
		Method:    method,
		FromState: from,
		ToState:   to,
		Reason:    reason,
		Timestamp: time.Now(),
	}

	if len(cl.entries) >= cl.maxSize {
		cl.entries = cl.entries[1:]
	}
	cl.entries = append(cl.entries, entry)
}

// Len returns the current number of log entries.
func (cl *CircuitLog) Len() int {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return len(cl.entries)
}

// Entries returns a snapshot of all current log entries, newest last.
func (cl *CircuitLog) Entries() []CircuitLogEntry {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	snap := make([]CircuitLogEntry, len(cl.entries))
	copy(snap, cl.entries)
	return snap
}

// Clear removes all entries from the log.
func (cl *CircuitLog) Clear() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.entries = cl.entries[:0]
}

// Summary returns a human-readable summary of the most recent transition for
// the given method, or an empty string if no entry exists.
func (cl *CircuitLog) Summary(method string) string {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	for i := len(cl.entries) - 1; i >= 0; i-- {
		e := cl.entries[i]
		if e.Method == method {
			return fmt.Sprintf("%s → %s (%s) at %s",
				e.FromState, e.ToState, e.Reason,
				e.Timestamp.Format(time.RFC3339))
		}
	}
	return ""
}
