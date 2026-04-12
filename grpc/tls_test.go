package grpc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultTLSConfig_Disabled(t *testing.T) {
	cfg := DefaultTLSConfig()
	if cfg.Enabled {
		t.Fatal("expected TLS to be disabled by default")
	}
}

func TestNewTLSBuilder_NotNil(t *testing.T) {
	b := NewTLSBuilder(DefaultTLSConfig())
	if b == nil {
		t.Fatal("expected non-nil TLSBuilder")
	}
}

func TestTLSBuilder_Build_Disabled(t *testing.T) {
	b := NewTLSBuilder(DefaultTLSConfig())
	tlsCfg, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tlsCfg != nil {
		t.Fatal("expected nil tls.Config when TLS disabled")
	}
}

func TestTLSBuilder_Build_Insecure(t *testing.T) {
	cfg := TLSConfig{Enabled: true, Insecure: true}
	b := NewTLSBuilder(cfg)
	tlsCfg, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tlsCfg == nil {
		t.Fatal("expected non-nil tls.Config")
	}
	if !tlsCfg.InsecureSkipVerify {
		t.Fatal("expected InsecureSkipVerify to be true")
	}
}

func TestTLSBuilder_Build_MissingCACert(t *testing.T) {
	cfg := TLSConfig{Enabled: true, CACert: "/nonexistent/ca.pem"}
	b := NewTLSBuilder(cfg)
	_, err := b.Build()
	if err == nil {
		t.Fatal("expected error for missing CA cert file")
	}
}

func TestTLSBuilder_Build_InvalidCACert(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad-ca.pem")
	_ = os.WriteFile(tmp, []byte("not-a-cert"), 0o600)
	cfg := TLSConfig{Enabled: true, CACert: tmp}
	b := NewTLSBuilder(cfg)
	_, err := b.Build()
	if err == nil {
		t.Fatal("expected error for invalid CA cert PEM")
	}
}

func TestTLSBuilder_Build_MissingClientKeyPair(t *testing.T) {
	cfg := TLSConfig{Enabled: true, Insecure: true, ClientCert: "/no/cert", ClientKey: "/no/key"}
	b := NewTLSBuilder(cfg)
	_, err := b.Build()
	if err == nil {
		t.Fatal("expected error for missing client cert/key")
	}
}

func TestTLSBuilder_Build_EnabledNoOptions(t *testing.T) {
	// TLS enabled with no CA cert, no insecure flag, and no client cert should
	// succeed and return a tls.Config that uses the system certificate pool.
	cfg := TLSConfig{Enabled: true}
	b := NewTLSBuilder(cfg)
	tlsCfg, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tlsCfg == nil {
		t.Fatal("expected non-nil tls.Config")
	}
	if tlsCfg.InsecureSkipVerify {
		t.Fatal("expected InsecureSkipVerify to be false")
	}
}
