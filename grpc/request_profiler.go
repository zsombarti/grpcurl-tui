package grpc

import (
	"sync"
	"time"
)

// ProfileEntry holds timing and size metrics for a single request.
type ProfileEntry struct {
	Method       string
	Duration     time.Duration
	RequestSize  int
	ResponseSize int
	Timestamp    time.Time
	Error        bool
}

// RequestProfiler collects per-method performance metrics.
type RequestProfiler struct {
	mu      sync.RWMutex
	entries []ProfileEntry
	maxSize int
}

// NewRequestProfiler creates a RequestProfiler with the given capacity.
func NewRequestProfiler(maxSize int) *RequestProfiler {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &RequestProfiler{
		entries: make([]ProfileEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Record appends a new ProfileEntry, evicting the oldest if at capacity.
func (p *RequestProfiler) Record(entry ProfileEntry) {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.entries) >= p.maxSize {
		p.entries = p.entries[1:]
	}
	p.entries = append(p.entries, entry)
}

// Len returns the number of recorded entries.
func (p *RequestProfiler) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.entries)
}

// All returns a snapshot of all entries.
func (p *RequestProfiler) All() []ProfileEntry {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]ProfileEntry, len(p.entries))
	copy(out, p.entries)
	return out
}

// Summary returns average duration and error rate across all entries.
func (p *RequestProfiler) Summary() (avgDuration time.Duration, errorRate float64) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if len(p.entries) == 0 {
		return 0, 0
	}
	var total time.Duration
	var errCount int
	for _, e := range p.entries {
		total += e.Duration
		if e.Error {
			errCount++
		}
	}
	avgDuration = total / time.Duration(len(p.entries))
	errorRate = float64(errCount) / float64(len(p.entries))
	return
}

// Clear removes all recorded entries.
func (p *RequestProfiler) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.entries = p.entries[:0]
}
