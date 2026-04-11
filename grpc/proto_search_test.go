package grpc

import (
	"testing"
)

func TestNewProtoSearcher_NotNil(t *testing.T) {
	s := NewProtoSearcher()
	if s == nil {
		t.Fatal("expected non-nil ProtoSearcher")
	}
}

func TestProtoSearcher_Len_Empty(t *testing.T) {
	s := NewProtoSearcher()
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
}

func TestProtoSearcher_Index_And_Len(t *testing.T) {
	s := NewProtoSearcher()
	s.Index(map[string][]string{
		"svc.Foo": {"Bar", "Baz"},
		"svc.Empty": {},
	})
	if s.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", s.Len())
	}
}

func TestProtoSearcher_Search_EmptyQuery_ReturnsAll(t *testing.T) {
	s := NewProtoSearcher()
	s.Index(map[string][]string{
		"svc.Foo": {"Bar", "Baz"},
	})
	results := s.Search("")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestProtoSearcher_Search_CaseInsensitive(t *testing.T) {
	s := NewProtoSearcher()
	s.Index(map[string][]string{
		"svc.UserService": {"GetUser", "ListUsers"},
	})
	results := s.Search("getuser")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Method != "GetUser" {
		t.Fatalf("expected GetUser, got %s", results[0].Method)
	}
}

func TestProtoSearcher_Search_NoMatch(t *testing.T) {
	s := NewProtoSearcher()
	s.Index(map[string][]string{
		"svc.Foo": {"Bar"},
	})
	results := s.Search("zzznomatch")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestProtoSearcher_Index_Reindex_Clears(t *testing.T) {
	s := NewProtoSearcher()
	s.Index(map[string][]string{"svc.A": {"M1"}})
	s.Index(map[string][]string{"svc.B": {"M2", "M3"}})
	if s.Len() != 2 {
		t.Fatalf("expected 2 after reindex, got %d", s.Len())
	}
}
