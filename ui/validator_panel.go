package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"google.golang.org/protobuf/reflect/protoreflect"

	grpcpkg "grpcurl-tui/grpc"
)

// ValidatorPanel displays request validation results in the TUI.
type ValidatorPanel struct {
	frame     *tview.Frame
	textView  *tview.TextView
	validator *grpcpkg.RequestValidator
}

// NewValidatorPanel constructs a ValidatorPanel.
func NewValidatorPanel() *ValidatorPanel {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetScrollable(true)
	tv.SetWrap(true)

	frame := tview.NewFrame(tv).
		SetBorders(0, 0, 0, 0, 0, 0)
	frame.SetBorder(true).SetTitle(" Validation ")

	return &ValidatorPanel{
		frame:     frame,
		textView:  tv,
		validator: grpcpkg.NewRequestValidator(),
	}
}

// Primitive returns the tview primitive for layout embedding.
func (p *ValidatorPanel) Primitive() tview.Primitive {
	return p.frame
}

// Validate runs validation on msg and renders the results.
func (p *ValidatorPanel) Validate(msg protoreflect.Message) bool {
	p.textView.Clear()

	if msg == nil {
		fmt.Fprintf(p.textView, "[red]✗ no message to validate[-]\n")
		return false
	}

	errs := p.validator.Validate(msg)
	if len(errs) == 0 {
		fmt.Fprintf(p.textView, "[green]✓ message is valid[-]\n")
		return true
	}

	lines := make([]string, 0, len(errs))
	for _, e := range errs {
		lines = append(lines, fmt.Sprintf("[red]✗[-] %s", e.Error()))
	}
	fmt.Fprint(p.textView, strings.Join(lines, "\n"))
	return false
}

// Clear resets the panel content.
func (p *ValidatorPanel) Clear() {
	p.textView.Clear()
}
