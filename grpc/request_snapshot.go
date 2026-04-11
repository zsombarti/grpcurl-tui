package grpc

import (
	"encoding/json"
	"fmt"
	"time"
)

// Snapshot captures a complete gRPC request configuration at a point in time.
type Snapshot struct {
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Method   string            `json:"method"`
	Payload  string            `json:"payload"`
	Metadata map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}

// SnapshotStore holds a bounded list of saved request snapshots.
type SnapshotStore struct {
	snapshots []Snapshot
	maxSize   int
}

// NewSnapshotStore creates a SnapshotStore with the given capacity.
func NewSnapshotStore(maxSize int) *SnapshotStore {
	if maxSize <= 0 {
		maxSize = 50
	}
	return &SnapshotStore{maxSize: maxSize}
}

// Save adds a snapshot, evicting the oldest if at capacity.
func (s *SnapshotStore) Save(snap Snapshot) {
	if snap.CreatedAt.IsZero() {
		snap.CreatedAt = time.Now()
	}
	if len(s.snapshots) >= s.maxSize {
		s.snapshots = s.snapshots[1:]
	}
	s.snapshots = append(s.snapshots, snap)
}

// All returns a copy of all stored snapshots.
func (s *SnapshotStore) All() []Snapshot {
	out := make([]Snapshot, len(s.snapshots))
	copy(out, s.snapshots)
	return out
}

// Len returns the number of stored snapshots.
func (s *SnapshotStore) Len() int { return len(s.snapshots) }

// Clear removes all snapshots.
func (s *SnapshotStore) Clear() { s.snapshots = s.snapshots[:0] }

// ExportJSON serialises all snapshots to a JSON byte slice.
func (s *SnapshotStore) ExportJSON() ([]byte, error) {
	data, err := json.MarshalIndent(s.snapshots, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("snapshot export: %w", err)
	}
	return data, nil
}
