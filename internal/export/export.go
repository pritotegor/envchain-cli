// Package export provides functionality to serialize environment chains
// into various shell-compatible formats for sourcing or inspection.
package export

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents an output format for exported variables.
type Format string

const (
	FormatShell  Format = "shell"  // export KEY=VALUE
	FormatDotenv Format = "dotenv" // KEY=VALUE
	FormatJSON   Format = "json"   // {"KEY": "VALUE"}
)

// Exporter writes environment variables in a specified format.
type Exporter struct {
	format Format
	writer io.Writer
}

// NewExporter creates a new Exporter writing to w in the given format.
func NewExporter(w io.Writer, format Format) *Exporter {
	return &Exporter{format: format, writer: w}
}

// Write serializes the provided key-value pairs to the writer.
func (e *Exporter) Write(vars map[string]string) error {
	switch e.format {
	case FormatShell:
		return e.writeShell(vars)
	case FormatDotenv:
		return e.writeDotenv(vars)
	case FormatJSON:
		return e.writeJSON(vars)
	default:
		return fmt.Errorf("export: unknown format %q", e.format)
	}
}

func sortedKeys(vars map[string]string) []string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (e *Exporter) writeShell(vars map[string]string) error {
	for _, k := range sortedKeys(vars) {
		_, err := fmt.Fprintf(e.writer, "export %s=%s\n", k, shellQuote(vars[k]))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeDotenv(vars map[string]string) error {
	for _, k := range sortedKeys(vars) {
		_, err := fmt.Fprintf(e.writer, "%s=%s\n", k, vars[k])
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeJSON(vars map[string]string) error {
	keys := sortedKeys(vars)
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		sb.WriteString(fmt.Sprintf("  %q: %q", k, vars[k]))
		if i < len(keys)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
	_, err := fmt.Fprint(e.writer, sb.String())
	return err
}

// shellQuote wraps a value in single quotes, escaping existing single quotes.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\'''") + "'"
}
