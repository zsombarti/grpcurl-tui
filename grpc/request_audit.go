package grpc

import (
	"sync"
	"time"
)

// AuditEntry records a single auditable gRPC request event.
type AuditEntry struct {
	Timestamp time.Time
	Method    string
	Address   string
	Status    string // "ok" | "error"
	Duration  time.Duration
	Note      string
}

// RequestAuditor maintains an append-only audit log of gRPC calls.
type RequestAuditor struct {
	mu      sync.RWMutex
	entries []AuditEntry
	maxSize int
}

const defaultAuditMaxSize = 500

// NewRequestAuditor creates a RequestAuditor with the given max log size.
// If maxSize <= 0 the default (500) is used.
func NewRequestAuditor(maxSize int) *RequestAuditor {
	if maxSize <= 0 {
		maxSize = defaultAuditMaxSize
	}
	return &RequestAuditor{maxSize: maxSize}
}

// Record appends an AuditEntry, evicting the oldest entry when full.
func (a *RequestAuditor) Record(entry AuditEntry) {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.entries) >= a.maxSize {
		a.entries = a.entries[1:]
	}
	a.entries = append(a.entries, entry)
}

// Entries returns a shallow copy of all audit entries.
func (a *RequestAuditor) Entries() []AuditEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]AuditEntry, len(a.entries))
	copy(out, a.entries)
	return out
}

// Len returns the number of recorded entries.
func (a *RequestAuditor) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.entries)
}

// Clear removes all audit entries.
func (a *RequestAuditor) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = nil
}
