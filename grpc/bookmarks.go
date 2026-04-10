package grpc

import (
	"errors"
	"sync"
	"time"
)

// Bookmark represents a saved gRPC request configuration.
type Bookmark struct {
	Name      string
	Address   string
	Method    string
	Payload   string
	Metadata  []string
	CreatedAt time.Time
}

// BookmarkStore manages saved request bookmarks.
type BookmarkStore struct {
	mu        sync.RWMutex
	bookmarks []Bookmark
	maxSize   int
}

// NewBookmarkStore creates a BookmarkStore with the given capacity.
func NewBookmarkStore(maxSize int) *BookmarkStore {
	if maxSize <= 0 {
		maxSize = 50
	}
	return &BookmarkStore{
		bookmarks: make([]Bookmark, 0, maxSize),
		maxSize:   maxSize,
	}
}

// Add saves a bookmark, evicting the oldest if at capacity.
func (s *BookmarkStore) Add(b Bookmark) error {
	if b.Name == "" {
		return errors.New("bookmark name must not be empty")
	}
	if b.Address == "" {
		return errors.New("bookmark address must not be empty")
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.bookmarks) >= s.maxSize {
		s.bookmarks = s.bookmarks[1:]
	}
	s.bookmarks = append(s.bookmarks, b)
	return nil
}

// All returns a copy of all stored bookmarks.
func (s *BookmarkStore) All() []Bookmark {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Bookmark, len(s.bookmarks))
	copy(out, s.bookmarks)
	return out
}

// Delete removes a bookmark by name. Returns error if not found.
func (s *BookmarkStore) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, b := range s.bookmarks {
		if b.Name == name {
			s.bookmarks = append(s.bookmarks[:i], s.bookmarks[i+1:]...)
			return nil
		}
	}
	return errors.New("bookmark not found: " + name)
}

// Len returns the number of stored bookmarks.
func (s *BookmarkStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.bookmarks)
}

// Clear removes all bookmarks.
func (s *BookmarkStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bookmarks = s.bookmarks[:0]
}
