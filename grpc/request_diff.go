package grpc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

// DiffResult holds the differences between two request payloads.
type DiffResult struct {
	Added   map[string]interface{}
	Removed map[string]interface{}
	Changed map[string][2]interface{}
}

// String returns a human-readable summary of the diff.
func (d *DiffResult) String() string {
	if len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0 {
		return "no differences"
	}
	out := ""
	for k, v := range d.Added {
		out += fmt.Sprintf("+ %s: %v\n", k, v)
	}
	for k, v := range d.Removed {
		out += fmt.Sprintf("- %s: %v\n", k, v)
	}
	keys := make([]string, 0, len(d.Changed))
	for k := range d.Changed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := d.Changed[k]
		out += fmt.Sprintf("~ %s: %v -> %v\n", k, v[0], v[1])
	}
	return out
}

// RequestDiffer compares two JSON request payloads.
type RequestDiffer struct{}

// NewRequestDiffer creates a new RequestDiffer.
func NewRequestDiffer() *RequestDiffer {
	return &RequestDiffer{}
}

// Diff compares two JSON strings and returns a DiffResult.
func (d *RequestDiffer) Diff(left, right string) (*DiffResult, error) {
	var lMap, rMap map[string]interface{}
	if err := json.Unmarshal([]byte(left), &lMap); err != nil {
		return nil, fmt.Errorf("left payload: %w", err)
	}
	if err := json.Unmarshal([]byte(right), &rMap); err != nil {
		return nil, fmt.Errorf("right payload: %w", err)
	}
	res := &DiffResult{
		Added:   make(map[string]interface{}),
		Removed: make(map[string]interface{}),
		Changed: make(map[string][2]interface{}),
	}
	for k, lv := range lMap {
		rv, ok := rMap[k]
		if !ok {
			res.Removed[k] = lv
			continue
		}
		if !reflect.DeepEqual(lv, rv) {
			res.Changed[k] = [2]interface{}{lv, rv}
		}
	}
	for k, rv := range rMap {
		if _, ok := lMap[k]; !ok {
			res.Added[k] = rv
		}
	}
	return res, nil
}
