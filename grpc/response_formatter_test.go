package grpc

import (
	"strings"
	"testing"
)

func TestNewResponseFormatter_NotNil(t *testing.T) {
	f := NewResponseFormatter(FormatJSON, "")
	if f == nil {
		t.Fatal("expected non-nil ResponseFormatter")
	}
}

func TestResponseFormatter_DefaultIndent(t *testing.T) {
	f := NewResponseFormatter(FormatJSON, "")
	if f.indent != "  " {
		t.Fatalf("expected default indent '  ', got %q", f.indent)
	}
}

func TestResponseFormatter_Style_RoundTrip(t *testing.T) {
	for _, style := range []FormatStyle{FormatJSON, FormatCompact, FormatText} {
		f := NewResponseFormatter(style, "")
		if f.Style() != style {
			t.Fatalf("expected style %d, got %d", style, f.Style())
		}
	}
}

func TestResponseFormatter_Format_NilMap(t *testing.T) {
	f := NewResponseFormatter(FormatJSON, "")
	out, err := f.Format(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Fatalf("expected empty string for nil map, got %q", out)
	}
}

func TestResponseFormatter_Format_PrettyJSON(t *testing.T) {
	f := NewResponseFormatter(FormatJSON, "  ")
	data := map[string]any{"name": "alice", "age": 30}
	out, err := f.Format(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\n") {
		t.Fatalf("expected pretty JSON with newlines, got %q", out)
	}
	if !strings.Contains(out, "\"name\"") {
		t.Fatalf("expected 'name' key in output, got %q", out)
	}
}

func TestResponseFormatter_Format_Compact(t *testing.T) {
	f := NewResponseFormatter(FormatCompact, "")
	data := map[string]any{"ok": true}
	out, err := f.Format(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "\n") {
		t.Fatalf("expected compact JSON without newlines, got %q", out)
	}
	if !strings.Contains(out, "ok") {
		t.Fatalf("expected 'ok' key in compact output, got %q", out)
	}
}

func TestResponseFormatter_Format_Text(t *testing.T) {
	f := NewResponseFormatter(FormatText, "")
	data := map[string]any{"status": "UP", "latency": 42}
	out, err := f.Format(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "=") {
		t.Fatalf("expected key=value format, got %q", out)
	}
}

func TestResponseFormatter_Format_Text_Nested(t *testing.T) {
	f := NewResponseFormatter(FormatText, "")
	data := map[string]any{
		"server": map[string]any{"host": "localhost", "port": 50051},
	}
	out, err := f.Format(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "server:") {
		t.Fatalf("expected nested key header, got %q", out)
	}
	if !strings.Contains(out, "host") {
		t.Fatalf("expected nested 'host' key, got %q", out)
	}
}
