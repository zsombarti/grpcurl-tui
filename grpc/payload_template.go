package grpc

import (
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// PayloadTemplateGenerator generates JSON template skeletons from proto descriptors.
type PayloadTemplateGenerator struct{}

// NewPayloadTemplateGenerator creates a new PayloadTemplateGenerator.
func NewPayloadTemplateGenerator() *PayloadTemplateGenerator {
	return &PayloadTemplateGenerator{}
}

// Generate returns a pretty-printed JSON template for the given message descriptor.
func (g *PayloadTemplateGenerator) Generate(md protoreflect.MessageDescriptor) (string, error) {
	if md == nil {
		return "", fmt.Errorf("message descriptor is nil")
	}
	m := g.buildMap(md)
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal template: %w", err)
	}
	return string(b), nil
}

// buildMap recursively constructs a map representing the message template.
func (g *PayloadTemplateGenerator) buildMap(md protoreflect.MessageDescriptor) map[string]any {
	result := make(map[string]any)
	fields := md.Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		name := string(fd.JSONName())
		result[name] = g.defaultValue(fd)
	}
	return result
}

// defaultValue returns a zero/placeholder value for a field descriptor.
func (g *PayloadTemplateGenerator) defaultValue(fd protoreflect.FieldDescriptor) any {
	if fd.IsList() {
		return []any{}
	}
	if fd.IsMap() {
		return map[string]any{}
	}
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return false
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return 0
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return 0.0
	case protoreflect.StringKind:
		return ""
	case protoreflect.BytesKind:
		return ""
	case protoreflect.EnumKind:
		return strings.ToLower(string(fd.Enum().Values().Get(0).Name()))
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return g.buildMap(fd.Message())
	default:
		return nil
	}
}
