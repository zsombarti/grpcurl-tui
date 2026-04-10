package ui

import (
	"strconv"

	"github.com/rivo/tview"

	grpcpkg "grpcurl-tui/grpc"
)

// TLSPanel is a UI panel for configuring TLS options.
type TLSPanel struct {
	form    *tview.Form
	enabled *tview.Checkbox
	insecure *tview.Checkbox
	caCert  *tview.InputField
	clientCert *tview.InputField
	clientKey  *tview.InputField
}

// NewTLSPanel creates and returns a new TLSPanel.
func NewTLSPanel() *TLSPanel {
	p := &TLSPanel{
		form:       tview.NewForm(),
		enabled:    tview.NewCheckbox(),
		insecure:   tview.NewCheckbox(),
		caCert:     tview.NewInputField(),
		clientCert: tview.NewInputField(),
		clientKey:  tview.NewInputField(),
	}

	p.enabled.SetLabel("Enable TLS")
	p.insecure.SetLabel("Skip Verify")
	p.caCert.SetLabel("CA Cert Path").SetFieldWidth(40)
	p.clientCert.SetLabel("Client Cert Path").SetFieldWidth(40)
	p.clientKey.SetLabel("Client Key Path").SetFieldWidth(40)

	p.form.SetBorder(true).SetTitle(" TLS ").SetTitleAlign(tview.AlignLeft)
	p.form.AddFormItem(p.enabled)
	p.form.AddFormItem(p.insecure)
	p.form.AddFormItem(p.caCert)
	p.form.AddFormItem(p.clientCert)
	p.form.AddFormItem(p.clientKey)

	return p
}

// GetConfig reads the form fields and returns a TLSConfig.
func (p *TLSPanel) GetConfig() grpcpkg.TLSConfig {
	return grpcpkg.TLSConfig{
		Enabled:    p.enabled.IsChecked(),
		Insecure:   p.insecure.IsChecked(),
		CACert:     p.caCert.GetText(),
		ClientCert: p.clientCert.GetText(),
		ClientKey:  p.clientKey.GetText(),
	}
}

// SetConfig populates the form fields from a TLSConfig.
func (p *TLSPanel) SetConfig(cfg grpcpkg.TLSConfig) {
	p.enabled.SetChecked(cfg.Enabled)
	p.insecure.SetChecked(cfg.Insecure)
	p.caCert.SetText(cfg.CACert)
	p.clientCert.SetText(cfg.ClientCert)
	p.clientKey.SetText(cfg.ClientKey)
	_ = strconv.FormatBool(cfg.Enabled) // satisfy import if needed
}

// Primitive returns the underlying tview.Form for embedding in layouts.
func (p *TLSPanel) Primitive() tview.Primitive {
	return p.form
}
