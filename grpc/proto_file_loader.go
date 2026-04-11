package grpc

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// ProtoFileLoader loads and parses .proto descriptor files from disk.
type ProtoFileLoader struct {
	files map[string]*descriptorpb.FileDescriptorProto
}

// NewProtoFileLoader returns a new ProtoFileLoader.
func NewProtoFileLoader() *ProtoFileLoader {
	return &ProtoFileLoader{
		files: make(map[string]*descriptorpb.FileDescriptorProto),
	}
}

// LoadFile reads a compiled .pb binary descriptor file and registers it.
func (p *ProtoFileLoader) LoadFile(path string) error {
	if path == "" {
		return errors.New("proto_file_loader: path must not be empty")
	}
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".pb" && ext != ".bin" {
		return errors.New("proto_file_loader: unsupported file extension, expected .pb or .bin")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fdp := &descriptorpb.FileDescriptorProto{}
	if err := proto.Unmarshal(data, fdp); err != nil {
		return err
	}
	name := fdp.GetName()
	if name == "" {
		name = filepath.Base(path)
	}
	p.files[name] = fdp
	return nil
}

// Get returns a registered descriptor by name, or nil if not found.
func (p *ProtoFileLoader) Get(name string) *descriptorpb.FileDescriptorProto {
	return p.files[name]
}

// Names returns all registered descriptor names.
func (p *ProtoFileLoader) Names() []string {
	names := make([]string, 0, len(p.files))
	for k := range p.files {
		names = append(names, k)
	}
	return names
}

// Len returns the number of loaded descriptors.
func (p *ProtoFileLoader) Len() int {
	return len(p.files)
}

// Clear removes all loaded descriptors.
func (p *ProtoFileLoader) Clear() {
	p.files = make(map[string]*descriptorpb.FileDescriptorProto)
}
