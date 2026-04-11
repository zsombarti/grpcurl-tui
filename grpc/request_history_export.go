package grpc

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ExportFormat defines the output format for exported history.
type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatText ExportFormat = "text"
)

// HistoryExporter writes history entries to a file or writer.
type HistoryExporter struct {
	history *History
}

// NewHistoryExporter creates a HistoryExporter backed by the given History.
func NewHistoryExporter(h *History) *HistoryExporter {
	if h == nil {
		h = NewHistory(0)
	}
	return &HistoryExporter{history: h}
}

// ExportToFile writes history entries to path using the specified format.
func (e *HistoryExporter) ExportToFile(path string, format ExportFormat) error {
	entries := e.history.Entries()
	switch format {
	case ExportFormatJSON:
		return e.writeJSON(path, entries)
	case ExportFormatText:
		return e.writeText(path, entries)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func (e *HistoryExporter) writeJSON(path string, entries []HistoryEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

func (e *HistoryExporter) writeText(path string, entries []HistoryEntry) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	for _, en := range entries {
		ts := en.Timestamp.Format(time.RFC3339)
		_, err = fmt.Fprintf(f, "[%s] %s %s\n  request:  %s\n  response: %s\n\n",
			ts, en.Address, en.Method, en.Request, en.Response)
		if err != nil {
			return err
		}
	}
	return nil
}
