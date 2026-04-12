package grpc

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// DeduplicatorPolicy controls deduplication window behaviour.
type DeduplicatorPolicy struct {
	WindowDuration time.Duration
	MaxSize        int
}

// DefaultDeduplicatorPolicy returns sensible defaults.
func DefaultDeduplicatorPolicy() DeduplicatorPolicy {
	return DeduplicatorPolicy{
		WindowDuration: 5 * time.Second,
		MaxSize:        256,
	}
}

type dedupEntry struct {
	key       string
	seenAt    time.Time
}

// RequestDeduplicator tracks recently seen request fingerprints and
// rejects duplicates that arrive within the configured window.
type RequestDeduplicator struct {
	mu     sync.Mutex
	policy DeduplicatorPolicy
	entries []dedupEntry
}

// NewRequestDeduplicator creates a RequestDeduplicator with the given policy.
// Falls back to defaults when the policy is zero-valued or invalid.
func NewRequestDeduplicator(p DeduplicatorPolicy) *RequestDeduplicator {
	def := DefaultDeduplicatorPolicy()
	if p.WindowDuration <= 0 {
		p.WindowDuration = def.WindowDuration
	}
	if p.MaxSize <= 0 {
		p.MaxSize = def.MaxSize
	}
	return &RequestDeduplicator{policy: p}
}

// fingerprint returns a SHA-256 hex digest of method+payload.
func fingerprint(method, payload string) string {
	h := sha256.Sum256([]byte(method + "\x00" + payload))
	return fmt.Sprintf("%x", h)
}

// IsDuplicate returns true when an identical (method, payload) pair was
// already seen within the deduplication window.
func (d *RequestDeduplicator) IsDuplicate(method, payload string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict()
	key := fingerprint(method, payload)
	for _, e := range d.entries {
		if e.key == key {
			return true
		}
	}
	return false
}

// Record stores a fingerprint so future calls to IsDuplicate can detect it.
func (d *RequestDeduplicator) Record(method, payload string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict()
	if len(d.entries) >= d.policy.MaxSize {
		d.entries = d.entries[1:]
	}
	d.entries = append(d.entries, dedupEntry{key: fingerprint(method, payload), seenAt: time.Now()})
}

// Len returns the number of active (non-expired) fingerprints.
func (d *RequestDeduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict()
	return len(d.entries)
}

// evict removes entries older than the window. Must be called with mu held.
func (d *RequestDeduplicator) evict() {
	cutoff := time.Now().Add(-d.policy.WindowDuration)
	start := 0
	for start < len(d.entries) && d.entries[start].seenAt.Before(cutoff) {
		start++
	}
	d.entries = d.entries[start:]
}
