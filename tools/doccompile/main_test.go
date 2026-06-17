package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestExtractBlocks(t *testing.T) {
	src := strings.Join([]string{
		"// Package p does things.",
		"//",
		"// # Quick Start",
		"//",
		"//\tx := p.New()",
		"//\ty := x.Run()",
		"//",
		"//\t_ = y",
		"//",
		"// Prose ends the block.",
		"//",
		"//   - a bullet, NOT code (space-indented)",
		"//",
		"// # Next",
		"//",
		"//\tp.Other()",
		"package p",
		"",
	}, "\n")

	blocks := extractBlocks(src)
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks, want 2: %+v", len(blocks), blocks)
	}
	if blocks[0].section != "Quick Start" {
		t.Errorf("block 0 section = %q, want %q", blocks[0].section, "Quick Start")
	}
	// The blank comment line inside the run is preserved; trailing blanks dropped.
	wantCode := "x := p.New()\ny := x.Run()\n\n_ = y"
	if blocks[0].code != wantCode {
		t.Errorf("block 0 code = %q, want %q", blocks[0].code, wantCode)
	}
	if blocks[1].section != "Next" || blocks[1].code != "p.Other()" {
		t.Errorf("block 1 = %+v, want section Next / code p.Other()", blocks[1])
	}
	// The space-indented bullet must NOT be treated as a code block.
	for _, b := range blocks {
		if strings.Contains(b.code, "bullet") {
			t.Errorf("bullet prose leaked into a code block: %q", b.code)
		}
	}
}

func TestCheckSymbols(t *testing.T) {
	api := map[string]map[string]bool{
		"collections": {"Pair": true, "NewDict": true},
	}
	b := block{
		section: "Demo",
		line:    10,
		code:    "a := collections.NewDict()\nb := collections.Pair{}\nc := collections.Nope{}\nd := local.Thing{}\ne := x.Field",
	}
	got := checkSymbols("collections/doc.go", b, api)
	if len(got) != 1 {
		t.Fatalf("got %d problems, want 1: %+v", len(got), got)
	}
	if !strings.Contains(got[0].msg, "collections.Nope") {
		t.Errorf("problem msg = %q, want it to name collections.Nope", got[0].msg)
	}
	if got[0].line != 12 { // collections.Nope is on the 3rd line of the block (10 + 2)
		t.Errorf("problem line = %d, want 12", got[0].line)
	}
}

func TestStripImports(t *testing.T) {
	in := "import \"fmt\"\nimport (\n\t\"context\"\n\t\"strings\"\n)\nx := 1"
	got := stripImports(in)
	if got != "x := 1" {
		t.Errorf("stripImports = %q, want %q", got, "x := 1")
	}
}

func TestWrapForCompile(t *testing.T) {
	cases := []struct {
		name   string
		code   string
		wantOK bool
	}{
		{"statements -> body", "x := 1\n_ = x", true},
		{"top-level decls -> file scope", "func F[T any](v T) T { return v }", true},
		{"composite-literal ellipsis pseudocode", "users := []User{...}", false},
		{"call ellipsis pseudocode", "x := New(...)", false},
		{"mixed decls and statements", "func F() {}\nx := 1", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, ok := wrapForCompile(c.code)
			if ok != c.wantOK {
				t.Errorf("wrapForCompile(%q) ok = %v, want %v", c.code, ok, c.wantOK)
			}
		})
	}
}

func TestIsLibraryError(t *testing.T) {
	re := regexp.MustCompile(`\b(channels|heaps|lru)\.[A-Za-z_]`)
	cases := []struct {
		msg  string
		want bool
	}{
		{"undefined: channels.Nope", true},
		{"x.Collect undefined (type channels.Pipeline[int,int] has no field or method Collect)", true},
		{"too many arguments in call to channels.Map", true},
		// Benign artifacts that mention a library package must NOT gate.
		{"heaps.Min[int] (value of type func(a int, b int) bool) is not used", false},
		{"cannot use lru.NewLRU[string, []byte](1000) (value of type *lru.LRU) as map[string][]byte value in assignment", false},
		{"\"github.com/x/heaps\" imported and not used", false},
		// An error about a non-library placeholder must NOT gate.
		{"undefined: someLocalVar", false},
	}
	for _, c := range cases {
		got := isLibraryError(c.msg, re)
		if got != c.want {
			t.Errorf("isLibraryError(%q) = %v, want %v", c.msg, got, c.want)
		}
	}
}

// TestRun_Acceptance is the issue-#151 acceptance check: a doc.go example that
// references a non-existent symbol must be reported, and an honest one must not.
// It builds a tiny throwaway module so the full pipeline (API index + symbol
// check + compile) runs end to end.
func TestRun_Acceptance(t *testing.T) {
	_, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain not available")
	}

	write := func(root, badOrGood string) {
		mustWrite(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\ngo 1.24\n")
		mustWrite(t, filepath.Join(root, "foo", "foo.go"),
			"package foo\n\n// Bar is a real exported function.\nfunc Bar() int { return 1 }\n")
		mustWrite(t, filepath.Join(root, "foo", "doc.go"), strings.Join([]string{
			"// Package foo demonstrates the guard.",
			"//",
			"// # Quick Start",
			"//",
			"//\tn := foo." + badOrGood + "()",
			"//\t_ = n",
			"package foo",
			"",
		}, "\n"))
	}

	t.Run("bad symbol fails", func(t *testing.T) {
		root := t.TempDir()
		write(root, "Baz") // foo.Baz does not exist
		problems, err := run(root, false)
		if err != nil {
			t.Fatalf("run: %v", err)
		}
		if len(problems) == 0 {
			t.Fatal("expected the non-existent symbol foo.Baz to be reported, got none")
		}
		var sawSymbol bool
		for _, p := range problems {
			if strings.Contains(p.msg, "Baz") {
				sawSymbol = true
			}
		}
		if !sawSymbol {
			t.Errorf("no problem named foo.Baz: %+v", problems)
		}
	})

	t.Run("honest doc passes", func(t *testing.T) {
		root := t.TempDir()
		write(root, "Bar") // foo.Bar exists
		problems, err := run(root, false)
		if err != nil {
			t.Fatalf("run: %v", err)
		}
		if len(problems) != 0 {
			t.Errorf("honest doc reported problems: %+v", problems)
		}
	})
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path, []byte(content), 0o644)
	if err != nil {
		t.Fatal(err)
	}
}
