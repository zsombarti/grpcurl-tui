package grpc

import (
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestNewProtoFileLoader_NotNil(t *testing.T) {
	l := NewProtoFileLoader()
	if l == nil {
		t.Fatal("expected non-nil ProtoFileLoader")
	}
}

func TestProtoFileLoader_Len_Empty(t *testing.T) {
	l := NewProtoFileLoader()
	if l.Len() != 0 {
		t.Fatalf("expected 0, got %d", l.Len())
	}
}

func TestProtoFileLoader_LoadFile_EmptyPath(t *testing.T) {
	l := NewProtoFileLoader()
	if err := l.LoadFile(""); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestProtoFileLoader_LoadFile_UnsupportedExtension(t *testing.T) {
	l := NewProtoFileLoader()
	if err := l.LoadFile("file.proto"); err == nil {
		t.Fatal("expected error for unsupported extension")
	}
}

func TestProtoFileLoader_LoadFile_MissingFile(t *testing.T) {
	l := NewProtoFileLoader()
	if err := l.LoadFile("/nonexistent/path/file.pb"); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestProtoFileLoader_LoadFile_ValidDescriptor(t *testing.T) {
	fdp := &descriptorpb.FileDescriptorProto{}
	fdp.Name = proto.String("test.proto")
	data, err := proto.Marshal(fdp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	tmp := filepath.Join(t.TempDir(), "test.pb")
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		t.Fatalf("write error: %v", err)
	}
	l := NewProtoFileLoader()
	if err := l.LoadFile(tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Len() != 1 {
		t.Fatalf("expected 1, got %d", l.Len())
	}
	if l.Get("test.proto") == nil {
		t.Fatal("expected descriptor to be retrievable by name")
	}
}

func TestProtoFileLoader_Clear(t *testing.T) {
	fdp := &descriptorpb.FileDescriptorProto{}
	fdp.Name = proto.String("a.proto")
	data, _ := proto.Marshal(fdp)
	tmp := filepath.Join(t.TempDir(), "a.pb")
	_ = os.WriteFile(tmp, data, 0o600)
	l := NewProtoFileLoader()
	_ = l.LoadFile(tmp)
	l.Clear()
	if l.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", l.Len())
	}
}
