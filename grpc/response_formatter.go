package grpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// FormatStyle controls the output format of a gRPC response.
type FormatStyle int

const (
	FormatJSON    FormatStyle = iota // pretty-printed JSON
	FormatCompact                    // single-line JSON
	FormatText                       // human-readable key=value
)

// ResponseFormatter converts a parsed gRPC response map into a
// display-ready string using the requested style.
type ResponseFormatter struct {
	style  FormatStyle
	indent string
}

// NewResponseFormatter creates a ResponseFormatter with the given style.
// indent is only used for FormatJSON; pass "" to use the default ("  ").
func NewResponseFormatter(style FormatStyle, indent string) *ResponseFormatter {
	if indent == "" {
		indent = "  "
	}
	return &ResponseFormatter{style: style, indent: indent}
}

// Format converts a map[string]any (as returned by ResponseParser.ToMap)
// into a formatted string. It returns an empty string for a nil map.
func (f *ResponseFormatter) Format(data map[string]any) (string, error) {
	if data == nil {
		return "", nil
	}
	switch f.style {
	case FormatCompact:
		return f.compact(data)
	case FormatText:
		return f.text(data, 0), nil
	default:
		return f.pretty(data)
	}
}

func (f *ResponseFormatter) pretty(data map[string]any) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", f.indent)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		return "", fmt.Errorf("response formatter: %w", err)
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

func (f *ResponseFormatter) compact(data map[string]any) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("response formatter: %w", err)
	}
	return string(b), nil
}

func (f *ResponseFormatter) text(data map[string]any, depth int) string {
	pad := strings.Repeat("  ", depth)
	var sb strings.Builder
	for k, v := range data {
		switch child := v.(type) {
		case map[string]any:
			sb.WriteString(fmt.Sprintf("%s%s:\n", pad, k))
			sb.WriteString(f.text(child, depth+1))
		default:
			sb.WriteString(fmt.Sprintf("%s%s = %v\n", pad, k, v))
		}
	}
	return sb.String()
}

// Style returns the current FormatStyle.
func (f *ResponseFormatter) Style() FormatStyle { return f.style }
