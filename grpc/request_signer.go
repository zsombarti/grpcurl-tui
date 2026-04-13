package grpc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// SignerPolicy holds configuration for request signing.
type SignerPolicy struct {
	Algorithm string
	HeaderName string
	IncludeTimestamp bool
}

// DefaultSignerPolicy returns a sensible default signing policy.
func DefaultSignerPolicy() SignerPolicy {
	return SignerPolicy{
		Algorithm: "hmac-sha256",
		HeaderName: "x-grpcurl-signature",
		IncludeTimestamp: true,
	}
}

// RequestSigner signs gRPC request payloads using HMAC-SHA256.
type RequestSigner struct {
	policy SignerPolicy
	secret []byte
}

// NewRequestSigner creates a new RequestSigner with the given secret and policy.
// Falls back to DefaultSignerPolicy if the algorithm is empty.
func NewRequestSigner(secret string, policy SignerPolicy) (*RequestSigner, error) {
	if secret == "" {
		return nil, errors.New("signer: secret must not be empty")
	}
	if policy.Algorithm == "" {
		policy = DefaultSignerPolicy()
	}
	if policy.HeaderName == "" {
		policy.HeaderName = DefaultSignerPolicy().HeaderName
	}
	return &RequestSigner{policy: policy, secret: []byte(secret)}, nil
}

// Sign computes an HMAC-SHA256 signature over the canonical form of the payload map.
// Returns the hex-encoded signature string.
func (s *RequestSigner) Sign(payload map[string]any) (string, error) {
	if payload == nil {
		return "", errors.New("signer: payload must not be nil")
	}
	canonical := s.canonicalize(payload)
	if s.policy.IncludeTimestamp {
		canonical = fmt.Sprintf("%s|ts=%d", canonical, time.Now().Unix())
	}
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(canonical))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// HeaderName returns the metadata header name to attach the signature to.
func (s *RequestSigner) HeaderName() string {
	return s.policy.HeaderName
}

// Policy returns the current signing policy.
func (s *RequestSigner) Policy() SignerPolicy {
	return s.policy
}

// canonicalize produces a deterministic string from a map by sorting keys.
func (s *RequestSigner) canonicalize(m map[string]any) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, m[k]))
	}
	return strings.Join(parts, "&")
}
