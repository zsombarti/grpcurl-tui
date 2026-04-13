package grpc

import (
	"testing"
	"time"
)

func TestNewRequestArchiver_NotNil(t *testing.T) {
	a := NewRequestArchiver(0)
	if a == nil {
		t.Fatal("expected non-nil RequestArchiver")
	}
}

func TestNewRequestArchiver_DefaultMaxSize(t *testing.T) {
	a := NewRequestArchiver(0)
	if a.maxSize != defaultArchiverMaxSize {
		t.Fatalf("expected maxSize %d, got %d", defaultArchiverMaxSize, a.maxSize)
	}
}

func TestRequestArchiver_Len_Empty(t *testing.T) {
	a := NewRequestArchiver(10)
	if a.Len() != 0 {
		t.Fatalf("expected 0, got %d", a.Len())
	}
}

func TestRequestArchiver_Archive_And_Len(t *testing.T) {
	a := NewRequestArchiver(10)
	err := a.Archive("SomeService/Method", map[string]interface{}{"key": "value"}, "v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Len() != 1 {
		t.Fatalf("expected 1, got %d", a.Len())
	}
}

func TestRequestArchiver_Archive_EmptyMethod_ReturnsError(t *testing.T) {
	a := NewRequestArchiver(10)
	err := a.Archive("", map[string]interface{}{"k": "v"}, "")
	if err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRequestArchiver_Archive_NilPayload_ReturnsError(t *testing.T) {
	a := NewRequestArchiver(10)
	err := a.Archive("SomeService/Method", nil, "")
	if err == nil {
		t.Fatal("expected error for nil payload")
	}
}

func TestRequestArchiver_Eviction(t *testing.T) {
	a := NewRequestArchiver(2)
	_ = a.Archive("M1", map[string]interface{}{"a": 1}, "")
	_ = a.Archive("M2", map[string]interface{}{"b": 2}, "")
	_ = a.Archive("M3", map[string]interface{}{"c": 3}, "")
	if a.Len() != 2 {
		t.Fatalf("expected 2 after eviction, got %d", a.Len())
	}
	if a.All()[0].Method != "M2" {
		t.Fatalf("expected oldest entry evicted, got %s", a.All()[0].Method)
	}
}

func TestRequestArchiver_TimestampAutoSet(t *testing.T) {
	before := time.Now()
	a := NewRequestArchiver(10)
	_ = a.Archive("M", map[string]interface{}{"x": 1}, "tag")
	after := time.Now()
	entries := a.All()
	if entries[0].ArchivedAt.Before(before) || entries[0].ArchivedAt.After(after) {
		t.Fatal("ArchivedAt timestamp out of expected range")
	}
}

func TestRequestArchiver_Clear(t *testing.T) {
	a := NewRequestArchiver(10)
	_ = a.Archive("M", map[string]interface{}{"x": 1}, "")
	a.Clear()
	if a.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", a.Len())
	}
}
