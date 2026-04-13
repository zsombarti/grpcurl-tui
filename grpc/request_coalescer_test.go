package grpc

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestNewRequestCoalescer_NotNil(t *testing.T) {
	c := NewRequestCoalescer(DefaultCoalescerPolicy())
	if c == nil {
		t.Fatal("expected non-nil coalescer")
	}
}

func TestDefaultCoalescerPolicy_Values(t *testing.T) {
	p := DefaultCoalescerPolicy()
	if p.Window <= 0 {
		t.Errorf("expected positive window, got %v", p.Window)
	}
	if p.MaxKeys <= 0 {
		t.Errorf("expected positive MaxKeys, got %d", p.MaxKeys)
	}
}

func TestNewRequestCoalescer_InvalidPolicy_FallsBackToDefault(t *testing.T) {
	c := NewRequestCoalescer(CoalescerPolicy{Window: -1, MaxKeys: 0})
	def := DefaultCoalescerPolicy()
	if c.Policy().Window != def.Window {
		t.Errorf("expected default window %v, got %v", def.Window, c.Policy().Window)
	}
}

func TestRequestCoalescer_Len_Empty(t *testing.T) {
	c := NewRequestCoalescer(DefaultCoalescerPolicy())
	if c.Len() != 0 {
		t.Errorf("expected 0, got %d", c.Len())
	}
}

func TestRequestCoalescer_Do_EmptyKey_ReturnsError(t *testing.T) {
	c := NewRequestCoalescer(DefaultCoalescerPolicy())
	_, err := c.Do("", func() (map[string]any, error) { return nil, nil })
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestRequestCoalescer_Do_ReturnsResult(t *testing.T) {
	c := NewRequestCoalescer(DefaultCoalescerPolicy())
	want := map[string]any{"k": "v"}
	got, err := c.Do("key1", func() (map[string]any, error) { return want, nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["k"] != "v" {
		t.Errorf("expected v, got %v", got["k"])
	}
}

func TestRequestCoalescer_Do_CoalescesParallelCalls(t *testing.T) {
	policy := CoalescerPolicy{Window: 200 * time.Millisecond, MaxKeys: 64}
	c := NewRequestCoalescer(policy)

	callCount := 0
	var mu sync.Mutex

	fn := func() (map[string]any, error) {
		time.Sleep(20 * time.Millisecond)
		mu.Lock()
		callCount++
		mu.Unlock()
		return map[string]any{"x": 1}, nil
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Do("shared", fn) //nolint:errcheck
		}()
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if callCount != 1 {
		t.Errorf("expected fn called once, got %d", callCount)
	}
}

func TestRequestCoalescer_Do_PropagatesError(t *testing.T) {
	c := NewRequestCoalescer(DefaultCoalescerPolicy())
	_, err := c.Do("errkey", func() (map[string]any, error) {
		return nil, errors.New("boom")
	})
	if err == nil || err.Error() != "boom" {
		t.Errorf("expected boom error, got %v", err)
	}
}
