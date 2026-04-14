package grpc

import (
	"testing"
)

func TestNewRequestBroadcaster_NotNil(t *testing.T) {
	b := NewRequestBroadcaster(DefaultBroadcasterPolicy())
	if b == nil {
		t.Fatal("expected non-nil RequestBroadcaster")
	}
}

func TestDefaultBroadcasterPolicy_Values(t *testing.T) {
	p := DefaultBroadcasterPolicy()
	if p.MaxTargets <= 0 {
		t.Fatalf("expected positive MaxTargets, got %d", p.MaxTargets)
	}
}

func TestNewRequestBroadcaster_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	b := NewRequestBroadcaster(BroadcasterPolicy{MaxTargets: -1})
	if b.policy.MaxTargets != DefaultBroadcasterPolicy().MaxTargets {
		t.Fatalf("expected fallback to default MaxTargets")
	}
}

func TestRequestBroadcaster_Len_Empty(t *testing.T) {
	b := NewRequestBroadcaster(DefaultBroadcasterPolicy())
	if b.Len() != 0 {
		t.Fatalf("expected 0, got %d", b.Len())
	}
}

func TestRequestBroadcaster_AddTarget_And_Len(t *testing.T) {
	b := NewRequestBroadcaster(DefaultBroadcasterPolicy())
	if err := b.AddTarget("service-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Len() != 1 {
		t.Fatalf("expected 1, got %d", b.Len())
	}
}

func TestRequestBroadcaster_AddTarget_EmptyName_ReturnsError(t *testing.T) {
	b := NewRequestBroadcaster(DefaultBroadcasterPolicy())
	if err := b.AddTarget(""); err == nil {
		t.Fatal("expected error for empty target name")
	}
}

func TestRequestBroadcaster_AddTarget_MaxReached_ReturnsError(t *testing.T) {
	b := NewRequestBroadcaster(BroadcasterPolicy{MaxTargets: 2})
	_ = b.AddTarget("a")
	_ = b.AddTarget("b")
	if err := b.AddTarget("c"); err == nil {
		t.Fatal("expected error when max targets reached")
	}
}

func TestRequestBroadcaster_Broadcast_ReturnsOneResultPerTarget(t *testing.T) {
	b := NewRequestBroadcaster(DefaultBroadcasterPolicy())
	_ = b.AddTarget("svc-1")
	_ = b.AddTarget("svc-2")
	payload := map[string]any{"key": "value"}
	results := b.Broadcast(payload)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Fatalf("unexpected error for target %s: %v", r.Target, r.Err)
		}
		if r.Payload["key"] != "value" {
			t.Fatalf("payload not propagated for target %s", r.Target)
		}
	}
}

func TestRequestBroadcaster_Clear_ResetsLen(t *testing.T) {
	b := NewRequestBroadcaster(DefaultBroadcasterPolicy())
	_ = b.AddTarget("svc-1")
	b.Clear()
	if b.Len() != 0 {
		t.Fatalf("expected 0 after Clear, got %d", b.Len())
	}
}
