package grpc

import (
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	"google.golang.org/protobuf/types/descriptorpb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
)

func TestNewRequestValidator_NotNil(t *testing.T) {
	v := NewRequestValidator()
	if v == nil {
		t.Fatal("expected non-nil RequestValidator")
	}
}

func TestRequestValidator_Validate_NilMessage(t *testing.T) {
	v := NewRequestValidator()
	errs := v.Validate(nil)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error for nil message, got %d", len(errs))
	}
	if errs[0].Field != "(root)" {
		t.Errorf("unexpected field: %s", errs[0].Field)
	}
}

func TestRequestValidator_Validate_EmptyMessage_NoErrors(t *testing.T) {
	msg := buildValidatorTestMessage(t)
	v := NewRequestValidator()
	errs := v.Validate(msg)
	if len(errs) != 0 {
		t.Fatalf("expected no errors for empty optional-field message, got %d", len(errs))
	}
}

func TestRequestValidator_IsValid_True(t *testing.T) {
	msg := buildValidatorTestMessage(t)
	v := NewRequestValidator()
	if !v.IsValid(msg) {
		t.Error("expected IsValid to return true")
	}
}

func TestRequestValidator_Summary_NilMessage_ReturnsError(t *testing.T) {
	v := NewRequestValidator()
	if err := v.Summary(nil); err == nil {
		t.Error("expected non-nil error for nil message")
	}
}

func TestRequestValidator_Summary_ValidMessage_ReturnsNil(t *testing.T) {
	msg := buildValidatorTestMessage(t)
	v := NewRequestValidator()
	if err := v.Summary(msg); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// buildValidatorTestMessage creates a simple dynamic message for testing.
func buildValidatorTestMessage(t *testing.T) protoreflect.Message {
	t.Helper()
	fdp := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("validator_test.proto"),
		Syntax:  proto.String("proto3"),
		Package: proto.String("validatortest"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("TestMsg"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:   proto.String("value"),
						Number: proto.Int32(1),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					},
				},
			},
		},
	}
	fd, err := protodesc.NewFile(fdp, nil)
	if err != nil {
		t.Fatalf("failed to build file descriptor: %v", err)
	}
	md := fd.Messages().Get(0)
	return dynamicpb.NewMessage(md)
}
