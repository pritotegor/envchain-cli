package envimport_test

import (
	"os"
	"strings"
	"testing"

	envimport "github.com/user/envchain-cli/internal/import"
	"github.com/user/envchain-cli/internal/chain"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "envimport-test-*")
	if err != nil {
		t.Fatalf("tempStore: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := chain.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestImport_Basic(t *testing.T) {
	s := tempStore(t)
	im := envimport.NewImporter(s)

	input := strings.NewReader("FOO=bar\nBAZ=qux\n")
	res, err := im.Import(input, "mychain", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Imported != 2 {
		t.Errorf("imported=%d, want 2", res.Imported)
	}
	v, _ := s.Get("mychain", "FOO")
	if v != "bar" {
		t.Errorf("FOO=%q, want \"bar\"", v)
	}
}

func TestImport_SkipsCommentAndBlank(t *testing.T) {
	s := tempStore(t)
	im := envimport.NewImporter(s)

	input := strings.NewReader("# comment\n\nKEY=val\n")
	res, err := im.Import(input, "chain", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Imported != 1 {
		t.Errorf("imported=%d, want 1", res.Imported)
	}
}

func TestImport_ExportPrefix(t *testing.T) {
	s := tempStore(t)
	im := envimport.NewImporter(s)

	input := strings.NewReader("export MY_VAR=hello\n")
	res, err := im.Import(input, "chain", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Imported != 1 {
		t.Errorf("imported=%d, want 1", res.Imported)
	}
	v, _ := s.Get("chain", "MY_VAR")
	if v != "hello" {
		t.Errorf("MY_VAR=%q, want \"hello\"", v)
	}
}

func TestImport_NoOverwrite(t *testing.T) {
	s := tempStore(t)
	_ = s.Add("chain", "KEY", "original")
	im := envimport.NewImporter(s)

	input := strings.NewReader("KEY=new\n")
	res, err := im.Import(input, "chain", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("skipped=%d, want 1", res.Skipped)
	}
	v, _ := s.Get("chain", "KEY")
	if v != "original" {
		t.Errorf("KEY=%q, want \"original\"", v)
	}
}

func TestImport_Overwrite(t *testing.T) {
	s := tempStore(t)
	_ = s.Add("chain", "KEY", "original")
	im := envimport.NewImporter(s)

	input := strings.NewReader("KEY=updated\n")
	res, err := im.Import(input, "chain", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Overwrote != 1 {
		t.Errorf("overwrote=%d, want 1", res.Overwrote)
	}
	v, _ := s.Get("chain", "KEY")
	if v != "updated" {
		t.Errorf("KEY=%q, want \"updated\"", v)
	}
}

func TestImport_QuotedValues(t *testing.T) {
	s := tempStore(t)
	im := envimport.NewImporter(s)

	input := strings.NewReader(`A="double"` + "\n" + `B='single'` + "\n")
	_, err := im.Import(input, "chain", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := s.Get("chain", "A"); v != "double" {
		t.Errorf("A=%q, want \"double\"", v)
	}
	if v, _ := s.Get("chain", "B"); v != "single" {
		t.Errorf("B=%q, want \"single\"", v)
	}
}
