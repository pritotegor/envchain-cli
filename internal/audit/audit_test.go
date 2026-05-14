package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain-cli/internal/audit"
)

func tempLog(t *testing.T) *audit.Logger {
	t.Helper()
	dir := t.TempDir()
	l, err := audit.NewLogger(filepath.Join(dir, "audit.log"))
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	return l
}

func TestRecord_And_ReadAll(t *testing.T) {
	l := tempLog(t)

	if err := l.Record("add", "myapp", "DB_URL", ""); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := l.Record("delete", "myapp", "OLD_KEY", ""); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Action != "add" || entries[0].Chain != "myapp" || entries[0].Key != "DB_URL" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Action != "delete" || entries[1].Key != "OLD_KEY" {
		t.Errorf("unexpected second entry: %+v", entries[1])
	}
}

func TestReadAll_EmptyFile(t *testing.T) {
	l := tempLog(t)
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecord_CreatesParentDir(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "nested", "deep", "audit.log")
	l, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	if err := l.Record("rename", "proj", "", "old->new"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("log file not created: %v", err)
	}
}

func TestRecord_AppendsPersists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	for i := 0; i < 5; i++ {
		l, _ := audit.NewLogger(path)
		_ = l.Record("add", "chain", "KEY", "")
	}

	l, _ := audit.NewLogger(path)
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 5 {
		t.Errorf("expected 5 entries after 5 appends, got %d", len(entries))
	}
}
