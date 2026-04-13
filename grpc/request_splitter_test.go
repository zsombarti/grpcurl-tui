package grpc

import (
	"testing"
)

func TestNewRequestSplitter_NotNil(t *testing.T) {
	s := NewRequestSplitter(DefaultSplitterPolicy())
	if s == nil {
		t.Fatal("expected non-nil RequestSplitter")
	}
}

func TestDefaultSplitterPolicy_Values(t *testing.T) {
	p := DefaultSplitterPolicy()
	if p.MaxTargets <= 0 {
		t.Fatalf("expected positive MaxTargets, got %d", p.MaxTargets)
	}
}

func TestNewRequestSplitter_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	s := NewRequestSplitter(SplitterPolicy{MaxTargets: -1})
	if s.policy.MaxTargets != DefaultSplitterPolicy().MaxTargets {
		t.Fatalf("expected fallback to default MaxTargets")
	}
}

func TestRequestSplitter_Len_Empty(t *testing.T) {
	s := NewRequestSplitter(DefaultSplitterPolicy())
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
}

func TestRequestSplitter_AddTarget_And_Len(t *testing.T) {
	s := NewRequestSplitter(DefaultSplitterPolicy())
	if err := s.AddTarget("localhost:50051"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1, got %d", s.Len())
	}
}

func TestRequestSplitter_AddTarget_EmptyAddress_ReturnsError(t *testing.T) {
	s := NewRequestSplitter(DefaultSplitterPolicy())
	if err := s.AddTarget(""); err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestRequestSplitter_AddTarget_ExceedsCap_ReturnsError(t *testing.T) {
	s := NewRequestSplitter(SplitterPolicy{MaxTargets: 2})
	_ = s.AddTarget("a:1")
	_ = s.AddTarget("b:2")
	if err := s.AddTarget("c:3"); err == nil {
		t.Fatal("expected error when max targets reached")
	}
}

func TestRequestSplitter_Split_ReturnsOneResultPerTarget(t *testing.T) {
	s := NewRequestSplitter(DefaultSplitterPolicy())
	_ = s.AddTarget("host1:50051")
	_ = s.AddTarget("host2:50051")
	results := s.Split(map[string]interface{}{"key": "value"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRequestSplitter_Split_NilPayload_Handled(t *testing.T) {
	s := NewRequestSplitter(DefaultSplitterPolicy())
	_ = s.AddTarget("host1:50051")
	results := s.Split(nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Payload == nil {
		t.Fatal("expected non-nil payload in result")
	}
}

func TestRequestSplitter_Split_PayloadIsCopied(t *testing.T) {
	s := NewRequestSplitter(DefaultSplitterPolicy())
	_ = s.AddTarget("host1:50051")
	_ = s.AddTarget("host2:50051")
	orig := map[string]interface{}{"x": 1}
	results := s.Split(orig)
	results[0].Payload["x"] = 99
	if orig["x"] == 99 {
		t.Fatal("split should not mutate original payload")
	}
	if results[1].Payload["x"] == 99 {
		t.Fatal("split should produce independent copies per target")
	}
}

func TestRequestSplitter_ClearTargets(t *testing.T) {
	s := NewRequestSplitter(DefaultSplitterPolicy())
	_ = s.AddTarget("host1:50051")
	s.ClearTargets()
	if s.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", s.Len())
	}
}
