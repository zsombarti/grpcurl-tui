package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// ResponseFormatterPanel lets the user choose a response display format
// (Pretty JSON / Compact JSON / Text) and renders formatted output.
type ResponseFormatterPanel struct {
	root    *tview.Flex
	dropdown *tview.DropDown
	output  *tview.TextView
	fmt     *grpcpkg.ResponseFormatter
}

// NewResponseFormatterPanel constructs the panel with a format selector
// and an output view.
func NewResponseFormatterPanel() *ResponseFormatterPanel {
	p := &ResponseFormatterPanel{
		fmt: grpcpkg.NewResponseFormatter(grpcpkg.FormatJSON, ""),
	}

	p.dropdown = tview.NewDropDown().
		SetLabel("Format: ").
		SetOptions([]string{"Pretty JSON", "Compact JSON", "Text"}, nil).
		SetCurrentOption(0)
	p.dropdown.SetBorder(false)

	p.output = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	p.output.SetBorder(true).SetTitle(" Response ")
	p.output.SetBorderColor(tcell.ColorDarkCyan)

	p.root = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(p.dropdown, 1, 0, false).
		AddItem(p.output, 0, 1, false)

	return p
}

// Primitive returns the root tview primitive for layout embedding.
func (p *ResponseFormatterPanel) Primitive() tview.Primitive { return p.root }

// SetData formats data with the currently selected style and updates
// the output view. Errors are displayed inline.
func (p *ResponseFormatterPanel) SetData(data map[string]any) {
	style := p.selectedStyle()
	f := grpcpkg.NewResponseFormatter(style, "")
	p.fmt = f

	out, err := f.Format(data)
	if err != nil {
		p.output.SetText("[red]Error: " + err.Error())
		return
	}
	if out == "" {
		p.output.SetText("[gray](empty response)")
		return
	}
	p.output.SetText(out)
	p.output.ScrollToBeginning()
}

// Clear wipes the output view.
func (p *ResponseFormatterPanel) Clear() {
	p.output.Clear()
}

func (p *ResponseFormatterPanel) selectedStyle() grpcpkg.FormatStyle {
	idx, _ := p.dropdown.GetCurrentOption()
	switch idx {
	case 1:
		return grpcpkg.FormatCompact
	case 2:
		return grpcpkg.FormatText
	default:
		return grpcpkg.FormatJSON
	}
}
