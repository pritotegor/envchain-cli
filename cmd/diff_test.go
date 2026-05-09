package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestDiffCmd_RequiresArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "diff")
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestDiffCmd_RequiresExactlyTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "diff", "only-one")
	if err == nil {
		t.Fatal("expected error for single arg")
	}
}

func TestDiffCmd_SourceNotFound(t *testing.T) {
	dir := setupTestStore(t)
	t.Setenv("ENVCHAIN_STORE", dir)

	_, err := executeCommand(rootCmd, "diff", "ghost", "also-ghost")
	if err == nil {
		t.Fatal("expected error for unknown chains")
	}
}

func TestDiffCmd_NoChanges(t *testing.T) {
	dir := setupTestStore(t)
	t.Setenv("ENVCHAIN_STORE", dir)

	// seed two identical chains via the set command
	_, _ = executeCommand(rootCmd, "set", "alpha", "FOO=bar")
	_, _ = executeCommand(rootCmd, "set", "beta", "FOO=bar")

	out, err := executeCommand(rootCmd, "diff", "alpha", "beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-differences message, got: %q", out)
	}
}

func TestDiffCmd_ShowsAddedKey(t *testing.T) {
	dir := setupTestStore(t)
	t.Setenv("ENVCHAIN_STORE", dir)

	_, _ = executeCommand(rootCmd, "set", "src", "FOO=1")
	_, _ = executeCommand(rootCmd, "set", "dst", "FOO=1", "BAR=2")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	_, err := executeCommand(rootCmd, "diff", "src", "dst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "BAR") {
		t.Errorf("expected BAR in diff output, got: %q", output)
	}
}
