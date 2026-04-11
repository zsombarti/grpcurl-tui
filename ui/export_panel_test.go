package ui

import (
	"path/filepath"
	"testing"

	"grpcurl-tui/grpc"
)

func newTestExportPanel(t *testing.T) *ExportPanel {
	t.Helper()
	h := grpc.NewHistory(10)
	e := grpc.NewHistoryExporter(h)
	return NewExportPanel(e)
}

func TestNewExportPanel_NotNil(t *testing.T) {
	if p := newTestExportPanel(t); p == nil {
		t.Fatal("expected non-nil ExportPanel")
	}
}

func TestExportPanel_Primitive_NotNil(t *testing.T) {
	p := newTestExportPanel(t)
	if p.Primitive() == nil {
		t.Fatal("expected non-nil Primitive")
	}
}

func TestExportPanel_DefaultPath(t *testing.T) {
	p := newTestExportPanel(t)
	if got := p.GetPath(); got == "" {
		t.Fatal("expected a default path")
	}
}

func TestExportPanel_SetAndGetPath(t *testing.T) {
	p := newTestExportPanel(t)
	p.SetPath("/tmp/test_export.json")
	if got := p.GetPath(); got != "/tmp/test_export.json" {
		t.Fatalf("expected /tmp/test_export.json, got %s", got)
	}
}

func TestExportPanel_Export_WritesFile(t *testing.T) {
	h := grpc.NewHistory(10)
	h.Add(grpc.HistoryEntry{
		Address: "localhost:50051",
		Method:  "/svc.Test/Call",
		Request: `{}`,
		Response: `{}`,
	})
	e := grpc.NewHistoryExporter(h)
	p := NewExportPanel(e)
	tmp := filepath.Join(t.TempDir(), "out.json")
	p.SetPath(tmp)
	p.Export() // should not panic
}
