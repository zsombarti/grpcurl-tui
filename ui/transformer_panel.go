package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// TransformerPanel displays registered transform steps and allows adding new ones.
type TransformerPanel struct {
	flex  *tview.Flex
	list  *tview.List
	input *tview.InputField
	steps []string
}

// NewTransformerPanel creates and wires up the transformer panel UI.
func NewTransformerPanel() *TransformerPanel {
	p := &TransformerPanel{
		list:  tview.NewList().ShowSecondaryText(false),
		input: tview.NewInputField(),
	}

	p.list.SetBorder(true).SetTitle(" Transform Steps ")

	p.input.SetLabel("Step name: ").
		SetFieldWidth(30).
		SetPlaceholder("e.g. redact-secrets")

	p.flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.list, 0, 1, false).
		AddItem(p.input, 3, 0, true)

	return p
}

// Primitive returns the root drawable for embedding in layouts.
func (p *TransformerPanel) Primitive() tview.Primitive {
	return p.flex
}

// AddStep registers a display entry for a named step.
func (p *TransformerPanel) AddStep(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("step name must not be empty")
	}
	p.steps = append(p.steps, name)
	p.list.AddItem(name, "", 0, nil)
	return nil
}

// Refresh rebuilds the list from the current steps slice.
func (p *TransformerPanel) Refresh() {
	p.list.Clear()
	for _, s := range p.steps {
		p.list.AddItem(s, "", 0, nil)
	}
}

// Clear removes all step entries from the panel.
func (p *TransformerPanel) Clear() {
	p.steps = nil
	p.list.Clear()
}

// Len returns the number of steps currently shown.
func (p *TransformerPanel) Len() int {
	return len(p.steps)
}
