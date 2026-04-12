package grpc

import (
	"testing"
)

func TestNewRequestLabeler_NotNil(t *testing.T) {
	l := NewRequestLabeler(0)
	if l == nil {
		t.Fatal("expected non-nil RequestLabeler")
	}
}

func TestNewRequestLabeler_DefaultMaxSize(t *testing.T) {
	l := NewRequestLabeler(0)
	if l.maxLen != 128 {
		t.Fatalf("expected default maxLen 128, got %d", l.maxLen)
	}
}

func TestRequestLabeler_Len_Empty(t *testing.T) {
	l := NewRequestLabeler(10)
	if l.Len() != 0 {
		t.Fatalf("expected 0, got %d", l.Len())
	}
}

func TestRequestLabeler_Set_And_Get(t *testing.T) {
	l := NewRequestLabeler(10)
	if err := l.Set("/pkg.Svc/Method", "My Label", "green"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := l.Get("/pkg.Svc/Method")
	if !ok {
		t.Fatal("expected label to be found")
	}
	if got.Label != "My Label" || got.Color != "green" {
		t.Fatalf("unexpected label: %+v", got)
	}
}

func TestRequestLabeler_Set_EmptyMethod_ReturnsError(t *testing.T) {
	l := NewRequestLabeler(10)
	if err := l.Set("", "lbl", ""); err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRequestLabeler_Set_EmptyLabel_ReturnsError(t *testing.T) {
	l := NewRequestLabeler(10)
	if err := l.Set("/pkg.Svc/Method", "", ""); err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestRequestLabeler_Remove(t *testing.T) {
	l := NewRequestLabeler(10)
	_ = l.Set("/pkg.Svc/Method", "lbl", "")
	l.Remove("/pkg.Svc/Method")
	if _, ok := l.Get("/pkg.Svc/Method"); ok {
		t.Fatal("expected label to be removed")
	}
}

func TestRequestLabeler_Eviction(t *testing.T) {
	l := NewRequestLabeler(2)
	_ = l.Set("/svc/A", "A", "")
	_ = l.Set("/svc/B", "B", "")
	_ = l.Set("/svc/C", "C", "") // should evict one
	if l.Len() > 2 {
		t.Fatalf("expected at most 2 entries, got %d", l.Len())
	}
}

func TestRequestLabeler_All_ReturnsSnapshot(t *testing.T) {
	l := NewRequestLabeler(10)
	_ = l.Set("/svc/A", "Alpha", "red")
	_ = l.Set("/svc/B", "Beta", "blue")
	all := l.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}
