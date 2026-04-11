package grpc

import (
	"sync"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// SchemaCacheEntry holds a cached service descriptor with its expiry time.
type SchemaCacheEntry struct {
	Descriptor protoreflect.ServiceDescriptor
	CachedAt   time.Time
	TTL        time.Duration
}

// IsExpired reports whether the cache entry has passed its TTL.
func (e SchemaCacheEntry) IsExpired() bool {
	return time.Since(e.CachedAt) > e.TTL
}

// SchemaCache is a thread-safe in-memory cache for protobuf service descriptors.
type SchemaCache struct {
	mu      sync.RWMutex
	entries map[string]SchemaCacheEntry
	ttl     time.Duration
}

// NewSchemaCache creates a SchemaCache with the given default TTL.
// If ttl is zero or negative, a default of 5 minutes is used.
func NewSchemaCache(ttl time.Duration) *SchemaCache {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &SchemaCache{
		entries: make(map[string]SchemaCacheEntry),
		ttl:     ttl,
	}
}

// Set stores a service descriptor under the given key.
func (c *SchemaCache) Set(key string, desc protoreflect.ServiceDescriptor) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = SchemaCacheEntry{
		Descriptor: desc,
		CachedAt:   time.Now(),
		TTL:        c.ttl,
	}
}

// Get retrieves a service descriptor by key. Returns nil and false if not found or expired.
func (c *SchemaCache) Get(key string) (protoreflect.ServiceDescriptor, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok || entry.IsExpired() {
		return nil, false
	}
	return entry.Descriptor, true
}

// Invalidate removes a single entry from the cache.
func (c *SchemaCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Flush removes all entries from the cache.
func (c *SchemaCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]SchemaCacheEntry)
}

// Len returns the number of entries currently in the cache (including expired).
func (c *SchemaCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
