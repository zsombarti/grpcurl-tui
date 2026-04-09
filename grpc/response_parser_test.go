package grpc

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewResponseParser_NotNil(t *testing.T) {
	p := NewResponseParser()
	if p == nil {
		t.Fatal("expected non-nil ResponseParser")
	}
}

func TestResponseParser_ToJSON_NilMessage(t *testing.T) {
	p := NewResponseParser()
	_, err := p.ToJSON(nil)
	if err == nil {
		t.Fatal("expected error for nil message, got nil")
	}
}

func TestResponseParser_ToJSON_EmptyMessage(t *testing.T) {
	p := NewResponseParser()
	msg := &emptypb.Empty{}

	result, err := p.ToJSON(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "{") {
		t.Errorf("expected JSON object, got: %s", result)
	}
}

func TestResponseParser_ToJSON_WithValue(t *testing.T) {
	p := NewResponseParser()
	msg := wrapperspb.String("hello")

	result, err := p.ToJSON(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "hello") {
		t.Errorf("expected 'hello' in JSON output, got: %s", result)
	}
}

func TestResponseParser_ToMap_NilMessage(t *testing.T) {
	p := NewResponseParser()
	_, err := p.ToMap(nil)
	if err == nil {
		t.Fatal("expected error for nil message")
	}
}

func TestResponseParser_ToMap_ValidMessage(t *testing.T) {
	p := NewResponseParser()
	var msg proto.Message = &emptypb.Empty{}

	result, err := p.ToMap(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil map")
	}
}

func TestResponseParser_FormatError_Nil(t *testing.T) {
	p := NewResponseParser()
	result := p.FormatError(nil)
	if result != "" {
		t.Errorf("expected empty string for nil error, got: %s", result)
	}
}

func TestResponseParser_FormatError_WithError(t *testing.T) {
	p := NewResponseParser()
	result := p.FormatError(fmt.Errorf("connection refused"))
	if !strings.Contains(result, "connection refused") {
		t.Errorf("expected error message in output, got: %s", result)
	}
}
