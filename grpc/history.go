package grpc

import (
	"sync"
	"time"
)

// HistoryEntry represents a single recorded gRPC call.
type HistoryEntry struct {
	Timestamp time.Time
	Address   string
	Service   string
	Method    string
	Request   string
	Response  string
	Error     string
}

// History stores a bounded list of past gRPC invocations.
type History struct {
	mu      sync.RWMutex
	entries []HistoryEntry
	maxSize int
}

// NewHistory creates a History with the given maximum capacity.
func NewHistory(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &History{
		entries: make([]HistoryEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add appends an entry, evicting the oldest if capacity is exceeded.
func (h *History) Add(entry HistoryEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	if len(h.entries) >= h.maxSize {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, entry)
}

// All returns a snapshot of all entries in chronological order.
func (h *History) All() []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	snap := make([]HistoryEntry, len(h.entries))
	copy(snap, h.entries)
	return snap
}

// Clear removes all stored entries.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = h.entries[:0]
}

// Len returns the current number of stored entries.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.entries)
}
