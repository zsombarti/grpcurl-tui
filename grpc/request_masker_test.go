package grpc

import (
	"testing"
)

func TestNewRequestMasker_NotNil(t *testing.T) {
	m := NewRequestMasker(DefaultMaskerPolicy())
	if m == nil {
		t.Fatal("expected non-nil RequestMasker")
	}
}

func TestDefaultMaskerPolicy_Values(t *testing.T) {
	p := DefaultMaskerPolicy()
	if p.MaskChar != "*" {
		t.Errorf("expected MaskChar '*', got %q", p.MaskChar)
	}
	if p.KeepFirst != 0 || p.KeepLast != 0 {
		t.Errorf("expected KeepFirst/KeepLast 0, got %d/%d", p.KeepFirst, p.KeepLast)
	}
}

func TestNewRequestMasker_EmptyMaskChar_FallsBackToDefault(t *testing.T) {
	m := NewRequestMasker(MaskerPolicy{MaskChar: ""})
	if m.policy.MaskChar != "*" {
		t.Errorf("expected fallback MaskChar '*', got %q", m.policy.MaskChar)
	}
}

func TestRequestMasker_Len_Empty(t *testing.T) {
	m := NewRequestMasker(DefaultMaskerPolicy())
	if m.Len() != 0 {
		t.Errorf("expected 0, got %d", m.Len())
	}
}

func TestRequestMasker_AddField_And_Len(t *testing.T) {
	m := NewRequestMasker(DefaultMaskerPolicy())
	if err := m.AddField("password"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 1 {
		t.Errorf("expected 1, got %d", m.Len())
	}
}

func TestRequestMasker_AddField_EmptyName_ReturnsError(t *testing.T) {
	m := NewRequestMasker(DefaultMaskerPolicy())
	if err := m.AddField(""); err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestRequestMasker_Mask_NilPayload(t *testing.T) {
	m := NewRequestMasker(DefaultMaskerPolicy())
	if m.Mask(nil) != nil {
		t.Error("expected nil for nil payload")
	}
}

func TestRequestMasker_Mask_MasksRegisteredField(t *testing.T) {
	m := NewRequestMasker(DefaultMaskerPolicy())
	_ = m.AddField("token")
	out := m.Mask(map[string]any{"token": "abc123", "user": "alice"})
	if out["token"] == "abc123" {
		t.Error("expected token to be masked")
	}
	if out["user"] != "alice" {
		t.Error("expected user to be unchanged")
	}
}

func TestRequestMasker_Mask_KeepFirstAndLast(t *testing.T) {
	p := MaskerPolicy{MaskChar: "#", KeepFirst: 2, KeepLast: 2}
	m := NewRequestMasker(p)
	_ = m.AddField("secret")
	out := m.Mask(map[string]any{"secret": "abcdefgh"})
	masked, _ := out["secret"].(string)
	if masked[:2] != "ab" || masked[len(masked)-2:] != "gh" {
		t.Errorf("unexpected masked value: %q", masked)
	}
}

func TestRequestMasker_Mask_NonStringValue(t *testing.T) {
	m := NewRequestMasker(DefaultMaskerPolicy())
	_ = m.AddField("count")
	out := m.Mask(map[string]any{"count": 42})
	if out["count"] == 42 {
		t.Error("expected non-string value to be masked")
	}
}
