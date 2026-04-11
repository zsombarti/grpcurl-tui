package grpc

import (
	"testing"
	"time"
)

func TestNewSchemaCache_NotNil(t *testing.T) {
	c := NewSchemaCache(0)
	if c == nil {
		t.Fatal("expected non-nil SchemaCache")
	}
}

func TestNewSchemaCache_DefaultTTL(t *testing.T) {
	c := NewSchemaCache(0)
	if c.ttl != 5*time.Minute {
		t.Fatalf("expected default TTL of 5m, got %v", c.ttl)
	}
}

func TestNewSchemaCache_CustomTTL(t *testing.T) {
	c := NewSchemaCache(10 * time.Second)
	if c.ttl != 10*time.Second {
		t.Fatalf("expected TTL of 10s, got %v", c.ttl)
	}
}

func TestSchemaCache_Len_Empty(t *testing.T) {
	c := NewSchemaCache(time.Minute)
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
}

func TestSchemaCache_Set_And_Get(t *testing.T) {
	c := NewSchemaCache(time.Minute)
	c.Set("svc.Foo", nil) // nil descriptor is valid for cache key testing
	if c.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", c.Len())
	}
	_, ok := c.Get("svc.Foo")
	if !ok {
		t.Fatal("expected to find cached entry")
	}
}

func TestSchemaCache_Get_MissingKey(t *testing.T) {
	c := NewSchemaCache(time.Minute)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected miss for nonexistent key")
	}
}

func TestSchemaCache_Get_ExpiredEntry(t *testing.T) {
	c := NewSchemaCache(1 * time.Millisecond)
	c.Set("svc.Bar", nil)
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("svc.Bar")
	if ok {
		t.Fatal("expected expired entry to be a cache miss")
	}
}

func TestSchemaCache_Invalidate(t *testing.T) {
	c := NewSchemaCache(time.Minute)
	c.Set("svc.Baz", nil)
	c.Invalidate("svc.Baz")
	if c.Len() != 0 {
		t.Fatalf("expected 0 after invalidate, got %d", c.Len())
	}
}

func TestSchemaCache_Invalidate_NonexistentKey(t *testing.T) {
	c := NewSchemaCache(time.Minute)
	c.Set("svc.Baz", nil)
	// Invalidating a key that doesn't exist should be a no-op.
	c.Invalidate("svc.DoesNotExist")
	if c.Len() != 1 {
		t.Fatalf("expected 1 entry after invalidating nonexistent key, got %d", c.Len())
	}
}

func TestSchemaCache_Flush(t *testing.T) {
	c := NewSchemaCache(time.Minute)
	c.Set("a", nil)
	c.Set("b", nil)
	c.Set("c", nil)
	c.Flush()
	if c.Len() != 0 {
		t.Fatalf("expected 0 after flush, got %d", c.Len())
	}
}

func TestSchemaCacheEntry_IsExpired_False(t *testing.T) {
	e := SchemaCacheEntry{CachedAt: time.Now(), TTL: time.Minute}
	if e.IsExpired() {
		t.Fatal("expected entry to not be expired")
	}
}

func TestSchemaCacheEntry_IsExpired_True(t *testing.T) {
	e := SchemaCacheEntry{CachedAt: time.Now().Add(-2 * time.Minute), TTL: time.Minute}
	if !e.IsExpired() {
		t.Fatal("expected entry to be expired")
	}
}
