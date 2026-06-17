// Command doccompile guards the package-level godoc examples (the indented code
// blocks in every doc.go) against drift from the real API. Confidently-wrong
// quick-starts are worse than missing docs for a library people copy-paste from,
// so this is a cheap PR-time gate that makes the docs provably real (issue #151).
//
// It runs two complementary checks over every indented code block found in a
// doc.go package comment:
//
//  1. Symbol existence (robust baseline). The library's exported API is parsed
//     straight from source; any `pkg.Symbol` reference whose pkg is a library
//     package but whose Symbol the package does not export is a failure. This is
//     immune to the pseudocode that fills illustrative blocks ([]User{...},
//     undefined helpers like isEven, top-level generic declarations) because it
//     only inspects package-qualified identifiers — exactly the class of bug the
//     reviews flagged (collections.Pair, slices.FlatMap, maps.Invert, …).
//
//  2. Compilation (catches the rest). Blocks that parse cleanly as Go — once the
//     intentional pseudocode is filtered out by the parser itself — are compiled
//     against the real packages. Only errors that implicate a library package
//     gate the build; benign fragment artifacts (an unused variable, an unused
//     import, an undefined ambient placeholder) are ignored. This additionally
//     catches method-chain typos on returned values (e.g. a renamed Pipeline
//     method) that a pure symbol check cannot see.
//
// The tool is dependency-free: it uses go/parser to read both the API and the
// doc blocks, and shells out to the installed `go` toolchain to compile the
// parseable blocks in a throwaway module that `replace`s the library with the
// local checkout.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Directories that are not part of the library's root module (nested modules,
// VCS metadata, fuzz corpora) and so never contribute API symbols or doc blocks.
var skipDirs = map[string]bool{".git": true, "testdata": true, "examples": true, "tools": true}

// stdlibImports maps the standard-library qualifiers the docs use to their
// import paths. Over-inclusion is harmless: an import that a block doesn't
// actually use surfaces only as a benign "imported and not used" error, which
// the gating filter ignores.
var stdlibImports = map[string]string{
	"context": "context", "fmt": "fmt", "strings": "strings", "strconv": "strconv",
	"sort": "sort", "errors": "errors", "sync": "sync", "time": "time", "math": "math",
	"os": "os", "io": "io", "bytes": "bytes", "bufio": "bufio", "regexp": "regexp",
	"unicode": "unicode",
}

// symbolRE matches a package-qualified, exported reference: a lowercase-led
// qualifier, a dot, then an exported (capitalised) identifier. Selectors on
// local variables (result.Value, u.Email) match too but are dropped later
// because their qualifier is not a known library package.
var symbolRE = regexp.MustCompile(`\b([a-z][A-Za-z0-9_]*)\.([A-Z][A-Za-z0-9_]*)`)

// qualRE finds every `qualifier.` lead-in, used to decide which imports a block
// needs before compiling it.
var qualRE = regexp.MustCompile(`\b([a-z][a-zA-Z0-9_]*)\.`)

// block is one indented code block lifted from a doc.go package comment.
type block struct {
	section string // the "# Heading" the block sits under, for reporting
	line    int    // 1-based line in the doc.go where the block starts
	code    string // de-indented source, lines joined by "\n"
}

// problem is a single doc-example defect to report.
type problem struct {
	file    string
	line    int
	section string
	kind    string // "symbol" or "compile"
	msg     string
}

func main() {
	root := flag.String("root", ".", "repository root containing the library packages")
	verbose := flag.Bool("v", false, "list blocks skipped by the compile check")
	flag.Parse()

	problems, err := run(*root, *verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "doccompile: %v\n", err)
		os.Exit(2)
	}
	if len(problems) == 0 {
		fmt.Println("doccompile: every doc.go example references real symbols and compiles ✔")
		return
	}
	for _, p := range problems {
		// GitHub annotates the PR from these; they also read fine in a plain log.
		fmt.Fprintf(os.Stderr, "::error file=%s,line=%d::[%s] %q: %s\n", p.file, p.line, p.kind, p.section, p.msg)
	}
	fmt.Fprintf(os.Stderr, "\ndoccompile: %d problem(s) found in doc.go examples\n", len(problems))
	os.Exit(1)
}

