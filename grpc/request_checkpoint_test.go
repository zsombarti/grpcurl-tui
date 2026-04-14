package grpc

import (
	"testing"
)

func TestNewRequestCheckpointer_NotNil(t *testing.T) {
	c := NewRequestCheckpointer(0)
	if c == nil {
		t.Fatal("expected non-nil checkpointer")
	}
}

func TestNewRequestCheckpointer_DefaultMaxSize(t *testing.T) {
	c := NewRequestCheckpointer(0)
	if c.maxSize != defaultCheckpointMaxSize {
		t.Fatalf("expected %d, got %d", defaultCheckpointMaxSize, c.maxSize)
	}
}

func TestRequestCheckpointer_Len_Empty(t *testing.T) {
	c := NewRequestCheckpointer(10)
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
}

func TestRequestCheckpointer_Save_And_Get(t *testing.T) {
	c := NewRequestCheckpointer(10)
	payload := map[string]interface{}{"key": "value"}
	if err := c.Save("cp1", "pkg.Service/Method", payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cp, ok := c.Get("cp1")
	if !ok {
		t.Fatal("expected checkpoint to be found")
	}
	if cp.Name != "cp1" || cp.Method != "pkg.Service/Method" {
		t.Fatalf("unexpected checkpoint: %+v", cp)
	}
	if cp.CreatedAt.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestRequestCheckpointer_Save_EmptyName_ReturnsError(t *testing.T) {
	c := NewRequestCheckpointer(10)
	if err := c.Save("", "Method", nil); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRequestCheckpointer_Save_EmptyMethod_ReturnsError(t *testing.T) {
	c := NewRequestCheckpointer(10)
	if err := c.Save("cp1", "", nil); err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestRequestCheckpointer_Save_Replace_ExistingName(t *testing.T) {
	c := NewRequestCheckpointer(10)
	_ = c.Save("cp1", "MethodA", nil)
	_ = c.Save("cp1", "MethodB", nil)
	if c.Len() != 1 {
		t.Fatalf("expected 1 after replace, got %d", c.Len())
	}
	cp, _ := c.Get("cp1")
	if cp.Method != "MethodB" {
		t.Fatalf("expected MethodB, got %s", cp.Method)
	}
}

func TestRequestCheckpointer_Eviction(t *testing.T) {
	c := NewRequestCheckpointer(2)
	_ = c.Save("a", "M", nil)
	_ = c.Save("b", "M", nil)
	_ = c.Save("c", "M", nil)
	if c.Len() != 2 {
		t.Fatalf("expected 2 after eviction, got %d", c.Len())
	}
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected 'a' to be evicted")
	}
}

func TestRequestCheckpointer_Delete(t *testing.T) {
	c := NewRequestCheckpointer(10)
	_ = c.Save("cp1", "M", nil)
	if !c.Delete("cp1") {
		t.Fatal("expected delete to return true")
	}
	if c.Len() != 0 {
		t.Fatalf("expected 0 after delete, got %d", c.Len())
	}
	if c.Delete("cp1") {
		t.Fatal("expected delete of missing key to return false")
	}
}

func TestRequestCheckpointer_All(t *testing.T) {
	c := NewRequestCheckpointer(10)
	_ = c.Save("x", "M", nil)
	_ = c.Save("y", "M", nil)
	all := c.All()
	if len(all) != 2 {
		t.Fatalf("expected 2, got %d", len(all))
	}
}
