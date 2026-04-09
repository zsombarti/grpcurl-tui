package grpc

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ResponseParser converts protobuf messages into human-readable formats.
type ResponseParser struct {
	indent string
}

// NewResponseParser creates a new ResponseParser with default indentation.
func NewResponseParser() *ResponseParser {
	return &ResponseParser{indent: "  "}
}

// ToJSON converts a proto.Message to a pretty-printed JSON string.
func (p *ResponseParser) ToJSON(msg proto.Message) (string, error) {
	if msg == nil {
		return "", fmt.Errorf("response message is nil")
	}

	m := protojson.MarshalOptions{
		EmitUnpopulated: true,
		Indent:          p.indent,
		UseProtoNames:   true,
	}

	b, err := m.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response to JSON: %w", err)
	}

	return string(b), nil
}

// ToMap converts a proto.Message to a map[string]interface{} for further processing.
func (p *ResponseParser) ToMap(msg proto.Message) (map[string]interface{}, error) {
	jsonStr, err := p.ToJSON(msg)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to map: %w", err)
	}

	return result, nil
}

// FormatError wraps a gRPC error into a structured JSON-like string for display.
func (p *ResponseParser) FormatError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("{\n%s\"error\": \"%s\"\n}", p.indent, err.Error())
}
