package grpc

import (
	"errors"
	"testing"
)

func TestNewRequestEnricher_NotNil(t *testing.T) {
	e := NewRequestEnricher()
	if e == nil {
		t.Fatal("expected non-nil RequestEnricher")
	}
}

func TestRequestEnricher_Len_Empty(t *testing.T) {
	e := NewRequestEnricher()
	if e.Len() != 0 {
		t.Fatalf("expected 0, got %d", e.Len())
	}
}

func TestRequestEnricher_AddStep_And_Len(t *testing.T) {
	e := NewRequestEnricher()
	err := e.AddStep("add-version", func(p map[string]interface{}) (map[string]interface{}, error) {
		p["version"] = "v1"
		return p, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Len() != 1 {
		t.Fatalf("expected 1, got %d", e.Len())
	}
}

func TestRequestEnricher_AddStep_EmptyName_ReturnsError(t *testing.T) {
	e := NewRequestEnricher()
	err := e.AddStep("", func(p map[string]interface{}) (map[string]interface{}, error) { return p, nil })
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRequestEnricher_AddStep_NilFn_ReturnsError(t *testing.T) {
	e := NewRequestEnricher()
	err := e.AddStep("step", nil)
	if err == nil {
		t.Fatal("expected error for nil function")
	}
}

func TestRequestEnricher_Enrich_NilPayload_ReturnsError(t *testing.T) {
	e := NewRequestEnricher()
	_, err := e.Enrich(nil)
	if err == nil {
		t.Fatal("expected error for nil payload")
	}
}

func TestRequestEnricher_Enrich_AppliesSteps(t *testing.T) {
	e := NewRequestEnricher()
	_ = e.AddStep("add-env", func(p map[string]interface{}) (map[string]interface{}, error) {
		p["env"] = "production"
		return p, nil
	})
	_ = e.AddStep("add-source", func(p map[string]interface{}) (map[string]interface{}, error) {
		p["source"] = "grpcurl-tui"
		return p, nil
	})
	result, err := e.Enrich(map[string]interface{}{"name": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["env"] != "production" || result["source"] != "grpcurl-tui" {
		t.Fatalf("enrichment steps not applied correctly: %v", result)
	}
}

func TestRequestEnricher_Enrich_StepError_Propagates(t *testing.T) {
	e := NewRequestEnricher()
	_ = e.AddStep("fail", func(p map[string]interface{}) (map[string]interface{}, error) {
		return nil, errors.New("enrichment failed")
	})
	_, err := e.Enrich(map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error from failing step")
	}
}

func TestRequestEnricher_StepNames(t *testing.T) {
	e := NewRequestEnricher()
	_ = e.AddStep("alpha", func(p map[string]interface{}) (map[string]interface{}, error) { return p, nil })
	_ = e.AddStep("beta", func(p map[string]interface{}) (map[string]interface{}, error) { return p, nil })
	names := e.StepNames()
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Fatalf("unexpected step names: %v", names)
	}
}

func TestRequestEnricher_Clear_ResetsLen(t *testing.T) {
	e := NewRequestEnricher()
	_ = e.AddStep("s1", func(p map[string]interface{}) (map[string]interface{}, error) { return p, nil })
	e.Clear()
	if e.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", e.Len())
	}
}
