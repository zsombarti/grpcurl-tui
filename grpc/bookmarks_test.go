package grpc

import (
	"testing"
)

func TestNewBookmarkStore_Defaults(t *testing.T) {
	bs := NewBookmarkStore(0)
	if bs == nil {
		t.Fatal("expected non-nil BookmarkStore")
	}
	if bs.maxSize != 50 {
		t.Errorf("expected default maxSize 50, got %d", bs.maxSize)
	}
}

func TestBookmarkStore_Add_And_Len(t *testing.T) {
	bs := NewBookmarkStore(10)
	err := bs.Add(Bookmark{Name: "local", Address: "localhost:50051", Method: "SayHello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bs.Len() != 1 {
		t.Errorf("expected len 1, got %d", bs.Len())
	}
}

func TestBookmarkStore_Add_MissingName(t *testing.T) {
	bs := NewBookmarkStore(10)
	err := bs.Add(Bookmark{Address: "localhost:50051"})
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestBookmarkStore_Add_MissingAddress(t *testing.T) {
	bs := NewBookmarkStore(10)
	err := bs.Add(Bookmark{Name: "test"})
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestBookmarkStore_Eviction(t *testing.T) {
	bs := NewBookmarkStore(2)
	_ = bs.Add(Bookmark{Name: "a", Address: "localhost:1"})
	_ = bs.Add(Bookmark{Name: "b", Address: "localhost:2"})
	_ = bs.Add(Bookmark{Name: "c", Address: "localhost:3"})
	if bs.Len() != 2 {
		t.Errorf("expected len 2 after eviction, got %d", bs.Len())
	}
	all := bs.All()
	if all[0].Name != "b" {
		t.Errorf("expected oldest evicted, got %s", all[0].Name)
	}
}

func TestBookmarkStore_Delete(t *testing.T) {
	bs := NewBookmarkStore(10)
	_ = bs.Add(Bookmark{Name: "keep", Address: "localhost:1"})
	_ = bs.Add(Bookmark{Name: "remove", Address: "localhost:2"})
	err := bs.Delete("remove")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bs.Len() != 1 {
		t.Errorf("expected len 1, got %d", bs.Len())
	}
}

func TestBookmarkStore_Delete_NotFound(t *testing.T) {
	bs := NewBookmarkStore(10)
	err := bs.Delete("ghost")
	if err == nil {
		t.Fatal("expected error for missing bookmark")
	}
}

func TestBookmarkStore_Clear(t *testing.T) {
	bs := NewBookmarkStore(10)
	_ = bs.Add(Bookmark{Name: "x", Address: "localhost:1"})
	bs.Clear()
	if bs.Len() != 0 {
		t.Errorf("expected len 0 after clear, got %d", bs.Len())
	}
}
