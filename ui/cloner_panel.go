package ui

import (
	"encoding/json"
	"fmt"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// ClonerPanel provides a UI for deep-cloning request payloads.
type ClonerPanel struct {
	flex    *tview.Flex
	input   *tview.TextArea
	output  *tview.TextView
	cloner  *grpcpkg.RequestCloner
}

// NewClonerPanel constructs a ClonerPanel backed by the given cloner.
func NewClonerPanel(cloner *grpcpkg.RequestCloner) *ClonerPanel {
	if cloner == nil {
		cloner = grpcpkg.NewRequestCloner(grpcpkg.DefaultClonerPolicy())
	}

	input := tview.NewTextArea().
		SetPlaceholder(`{"key": "value"}`)
	input.SetTitle(" Source Payload ").SetBorder(true)

	output := tview.NewTextView().SetDynamicColors(true)
	output.SetTitle(" Cloned Payload ").SetBorder(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(input, 0, 1, true).
		AddItem(output, 0, 1, false)

	return &ClonerPanel{
		flex:   flex,
		input:  input,
		output: output,
		cloner: cloner,
	}
}

// Primitive returns the root tview primitive.
func (p *ClonerPanel) Primitive() tview.Primitive { return p.flex }

// Clone reads the input text area, clones the payload, and displays the result.
func (p *ClonerPanel) Clone() error {
	raw := p.input.GetText()
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		p.output.SetText(fmt.Sprintf("[red]parse error: %v[-]", err))
		return err
	}
	cloned, err := p.cloner.CloneJSON(payload)
	if err != nil {
		p.output.SetText(fmt.Sprintf("[red]clone error: %v[-]", err))
		return err
	}
	b, _ := json.MarshalIndent(cloned, "", "  ")
	p.output.SetText(string(b))
	return nil
}

// Clearsets both input and output.
func (p *ClonerPanel) Clear() {
	p.input.SetText("", false)
	p.output.SetText("")
}

// SetInput populates the input text area programmatically.
func (p *ClonerPanel) SetInput(raw string) {
	p.input.SetText(raw, false)
}

// GetOutput returns the current output text.
func (p *ClonerPanel) GetOutput() string {
	return p.output.GetText(true)
}
