// Command cihealth computes the main-branch CI health badges (issue #209): the
// percentage of healthy `CI` runs on `main` over rolling 7-, 30- and 90-day
// windows, rendered as shields.io endpoint JSON.
//
// It is a pure function of its inputs. GitHub's `actions/runs` API is fetched by
// the caller (the scheduled ci-health-badges workflow, via `gh api`) and piped
// in as newline-delimited JSON; this tool merges those runs into the committed
// tally store, prunes the aged-out tail, recomputes the three windows from the
// store, and writes the store back plus one badge JSON per window. Persisting
// the tally is what keeps the 90-day badge honest once runs age past GitHub's
// Actions retention limit — see store.go.
//
// Each NDJSON run object has the shape emitted by the workflow's `gh api --jq`:
//
//	{"id":123,"sha":"a894dc3","conclusion":"success","timestamp":"2026-06-16T19:23:20Z"}
//
// Re-running with unchanged inputs and flags produces byte-identical output, so
// the badge job's "commit only if changed" step stays a no-op on a quiet day.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "cihealth:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("cihealth", flag.ContinueOnError)
	var (
		runsPath  = fs.String("runs", "-", "newly-fetched runs as NDJSON (\"-\" for stdin)")
		storePath = fs.String("store", "docs/ci-health/history.csv", "persisted tally store CSV")
		outDir    = fs.String("out", "docs/ci-health", "directory to write the per-window badge JSON into")
		nowStr    = fs.String("now", "", "the \"now\" instant as RFC-3339 (defaults to the current time)")
		retention = fs.Int("retention-days", 100, "drop store rows older than this many days (a margin past the 90d window; <1 keeps all)")
		minSample = fs.Int("min-sample", 3, "grey a window with fewer than this many counted runs (too few to colour-code)")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}

	now := time.Now().UTC()
	if *nowStr != "" {
		parsed, err := time.Parse(time.RFC3339, *nowStr)
		if err != nil {
			return fmt.Errorf("-now: %w", err)
		}
		now = parsed.UTC()
	}

	raw, err := readInput(*runsPath)
	if err != nil {
		return err
	}
	fetched, err := ParseRuns(raw)
	if err != nil {
		return err
	}

	stored, err := LoadStore(*storePath)
	if err != nil {
		return err
	}

	records := Prune(Merge(stored, fetched), now, *retention)
	if err := SaveStore(*storePath, records); err != nil {
		return err
	}

	for _, w := range windows {
		badge := Render(w.Label, Count(records, now, w.Days), *minSample)
		if err := writeBadge(filepath.Join(*outDir, w.File), badge); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "%s: %s\n", w.Label, badge.Message)
	}
	fmt.Fprintf(os.Stderr, "store: %d run(s) retained (%d fetched, pruned to %dd) → %s\n",
		len(records), len(fetched), *retention, *storePath)
	return nil
}

// writeBadge serialises one shields endpoint badge to path (creating parents),
// pretty-printed with a trailing newline so the committed file is diff-friendly.
func writeBadge(path string, b Badge) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func readInput(path string) (io.Reader, error) {
	if path == "-" {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(b), nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
