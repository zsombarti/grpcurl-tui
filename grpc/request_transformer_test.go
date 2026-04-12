package grpc

import (
	"errors"
	"testing"
)

func TestNewRequestTransformer_NotNil(t *testing.T) {
	tr := NewRequestTransformer()
	if tr == nil {
		t.Fatal("expected non-nil RequestTransformer")
	}
}

func TestRequestTransformer_Len_Empty(t *testing.T) {
	tr := NewRequestTransformer()
	if tr.Len() != 0 {
		t.Fatalf("expected 0, got %d", tr.Len())
	}
}

func TestRequestTransformer_AddStep_And_Len(t *testing.T) {
	tr := NewRequestTransformer()
	err := tr.AddStep("noop", func(p map[string]interface{}) (map[string]interface{}, error) { return p, nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Len() != 1 {
		t.Fatalf("expected 1, got %d", tr.Len())
	}
}

func TestRequestTransformer_AddStep_EmptyName_ReturnsError(t *testing.T) {
	tr := NewRequestTransformer()
	err := tr.AddStep("", func(p map[string]interface{}) (map[string]interface{}, error) { return p, nil })
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRequestTransformer_AddStep_NilFn_ReturnsError(t *testing.T) {
	tr := NewRequestTransformer()
	err := tr.AddStep("step1", nil)
	if err == nil {
		t.Fatal("expected error for nil function")
	}
}

func TestRequestTransformer_Transform_NilPayload_ReturnsError(t *testing.T) {
	tr := NewRequestTransformer()
	_, err := tr.Transform(nil)
	if err == nil {
		t.Fatal("expected error for nil payload")
	}
}

func TestRequestTransformer_Transform_ModifiesPayload(t *testing.T) {
	tr := NewRequestTransformer()
	_ = tr.AddStep("add-key", func(p map[string]interface{}) (map[string]interface{}, error) {
		p["injected"] = "value"
		return p, nil
	})
	result, err := tr.Transform(map[string]interface{}{"original": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["injected"] != "value" {
		t.Fatalf("expected injected key, got %v", result)
	}
}

func TestRequestTransformer_Transform_StepError_Propagates(t *testing.T) {
	tr := NewRequestTransformer()
	_ = tr.AddStep("fail", func(p map[string]interface{}) (map[string]interface{}, error) {
		return nil, errors.New("step failed")
	})
	_, err := tr.Transform(map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error from step")
	}
}

func TestRequestTransformer_Names_Order(t *testing.T) {
	tr := NewRequestTransformer()
	noop := func(p map[string]interface{}) (map[string]interface{}, error) { return p, nil }
	_ = tr.AddStep("alpha", noop)
	_ = tr.AddStep("beta", noop)
	names := tr.Names()
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Fatalf("unexpected names: %v", names)
	}
}

func TestRequestTransformer_Clear_ResetsLen(t *testing.T) {
	tr := NewRequestTransformer()
	noop := func(p map[string]interface{}) (map[string]interface{}, error) { return p, nil }
	_ = tr.AddStep("s1", noop)
	tr.Clear()
	if tr.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", tr.Len())
	}
}
