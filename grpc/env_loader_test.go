package grpc

import (
	"os"
	"testing"
)

func TestNewEnvLoader_NotNil(t *testing.T) {
	el := NewEnvLoader()
	if el == nil {
		t.Fatal("expected non-nil EnvLoader")
	}
}

func TestEnvLoader_Len_Empty(t *testing.T) {
	el := NewEnvLoader()
	if el.Len() != 0 {
		t.Fatalf("expected 0, got %d", el.Len())
	}
}

func TestEnvLoader_LoadFile_ValidPairs(t *testing.T) {
	f := writeTempEnv(t, "KEY=value\n# comment\n\nFOO=bar\n")
	el := NewEnvLoader()
	if err := el.LoadFile(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if el.Len() != 2 {
		t.Fatalf("expected 2 pairs, got %d", el.Len())
	}
	pairs := el.Pairs()
	if pairs[0] != "KEY=value" || pairs[1] != "FOO=bar" {
		t.Fatalf("unexpected pairs: %v", pairs)
	}
}

func TestEnvLoader_LoadFile_MissingFile(t *testing.T) {
	el := NewEnvLoader()
	err := el.LoadFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestEnvLoader_LoadFile_MalformedLine(t *testing.T) {
	f := writeTempEnv(t, "BADLINE\n")
	el := NewEnvLoader()
	err := el.LoadFile(f)
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestEnvLoader_Clear(t *testing.T) {
	f := writeTempEnv(t, "A=1\nB=2\n")
	el := NewEnvLoader()
	_ = el.LoadFile(f)
	el.Clear()
	if el.Len() != 0 {
		t.Fatalf("expected 0 after Clear, got %d", el.Len())
	}
}

func TestEnvLoader_Pairs_ReturnsCopy(t *testing.T) {
	f := writeTempEnv(t, "X=1\n")
	el := NewEnvLoader()
	_ = el.LoadFile(f)
	p := el.Pairs()
	p[0] = "mutated"
	if el.Pairs()[0] == "mutated" {
		t.Fatal("Pairs should return a copy, not a reference")
	}
}

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := tmp.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	_ = tmp.Close()
	return tmp.Name()
}