func run(root string, verbose bool) ([]problem, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	modPath, goVersion, err := readGoMod(absRoot)
	if err != nil {
		return nil, err
	}
	api, libImports, err := indexAPI(absRoot, modPath)
	if err != nil {
		return nil, fmt.Errorf("indexing library API: %w", err)
	}
	if len(libImports) == 0 {
		// No packages were indexed — almost certainly a wrong -root. Fail loudly
		// rather than silently passing every doc (and rather than letting the
		// compile gate build an over-broad `\b()\.` matcher, see compileEntries).
		return nil, fmt.Errorf("no library packages found under %s — is -root correct?", absRoot)
	}
	docFiles, err := findDocFiles(absRoot)
	if err != nil {
		return nil, err
	}

	var problems []problem
	var entries []compileEntry
	for _, f := range docFiles {
		src, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		rel, _ := filepath.Rel(absRoot, f)
		rel = filepath.ToSlash(rel)
		for _, b := range extractBlocks(string(src)) {
			problems = append(problems, checkSymbols(rel, b, api)...)

			stripped := stripImports(b.code)
			body, ok := wrapForCompile(stripped)
			if ok {
				entries = append(entries, compileEntry{
					file:    rel,
					block:   b,
					body:    body,
					imports: detectImports(stripped, libImports),
				})
			} else if verbose {
				fmt.Printf("skip-compile %s:%d (%q): not parseable as decls or a statement body\n", rel, b.line, b.section)
			}
		}
	}

	cprobs, err := compileEntries(modPath, goVersion, absRoot, libImports, entries)
	if err != nil {
		return nil, err
	}
	problems = append(problems, cprobs...)

	sort.SliceStable(problems, func(i, j int) bool {
		if problems[i].file != problems[j].file {
			return problems[i].file < problems[j].file
		}
		return problems[i].line < problems[j].line
	})
	return problems, nil
}

// readGoMod returns the module path and the `go` directive declared in
// <root>/go.mod. The go version is reused for the synthesized doc-check module
// so it tracks the repo's toolchain instead of a hard-coded literal that could
// drift on a future Go bump.
func readGoMod(absRoot string) (modPath, goVersion string, err error) {
	data, err := os.ReadFile(filepath.Join(absRoot, "go.mod"))
	if err != nil {
		return "", "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "module "):
			modPath = strings.TrimSpace(strings.TrimPrefix(line, "module "))
		case strings.HasPrefix(line, "go "):
			goVersion = strings.TrimSpace(strings.TrimPrefix(line, "go "))
		}
	}
	if modPath == "" {
		return "", "", fmt.Errorf("no module path found in go.mod")
	}
	if goVersion == "" {
		return "", "", fmt.Errorf("no go directive found in go.mod")
	}
	return modPath, goVersion, nil
}

// indexAPI walks the root module and records, per package, the set of exported
// top-level identifiers (funcs without a receiver, types, vars, consts). It also
// returns each package's import path, so the compile step knows how to import it.
func indexAPI(absRoot, modPath string) (map[string]map[string]bool, map[string]string, error) {
	api := map[string]map[string]bool{}
	libImports := map[string]string{}
	fset := token.NewFileSet()

	err := filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != absRoot && skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, perr := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if perr != nil {
			return nil // a file that doesn't parse can't define API we can rely on
		}
		pkg := f.Name.Name
		if pkg == "main" || strings.HasSuffix(pkg, "_test") {
			return nil
		}
		relDir, _ := filepath.Rel(absRoot, filepath.Dir(path))
		importPath := modPath
		if relDir != "." {
			importPath = modPath + "/" + filepath.ToSlash(relDir)
		}
		if api[pkg] == nil {
			api[pkg] = map[string]bool{}
		}
		libImports[pkg] = importPath
		for _, decl := range f.Decls {
			switch dd := decl.(type) {
			case *ast.FuncDecl:
				if dd.Recv == nil && dd.Name.IsExported() {
					api[pkg][dd.Name.Name] = true
				}
			case *ast.GenDecl:
				for _, spec := range dd.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						if s.Name.IsExported() {
							api[pkg][s.Name.Name] = true
						}
					case *ast.ValueSpec:
						for _, n := range s.Names {
							if n.IsExported() {
								api[pkg][n.Name] = true
							}
						}
					}
				}
			}
		}
		return nil
	})
	return api, libImports, err
}

// findDocFiles returns every doc.go in the root module (nested modules skipped).
func findDocFiles(absRoot string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != absRoot && skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == "doc.go" {
			files = append(files, path)
		}
		return nil
	})
	sort.Strings(files)
	return files, err
}

