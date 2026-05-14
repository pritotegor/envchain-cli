package lint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/envchain-cli/internal/chain"
	"github.com/envchain-cli/envchain-cli/internal/lint"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "lint-test-*")
	if err != nil {
		t.Fatalf("tempStore: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := chain.NewStore(filepath.Join(dir, "chains.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestLint_NoIssues(t *testing.T) {
	s := tempStore(t)
	_ = s.Add("prod", map[string]string{"DATABASE_URL": "postgres://", "PORT": "5432"})

	l := lint.NewLinter(s)
	issues, err := l.Lint()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	s := tempStore(t)
	_ = s.Add("dev", map[string]string{"db_host": "localhost"})

	l := lint.NewLinter(s)
	issues, err := l.Lint()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %v", len(issues), issues)
	}
	if issues[0].Key != "db_host" {
		t.Errorf("unexpected key in issue: %q", issues[0].Key)
	}
}

func TestLint_InvalidKey(t *testing.T) {
	s := tempStore(t)
	_ = s.Add("ci", map[string]string{"1INVALID": "value"})

	l := lint.NewLinter(s)
	issues, err := l.Lint()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected at least one issue for invalid key")
	}
}

func TestLint_FilterByChainName(t *testing.T) {
	s := tempStore(t)
	_ = s.Add("prod", map[string]string{"GOOD_KEY": "v"})
	_ = s.Add("dev", map[string]string{"bad_key": "v"})

	l := lint.NewLinter(s)
	issues, err := l.Lint("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected no issues for prod chain, got %v", issues)
	}
}

func TestLint_IssueString(t *testing.T) {
	i := lint.Issue{Chain: "dev", Key: "foo", Message: "key is not upper-case"}
	want := "dev.foo: key is not upper-case"
	if i.String() != want {
		t.Errorf("String() = %q, want %q", i.String(), want)
	}
}
