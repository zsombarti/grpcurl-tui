package grpc

import (
	"regexp"
	"strings"
)

// EnvSubstitutor replaces ${VAR} or $VAR placeholders in strings using a
// provided environment map, falling back to a default value when a variable
// is not found.
type EnvSubstitutor struct {
	env      map[string]string
	pattern  *regexp.Regexp
	fallback string
}

// NewEnvSubstitutor creates an EnvSubstitutor backed by the given env map.
// An empty fallback string is used when a variable is not present in the map.
func NewEnvSubstitutor(env map[string]string, fallback string) *EnvSubstitutor {
	if env == nil {
		env = make(map[string]string)
	}
	return &EnvSubstitutor{
		env:      env,
		pattern:  regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`),
		fallback: fallback,
	}
}

// Substitute replaces all environment variable references in s and returns
// the resulting string.
func (e *EnvSubstitutor) Substitute(s string) string {
	return e.pattern.ReplaceAllStringFunc(s, func(match string) string {
		key := strings.TrimPrefix(match, "$")
		key = strings.TrimPrefix(key, "{")
		key = strings.TrimSuffix(key, "}")
		if val, ok := e.env[key]; ok {
			return val
		}
		return e.fallback
	})
}

// SubstituteMap applies Substitute to every value in the provided map and
// returns a new map with the substituted values.
func (e *EnvSubstitutor) SubstituteMap(pairs map[string]string) map[string]string {
	out := make(map[string]string, len(pairs))
	for k, v := range pairs {
		out[k] = e.Substitute(v)
	}
	return out
}

// Len returns the number of variables available for substitution.
func (e *EnvSubstitutor) Len() int {
	return len(e.env)
}
