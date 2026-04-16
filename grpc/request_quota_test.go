package grpc

import (
	"testing"
	"time"
)

func TestNewRequestQuota_NotNil(t *testing.T) {
	q := NewRequestQuota(DefaultQuotaPolicy())
	if q == nil {
		t.Fatal("expected non-nil RequestQuota")
	}
}

func TestDefaultQuotaPolicy_Values(t *testing.T) {
	p := DefaultQuotaPolicy()
	if p.MaxRequests != 100 {
		t.Errorf("expected 100, got %d", p.MaxRequests)
	}
	if p.WindowSize != time.Minute {
		t.Errorf("expected 1m, got %v", p.WindowSize)
	}
}

func TestNewRequestQuota_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	q := NewRequestQuota(QuotaPolicy{MaxRequests: -1, WindowSize: 0})
	if q.policy.MaxRequests != 100 {
		t.Errorf("expected fallback to 100, got %d", q.policy.MaxRequests)
	}
}

func TestRequestQuota_Len_Empty(t *testing.T) {
	q := NewRequestQuota(DefaultQuotaPolicy())
	if q.Len() != 0 {
		t.Errorf("expected 0, got %d", q.Len())
	}
}

func TestRequestQuota_Allow_EmptyMethod_ReturnsError(t *testing.T) {
	q := NewRequestQuota(DefaultQuotaPolicy())
	if err := q.Allow(""); err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRequestQuota_Allow_WithinLimit(t *testing.T) {
	q := NewRequestQuota(DefaultQuotaPolicy())
	if err := q.Allow("/svc/Method"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Len() != 1 {
		t.Errorf("expected 1, got %d", q.Len())
	}
}

func TestRequestQuota_Allow_ExceedsLimit(t *testing.T) {
	p := QuotaPolicy{MaxRequests: 2, WindowSize: time.Minute}
	q := NewRequestQuota(p)
	_ = q.Allow("m")
	_ = q.Allow("m")
	if err := q.Allow("m"); err == nil {
		t.Fatal("expected quota exceeded error")
	}
}

func TestRequestQuota_Reset_ClearsUsage(t *testing.T) {
	q := NewRequestQuota(DefaultQuotaPolicy())
	_ = q.Allow("m")
	q.Reset()
	if q.Len() != 0 {
		t.Errorf("expected 0 after reset, got %d", q.Len())
	}
}

func TestRequestQuota_Usage_ReturnsCopy(t *testing.T) {
	q := NewRequestQuota(DefaultQuotaPolicy())
	_ = q.Allow("/svc/A")
	_ = q.Allow("/svc/B")
	if len(q.Usage()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(q.Usage()))
	}
}
