package grpc

import (
	"testing"
)

func TestNewRequestFilter_NotNil(t *testing.T) {
	f := NewRequestFilter()
	if f == nil {
		t.Fatal("expected non-nil RequestFilter")
	}
}

func TestRequestFilter_Len_Empty(t *testing.T) {
	f := NewRequestFilter()
	if f.Len() != 0 {
		t.Fatalf("expected 0, got %d", f.Len())
	}
}

func TestRequestFilter_AddRule_And_Len(t *testing.T) {
	f := NewRequestFilter()
	if err := f.AddRule("method", "eq", "SayHello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Len() != 1 {
		t.Fatalf("expected 1, got %d", f.Len())
	}
}

func TestRequestFilter_AddRule_EmptyField_ReturnsError(t *testing.T) {
	f := NewRequestFilter()
	if err := f.AddRule("", "eq", "val"); err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestRequestFilter_AddRule_InvalidOperator_ReturnsError(t *testing.T) {
	f := NewRequestFilter()
	if err := f.AddRule("field", "regex", ".*"); err == nil {
		t.Fatal("expected error for unsupported operator")
	}
}

func TestRequestFilter_Match_Eq_True(t *testing.T) {
	f := NewRequestFilter()
	_ = f.AddRule("service", "eq", "Greeter")
	if !f.Match(map[string]string{"service": "Greeter"}) {
		t.Fatal("expected match")
	}
}

func TestRequestFilter_Match_Eq_False(t *testing.T) {
	f := NewRequestFilter()
	_ = f.AddRule("service", "eq", "Greeter")
	if f.Match(map[string]string{"service": "Other"}) {
		t.Fatal("expected no match")
	}
}

func TestRequestFilter_Match_Contains(t *testing.T) {
	f := NewRequestFilter()
	_ = f.AddRule("body", "contains", "hello")
	if !f.Match(map[string]string{"body": "say hello world"}) {
		t.Fatal("expected match")
	}
}

func TestRequestFilter_Match_Prefix(t *testing.T) {
	f := NewRequestFilter()
	_ = f.AddRule("method", "prefix", "Get")
	if !f.Match(map[string]string{"method": "GetUser"}) {
		t.Fatal("expected match")
	}
}

func TestRequestFilter_Match_MissingField_ReturnsFalse(t *testing.T) {
	f := NewRequestFilter()
	_ = f.AddRule("method", "eq", "Foo")
	if f.Match(map[string]string{"other": "Foo"}) {
		t.Fatal("expected no match on missing field")
	}
}

func TestRequestFilter_Clear_ResetsLen(t *testing.T) {
	f := NewRequestFilter()
	_ = f.AddRule("a", "eq", "b")
	f.Clear()
	if f.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", f.Len())
	}
}
