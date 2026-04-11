package grpc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvLoader loads environment variables from a .env-style file
// and exposes them as gRPC metadata key-value pairs.
type EnvLoader struct {
	pairs []string
}

// NewEnvLoader returns a new EnvLoader instance.
func NewEnvLoader() *EnvLoader {
	return &EnvLoader{}
}

// LoadFile reads key=value pairs from the given file path.
// Lines starting with '#' and blank lines are ignored.
func (e *EnvLoader) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("env_loader: open %q: %w", path, err)
	}
	defer f.Close()

	var pairs []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "=") {
			return fmt.Errorf("env_loader: malformed line %q: missing '='" , line)
		}
		pairs = append(pairs, line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("env_loader: scan: %w", err)
	}

	e.pairs = pairs
	return nil
}

// Pairs returns the loaded key=value pairs.
func (e *EnvLoader) Pairs() []string {
	out := make([]string, len(e.pairs))
	copy(out, e.pairs)
	return out
}

// Len returns the number of loaded pairs.
func (e *EnvLoader) Len() int {
	return len(e.pairs)
}

// Clear removes all loaded pairs.
func (e *EnvLoader) Clear() {
	e.pairs = nil
}
