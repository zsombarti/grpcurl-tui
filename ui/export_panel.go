package ui

import (
	"path/filepath"

	"github.com/rivo/tview"

	"grpcurl-tui/grpc"
)

// ExportPanel provides a UI for exporting request history to a file.
type ExportPanel struct {
	flex     *tview.Flex
	pathInput *tview.InputField
	formatDrop *tview.DropDown
	status   *tview.TextView
	exporter *grpc.HistoryExporter
}

// NewExportPanel constructs an ExportPanel backed by the given HistoryExporter.
func NewExportPanel(exporter *grpc.HistoryExporter) *ExportPanel {
	p := &ExportPanel{exporter: exporter}

	p.pathInput = tview.NewInputField().
		SetLabel("Output path: ").
		SetText("history_export.json").
		SetFieldWidth(40)

	p.formatDrop = tview.NewDropDown().
		SetLabel("Format: ").
		SetOptions([]string{"json", "text"}, nil).
		SetCurrentOption(0)

	p.status = tview.NewTextView().SetDynamicColors(true)

	p.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(p.pathInput, 1, 0, true).
		AddItem(p.formatDrop, 1, 0, false).
		AddItem(p.status, 0, 1, false)

	return p
}

// Primitive returns the root tview primitive for layout embedding.
func (p *ExportPanel) Primitive() tview.Primitive { return p.flex }

// Export triggers the export using current panel values.
func (p *ExportPanel) Export() {
	path := filepath.Clean(p.pathInput.GetText())
	_, fmtLabel := p.formatDrop.GetCurrentOption()
	fmt := grpc.ExportFormat(fmtLabel)

	if err := p.exporter.ExportToFile(path, fmt); err != nil {
		p.status.SetText("[red]Error: " + err.Error())
		return
	}
	p.status.SetText("[green]Exported to " + path)
}

// SetPath sets the output file path in the input field.
func (p *ExportPanel) SetPath(path string) { p.pathInput.SetText(path) }

// GetPath returns the current output file path.
func (p *ExportPanel) GetPath() string { return p.pathInput.GetText() }
