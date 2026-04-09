package ui

import (
	"strings"

	"github.com/rivo/tview"
)

// MetadataPanel is a UI component for editing gRPC request metadata as key-value pairs.
type MetadataPanel struct {
	*tview.TextArea
}

// NewMetadataPanel creates and returns a new MetadataPanel.
func NewMetadataPanel() *MetadataPanel {
	ta := tview.NewTextArea()
	ta.SetTitle(" Metadata (key: value) ").SetBorder(true)
	ta.SetPlaceholder("authorization: Bearer <token>\nx-request-id: abc-123")
	return &MetadataPanel{TextArea: ta}
}

// GetPairs returns the current text content as a slice of non-empty trimmed lines.
func (m *MetadataPanel) GetPairs() []string {
	raw := m.GetText()
	lines := strings.Split(raw, "\n")
	var pairs []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			pairs = append(pairs, trimmed)
		}
	}
	return pairs
}

// SetPairs replaces the panel content with the provided key-value pairs.
func (m *MetadataPanel) SetPairs(pairs []string) {
	m.SetText(strings.Join(pairs, "\n"), true)
}

// Clear removes all text from the panel.
func (m *MetadataPanel) Clear() {
	m.SetText("", true)
}
