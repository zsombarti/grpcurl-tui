package ui

import (
	"errors"
	"testing"
)

func TestNewStreamPanel_NotNil(t *testing.T) {
	p := NewStreamPanel()
	if p == nil {
		t.Fatal("expected non-nil StreamPanel")
	}
}

func TestStreamPanel_InitialLen_Zero(t *testing.T) {
	p := NewStreamPanel()
	if p.Len() != 0 {
		t.Fatalf("expected 0 lines, got %d", p.Len())
	}
}

func TestStreamPanel_AppendMessage_IncrementsLen(t *testing.T) {
	p := NewStreamPanel()
	p.AppendMessage(0, `{"field":"value"}`)
	p.AppendMessage(1, `{"field":"other"}`)
	if p.Len() != 2 {
		t.Fatalf("expected 2 lines, got %d", p.Len())
	}
}

func TestStreamPanel_AppendError_IncrementsLen(t *testing.T) {
	p := NewStreamPanel()
	p.AppendError(errors.New("connection reset"))
	if p.Len() != 1 {
		t.Fatalf("expected 1 line after error, got %d", p.Len())
	}
}

func TestStreamPanel_Clear_ResetsLen(t *testing.T) {
	p := NewStreamPanel()
	p.AppendMessage(0, "hello")
	p.AppendMessage(1, "world")
	p.Clear()
	if p.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", p.Len())
	}
}

func TestStreamPanel_Primitive_NotNil(t *testing.T) {
	p := NewStreamPanel()
	if p.TextView == nil {
		t.Fatal("expected non-nil underlying TextView")
	}
}
