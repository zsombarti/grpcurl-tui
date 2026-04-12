package grpc

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// TagStore manages key-value tags attached to named requests.
type TagStore struct {
	mu      sync.RWMutex
	entries map[string][]string
	maxTags int
}

// NewTagStore creates a TagStore with a maximum number of tags per entry.
func NewTagStore(maxTags int) *TagStore {
	if maxTags <= 0 {
		maxTags = 20
	}
	return &TagStore{
		entries: make(map[string][]string),
		maxTags: maxTags,
	}
}

// Add attaches one or more tags to a named request key.
// Duplicate tags are ignored. Excess tags are silently dropped.
func (t *TagStore) Add(key string, tags ...string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("tag key must not be empty")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	existing := t.entries[key]
	set := make(map[string]struct{}, len(existing))
	for _, v := range existing {
		set[v] = struct{}{}
	}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, dup := set[tag]; dup {
			continue
		}
		if len(existing) >= t.maxTags {
			break
		}
		existing = append(existing, tag)
		set[tag] = struct{}{}
	}
	t.entries[key] = existing
	return nil
}

// Get returns a sorted copy of tags for the given key.
func (t *TagStore) Get(key string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	src := t.entries[key]
	out := make([]string, len(src))
	copy(out, src)
	sort.Strings(out)
	return out
}

// Remove deletes all tags for the given key.
func (t *TagStore) Remove(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// Keys returns all keys that have at least one tag.
func (t *TagStore) Keys() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	keys := make([]string, 0, len(t.entries))
	for k := range t.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Len returns the number of tagged keys.
func (t *TagStore) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.entries)
}
