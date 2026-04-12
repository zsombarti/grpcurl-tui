package grpc

import (
	"errors"
	"sync"
	"time"
)

const defaultRequestCacheMaxSize = 128

// RequestCacheEntry holds a cached response payload with metadata.
type RequestCacheEntry struct {
	Method    string
	Payload   string
	Response  string
	CachedAt  time.Time
	ExpiresAt time.Time
}

// RequestCache stores responses keyed by method+payload fingerprint.
type RequestCache struct {
	mu      sync.RWMutex
	entries map[string]*RequestCacheEntry
	ttl     time.Duration
	maxSize int
	order   []string
}

// NewRequestCache creates a RequestCache with the given TTL and max size.
// If ttl <= 0 it defaults to 30 seconds; if maxSize <= 0 it defaults to 128.
func NewRequestCache(ttl time.Duration, maxSize int) *RequestCache {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	if maxSize <= 0 {
		maxSize = defaultRequestCacheMaxSize
	}
	return &RequestCache{
		entries: make(map[string]*RequestCacheEntry),
		ttl:     ttl,
		maxSize: maxSize,
	}
}

func cacheKey(method, payload string) string {
	return method + "::" + payload
}

// Set stores a response in the cache for the given method and payload.
func (c *RequestCache) Set(method, payload, response string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := cacheKey(method, payload)
	if _, exists := c.entries[key]; !exists {
		if len(c.order) >= c.maxSize {
			oldest := c.order[0]
			c.order = c.order[1:]
			delete(c.entries, oldest)
		}
		c.order = append(c.order, key)
	}
	now := time.Now()
	c.entries[key] = &RequestCacheEntry{
		Method:    method,
		Payload:   payload,
		Response:  response,
		CachedAt:  now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Get retrieves a cached response. Returns error if not found or expired.
func (c *RequestCache) Get(method, payload string) (*RequestCacheEntry, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[cacheKey(method, payload)]
	if !ok {
		return nil, errors.New("cache miss")
	}
	if time.Now().After(entry.ExpiresAt) {
		return nil, errors.New("cache entry expired")
	}
	return entry, nil
}

// Len returns the number of entries currently in the cache.
func (c *RequestCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// Clear removes all entries from the cache.
func (c *RequestCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*RequestCacheEntry)
	c.order = nil
}
