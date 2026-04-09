package grpc

import (
	"testing"
)

func TestNewMetadataBuilder_NotNil(t *testing.T) {
	b := NewMetadataBuilder()
	if b == nil {
		t.Fatal("expected non-nil MetadataBuilder")
	}
}

func TestMetadataBuilder_Build_ColonSeparator(t *testing.T) {
	b := NewMetadataBuilder()
	md, err := b.Build([]string{"authorization: Bearer token123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vals := md.Get("authorization")
	if len(vals) == 0 || vals[0] != "Bearer token123" {
		t.Fatalf("expected 'Bearer token123', got %v", vals)
	}
}

func TestMetadataBuilder_Build_EqualsSeparator(t *testing.T) {
	b := NewMetadataBuilder()
	md, err := b.Build([]string{"x-request-id=abc-123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vals := md.Get("x-request-id")
	if len(vals) == 0 || vals[0] != "abc-123" {
		t.Fatalf("expected 'abc-123', got %v", vals)
	}
}

func TestMetadataBuilder_Build_MalformedPair(t *testing.T) {
	b := NewMetadataBuilder()
	_, err := b.Build([]string{"badentry"})
	if err == nil {
		t.Fatal("expected error for malformed pair")
	}
}

func TestMetadataBuilder_Build_KeyLowercased(t *testing.T) {
	b := NewMetadataBuilder()
	md, err := b.Build([]string{"Content-Type: application/json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(md.Get("content-type")) == 0 {
		t.Fatal("expected key to be lowercased")
	}
}

func TestMetadataBuilder_FromMap(t *testing.T) {
	b := NewMetadataBuilder()
	md := b.FromMap(map[string]string{"X-Trace-Id": "trace-999"})
	vals := md.Get("x-trace-id")
	if len(vals) == 0 || vals[0] != "trace-999" {
		t.Fatalf("expected 'trace-999', got %v", vals)
	}
}

func TestMetadataBuilder_Build_EmptySlice(t *testing.T) {
	b := NewMetadataBuilder()
	md, err := b.Build([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(md) != 0 {
		t.Fatalf("expected empty metadata, got %v", md)
	}
}
