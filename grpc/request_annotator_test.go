package grpc

import (
	"testing"
)

func TestNewRequestAnnotator_NotNil(t *testing.T) {
	a := NewRequestAnnotator(0)
	if a == nil {
		t.Fatal("expected non-nil annotator")
	}
}

func TestNewRequestAnnotator_DefaultMaxSize(t *testing.T) {
	a := NewRequestAnnotator(0)
	if a.maxSize != defaultAnnotatorMaxSize {
		t.Fatalf("expected maxSize %d, got %d", defaultAnnotatorMaxSize, a.maxSize)
	}
}

func TestRequestAnnotator_Len_Empty(t *testing.T) {
	a := NewRequestAnnotator(10)
	if a.Len() != 0 {
		t.Fatalf("expected 0, got %d", a.Len())
	}
}

func TestRequestAnnotator_Annotate_And_Len(t *testing.T) {
	a := NewRequestAnnotator(10)
	if err := a.Annotate("/pkg.Svc/Method", "note", "hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Len() != 1 {
		t.Fatalf("expected 1, got %d", a.Len())
	}
}

func TestRequestAnnotator_Annotate_EmptyMethod_ReturnsError(t *testing.T) {
	a := NewRequestAnnotator(10)
	if err := a.Annotate("", "key", "val"); err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRequestAnnotator_Annotate_EmptyKey_ReturnsError(t *testing.T) {
	a := NewRequestAnnotator(10)
	if err := a.Annotate("/pkg.Svc/Method", "", "val"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestRequestAnnotator_Get_ExistingKey(t *testing.T) {
	a := NewRequestAnnotator(10)
	_ = a.Annotate("/pkg.Svc/Method", "env", "production")
	v, ok := a.Get("/pkg.Svc/Method", "env")
	if !ok {
		t.Fatal("expected annotation to be found")
	}
	if v != "production" {
		t.Fatalf("expected 'production', got %q", v)
	}
}

func TestRequestAnnotator_Get_MissingKey(t *testing.T) {
	a := NewRequestAnnotator(10)
	_, ok := a.Get("/pkg.Svc/Method", "missing")
	if ok {
		t.Fatal("expected annotation to be absent")
	}
}

func TestRequestAnnotator_Annotate_Overwrite(t *testing.T) {
	a := NewRequestAnnotator(10)
	_ = a.Annotate("/pkg.Svc/Method", "env", "staging")
	_ = a.Annotate("/pkg.Svc/Method", "env", "production")
	if a.Len() != 1 {
		t.Fatalf("expected 1 after overwrite, got %d", a.Len())
	}
	v, _ := a.Get("/pkg.Svc/Method", "env")
	if v != "production" {
		t.Fatalf("expected overwritten value 'production', got %q", v)
	}
}

func TestRequestAnnotator_Eviction(t *testing.T) {
	a := NewRequestAnnotator(2)
	_ = a.Annotate("/svc/A", "k", "1")
	_ = a.Annotate("/svc/B", "k", "2")
	_ = a.Annotate("/svc/C", "k", "3")
	if a.Len() != 2 {
		t.Fatalf("expected 2 after eviction, got %d", a.Len())
	}
}

func TestRequestAnnotator_Clear(t *testing.T) {
	a := NewRequestAnnotator(10)
	_ = a.Annotate("/svc/A", "k", "v")
	a.Clear()
	if a.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", a.Len())
	}
}
