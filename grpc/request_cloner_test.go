package grpc

import (
	"testing"
)

func TestNewRequestCloner_NotNil(t *testing.T) {
	c := NewRequestCloner(DefaultClonerPolicy())
	if c == nil {
		t.Fatal("expected non-nil cloner")
	}
}

func TestDefaultClonerPolicy_Values(t *testing.T) {
	p := DefaultClonerPolicy()
	if p.MaxDepth <= 0 {
		t.Fatalf("expected positive MaxDepth, got %d", p.MaxDepth)
	}
}

func TestNewRequestCloner_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	c := NewRequestCloner(ClonerPolicy{MaxDepth: -1})
	if c.Policy().MaxDepth != DefaultClonerPolicy().MaxDepth {
		t.Fatal("expected fallback to default policy")
	}
}

func TestRequestCloner_Clone_NilPayload(t *testing.T) {
	c := NewRequestCloner(DefaultClonerPolicy())
	_, err := c.Clone(nil)
	if err == nil {
		t.Fatal("expected error for nil payload")
	}
}

func TestRequestCloner_Clone_SimpleMap(t *testing.T) {
	c := NewRequestCloner(DefaultClonerPolicy())
	orig := map[string]any{"key": "value", "num": 42.0}
	cloned, err := c.Clone(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cloned["key"] != "value" || cloned["num"] != 42.0 {
		t.Fatal("cloned map does not match original")
	}
}

func TestRequestCloner_Clone_Isolation(t *testing.T) {
	c := NewRequestCloner(DefaultClonerPolicy())
	orig := map[string]any{"nested": map[string]any{"x": "original"}}
	cloned, _ := c.Clone(orig)
	// mutate nested in original
	orig["nested"].(map[string]any)["x"] = "mutated"
	if cloned["nested"].(map[string]any)["x"] != "original" {
		t.Fatal("clone is not isolated from original")
	}
}

func TestRequestCloner_CloneJSON_NilPayload(t *testing.T) {
	c := NewRequestCloner(DefaultClonerPolicy())
	_, err := c.CloneJSON(nil)
	if err == nil {
		t.Fatal("expected error for nil payload")
	}
}

func TestRequestCloner_CloneJSON_RoundTrip(t *testing.T) {
	c := NewRequestCloner(DefaultClonerPolicy())
	orig := map[string]any{"a": "b", "c": 1.0}
	cloned, err := c.CloneJSON(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cloned["a"] != "b" {
		t.Fatal("JSON clone does not match original")
	}
}

func TestRequestCloner_Clone_MaxDepthExceeded(t *testing.T) {
	c := NewRequestCloner(ClonerPolicy{MaxDepth: 1})
	deep := map[string]any{"l1": map[string]any{"l2": map[string]any{"l3": "v"}}}
	_, err := c.Clone(deep)
	if err == nil {
		t.Fatal("expected max depth error")
	}
}
