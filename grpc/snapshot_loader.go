package grpc

import (
	"encoding/json"
	"fmt"
	"os"
)

// SnapshotLoader handles persisting and restoring a SnapshotStore to/from disk.
type SnapshotLoader struct {
	store *SnapshotStore
}

// NewSnapshotLoader wraps a SnapshotStore with file I/O capabilities.
func NewSnapshotLoader(store *SnapshotStore) *SnapshotLoader {
	if store == nil {
		store = NewSnapshotStore(50)
	}
	return &SnapshotLoader{store: store}
}

// SaveToFile writes all snapshots as JSON to the given path.
func (l *SnapshotLoader) SaveToFile(path string) error {
	if path == "" {
		return fmt.Errorf("snapshot save: path must not be empty")
	}
	data, err := l.store.ExportJSON()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot save: %w", err)
	}
	return nil
}

// LoadFromFile reads snapshots from a JSON file and appends them to the store.
func (l *SnapshotLoader) LoadFromFile(path string) error {
	if path == "" {
		return fmt.Errorf("snapshot load: path must not be empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("snapshot load: %w", err)
	}
	var snaps []Snapshot
	if err := json.Unmarshal(data, &snaps); err != nil {
		return fmt.Errorf("snapshot load: invalid JSON: %w", err)
	}
	for _, s := range snaps {
		l.store.Save(s)
	}
	return nil
}

// Store returns the underlying SnapshotStore.
func (l *SnapshotLoader) Store() *SnapshotStore { return l.store }
