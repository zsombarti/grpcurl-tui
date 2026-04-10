package ui

import (
	"testing"

	grpcpkg "grpcurl-tui/grpc"
)

func TestNewTLSPanel_NotNil(t *testing.T) {
	p := NewTLSPanel()
	if p == nil {
		t.Fatal("expected non-nil TLSPanel")
	}
}

func TestTLSPanel_Primitive_NotNil(t *testing.T) {
	p := NewTLSPanel()
	if p.Primitive() == nil {
		t.Fatal("expected non-nil primitive")
	}
}

func TestTLSPanel_GetConfig_Defaults(t *testing.T) {
	p := NewTLSPanel()
	cfg := p.GetConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.Insecure {
		t.Error("expected Insecure to be false by default")
	}
	if cfg.CACert != "" || cfg.ClientCert != "" || cfg.ClientKey != "" {
		t.Error("expected empty cert paths by default")
	}
}

func TestTLSPanel_SetAndGetConfig(t *testing.T) {
	p := NewTLSPanel()
	input := grpcpkg.TLSConfig{
		Enabled:    true,
		Insecure:   false,
		CACert:     "/etc/certs/ca.pem",
		ClientCert: "/etc/certs/client.crt",
		ClientKey:  "/etc/certs/client.key",
	}
	p.SetConfig(input)
	out := p.GetConfig()

	if out.Enabled != input.Enabled {
		t.Errorf("Enabled: got %v, want %v", out.Enabled, input.Enabled)
	}
	if out.CACert != input.CACert {
		t.Errorf("CACert: got %q, want %q", out.CACert, input.CACert)
	}
	if out.ClientCert != input.ClientCert {
		t.Errorf("ClientCert: got %q, want %q", out.ClientCert, input.ClientCert)
	}
	if out.ClientKey != input.ClientKey {
		t.Errorf("ClientKey: got %q, want %q", out.ClientKey, input.ClientKey)
	}
}

func TestTLSPanel_SetConfig_InsecureFlag(t *testing.T) {
	p := NewTLSPanel()
	p.SetConfig(grpcpkg.TLSConfig{Enabled: true, Insecure: true})
	cfg := p.GetConfig()
	if !cfg.Insecure {
		t.Error("expected Insecure to be true after SetConfig")
	}
}
