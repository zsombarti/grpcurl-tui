package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// TLSConfig holds TLS configuration for a gRPC connection.
type TLSConfig struct {
	Enabled    bool
	Insecure   bool   // skip server certificate verification
	CACert     string // path to CA certificate file
	ClientCert string // path to client certificate file
	ClientKey  string // path to client key file
}

// DefaultTLSConfig returns a TLSConfig with TLS disabled.
func DefaultTLSConfig() TLSConfig {
	return TLSConfig{Enabled: false}
}

// NewTLSBuilder creates a TLSBuilder from the given config.
func NewTLSBuilder(cfg TLSConfig) *TLSBuilder {
	return &TLSBuilder{cfg: cfg}
}

// TLSBuilder constructs a *tls.Config from a TLSConfig.
type TLSBuilder struct {
	cfg TLSConfig
}

// Build returns a *tls.Config or nil if TLS is disabled.
func (b *TLSBuilder) Build() (*tls.Config, error) {
	if !b.cfg.Enabled {
		return nil, nil
	}

	tlsCfg := &tls.Config{
		InsecureSkipVerify: b.cfg.Insecure, //nolint:gosec
	}

	if b.cfg.CACert != "" {
		pem, err := os.ReadFile(b.cfg.CACert)
		if err != nil {
			return nil, fmt.Errorf("reading CA cert: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("failed to append CA cert from %s", b.cfg.CACert)
		}
		tlsCfg.RootCAs = pool
	}

	if b.cfg.ClientCert != "" && b.cfg.ClientKey != "" {
		cert, err := tls.LoadX509KeyPair(b.cfg.ClientCert, b.cfg.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("loading client key pair: %w", err)
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
	}

	return tlsCfg, nil
}
