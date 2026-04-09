package ui

import (
	"testing"
)

func TestNewMetadataPanel_NotNil(t *testing.T) {
	p := NewMetadataPanel()
	if p == nil {
		t.Fatal("expected non-nil MetadataPanel")
	}
}

func TestMetadataPanel_GetPairs_Empty(t *testing.T) {
	p := NewMetadataPanel()
	pairs := p.GetPairs()
	if len(pairs) != 0 {
		t.Fatalf("expected 0 pairs, got %d", len(pairs))
	}
}

func TestMetadataPanel_SetAndGetPairs(t *testing.T) {
	p := NewMetadataPanel()
	input := []string{"authorization: Bearer tok", "x-id: 42"}
	p.SetPairs(input)
	got := p.GetPairs()
	if len(got) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(got))
	}
	if got[0] != "authorization: Bearer tok" {
		t.Errorf("unexpected first pair: %q", got[0])
	}
	if got[1] != "x-id: 42" {
		t.Errorf("unexpected second pair: %q", got[1])
	}
}

func TestMetadataPanel_GetPairs_SkipsBlankLines(t *testing.T) {
	p := NewMetadataPanel()
	p.SetText("key: val\n\n   \nanother: one", true)
	pairs := p.GetPairs()
	if len(pairs) != 2 {
		t.Fatalf("expected 2 non-empty pairs, got %d: %v", len(pairs), pairs)
	}
}

func TestMetadataPanel_Clear(t *testing.T) {
	p := NewMetadataPanel()
	p.SetPairs([]string{"k: v"})
	p.Clear()
	pairs := p.GetPairs()
	if len(pairs) != 0 {
		t.Fatalf("expected 0 pairs after clear, got %d", len(pairs))
	}
}
