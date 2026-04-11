package grpc

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewSnapshotStore_NotNil(t *testing.T) {
	store := NewSnapshotStore(10)
	if store == nil {
		t.Fatal("expected non-nil SnapshotStore")
	}
}

func TestSnapshotStore_DefaultMaxSize(t *testing.T) {
	store := NewSnapshotStore(0)
	if store.maxSize != 50 {
		t.Fatalf("expected default maxSize 50, got %d", store.maxSize)
	}
}

func TestSnapshotStore_Save_And_Len(t *testing.T) {
	store := NewSnapshotStore(5)
	store.Save(Snapshot{Name: "s1", Address: "localhost:50051", Method: "Foo"})
	store.Save(Snapshot{Name: "s2", Address: "localhost:50051", Method: "Bar"})
	if store.Len() != 2 {
		t.Fatalf("expected 2, got %d", store.Len())
	}
}

func TestSnapshotStore_Eviction(t *testing.T) {
	store := NewSnapshotStore(2)
	store.Save(Snapshot{Name: "a"})
	store.Save(Snapshot{Name: "b"})
	store.Save(Snapshot{Name: "c"})
	if store.Len() != 2 {
		t.Fatalf("expected 2 after eviction, got %d", store.Len())
	}
	if store.All()[0].Name != "b" {
		t.Fatalf("expected oldest evicted, got %s", store.All()[0].Name)
	}
}

func TestSnapshotStore_TimestampAutoSet(t *testing.T) {
	store := NewSnapshotStore(5)
	before := time.Now()
	store.Save(Snapshot{Name: "ts"})
	after := time.Now()
	snap := store.All()[0]
	if snap.CreatedAt.Before(before) || snap.CreatedAt.After(after) {
		t.Fatal("timestamp not set correctly")
	}
}

func TestSnapshotStore_Clear(t *testing.T) {
	store := NewSnapshotStore(5)
	store.Save(Snapshot{Name: "x"})
	store.Clear()
	if store.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", store.Len())
	}
}

func TestSnapshotStore_ExportJSON(t *testing.T) {
	store := NewSnapshotStore(5)
	store.Save(Snapshot{Name: "export", Address: "localhost:9090", Method: "SayHello"})
	data, err := store.ExportJSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var snaps []Snapshot
	if err := json.Unmarshal(data, &snaps); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(snaps) != 1 || snaps[0].Name != "export" {
		t.Fatal("exported snapshot mismatch")
	}
}
