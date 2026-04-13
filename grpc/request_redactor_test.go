package grpc

import (
	"testing"
)

func TestNewRequestRedactor_NotNil(t *testing.T) {
	r := NewRequestRedactor(DefaultRedactorPolicy())
	if r == nil {
		t.Fatal("expected non-nil RequestRedactor")
	}
}

func TestDefaultRedactorPolicy_Values(t *testing.T) {
	p := DefaultRedactorPolicy()
	if p.Placeholder != "[REDACTED]" {
		t.Errorf("unexpected placeholder: %s", p.Placeholder)
	}
}

func TestNewRequestRedactor_EmptyPlaceholder_FallsBackToDefault(t *testing.T) {
	r := NewRequestRedactor(RedactorPolicy{Placeholder: ""})
	if r.Policy().Placeholder != "[REDACTED]" {
		t.Errorf("expected default placeholder, got %s", r.Policy().Placeholder)
	}
}

func TestRequestRedactor_Len_Empty(t *testing.T) {
	r := NewRequestRedactor(DefaultRedactorPolicy())
	if r.Len() != 0 {
		t.Errorf("expected 0, got %d", r.Len())
	}
}

func TestRequestRedactor_AddField_And_Len(t *testing.T) {
	r := NewRequestRedactor(DefaultRedactorPolicy())
	if err := r.AddField("password"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 1 {
		t.Errorf("expected 1, got %d", r.Len())
	}
}

func TestRequestRedactor_AddField_EmptyName_ReturnsError(t *testing.T) {
	r := NewRequestRedactor(DefaultRedactorPolicy())
	if err := r.AddField(""); err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestRequestRedactor_Redact_NilPayload(t *testing.T) {
	r := NewRequestRedactor(DefaultRedactorPolicy())
	if r.Redact(nil) != nil {
		t.Fatal("expected nil result for nil payload")
	}
}

func TestRequestRedactor_Redact_ReplacesFields(t *testing.T) {
	r := NewRequestRedactor(DefaultRedactorPolicy())
	_ = r.AddField("token")
	_ = r.AddField("secret")

	payload := map[string]any{
		"username": "alice",
		"token":    "abc123",
		"secret":   "s3cr3t",
	}
	out := r.Redact(payload)

	if out["username"] != "alice" {
		t.Errorf("expected username preserved, got %v", out["username"])
	}
	if out["token"] != "[REDACTED]" {
		t.Errorf("expected token redacted, got %v", out["token"])
	}
	if out["secret"] != "[REDACTED]" {
		t.Errorf("expected secret redacted, got %v", out["secret"])
	}
}

func TestRequestRedactor_Redact_CaseInsensitive(t *testing.T) {
	r := NewRequestRedactor(DefaultRedactorPolicy())
	_ = r.AddField("Password")

	payload := map[string]any{"password": "hunter2"}
	out := r.Redact(payload)
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected redaction, got %v", out["password"])
	}
}
