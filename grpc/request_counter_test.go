package grpc

import (
	"testing"
	"time"
)

func TestNewRequestCounter_NotNil(t *testing.T) {
	c := NewRequestCounter()
	if c == nil {
		t.Fatal("expected non-nil RequestCounter")
	}
}

func TestRequestCounter_Count_Empty(t *testing.T) {
	c := NewRequestCounter()
	if got := c.Count("/pkg.Service/Method"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestRequestCounter_Increment_And_Count(t *testing.T) {
	c := NewRequestCounter()
	method := "/pkg.Service/Method"
	c.Increment(method)
	c.Increment(method)
	if got := c.Count(method); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestRequestCounter_IncrementError_And_ErrorCount(t *testing.T) {
	c := NewRequestCounter()
	method := "/pkg.Service/Method"
	c.IncrementError(method)
	if got := c.ErrorCount(method); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestRequestCounter_Methods_ReturnsTracked(t *testing.T) {
	c := NewRequestCounter()
	c.Increment("/a/A")
	c.Increment("/b/B")
	if got := len(c.Methods()); got != 2 {
		t.Fatalf("expected 2 methods, got %d", got)
	}
}

func TestRequestCounter_Reset_ClearsAll(t *testing.T) {
	c := NewRequestCounter()
	method := "/pkg.Service/Method"
	c.Increment(method)
	c.IncrementError(method)
	c.Reset()
	if got := c.Count(method); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	if got := c.ErrorCount(method); got != 0 {
		t.Fatalf("expected 0 errors after reset, got %d", got)
	}
	if got := len(c.Methods()); got != 0 {
		t.Fatalf("expected 0 methods after reset, got %d", got)
	}
}

func TestRequestCounter_Since_Positive(t *testing.T) {
	c := NewRequestCounter()
	time.Sleep(2 * time.Millisecond)
	if c.Since() <= 0 {
		t.Fatal("expected positive duration since reset")
	}
}
