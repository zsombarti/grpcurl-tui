package ui

import (
	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// SignerPanel provides a UI form for configuring request signing.
type SignerPanel struct {
	form       *tview.Form
	secretInput *tview.InputField
	headerInput *tview.InputField
	timestamp  *tview.Checkbox
}

// NewSignerPanel constructs a SignerPanel with default values pre-filled.
func NewSignerPanel() *SignerPanel {
	p := &SignerPanel{}
	defaults := grpcpkg.DefaultSignerPolicy()

	p.secretInput = tview.NewInputField().
		SetLabel("Secret Key").
		SetFieldWidth(32).
		SetMaskCharacter('*')

	p.headerInput = tview.NewInputField().
		SetLabel("Header Name").
		SetText(defaults.HeaderName).
		SetFieldWidth(28)

	p.timestamp = tview.NewCheckbox().
		SetLabel("Include Timestamp").
		SetChecked(defaults.IncludeTimestamp)

	p.form = tview.NewForm().
		AddFormItem(p.secretInput).
		AddFormItem(p.headerInput).
		AddFormItem(p.timestamp)

	p.form.SetBorder(true).SetTitle(" Request Signer ")
	return p
}

// Primitive returns the underlying tview primitive for layout embedding.
func (p *SignerPanel) Primitive() tview.Primitive {
	return p.form
}

// GetSigner builds a RequestSigner from the current panel inputs.
// Returns nil and an error if the secret is empty or construction fails.
func (p *SignerPanel) GetSigner() (*grpcpkg.RequestSigner, error) {
	policy := grpcpkg.SignerPolicy{
		Algorithm: "hmac-sha256",
		HeaderName: p.headerInput.GetText(),
		IncludeTimestamp: p.timestamp.IsChecked(),
	}
	return grpcpkg.NewRequestSigner(p.secretInput.GetText(), policy)
}

// Clear resets all input fields to their default values.
func (p *SignerPanel) Clear() {
	defaults := grpcpkg.DefaultSignerPolicy()
	p.secretInput.SetText("")
	p.headerInput.SetText(defaults.HeaderName)
	p.timestamp.SetChecked(defaults.IncludeTimestamp)
}
