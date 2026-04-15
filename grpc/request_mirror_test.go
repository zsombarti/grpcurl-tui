package grpc

import (
	"testing"
)

func TestNewRequestMirror_NotNil(t *testing.T) {
	m := NewRequestMirror(4)
	if m == nil {
		t.Fatal("expected non-nil RequestMirror")
	}
}

func TestNewRequestMirror_DefaultMaxSize(t *testing.T) {
	m := NewRequestMirror(0)
	for i := 0; i < 8; i++ {
		if err := m.AddTarget(MirrorTarget{Address: "localhost:500" + string(rune('0'+i)), Weight: 10}); err != nil {
			t.Fatalf("unexpected error on target %d: %v", i, err)
		}
	}
	if m.Len() != 8 {
		t.Fatalf("expected 8 targets, got %d", m.Len())
	}
}

func TestRequestMirror_Len_Empty(t *testing.T) {
	m := NewRequestMirror(4)
	if m.Len() != 0 {
		t.Fatalf("expected 0, got %d", m.Len())
	}
}

func TestRequestMirror_AddTarget_And_Len(t *testing.T) {
	m := NewRequestMirror(4)
	if err := m.AddTarget(MirrorTarget{Address: "localhost:9090", Weight: 50}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 1 {
		t.Fatalf("expected 1, got %d", m.Len())
	}
}

func TestRequestMirror_AddTarget_EmptyAddress_ReturnsError(t *testing.T) {
	m := NewRequestMirror(4)
	if err := m.AddTarget(MirrorTarget{Address: "", Weight: 10}); err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestRequestMirror_AddTarget_InvalidWeight_ReturnsError(t *testing.T) {
	m := NewRequestMirror(4)
	if err := m.AddTarget(MirrorTarget{Address: "localhost:9090", Weight: 0}); err == nil {
		t.Fatal("expected error for weight 0")
	}
	if err := m.AddTarget(MirrorTarget{Address: "localhost:9090", Weight: 101}); err == nil {
		t.Fatal("expected error for weight 101")
	}
}

func TestRequestMirror_Full_ReturnsError(t *testing.T) {
	m := NewRequestMirror(2)
	_ = m.AddTarget(MirrorTarget{Address: "localhost:9001", Weight: 50})
	_ = m.AddTarget(MirrorTarget{Address: "localhost:9002", Weight: 50})
	if err := m.AddTarget(MirrorTarget{Address: "localhost:9003", Weight: 50}); err == nil {
		t.Fatal("expected error when mirror is full")
	}
}

func TestRequestMirror_Targets_ReturnsCopy(t *testing.T) {
	m := NewRequestMirror(4)
	_ = m.AddTarget(MirrorTarget{Address: "localhost:9090", Weight: 100})
	targets := m.Targets()
	targets[0].Address = "mutated"
	if m.Targets()[0].Address == "mutated" {
		t.Fatal("Targets should return a copy, not a reference")
	}
}

func TestRequestMirror_Clear_ResetsLen(t *testing.T) {
	m := NewRequestMirror(4)
	_ = m.AddTarget(MirrorTarget{Address: "localhost:9090", Weight: 50})
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("expected 0 after Clear, got %d", m.Len())
	}
}
