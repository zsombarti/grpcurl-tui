package grpc

import (
	"strings"
	"testing"
)

func TestNewRequestDiffer_NotNil(t *testing.T) {
	d := NewRequestDiffer()
	if d == nil {
		t.Fatal("expected non-nil RequestDiffer")
	}
}

func TestRequestDiffer_Diff_NoDifferences(t *testing.T) {
	d := NewRequestDiffer()
	res, err := d.Diff(`{"a":1}`, `{"a":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 0 || len(res.Removed) != 0 || len(res.Changed) != 0 {
		t.Errorf("expected empty diff, got %+v", res)
	}
	if res.String() != "no differences" {
		t.Errorf("expected 'no differences', got %q", res.String())
	}
}

func TestRequestDiffer_Diff_Added(t *testing.T) {
	d := NewRequestDiffer()
	res, err := d.Diff(`{"a":1}`, `{"a":1,"b":2}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Added["b"]; !ok {
		t.Errorf("expected 'b' in Added")
	}
}

func TestRequestDiffer_Diff_Removed(t *testing.T) {
	d := NewRequestDiffer()
	res, err := d.Diff(`{"a":1,"b":2}`, `{"a":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Removed["b"]; !ok {
		t.Errorf("expected 'b' in Removed")
	}
}

func TestRequestDiffer_Diff_Changed(t *testing.T) {
	d := NewRequestDiffer()
	res, err := d.Diff(`{"a":1}`, `{"a":2}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Changed["a"]; !ok {
		t.Errorf("expected 'a' in Changed")
	}
	summary := res.String()
	if !strings.Contains(summary, "~ a:") {
		t.Errorf("expected changed field in summary, got %q", summary)
	}
}

func TestRequestDiffer_Diff_InvalidLeft(t *testing.T) {
	d := NewRequestDiffer()
	_, err := d.Diff(`not json`, `{"a":1}`)
	if err == nil {
		t.Fatal("expected error for invalid left payload")
	}
}

func TestRequestDiffer_Diff_InvalidRight(t *testing.T) {
	d := NewRequestDiffer()
	_, err := d.Diff(`{"a":1}`, `not json`)
	if err == nil {
		t.Fatal("expected error for invalid right payload")
	}
}
