package grpc

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// RequestBuilder builds a proto message from a JSON payload and a method descriptor.
type RequestBuilder struct{}

// NewRequestBuilder returns a new RequestBuilder.
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{}
}

// Build constructs a proto.Message from a JSON string using the given method descriptor.
// It returns an error if the JSON is invalid or does not match the input type.
func (rb *RequestBuilder) Build(md protoreflect.MethodDescriptor, jsonPayload string) (proto.Message, error) {
	if md == nil {
		return nil, fmt.Errorf("method descriptor must not be nil")
	}

	inputDesc := md.Input()
	if inputDesc == nil {
		return nil, fmt.Errorf("method has no input descriptor")
	}

	// Validate that the payload is valid JSON.
	if !json.Valid([]byte(jsonPayload)) {
		return nil, fmt.Errorf("invalid JSON payload: %q", jsonPayload)
	}

	// Build a dynamic message from the descriptor.
	msgDesc, err := protodesc.NewFile(inputDesc.ParentFile().ParentFile(), nil)
	if err != nil {
		// Fall back to using the descriptor directly.
		_ = msgDesc
	}

	msg := dynamicpb.NewMessage(inputDesc)

	if err := unmarshalJSON([]byte(jsonPayload), msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON into %s: %w", inputDesc.FullName(), err)
	}

	return msg, nil
}

// unmarshalJSON unmarshals JSON bytes into a proto.Message using the protojson package.
func unmarshalJSON(data []byte, msg proto.Message) error {
	// Use a raw map to populate dynamic message fields.
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	// For dynamic messages we rely on jsonpb / protojson; here we do a best-effort
	// field-by-field population for primitive types.
	dynMsg, ok := msg.(*dynamicpb.Message)
	if !ok {
		return fmt.Errorf("expected *dynamicpb.Message, got %T", msg)
	}
	return populateDynamicMessage(dynMsg, raw)
}

// populateDynamicMessage sets fields on a dynamic proto message from a Go map.
func populateDynamicMessage(msg *dynamicpb.Message, raw map[string]interface{}) error {
	fields := msg.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		val, ok := raw[string(fd.Name())]
		if !ok {
			continue
		}
		pv, err := toProtoValue(fd, val)
		if err != nil {
			return fmt.Errorf("field %s: %w", fd.Name(), err)
		}
		msg.Set(fd, pv)
	}
	return nil
}

// toProtoValue converts a Go interface{} value to a protoreflect.Value.
func toProtoValue(fd protoreflect.FieldDescriptor, val interface{}) (protoreflect.Value, error) {
	switch fd.Kind() {
	case protoreflect.StringKind:
		s, ok := val.(string)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected string, got %T", val)
		}
		return protoreflect.ValueOfString(s), nil
	case protoreflect.BoolKind:
		b, ok := val.(bool)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected bool, got %T", val)
		}
		return protoreflect.ValueOfBool(b), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		f, ok := val.(float64 !ok {
			return protoreflect.Value{}, fmt.Errorf("expected number, got %T", val)
		}
		return protoreflect.ValueOfInt32
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		f, ok := val.(float64)
		if !ok {
	flect.Value{}, fmt.Errorf("expected number, got %T", val)
		}
		return protoreflect.ValueOfInt64(int64(f)), nil
	case protoreflect.FloatKind:
	float64)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected number, got %T", val)
		}
		return protoreflect.(f)), nil
	case protoreflect.DoubleKind:
		f, ok := val.(float64)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected number, got %T", val)
		}
		return protoreflect.ValueOfFloat64(f), nil
	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported field kind: %s", fd.Kind())
	}
}
