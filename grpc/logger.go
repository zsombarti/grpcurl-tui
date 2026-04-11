package grpc

import (
	"fmt"
	"sync"
	"time"
)

// LogLevel represents the severity of a log entry.
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogEntry holds a single structured log record.
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
}

// Logger captures structured log entries in memory for display in the TUI.
type Logger struct {
	mu      sync.RWMutex
	entries []LogEntry
	maxSize int
	minLevel LogLevel
}

// NewLogger creates a Logger that retains up to maxSize entries.
func NewLogger(maxSize int, minLevel LogLevel) *Logger {
	if maxSize <= 0 {
		maxSize = 200
	}
	return &Logger{maxSize: maxSize, minLevel: minLevel}
}

// log appends an entry if its level meets the minimum threshold.
func (l *Logger) log(level LogLevel, format string, args ...any) {
	if level < l.minLevel {
		return
	}
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   fmt.Sprintf(format, args...),
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.entries) >= l.maxSize {
		l.entries = l.entries[1:]
	}
	l.entries = append(l.entries, entry)
}

func (l *Logger) Debug(format string, args ...any) { l.log(LogLevelDebug, format, args...) }
func (l *Logger) Info(format string, args ...any)  { l.log(LogLevelInfo, format, args...) }
func (l *Logger) Warn(format string, args ...any)  { l.log(LogLevelWarn, format, args...) }
func (l *Logger) Error(format string, args ...any) { l.log(LogLevelError, format, args...) }

// Entries returns a snapshot of all retained log entries.
func (l *Logger) Entries() []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	copy := make([]LogEntry, len(l.entries))
	for i, e := range l.entries {
		copy[i] = e
	}
	return copy
}

// Len returns the number of retained entries.
func (l *Logger) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.entries)
}

// Clear removes all retained entries.
func (l *Logger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = nil
}
