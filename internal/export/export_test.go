package export_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain-cli/internal/export"
)

func TestWrite_Shell(t *testing.T) {
	var buf strings.Builder
	ex := export.NewExporter(&buf, export.FormatShell)
	vars := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	if err := ex.Write(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export BAZ='qux'") {
		t.Errorf("expected shell export for BAZ, got:\n%s", out)
	}
	if !strings.Contains(out, "export FOO='bar'") {
		t.Errorf("expected shell export for FOO, got:\n%s", out)
	}
}

func TestWrite_Dotenv(t *testing.T) {
	var buf strings.Builder
	ex := export.NewExporter(&buf, export.FormatDotenv)
	vars := map[string]string{"TOKEN": "abc123"}
	if err := ex.Write(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "TOKEN=abc123\n" {
		t.Errorf("expected 'TOKEN=abc123\\n', got %q", got)
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf strings.Builder
	ex := export.NewExporter(&buf, export.FormatJSON)
	vars := map[string]string{"KEY": "val"}
	if err := ex.Write(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"KEY": "val"`) {
		t.Errorf("expected JSON entry, got:\n%s", out)
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf strings.Builder
	ex := export.NewExporter(&buf, export.Format("xml"))
	err := ex.Write(map[string]string{"A": "B"})
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestWrite_ShellQuoteSpecialChars(t *testing.T) {
	var buf strings.Builder
	ex := export.NewExporter(&buf, export.FormatShell)
	vars := map[string]string{"PW": "it's a secret"}
	if err := ex.Write(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	// single quote inside value must be escaped
	if !strings.Contains(out, `'\''`) {
		t.Errorf("expected escaped single quote in output, got:\n%s", out)
	}
}
