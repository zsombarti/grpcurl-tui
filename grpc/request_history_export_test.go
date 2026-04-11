package grpc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewHistoryExporter_NotNil(t *testing.T) {
	e := NewHistoryExporter(NewHistory(10))
	if e == nil {
		t.Fatal("expected non-nil HistoryExporter")
	}
}

func TestNewHistoryExporter_NilHistory(t *testing.T) {
	e := NewHistoryExporter(nil)
	if e == nil {
		t.Fatal("expected non-nil HistoryExporter even with nil history")
	}
}

func TestHistoryExporter_ExportJSON(t *testing.T) {
	h := NewHistory(10)
	h.Add(HistoryEntry{Address: "localhost:50051", Method: "/pkg.Svc/Hello",
		Request: `{}`, Response: `{"msg":"hi"}`, Timestamp: time.Now()})
	e := NewHistoryExporter(h)
	tmp := filepath.Join(t.TempDir(), "out.json")
	if err := e.ExportToFile(tmp, ExportFormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestHistoryExporter_ExportText(t *testing.T) {
	h := NewHistory(10)
	h.Add(HistoryEntry{Address: "localhost:50051", Method: "/pkg.Svc/Hello",
		Request: `{}`, Response: `{}`, Timestamp: time.Now()})
	e := NewHistoryExporter(h)
	tmp := filepath.Join(t.TempDir(), "out.txt")
	if err := e.ExportToFile(tmp, ExportFormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if len(data) == 0 {
		t.Fatal("expected non-empty text output")
	}
}

func TestHistoryExporter_ExportUnsupportedFormat(t *testing.T) {
	e := NewHistoryExporter(NewHistory(10))
	err := e.ExportToFile("/tmp/noop", ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