// extractBlocks returns the indented code blocks from a Go file's package doc
// comment. Godoc treats a tab-indented run of comment lines (blank comment lines
// inside it included) as preformatted code; a non-indented prose line ends the
// block. Each block is de-indented by one tab and tagged with the "# Heading" it
// sits under.
func extractBlocks(src string) []block {
	lines := strings.Split(src, "\n")

	pkgIdx := -1
	for i, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "package ") {
			pkgIdx = i
			break
		}
	}
	if pkgIdx < 0 {
		return nil
	}
	// The doc comment is the contiguous run of //-lines directly above `package`.
	start := pkgIdx
	for start-1 >= 0 && strings.HasPrefix(strings.TrimSpace(lines[start-1]), "//") {
		start--
	}

	var blocks []block
	section := ""
	var cur []string
	blankHold := 0
	curStart := 0
	flush := func() {
		for len(cur) > 0 && strings.TrimSpace(cur[len(cur)-1]) == "" {
			cur = cur[:len(cur)-1]
		}
		if len(cur) > 0 {
			blocks = append(blocks, block{section: section, line: curStart, code: strings.Join(cur, "\n")})
		}
		cur = nil
		blankHold = 0
	}

	for i := start; i < pkgIdx; i++ {
		raw := lines[i]
		marker := strings.Index(raw, "//")
		if marker < 0 {
			continue
		}
		content := raw[marker+2:]
		switch {
		case strings.TrimSpace(content) == "":
			if len(cur) > 0 {
				blankHold++ // a blank inside a code run is held, not a terminator
			}
		case strings.HasPrefix(content, "\t"):
			if len(cur) == 0 {
				curStart = i + 1 // 1-based file line of the block's first code line
			}
			for ; blankHold > 0; blankHold-- {
				cur = append(cur, "")
			}
			cur = append(cur, strings.TrimPrefix(content, "\t"))
		default:
			flush()
			t := strings.TrimSpace(content)
			if strings.HasPrefix(t, "# ") {
				section = strings.TrimSpace(t[2:])
			}
		}
	}
	flush()
	return blocks
}

// checkSymbols flags every `pkg.Symbol` reference in a block whose pkg is a
// library package but whose Symbol that package does not export.
func checkSymbols(file string, b block, api map[string]map[string]bool) []problem {
	var ps []problem
	for _, m := range symbolRE.FindAllStringSubmatchIndex(b.code, -1) {
		pkg := b.code[m[2]:m[3]]
		sym := b.code[m[4]:m[5]]
		exported, ok := api[pkg]
		if !ok || exported[sym] {
			continue // not a library package, or the symbol really exists
		}
		ps = append(ps, problem{
			file:    file,
			line:    b.line + strings.Count(b.code[:m[0]], "\n"),
			section: b.section,
			kind:    "symbol",
			msg:     fmt.Sprintf("references %s.%s, which package %s does not export", pkg, sym, pkg),
		})
	}
	return ps
}

// stripImports removes import declarations from a block so the compile step can
// supply its own (a block's import of a package we also import would collide).
func stripImports(code string) string {
	var out []string
	inBlock := false
	for _, l := range strings.Split(code, "\n") {
		t := strings.TrimSpace(l)
		switch {
		case inBlock:
			if t == ")" {
				inBlock = false
			}
		case strings.HasPrefix(t, "import ("):
			inBlock = true
		case strings.HasPrefix(t, "import "):
			// single-line import — drop it
		default:
			out = append(out, l)
		}
	}
	return strings.Join(out, "\n")
}

// wrapForCompile decides how a block should be placed in a synthesized file and
// returns that placement, or ok=false if the block is not compilable Go (the
// common case for intentional pseudocode: []User{...}, f(...), bare `...`, or a
// mix of top-level generic decls and statements). It tries file-scope first
// (declarations), then a statement body.
func wrapForCompile(code string) (string, bool) {
	if parses("package doccheck\n" + code + "\n") {
		return code, true
	}
	body := "func _docblock() {\n" + code + "\n}"
	if parses("package doccheck\n" + body + "\n") {
		return body, true
	}
	return "", false
}

func parses(src string) bool {
	_, err := parser.ParseFile(token.NewFileSet(), "", src, parser.SkipObjectResolution)
	return err == nil
}

// detectImports returns the import paths a block needs, inferred from the
// qualifiers it uses. Unknown qualifiers (local placeholders) are ignored.
func detectImports(code string, libImports map[string]string) []string {
	seen := map[string]bool{}
	var imports []string
	for _, m := range qualRE.FindAllStringSubmatch(code, -1) {
		q := m[1]
		path, ok := libImports[q]
		if !ok {
			path, ok = stdlibImports[q]
			if !ok {
				continue
			}
		}
		if !seen[path] {
			seen[path] = true
			imports = append(imports, path)
		}
	}
	sort.Strings(imports)
	return imports
}

