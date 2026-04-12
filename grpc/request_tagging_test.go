package grpc

import (
	"testing"
)

func TestNewTagStore_NotNil(t *testing.T) {
	s := NewTagStore(10)
	if s == nil {
		t.Fatal("expected non-nil TagStore")
	}
}

func TestNewTagStore_DefaultMaxTags(t *testing.T) {
	s := NewTagStore(0)
	if s.maxTags != 20 {
		t.Fatalf("expected default maxTags=20, got %d", s.maxTags)
	}
}

func TestTagStore_Add_And_Get(t *testing.T) {
	s := NewTagStore(10)
	if err := s.Add("req1", "alpha", "beta"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags := s.Get("req1")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestTagStore_Add_EmptyKey_ReturnsError(t *testing.T) {
	s := NewTagStore(10)
	if err := s.Add("", "tag"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestTagStore_Add_Deduplication(t *testing.T) {
	s := NewTagStore(10)
	_ = s.Add("req1", "alpha", "alpha", "beta")
	if len(s.Get("req1")) != 2 {
		t.Fatalf("expected 2 unique tags")
	}
}

func TestTagStore_Add_MaxTagsEnforced(t *testing.T) {
	s := NewTagStore(3)
	_ = s.Add("req1", "a", "b", "c", "d", "e")
	if len(s.Get("req1")) != 3 {
		t.Fatalf("expected at most 3 tags")
	}
}

func TestTagStore_Remove(t *testing.T) {
	s := NewTagStore(10)
	_ = s.Add("req1", "alpha")
	s.Remove("req1")
	if s.Len() != 0 {
		t.Fatalf("expected 0 entries after Remove")
	}
}

func TestTagStore_Keys_Sorted(t *testing.T) {
	s := NewTagStore(10)
	_ = s.Add("z", "tag")
	_ = s.Add("a", "tag")
	keys := s.Keys()
	if len(keys) != 2 || keys[0] != "a" || keys[1] != "z" {
		t.Fatalf("expected sorted keys [a z], got %v", keys)
	}
}

func TestTagStore_Len_Empty(t *testing.T) {
	s := NewTagStore(10)
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
}
