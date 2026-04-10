package grpc

import (
	"testing"
	"time"
)

func TestDefaultTimeoutPolicy_Values(t *testing.T) {
	p := DefaultTimeoutPolicy()
	if p.DialTimeout != 5*time.Second {
		t.Errorf("expected dial timeout 5s, got %v", p.DialTimeout)
	}
	if p.RequestTimeout != 30*time.Second {
		t.Errorf("expected request timeout 30s, got %v", p.RequestTimeout)
	}
}

func TestNewTimeoutManager_NotNil(t *testing.T) {
	m, err := NewTimeoutManager(DefaultTimeoutPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil TimeoutManager")
	}
}

func TestNewTimeoutManager_InvalidDialTimeout(t *testing.T) {
	p := TimeoutPolicy{DialTimeout: 0, RequestTimeout: 10 * time.Second}
	_, err := NewTimeoutManager(p)
	if err == nil {
		t.Fatal("expected error for zero dial timeout")
	}
}

func TestNewTimeoutManager_InvalidRequestTimeout(t *testing.T) {
	p := TimeoutPolicy{DialTimeout: 5 * time.Second, RequestTimeout: -1}
	_, err := NewTimeoutManager(p)
	if err == nil {
		t.Fatal("expected error for negative request timeout")
	}
}

func TestTimeoutManager_Policy_RoundTrip(t *testing.T) {
	expected := TimeoutPolicy{
		DialTimeout:    2 * time.Second,
		RequestTimeout: 15 * time.Second,
	}
	m, err := NewTimeoutManager(expected)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Policy()
	if got.DialTimeout != expected.DialTimeout || got.RequestTimeout != expected.RequestTimeout {
		t.Errorf("policy mismatch: got %+v, want %+v", got, expected)
	}
}

func TestTimeoutManager_Accessors(t *testing.T) {
	p := TimeoutPolicy{DialTimeout: 3 * time.Second, RequestTimeout: 20 * time.Second}
	m, _ := NewTimeoutManager(p)
	if m.DialTimeout() != 3*time.Second {
		t.Errorf("unexpected dial timeout: %v", m.DialTimeout())
	}
	if m.RequestTimeout() != 20*time.Second {
		t.Errorf("unexpected request timeout: %v", m.RequestTimeout())
	}
}
