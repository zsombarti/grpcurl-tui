package grpc

import (
	"errors"
	"testing"
	"time"
)

func TestNewRequestLogger_NotNil(t *testing.T) {
	l := NewRequestLogger(0)
	if l == nil {
		t.Fatal("expected non-nil RequestLogger")
	}
}

func TestNewRequestLogger_DefaultMaxSize(t *testing.T) {
	l := NewRequestLogger(0)
	if l.maxSize != defaultRequestLogMaxSize {
		t.Fatalf("expected maxSize %d, got %d", defaultRequestLogMaxSize, l.maxSize)
	}
}

func TestRequestLogger_Len_Empty(t *testing.T) {
	l := NewRequestLogger(10)
	if l.Len() != 0 {
		t.Fatalf("expected 0, got %d", l.Len())
	}
}

func TestRequestLogger_Log_And_Len(t *testing.T) {
	l := NewRequestLogger(10)
	l.Log("pkg.Svc/Method", `{}`, `{"ok":true}`, nil, 5*time.Millisecond)
	if l.Len() != 1 {
		t.Fatalf("expected 1, got %d", l.Len())
	}
}

func TestRequestLogger_Eviction(t *testing.T) {
	l := NewRequestLogger(3)
	for i := 0; i < 5; i++ {
		l.Log("m", "", "", nil, 0)
	}
	if l.Len() != 3 {
		t.Fatalf("expected 3 after eviction, got %d", l.Len())
	}
}

func TestRequestLogger_Log_ErrorEntry(t *testing.T) {
	l := NewRequestLogger(10)
	l.Log("m", "", "", errors.New("boom"), 0)
	entries := l.Entries()
	if entries[0].Error != "boom" {
		t.Fatalf("expected error string 'boom', got %q", entries[0].Error)
	}
}

func TestRequestLogger_Clear(t *testing.T) {
	l := NewRequestLogger(10)
	l.Log("m", "", "", nil, 0)
	l.Clear()
	if l.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", l.Len())
	}
}

func TestRequestLogger_Summary_OutOfRange(t *testing.T) {
	l := NewRequestLogger(10)
	_, err := l.Summary(0)
	if err == nil {
		t.Fatal("expected error for out-of-range index")
	}
}

func TestRequestLogger_Summary_Format(t *testing.T) {
	l := NewRequestLogger(10)
	l.Log("pkg.Svc/Hello", `{}`, `{}`, nil, 12*time.Millisecond)
	s, err := l.Summary(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
