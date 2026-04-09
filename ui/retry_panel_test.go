package ui

import (
	"testing"
	"time"
)

func TestNewRetryPanel_NotNil(t *testing.T) {
	p := NewRetryPanel()
	if p == nil {
		t.Fatal("expected non-nil RetryPanel")
	}
}

func TestRetryPanel_GetPolicy_Defaults(t *testing.T) {
	p := NewRetryPanel()
	policy := p.GetPolicy()

	if policy.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", policy.MaxAttempts)
	}
	if policy.InitialDelay != 100*time.Millisecond {
		t.Errorf("expected InitialDelay=100ms, got %v", policy.InitialDelay)
	}
	if policy.MaxDelay != 2*time.Second {
		t.Errorf("expected MaxDelay=2s, got %v", policy.MaxDelay)
	}
	if policy.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", policy.Multiplier)
	}
}

func TestRetryPanel_GetPolicy_InvalidMaxAttempts_FallsBackToDefault(t *testing.T) {
	p := NewRetryPanel()
	p.maxAttemptsField.SetText("not-a-number")
	policy := p.GetPolicy()
	if policy.MaxAttempts != 3 {
		t.Errorf("expected fallback MaxAttempts=3, got %d", policy.MaxAttempts)
	}
}

func TestRetryPanel_GetPolicy_ZeroMaxAttempts_FallsBackToDefault(t *testing.T) {
	p := NewRetryPanel()
	p.maxAttemptsField.SetText("0")
	policy := p.GetPolicy()
	if policy.MaxAttempts != 3 {
		t.Errorf("expected fallback MaxAttempts=3, got %d", policy.MaxAttempts)
	}
}

func TestRetryPanel_GetPolicy_InvalidMultiplier_FallsBackToDefault(t *testing.T) {
	p := NewRetryPanel()
	p.multiplierField.SetText("0.5")
	policy := p.GetPolicy()
	if policy.Multiplier != 2.0 {
		t.Errorf("expected fallback Multiplier=2.0, got %f", policy.Multiplier)
	}
}

func TestRetryPanel_GetPolicy_CustomValues(t *testing.T) {
	p := NewRetryPanel()
	p.maxAttemptsField.SetText("5")
	p.initialDelayField.SetText("200")
	p.maxDelayField.SetText("5000")
	p.multiplierField.SetText("3.0")

	policy := p.GetPolicy()
	if policy.MaxAttempts != 5 {
		t.Errorf("expected MaxAttempts=5, got %d", policy.MaxAttempts)
	}
	if policy.InitialDelay != 200*time.Millisecond {
		t.Errorf("expected InitialDelay=200ms, got %v", policy.InitialDelay)
	}
	if policy.MaxDelay != 5*time.Second {
		t.Errorf("expected MaxDelay=5s, got %v", policy.MaxDelay)
	}
	if policy.Multiplier != 3.0 {
		t.Errorf("expected Multiplier=3.0, got %f", policy.Multiplier)
	}
}
