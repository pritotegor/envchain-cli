package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/envchain-cli/internal/chain"
	"github.com/spf13/cobra"
)

func setupCompareStore(t *testing.T) (string, *chain.Store) {
	t.Helper()
	dir, err := os.MkdirTemp("", "compare-cmd-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := chain.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	return dir, s
}

func runCompareCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs(args)
	_, err := rootCmd.ExecuteC()
	return buf.String(), err
}

func TestCompareCmd_RequiresArgs(t *testing.T) {
	_, err := runCompareCmd(t, "compare")
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestCompareCmd_RequiresExactlyTwoArgs(t *testing.T) {
	_, err := runCompareCmd(t, "compare", "only-one")
	if err == nil {
		t.Fatal("expected error with one arg")
	}
}

func TestCompareCmd_SourceNotFound(t *testing.T) {
	dir, _ := setupCompareStore(t)
	_, err := runCompareCmd(t, "compare", "--store", dir, "ghost", "also-ghost")
	if err == nil {
		t.Fatal("expected error for missing chain")
	}
}

func TestCompareCmd_Identical(t *testing.T) {
	dir, s := setupCompareStore(t)
	_ = s.Add("x", "KEY", "val")
	_ = s.Add("y", "KEY", "val")

	out, err := runCompareCmd(t, "compare", "--store", dir, "x", "y")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "identical") {
		t.Errorf("expected 'identical' in output, got: %s", out)
	}
}

func TestCompareCmd_ShowsDiff(t *testing.T) {
	dir, s := setupCompareStore(t)
	_ = s.Add("x", "FOO", "old")
	_ = s.Add("y", "FOO", "new")
	_ = s.Add("y", "EXTRA", "e")

	out, err := runCompareCmd(t, "compare", "--store", dir, "x", "y")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "~") {
		t.Errorf("expected diff marker '~', got: %s", out)
	}
	if !strings.Contains(out, ">") {
		t.Errorf("expected only-in-B marker '>', got: %s", out)
	}
}

var _ = cobra.Command{} // ensure cobra import used
