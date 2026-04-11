package grpc

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the result of comparing two JSON request payloads.
type CompareResult struct {
	Match      bool
	Similarity float64 // 0.0 to 1.0
	Differences []string
}

// RequestComparator compares two JSON request payloads and scores their similarity.
type RequestComparator struct{}

// NewRequestComparator returns a new RequestComparator.
func NewRequestComparator() *RequestComparator {
	return &RequestComparator{}
}

// Compare compares two JSON strings and returns a CompareResult.
func (c *RequestComparator) Compare(a, b string) (*CompareResult, error) {
	var mapA, mapB map[string]interface{}
	if err := json.Unmarshal([]byte(a), &mapA); err != nil {
		return nil, fmt.Errorf("invalid JSON (a): %w", err)
	}
	if err := json.Unmarshal([]byte(b), &mapB); err != nil {
		return nil, fmt.Errorf("invalid JSON (b): %w", err)
	}

	diffs := c.diffMaps("", mapA, mapB)
	total := c.countKeys(mapA) + c.countKeys(mapB)
	var similarity float64
	if total == 0 {
		similarity = 1.0
	} else {
		similarity = 1.0 - float64(len(diffs))/float64(total)
		if similarity < 0 {
			similarity = 0
		}
	}

	return &CompareResult{
		Match:       len(diffs) == 0,
		Similarity:  similarity,
		Differences: diffs,
	}, nil
}

func (c *RequestComparator) diffMaps(prefix string, a, b map[string]interface{}) []string {
	var diffs []string
	keys := mergeKeys(a, b)
	sort.Strings(keys)
	for _, k := range keys {
		path := k
		if prefix != "" {
			path = prefix + "." + k
		}
		va, aOk := a[k]
		vb, bOk := b[k]
		if !aOk {
			diffs = append(diffs, fmt.Sprintf("added: %s = %v", path, vb))
		} else if !bOk {
			diffs = append(diffs, fmt.Sprintf("removed: %s = %v", path, va))
		} else {
			subA, aIsMap := va.(map[string]interface{})
			subB, bIsMap := vb.(map[string]interface{})
			if aIsMap && bIsMap {
				diffs = append(diffs, c.diffMaps(path, subA, subB)...)
			} else if fmt.Sprintf("%v", va) != fmt.Sprintf("%v", vb) {
				diffs = append(diffs, fmt.Sprintf("changed: %s: %v -> %v", path, va, vb))
			}
		}
	}
	return diffs
}

func (c *RequestComparator) countKeys(m map[string]interface{}) int {
	count := 0
	for _, v := range m {
		count++
		if sub, ok := v.(map[string]interface{}); ok {
			count += c.countKeys(sub)
		}
	}
	return count
}

func mergeKeys(a, b map[string]interface{}) []string {
	seen := map[string]bool{}
	for k := range a {
		seen[k] = true
	}
	for k := range b {
		seen[k] = true
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}

// SummaryLine returns a one-line human-readable summary.
func (r *CompareResult) SummaryLine() string {
	if r.Match {
		return "requests are identical"
	}
	return fmt.Sprintf("%.0f%% similar, %d difference(s): %s",
		r.Similarity*100, len(r.Differences), strings.Join(r.Differences[:min(3, len(r.Differences))], "; "))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
