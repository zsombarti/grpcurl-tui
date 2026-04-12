package grpc

import (
	"strings"
	"testing"
)

func TestNewRequestSanitizer_NotNil(t *testing.T) {
	s := NewRequestSanitizer(DefaultSanitizerPolicy())
	if s == nil {
		t.Fatal("expected non-nil sanitizer")
	}
}

func TestDefaultSanitizerPolicy_Values(t *testing.T) {
	p := DefaultSanitizerPolicy()
	if p.MaxValueLength <= 0 {
		t.Errorf("expected positive MaxValueLength, got %d", p.MaxValueLength)
	}
	if p.RedactedValue == "" {
		t.Error("expected non-empty RedactedValue")
	}
	if len(p.RedactedKeys) == 0 {
		t.Error("expected at least one redacted key")
	}
}

func TestNewRequestSanitizer_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	s := NewRequestSanitizer(SanitizerPolicy{MaxValueLength: 0})
	if s.Policy().MaxValueLength != DefaultSanitizerPolicy().MaxValueLength {
		t.Error("expected fallback to default policy")
	}
}

func TestRequestSanitizer_Sanitize_NilPayload(t *testing.T) {
	s := NewRequestSanitizer(DefaultSanitizerPolicy())
	if s.Sanitize(nil) != nil {
		t.Error("expected nil result for nil payload")
	}
}

func TestRequestSanitizer_Sanitize_RedactsSensitiveKeys(t *testing.T) {
	s := NewRequestSanitizer(DefaultSanitizerPolicy())
	payload := map[string]interface{}{
		"username": "alice",
		"password": "s3cr3t",
		"token":    "abc123",
	}
	out := s.Sanitize(payload)
	if out["username"] != "alice" {
		t.Errorf("expected username preserved, got %v", out["username"])
	}
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected password redacted, got %v", out["password"])
	}
	if out["token"] != "[REDACTED]" {
		t.Errorf("expected token redacted, got %v", out["token"])
	}
}

func TestRequestSanitizer_Sanitize_TruncatesLongValues(t *testing.T) {
	policy := DefaultSanitizerPolicy()
	policy.MaxValueLength = 10
	s := NewRequestSanitizer(policy)
	payload := map[string]interface{}{"data": strings.Repeat("x", 50)}
	out := s.Sanitize(payload)
	val, _ := out["data"].(string)
	if len(val) != 10 {
		t.Errorf("expected truncated value of length 10, got %d", len(val))
	}
}

func TestRequestSanitizer_Sanitize_NestedMap(t *testing.T) {
	s := NewRequestSanitizer(DefaultSanitizerPolicy())
	payload := map[string]interface{}{
		"user": map[string]interface{}{
			"name":   "bob",
			"secret": "topsecret",
		},
	}
	out := s.Sanitize(payload)
	nested, _ := out["user"].(map[string]interface{})
	if nested == nil {
		t.Fatal("expected nested map")
	}
	if nested["name"] != "bob" {
		t.Errorf("expected name preserved, got %v", nested["name"])
	}
	if nested["secret"] != "[REDACTED]" {
		t.Errorf("expected secret redacted, got %v", nested["secret"])
	}
}
