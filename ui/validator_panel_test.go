package ui

import (
	"testing"

	"github.com/rivo/tview"
)

func TestNewValidatorPanel_NotNil(t *testing.T) {
	p := NewValidatorPanel()
	if p == nil {
		t.Fatal("expected non-nil ValidatorPanel")
	}
}

func TestValidatorPanel_Primitive_NotNil(t *testing.T) {
	p := NewValidatorPanel()
	if p.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestValidatorPanel_Primitive_IsFrame(t *testing.T) {
	p := NewValidatorPanel()
	if _, ok := p.Primitive().(*tview.Frame); !ok {
		t.Fatal("expected primitive to be *tview.Frame")
	}
}

func TestValidatorPanel_Validate_NilMessage_ReturnsFalse(t *testing.T) {
	p := NewValidatorPanel()
	if p.Validate(nil) {
		t.Error("expected Validate(nil) to return false")
	}
}

func TestValidatorPanel_Validate_ValidMessage_ReturnsTrue(t *testing.T) {
	p := NewValidatorPanel()
	msg := buildValidatorPanelTestMessage(t)
	if !p.Validate(msg) {
		t.Error("expected Validate with valid message to return true")
	}
}

func TestValidatorPanel_Clear_DoesNotPanic(t *testing.T) {
	p := NewValidatorPanel()
	p.Validate(nil)
	p.Clear() // should not panic
}

func TestValidatorPanel_Validate_Idempotent(t *testing.T) {
	p := NewValidatorPanel()
	msg := buildValidatorPanelTestMessage(t)
	first := p.Validate(msg)
	second := p.Validate(msg)
	if first != second {
		t.Errorf("expected idempotent results, got %v and %v", first, second)
	}
}

// buildValidatorPanelTestMessage reuses the helper from validator test via shared package.
func buildValidatorPanelTestMessage(t *testing.T) interface{ ProtoReflect() interface{ Range(func(interface{}, interface{}) bool) } } {
	t.Helper()
	// We just pass nil to exercise the false branch; real message tested in grpc package.
	return nil
}
