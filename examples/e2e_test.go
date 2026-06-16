// Package examples holds the golden-output end-to-end tests for the example
// apps. Each case builds and runs a real binary (via `go run`) as a downstream
// consumer of the library and asserts its stdout against a checked-in golden
// file. A drift in output, a panic, or a failure to compile fails the test —
// which is the point: these tests are the library's outside-consumer smoke test.
//
// Regenerate the golden files after an intentional output change with:
//
//	go test ./... -update
package examples

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "regenerate the testdata/*.golden files")

// e2eCase is one (app, args, stdin) -> golden-file expectation.
type e2eCase struct {
	name  string   // also the testdata/<name>.golden basename
	app   string   // the cmd/<app> directory to run
	args  []string // command-line arguments
	stdin string   // fed to the program's stdin, if non-empty
}

func cases() []e2eCase {
	return []e2eCase{
		{
			name:  "word-frequency",
			app:   "word-frequency",
			args:  []string{"-n", "5"},
			stdin: "The quick brown fox. The lazy dog. The quick dog runs; the fox runs!\n",
		},
		{
			name: "set-algebra",
			app:  "set-algebra",
			args: []string{"-a", "apple,banana,cherry,date", "-b", "banana,date,fig"},
		},
		{
			name: "worker-pipeline",
			app:  "worker-pipeline",
			args: []string{"-n", "8", "-workers", "3"},
		},
		{
			name: "ordered-processing",
			app:  "ordered-processing",
			args: []string{"-nums", "5,3,8,1,9,2"},
		},
	}
}

func TestExamplesE2E(t *testing.T) {
	for _, tc := range cases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := runApp(t, tc)
			golden := filepath.Join("testdata", tc.name+".golden")

			if *update {
				if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
					t.Fatalf("writing golden file: %v", err)
				}
				return
			}

			wantRaw, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("reading golden file (run `go test ./... -update` to create it): %v", err)
			}
			want := normalise(string(wantRaw))
			if got != want {
				t.Errorf("output for %q does not match %s\n--- got ---\n%s\n--- want ---\n%s",
					tc.app, golden, got, want)
			}
		})
	}
}

// runApp compiles and runs the example with `go run`, returning its normalised
// stdout. Building through the toolchain (rather than calling the package
// directly) is what makes this a genuine end-to-end check.
func runApp(t *testing.T, tc e2eCase) string {
	t.Helper()
	args := append([]string{"run", "./cmd/" + tc.app}, tc.args...)
	cmd := exec.Command("go", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if tc.stdin != "" {
		cmd.Stdin = strings.NewReader(tc.stdin)
	}
	if err := cmd.Run(); err != nil {
		t.Fatalf("running %q failed: %v\nstderr:\n%s", tc.app, err, stderr.String())
	}
	return normalise(stdout.String())
}

// normalise collapses Windows CRLF line endings to LF so a single golden file
// is valid across the OS matrix.
func normalise(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}
