package grpc

import (
	"testing"
)

func TestNewEnvSubstitutor_NotNil(t *testing.T) {
	s := NewEnvSubstitutor(nil, "")
	if s == nil {
		t.Fatal("expected non-nil EnvSubstitutor")
	}
}

func TestEnvSubstitutor_Len_Empty(t *testing.T) {
	s := NewEnvSubstitutor(nil, "")
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
}

func TestEnvSubstitutor_Len_WithEntries(t *testing.T) {
	s := NewEnvSubstitutor(map[string]string{"A": "1", "B": "2"}, "")
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestEnvSubstitutor_Substitute_BraceStyle(t *testing.T) {
	s := NewEnvSubstitutor(map[string]string{"HOST": "localhost"}, "")
	got := s.Substitute("grpc://${HOST}:50051")
	if got != "grpc://localhost:50051" {
		t.Fatalf("unexpected result: %s", got)
	}
}

func TestEnvSubstitutor_Substitute_DollarStyle(t *testing.T) {
	s := NewEnvSubstitutor(map[string]string{"PORT": "9090"}, "")
	got := s.Substitute("host:$PORT")
	if got != "host:9090" {
		t.Fatalf("unexpected result: %s", got)
	}
}

func TestEnvSubstitutor_Substitute_MissingVar_UsesFallback(t *testing.T) {
	s := NewEnvSubstitutor(map[string]string{}, "MISSING")
	got := s.Substitute("${UNDEFINED}")
	if got != "MISSING" {
		t.Fatalf("expected MISSING, got %s", got)
	}
}

func TestEnvSubstitutor_Substitute_NoPlaceholders(t *testing.T) {
	s := NewEnvSubstitutor(map[string]string{"A": "1"}, "")
	input := "plain string"
	got := s.Substitute(input)
	if got != input {
		t.Fatalf("expected unchanged string, got %s", got)
	}
}

func TestEnvSubstitutor_SubstituteMap(t *testing.T) {
	s := NewEnvSubstitutor(map[string]string{"TOKEN": "abc123"}, "")
	pairs := map[string]string{
		"authorization": "Bearer ${TOKEN}",
		"x-custom":      "static",
	}
	out := s.SubstituteMap(pairs)
	if out["authorization"] != "Bearer abc123" {
		t.Fatalf("unexpected authorization: %s", out["authorization"])
	}
	if out["x-custom"] != "static" {
		t.Fatalf("unexpected x-custom: %s", out["x-custom"])
	}
}
