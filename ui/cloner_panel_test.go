package ui

import (
	"strings"
	"testing"

	grpcpkg "grpcurl-tui/grpc"
)

func newTestClonerPanel() *ClonerPanel {
	return NewClonerPanel(grpcpkg.NewRequestCloner(grpcpkg.DefaultClonerPolicy()))
}

func TestNewClonerPanel_NotNil(t *testing.T) {
	p := newTestClonerPanel()
	if p == nil {
		t.Fatal("expected non-nil panel")
	}
}

func TestClonerPanel_Primitive_NotNil(t *testing.T) {
	p := newTestClonerPanel()
	if p.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestClonerPanel_NilCloner_FallsBackToDefault(t *testing.T) {
	p := NewClonerPanel(nil)
	if p == nil {
		t.Fatal("expected non-nil panel even with nil cloner")
	}
}

func TestClonerPanel_Clone_ValidPayload(t *testing.T) {
	p := newTestClonerPanel()
	p.SetInput(`{"hello":"world"}`)
	if err := p.Clone(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := p.GetOutput()
	if !strings.Contains(out, "hello") {
		t.Fatalf("expected output to contain key, got: %s", out)
	}
}

func TestClonerPanel_Clone_InvalidJSON_ReturnsError(t *testing.T) {
	p := newTestClonerPanel()
	p.SetInput(`not json`)
	if err := p.Clone(); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestClonerPanel_Clear_ResetsOutput(t *testing.T) {
	p := newTestClonerPanel()
	p.SetInput(`{"a":1}`)
	_ = p.Clone()
	p.Clear()
	if p.GetOutput() != "" {
		t.Fatal("expected empty output after clear")
	}
}

func TestClonerPanel_SetInput_And_GetOutput_Roundtrip(t *testing.T) {
	p := newTestClonerPanel()
	p.SetInput(`{"x":"y"}`)
	_ = p.Clone()
	out := p.GetOutput()
	if !strings.Contains(out, "\"x\"") {
		t.Fatalf("expected key in output, got: %s", out)
	}
}
