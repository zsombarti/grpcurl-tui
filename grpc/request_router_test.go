package grpc

import (
	"testing"
)

func TestNewRequestRouter_NotNil(t *testing.T) {
	r := NewRequestRouter("localhost:50051")
	if r == nil {
		t.Fatal("expected non-nil RequestRouter")
	}
}

func TestRequestRouter_Fallback_Default(t *testing.T) {
	r := NewRequestRouter("default:9090")
	if r.Fallback() != "default:9090" {
		t.Fatalf("expected fallback 'default:9090', got %q", r.Fallback())
	}
}

func TestRequestRouter_Len_Empty(t *testing.T) {
	r := NewRequestRouter("")
	if r.Len() != 0 {
		t.Fatalf("expected 0 rules, got %d", r.Len())
	}
}

func TestRequestRouter_AddRule_And_Len(t *testing.T) {
	r := NewRequestRouter("")
	if err := r.AddRule("/pkg.Service/Method", "host:1234"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1 rule, got %d", r.Len())
	}
}

func TestRequestRouter_AddRule_EmptyMethod_ReturnsError(t *testing.T) {
	r := NewRequestRouter("")
	if err := r.AddRule("", "host:1234"); err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRequestRouter_AddRule_EmptyAddress_ReturnsError(t *testing.T) {
	r := NewRequestRouter("")
	if err := r.AddRule("/pkg.Service/Method", ""); err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestRequestRouter_Route_MatchesRule(t *testing.T) {
	r := NewRequestRouter("fallback:9999")
	_ = r.AddRule("/pkg.Service/Foo", "primary:1111")
	got := r.Route("/pkg.Service/Foo")
	if got != "primary:1111" {
		t.Fatalf("expected 'primary:1111', got %q", got)
	}
}

func TestRequestRouter_Route_FallsBackWhenNoMatch(t *testing.T) {
	r := NewRequestRouter("fallback:9999")
	_ = r.AddRule("/pkg.Service/Foo", "primary:1111")
	got := r.Route("/pkg.Service/Bar")
	if got != "fallback:9999" {
		t.Fatalf("expected fallback 'fallback:9999', got %q", got)
	}
}

func TestRequestRouter_Clear_ResetsLen(t *testing.T) {
	r := NewRequestRouter("")
	_ = r.AddRule("/a/B", "host:1")
	_ = r.AddRule("/a/C", "host:2")
	r.Clear()
	if r.Len() != 0 {
		t.Fatalf("expected 0 after Clear, got %d", r.Len())
	}
}

func TestRequestRouter_SetFallback(t *testing.T) {
	r := NewRequestRouter("old:1")
	r.SetFallback("new:2")
	if r.Fallback() != "new:2" {
		t.Fatalf("expected 'new:2', got %q", r.Fallback())
	}
}