// compileEntry is a block staged for the compile check.
type compileEntry struct {
	file    string
	block   block
	body    string   // file-scope decls or a func _docblock(){...} wrapper
	imports []string // import paths to supply
}

// compileEntries builds every staged block in a throwaway module that replaces
// the library with the local checkout, runs `go build ./...`, and returns only
// the errors that implicate a library package. Each block is its own package
// (its own subdirectory) so a name clash or an undefined identifier in one block
// can never affect another.
func compileEntries(modPath, goVersion, absRoot string, libImports map[string]string, entries []compileEntry) ([]problem, error) {
	keys := sortedKeys(libImports)
	if len(entries) == 0 || len(keys) == 0 {
		// No blocks to compile, or no library packages to attribute errors to —
		// the latter would make `\b()\.` match broadly, so don't gate at all.
		return nil, nil
	}
	tmp, err := os.MkdirTemp("", "doccompile-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)

	gomod := fmt.Sprintf("module doccheck\n\ngo %s\n\nrequire %s v0.0.0\n\nreplace %s => %s\n", goVersion, modPath, modPath, absRoot)
	err = os.WriteFile(filepath.Join(tmp, "go.mod"), []byte(gomod), 0o644)
	if err != nil {
		return nil, err
	}

	byDir := make(map[string]compileEntry, len(entries))
	for i, e := range entries {
		dirName := fmt.Sprintf("b%04d", i)
		byDir[dirName] = e
		dir := filepath.Join(tmp, dirName)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return nil, err
		}
		var sb strings.Builder
		sb.WriteString("package doccheck\n\n")
		if len(e.imports) > 0 {
			sb.WriteString("import (\n")
			for _, p := range e.imports {
				fmt.Fprintf(&sb, "\t%q\n", p)
			}
			sb.WriteString(")\n\n")
		}
		sb.WriteString(e.body)
		sb.WriteString("\n")
		err = os.WriteFile(filepath.Join(dir, "doc_block.go"), []byte(sb.String()), 0o644)
		if err != nil {
			return nil, err
		}
	}

	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = tmp
	cmd.Env = append(os.Environ(), "GOWORK=off") // ignore any parent go.work
	out, _ := cmd.CombinedOutput()               // non-zero is expected (benign fragment errors)

	// A module/toolchain failure (as opposed to per-block compile errors) prints
	// `go:`-prefixed lines and no block diagnostics — surface that loudly rather
	// than silently passing.
	if len(out) > 0 && !errLineRE.Match(out) {
		for _, l := range strings.Split(string(out), "\n") {
			if strings.HasPrefix(l, "go: ") {
				return nil, fmt.Errorf("go build failed in the doc-check module:\n%s", out)
			}
		}
	}

	libQualRE := regexp.MustCompile(`\b(` + strings.Join(keys, "|") + `)\.[A-Za-z_]`)
	var ps []problem
	for _, m := range errLineRE.FindAllStringSubmatch(string(out), -1) {
		dirName, msg := m[1], m[2]
		if !isLibraryError(msg, libQualRE) {
			continue // benign fragment artifact, or an error about a placeholder
		}
		e, ok := byDir[dirName]
		if !ok {
			continue
		}
		ps = append(ps, problem{
			file:    e.file,
			line:    e.block.line,
			section: e.block.section,
			kind:    "compile",
			msg:     msg,
		})
	}
	return ps, nil
}

// isLibraryError reports whether a compiler message is a high-confidence
// library-API defect, as opposed to a benign artifact of compiling an isolated
// doc fragment. We gate only on the error *kinds* that unambiguously mean "this
// doc references the library wrongly" AND that name a library package:
//
//   - a missing symbol or method — "undefined: collections.Pair",
//     "x.Collect undefined (type channels.Pipeline has no field or method Collect)"
//   - a call-arity mismatch — "too many arguments in call to slices.Map"
//
// Everything else is ignored on purpose: "declared/imported and not used" and
// "X is not used" are fragment artifacts; "cannot use … in assignment" fires
// when two independent illustrative snippets in one godoc block redeclare a
// variable with different types, which is not a real defect. Those classes would
// otherwise produce false positives even though they mention a library package.
func isLibraryError(msg string, libQualRE *regexp.Regexp) bool {
	if !libQualRE.MatchString(msg) {
		return false
	}
	return strings.Contains(msg, "undefined") || strings.Contains(msg, "arguments in call to")
}

// errLineRE matches a `go build` diagnostic line, capturing the block's
// subdirectory and the compiler message.
var errLineRE = regexp.MustCompile(`(?m)^(b\d+)[/\\][^:]+:\d+:\d+: (.+)$`)

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
