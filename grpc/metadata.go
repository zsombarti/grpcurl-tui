package grpc

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

// MetadataBuilder constructs gRPC metadata from key-value string pairs.
type MetadataBuilder struct{}

// NewMetadataBuilder returns a new MetadataBuilder.
func NewMetadataBuilder() *MetadataBuilder {
	return &MetadataBuilder{}
}

// Build parses a slice of "key: value" or "key=value" strings into gRPC metadata.
// Returns an error if any entry is malformed.
func (m *MetadataBuilder) Build(pairs []string) (metadata.MD, error) {
	md := metadata.New(nil)
	for _, pair := range pairs {
		key, value, err := splitPair(pair)
		if err != nil {
			return nil, err
		}
		md.Append(strings.ToLower(strings.TrimSpace(key)), strings.TrimSpace(value))
	}
	return md, nil
}

// FromMap converts a map[string]string into gRPC metadata.
func (m *MetadataBuilder) FromMap(kv map[string]string) metadata.MD {
	md := metadata.New(nil)
	for k, v := range kv {
		md.Append(strings.ToLower(strings.TrimSpace(k)), strings.TrimSpace(v))
	}
	return md
}

// splitPair splits a string by ":" or "=" into key and value.
func splitPair(s string) (string, string, error) {
	for _, sep := range []string{":", "="} {
		if idx := strings.Index(s, sep); idx > 0 {
			return s[:idx], s[idx+1:], nil
		}
	}
	return "", "", fmt.Errorf("metadata: malformed pair %q, expected 'key: value' or 'key=value'", s)
}
