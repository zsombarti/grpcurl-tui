package grpc

import (
	"testing"
)

func TestNewRequestNormalizer_NotNil(t *testing.T) {
	n := NewRequestNormalizer(DefaultNormalizerPolicy())
	if n == nil {
		t.Fatal("expected non-nil RequestNormalizer")
	}
}

func TestDefaultNormalizerPolicy_Values(t *testing.T) {
	p := DefaultNormalizerPolicy()
	if !p.TrimStrings {
		t.Error("expected TrimStrings to be true")
	}
	if p.LowercaseKeys {
		t.Error("expected LowercaseKeys to be false")
	}
	if !p.RemoveNullKeys {
		t.Error("expected RemoveNullKeys to be true")
	}
}

func TestNewRequestNormalizer_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	n := NewRequestNormalizer(NormalizerPolicy{})
	p := n.Policy()
	if !p.TrimStrings || !p.RemoveNullKeys {
		t.Error("expected fallback to default policy")
	}
}

func TestRequestNormalizer_Normalize_NilPayload(t *testing.T) {
	n := NewRequestNormalizer(DefaultNormalizerPolicy())
	_, err := n.Normalize(nil)
	if err == nil {
		t.Fatal("expected error for nil payload")
	}
}

func TestRequestNormalizer_Normalize_TrimStrings(t *testing.T) {
	n := NewRequestNormalizer(NormalizerPolicy{TrimStrings: true})
	out, err := n.Normalize(map[string]any{"name": "  hello  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["name"] != "hello" {
		t.Errorf("expected trimmed string, got %q", out["name"])
	}
}

func TestRequestNormalizer_Normalize_RemovesNullKeys(t *testing.T) {
	n := NewRequestNormalizer(NormalizerPolicy{RemoveNullKeys: true})
	out, err := n.Normalize(map[string]any{"key": nil, "other": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["key"]; ok {
		t.Error("expected null key to be removed")
	}
	if out["other"] != "val" {
		t.Error("expected non-null key to be preserved")
	}
}

func TestRequestNormalizer_Normalize_LowercaseKeys(t *testing.T) {
	n := NewRequestNormalizer(NormalizerPolicy{LowercaseKeys: true})
	out, err := n.Normalize(map[string]any{"Name": "Alice", "AGE": 30})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["name"]; !ok {
		t.Error("expected 'name' key after lowercasing")
	}
	if _, ok := out["age"]; !ok {
		t.Error("expected 'age' key after lowercasing")
	}
}
