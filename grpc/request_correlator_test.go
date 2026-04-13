package grpc

import (
	"testing"
)

func TestNewRequestCorrelator_NotNil(t *testing.T) {
	c := NewRequestCorrelator(0)
	if c == nil {
		t.Fatal("expected non-nil correlator")
	}
}

func TestNewRequestCorrelator_DefaultMaxSize(t *testing.T) {
	c := NewRequestCorrelator(0)
	if c.maxSize != 200 {
		t.Fatalf("expected default maxSize 200, got %d", c.maxSize)
	}
}

func TestRequestCorrelator_Len_Empty(t *testing.T) {
	c := NewRequestCorrelator(10)
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
}

func TestRequestCorrelator_Assign_And_Len(t *testing.T) {
	c := NewRequestCorrelator(10)
	_, err := c.Assign("/pkg.Service/Method")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 1 {
		t.Fatalf("expected 1, got %d", c.Len())
	}
}

func TestRequestCorrelator_Assign_EmptyMethod_ReturnsError(t *testing.T) {
	c := NewRequestCorrelator(10)
	_, err := c.Assign("")
	if err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRequestCorrelator_Lookup_Found(t *testing.T) {
	c := NewRequestCorrelator(10)
	id, _ := c.Assign("/pkg.Service/Foo")
	entry, err := c.Lookup(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.ID != id {
		t.Fatalf("expected id %s, got %s", id, entry.ID)
	}
	if entry.Method != "/pkg.Service/Foo" {
		t.Fatalf("unexpected method: %s", entry.Method)
	}
}

func TestRequestCorrelator_Lookup_NotFound(t *testing.T) {
	c := NewRequestCorrelator(10)
	_, err := c.Lookup("nonexistent-id")
	if err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestRequestCorrelator_Eviction(t *testing.T) {
	c := NewRequestCorrelator(2)
	id1, _ := c.Assign("/svc/A")
	c.Assign("/svc/B")
	c.Assign("/svc/C")
	if c.Len() != 2 {
		t.Fatalf("expected 2 after eviction, got %d", c.Len())
	}
	_, err := c.Lookup(id1)
	if err == nil {
		t.Fatal("expected evicted entry to be gone")
	}
}

func TestRequestCorrelator_Clear(t *testing.T) {
	c := NewRequestCorrelator(10)
	c.Assign("/svc/A")
	c.Assign("/svc/B")
	c.Clear()
	if c.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", c.Len())
	}
}
