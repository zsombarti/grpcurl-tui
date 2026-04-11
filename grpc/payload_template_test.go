package grpc

import (
	"encoding/json"
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/proto"
)

func TestNewPayloadTemplateGenerator_NotNil(t *testing.T) {
	g := NewPayloadTemplateGenerator()
	if g == nil {
		t.Fatal("expected non-nil generator")
	}
}

func TestPayloadTemplateGenerator_Generate_NilDescriptor(t *testing.T) {
	g := NewPayloadTemplateGenerator()
	_, err := g.Generate(nil)
	if err == nil {
		t.Fatal("expected error for nil descriptor")
	}
}

func TestPayloadTemplateGenerator_Generate_SimpleMessage(t *testing.T) {
	md := buildSimpleMessageDescriptor(t)
	g := NewPayloadTemplateGenerator()
	out, err := g.Generate(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "name") {
		t.Errorf("expected 'name' field in output, got: %s", out)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestPayloadTemplateGenerator_Generate_IsValidJSON(t *testing.T) {
	md := buildSimpleMessageDescriptor(t)
	g := NewPayloadTemplateGenerator()
	out, err := g.Generate(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var v any
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("generated template is not valid JSON: %v", err)
	}
}

// buildSimpleMessageDescriptor creates a minimal proto message descriptor for testing.
func buildSimpleMessageDescriptor(t *testing.T) protoreflect.MessageDescriptor {
	t.Helper()
	fdp := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("test.proto"),
		Syntax:  proto.String("proto3"),
		Package: proto.String("test"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("TestMsg"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:   proto.String("name"),
						Number: proto.Int32(1),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					},
					{
						Name:   proto.String("count"),
						Number: proto.Int32(2),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
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
	_ = dynamicpb.NewMessage(fd.Messages().Get(0))
	return fd.Messages().Get(0)
}
