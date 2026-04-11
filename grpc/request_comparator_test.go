package grpc

import (
	"testing"
)

func TestNewRequestComparator_NotNil(t *testing.T) {
	c := NewRequestComparator()
	if c == nil {
		t.Fatal("expected non-nil comparator")
	}
}

func TestRequestComparator_Compare_IdenticalPayloads(t *testing.T) {
	c := NewRequestComparator()
	result, err := c.Compare(`{"name":"alice","age":30}`, `{"name":"alice","age":30}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Match {
		t.Errorf("expected match, got differences: %v", result.Differences)
	}
	if result.Similarity != 1.0 {
		t.Errorf("expected similarity 1.0, got %f", result.Similarity)
	}
}

func TestRequestComparator_Compare_DifferentValues(t *testing.T) {
	c := NewRequestComparator()
	result, err := c.Compare(`{"name":"alice"}`, `{"name":"bob"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Match {
		t.Error("expected no match")
	}
	if len(result.Differences) != 1 {
		t.Errorf("expected 1 difference, got %d", len(result.Differences))
	}
}

func TestRequestComparator_Compare_AddedKey(t *testing.T) {
	c := NewRequestComparator()
	result, err := c.Compare(`{"name":"alice"}`, `{"name":"alice","age":30}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Match {
		t.Error("expected no match due to added key")
	}
	if len(result.Differences) == 0 {
		t.Error("expected at least one difference")
	}
}

func TestRequestComparator_Compare_RemovedKey(t *testing.T) {
	c := NewRequestComparator()
	result, err := c.Compare(`{"name":"alice","age":30}`, `{"name":"alice"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Match {
		t.Error("expected no match due to removed key")
	}
}

func TestRequestComparator_Compare_EmptyObjects(t *testing.T) {
	c := NewRequestComparator()
	result, err := c.Compare(`{}`, `{}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Match {
		t.Error("expected empty objects to match")
	}
	if result.Similarity != 1.0 {
		t.Errorf("expected similarity 1.0, got %f", result.Similarity)
	}
}

func TestRequestComparator_Compare_InvalidJSON(t *testing.T) {
	c := NewRequestComparator()
	_, err := c.Compare(`not-json`, `{}`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestCompareResult_SummaryLine_Match(t *testing.T) {
	r := &CompareResult{Match: true, Similarity: 1.0}
	if r.SummaryLine() != "requests are identical" {
		t.Errorf("unexpected summary: %s", r.SummaryLine())
	}
}

func TestCompareResult_SummaryLine_Mismatch(t *testing.T) {
	r := &CompareResult{Match: false, Similarity: 0.5, Differences: []string{"changed: name: alice -> bob"}}
	s := r.SummaryLine()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
