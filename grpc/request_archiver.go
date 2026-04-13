package grpc

import (
	"errors"
	"sync"
	"time"
)

const defaultArchiverMaxSize = 200

// ArchiveEntry holds a single archived request payload with metadata.
type ArchiveEntry struct {
	Method    string
	Payload   map[string]interface{}
	Tag       string
	ArchivedAt time.Time
}

// RequestArchiver stores request payloads for long-term reference.
type RequestArchiver struct {
	mu      sync.RWMutex
	entries []ArchiveEntry
	maxSize int
}

// NewRequestArchiver creates a new RequestArchiver with the given max size.
// If maxSize <= 0, defaultArchiverMaxSize is used.
func NewRequestArchiver(maxSize int) *RequestArchiver {
	if maxSize <= 0 {
		maxSize = defaultArchiverMaxSize
	}
	return &RequestArchiver{
		entries: make([]ArchiveEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Archive saves a request payload under the given method and optional tag.
func (a *RequestArchiver) Archive(method string, payload map[string]interface{}, tag string) error {
	if method == "" {
		return errors.New("request archiver: method must not be empty")
	}
	if payload == nil {
		return errors.New("request archiver: payload must not be nil")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.entries) >= a.maxSize {
		a.entries = a.entries[1:]
	}
	a.entries = append(a.entries, ArchiveEntry{
		Method:     method,
		Payload:    payload,
		Tag:        tag,
		ArchivedAt: time.Now(),
	})
	return nil
}

// Len returns the number of archived entries.
func (a *RequestArchiver) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.entries)
}

// All returns a copy of all archived entries.
func (a *RequestArchiver) All() []ArchiveEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]ArchiveEntry, len(a.entries))
	copy(result, a.entries)
	return result
}

// Clear removes all archived entries.
func (a *RequestArchiver) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = a.entries[:0]
}
