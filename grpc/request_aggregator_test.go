package grpc

import (
	"testing"
)

func TestNewRequestAggregator_NotNil(t *testing.T) {
	a := NewRequestAggregator(0)
	if a == nil {
		t.Fatal("expected non-nil aggregator")
	}
}

func TestNewRequestAggregator_DefaultMaxSize(t *testing.T) {
	a := NewRequestAggregator(0)
	if a.maxSize != defaultAggregatorMaxSize {
		t.Fatalf("expected %d, got %d", defaultAggregatorMaxSize, a.maxSize)
	}
}

func TestRequestAggregator_Len_Empty(t *testing.T) {
	a := NewRequestAggregator(10)
	if a.Len() != 0 {
		t.Fatalf("expected 0, got %d", a.Len())
	}
}

func TestRequestAggregator_Aggregate_And_Len(t *testing.T) {
	a := NewRequestAggregator(10)
	_, err := a.Aggregate("SomeMethod", []map[string]interface{}{{"key": "val"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Len() != 1 {
		t.Fatalf("expected 1, got %d", a.Len())
	}
}

func TestRequestAggregator_Aggregate_EmptyMethod_ReturnsError(t *testing.T) {
	a := NewRequestAggregator(10)
	_, err := a.Aggregate("", []map[string]interface{}{{"k": "v"}})
	if err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRequestAggregator_Aggregate_EmptyPayloads_ReturnsError(t *testing.T) {
	a := NewRequestAggregator(10)
	_, err := a.Aggregate("Method", nil)
	if err == nil {
		t.Fatal("expected error for empty payloads")
	}
}

func TestRequestAggregator_Aggregate_MergesKeys(t *testing.T) {
	a := NewRequestAggregator(10)
	result, err := a.Aggregate("Method", []map[string]interface{}{
		{"a": 1, "b": 2},
		{"b": 99, "c": 3},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Merged["b"] != 99 {
		t.Fatalf("expected b=99, got %v", result.Merged["b"])
	}
	if result.Merged["a"] != 1 {
		t.Fatalf("expected a=1, got %v", result.Merged["a"])
	}
	if result.Merged["c"] != 3 {
		t.Fatalf("expected c=3, got %v", result.Merged["c"])
	}
}

func TestRequestAggregator_Eviction(t *testing.T) {
	a := NewRequestAggregator(2)
	for i := 0; i < 3; i++ {
		a.Aggregate("M", []map[string]interface{}{{"i": i}})
	}
	if a.Len() != 2 {
		t.Fatalf("expected 2, got %d", a.Len())
	}
}

func TestRequestAggregator_Clear(t *testing.T) {
	a := NewRequestAggregator(10)
	a.Aggregate("M", []map[string]interface{}{{"x": 1}})
	a.Clear()
	if a.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", a.Len())
	}
}

func TestRequestAggregator_All_ReturnsCopy(t *testing.T) {
	a := NewRequestAggregator(10)
	a.Aggregate("M", []map[string]interface{}{{"x": 1}})
	all := a.All()
	if len(all) != 1 {
		t.Fatalf("expected 1, got %d", len(all))
	}
}
