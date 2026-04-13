package grpc

import (
	"strings"
	"testing"
)

func TestNewRequestSigner_NotNil(t *testing.T) {
	s, err := NewRequestSigner("mysecret", DefaultSignerPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil signer")
	}
}

func TestDefaultSignerPolicy_Values(t *testing.T) {
	p := DefaultSignerPolicy()
	if p.Algorithm != "hmac-sha256" {
		t.Errorf("expected hmac-sha256, got %s", p.Algorithm)
	}
	if p.HeaderName != "x-grpcurl-signature" {
		t.Errorf("unexpected header name: %s", p.HeaderName)
	}
	if !p.IncludeTimestamp {
		t.Error("expected IncludeTimestamp to be true")
	}
}

func TestNewRequestSigner_EmptySecret_ReturnsError(t *testing.T) {
	_, err := NewRequestSigner("", DefaultSignerPolicy())
	if err == nil {
		t.Fatal("expected error for empty secret")
	}
}

func TestNewRequestSigner_EmptyAlgorithm_FallsBackToDefault(t *testing.T) {
	s, err := NewRequestSigner("secret", SignerPolicy{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Policy().Algorithm != DefaultSignerPolicy().Algorithm {
		t.Errorf("expected fallback algorithm, got %s", s.Policy().Algorithm)
	}
}

func TestRequestSigner_Sign_NilPayload_ReturnsError(t *testing.T) {
	s, _ := NewRequestSigner("secret", DefaultSignerPolicy())
	_, err := s.Sign(nil)
	if err == nil {
		t.Fatal("expected error for nil payload")
	}
}

func TestRequestSigner_Sign_ReturnsHexString(t *testing.T) {
	s, _ := NewRequestSigner("secret", DefaultSignerPolicy())
	sig, err := s.Sign(map[string]any{"method": "SayHello", "body": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sig) != 64 {
		t.Errorf("expected 64 hex chars, got %d", len(sig))
	}
}

func TestRequestSigner_Sign_Deterministic_WithoutTimestamp(t *testing.T) {
	policy := DefaultSignerPolicy()
	policy.IncludeTimestamp = false
	s, _ := NewRequestSigner("secret", policy)
	payload := map[string]any{"a": "1", "b": "2"}
	sig1, _ := s.Sign(payload)
	sig2, _ := s.Sign(payload)
	if sig1 != sig2 {
		t.Error("expected deterministic signatures without timestamp")
	}
}

func TestRequestSigner_HeaderName(t *testing.T) {
	s, _ := NewRequestSigner("secret", DefaultSignerPolicy())
	if !strings.HasPrefix(s.HeaderName(), "x-") {
		t.Errorf("unexpected header name: %s", s.HeaderName())
	}
}
