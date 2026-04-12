package grpc

import (
	"testing"
	"time"
)

func TestNewRequestCache_NotNil(t *testing.T) {
	c := NewRequestCache(0, 0)
	if c == nil {
		t.Fatal("expected non-nil RequestCache")
	}
}

func TestNewRequestCache_DefaultTTL(t *testing.T) {
	c := NewRequestCache(0, 0)
	if c.ttl != 30*time.Second {
		t.Fatalf("expected default TTL 30s, got %v", c.ttl)
	}
}

func TestNewRequestCache_DefaultMaxSize(t *testing.T) {
	c := NewRequestCache(0, 0)
	if c.maxSize != defaultRequestCacheMaxSize {
		t.Fatalf("expected max size %d, got %d", defaultRequestCacheMaxSize, c.maxSize)
	}
}

func TestRequestCache_Len_Empty(t *testing.T) {
	c := NewRequestCache(time.Minute, 10)
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
}

func TestRequestCache_Set_And_Get(t *testing.T) {
	c := NewRequestCache(time.Minute, 10)
	c.Set("pkg.Service/Method", `{"id":1}`, `{"result":"ok"}`)
	entry, err := c.Get("pkg.Service/Method", `{"id":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Response != `{"result":"ok"}` {
		t.Fatalf("unexpected response: %s", entry.Response)
	}
}

func TestRequestCache_Get_Miss(t *testing.T) {
	c := NewRequestCache(time.Minute, 10)
	_, err := c.Get("missing", "{}")
	if err == nil {
		t.Fatal("expected cache miss error")
	}
}

func TestRequestCache_Get_Expired(t *testing.T) {
	c := NewRequestCache(time.Millisecond, 10)
	c.Set("m", "p", "r")
	time.Sleep(5 * time.Millisecond)
	_, err := c.Get("m", "p")
	if err == nil {
		t.Fatal("expected expiry error")
	}
}

func TestRequestCache_Eviction(t *testing.T) {
	c := NewRequestCache(time.Minute, 3)
	c.Set("m", "a", "r")
	c.Set("m", "b", "r")
	c.Set("m", "c", "r")
	c.Set("m", "d", "r") // should evict "a"
	if c.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", c.Len())
	}
	_, err := c.Get("m", "a")
	if err == nil {
		t.Fatal("expected evicted entry to be missing")
	}
}

func TestRequestCache_Clear(t *testing.T) {
	c := NewRequestCache(time.Minute, 10)
	c.Set("m", "p", "r")
	c.Clear()
	if c.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", c.Len())
	}
}
