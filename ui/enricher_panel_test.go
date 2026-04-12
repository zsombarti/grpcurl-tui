package ui

import (
	"testing"

	"github.com/rivo/tview"
)

func newTestEnricherPanel() *EnricherPanel {
	return NewEnricherPanel()
}

func TestNewEnricherPanel_NotNil(t *testing.T) {
	p := newTestEnricherPanel()
	if p == nil {
		t.Fatal("expected non-nil EnricherPanel")
	}
}

func TestEnricherPanel_Primitive_NotNil(t *testing.T) {
	p := newTestEnricherPanel()
	if p.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestEnricherPanel_Primitive_IsFlex(t *testing.T) {
	p := newTestEnricherPanel()
	if _, ok := p.Primitive().(*tview.Flex); !ok {
		t.Fatal("expected primitive to be *tview.Flex")
	}
}

func TestEnricherPanel_StepCount_Initial_Zero(t *testing.T) {
	p := newTestEnricherPanel()
	if p.StepCount() != 0 {
		t.Fatalf("expected 0, got %d", p.StepCount())
	}
}

func TestEnricherPanel_Refresh_Empty(t *testing.T) {
	p := newTestEnricherPanel()
	p.Refresh([]string{})
	if p.StepCount() != 0 {
		t.Fatalf("expected 0 after empty refresh, got %d", p.StepCount())
	}
}

func TestEnricherPanel_Refresh_WithSteps(t *testing.T) {
	p := newTestEnricherPanel()
	p.Refresh([]string{"add-env", "add-source", "add-trace-id"})
	if p.StepCount() != 3 {
		t.Fatalf("expected 3, got %d", p.StepCount())
	}
}

func TestEnricherPanel_Refresh_Idempotent(t *testing.T) {
	p := newTestEnricherPanel()
	p.Refresh([]string{"step-a", "step-b"})
	p.Refresh([]string{"step-a", "step-b"})
	if p.StepCount() != 2 {
		t.Fatalf("expected 2 after idempotent refresh, got %d", p.StepCount())
	}
}

func TestEnricherPanel_Refresh_ClearsOldSteps(t *testing.T) {
	p := newTestEnricherPanel()
	p.Refresh([]string{"old-step"})
	p.Refresh([]string{})
	if p.StepCount() != 0 {
		t.Fatalf("expected 0 after clear refresh, got %d", p.StepCount())
	}
}
