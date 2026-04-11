package ui

import (
	"github.com/rivo/tview"

	"github.com/user/grpcurl-tui/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// PayloadTemplatePanel is a UI panel that displays a generated JSON payload template
// for a selected proto message descriptor and allows the user to edit it.
type PayloadTemplatePanel struct {
	root      *tview.Flex
	textArea  *tview.TextArea
	generator *grpc.PayloadTemplateGenerator
}

// NewPayloadTemplatePanel creates a new PayloadTemplatePanel.
func NewPayloadTemplatePanel() *PayloadTemplatePanel {
	textArea := tview.NewTextArea().
		SetPlaceholder("Select a method to generate a payload template...")
	textArea.SetBorder(true).SetTitle(" Payload (JSON) ")

	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textArea, 0, 1, true)

	return &PayloadTemplatePanel{
		root:      root,
		textArea:  textArea,
		generator: grpc.NewPayloadTemplateGenerator(),
	}
}

// Primitive returns the root tview primitive for layout embedding.
func (p *PayloadTemplatePanel) Primitive() tview.Primitive {
	return p.root
}

// LoadDescriptor generates a template from the given message descriptor and
// populates the text area. Existing content is replaced.
func (p *PayloadTemplatePanel) LoadDescriptor(md protoreflect.MessageDescriptor) error {
	template, err := p.generator.Generate(md)
	if err != nil {
		return err
	}
	p.textArea.SetText(template, true)
	return nil
}

// GetPayload returns the current text content of the payload editor.
func (p *PayloadTemplatePanel) GetPayload() string {
	return p.textArea.GetText()
}

// Clear resets the text area to empty.
func (p *PayloadTemplatePanel) Clear() {
	p.textArea.SetText("", true)
}
