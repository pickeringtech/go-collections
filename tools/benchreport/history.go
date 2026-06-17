package main

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// The long-term trend store (issue #51) is a directory of per-commit snapshots,
// one file per push to `main`, named "<timestamp>_<sha>.txt":
//
//	docs/bench/history/2026-06-16T19-23-20Z_a894dc3.txt
//
// Each file is *raw* `go test -bench` output — the same multi-sample text the
// snapshot pipeline already produces (build/bench.txt) — with no provenance
// preamble. Two deliberate consequences:
//
//   - The timestamp prefix is filename-safe (the ISO colons become dashes) and
//     lexically sortable, so a plain name sort is chronological; provenance
//     (when, which commit) lives entirely in the name.
//   - Because the files hold nothing but benchmark output, benchstat reads them
//     directly — `benchstat <old>.txt <new>.txt` is the documented manual recipe
//     for comparing any two commits, with no preprocessing.
//
// Keeping the *multi-sample* output (not just benchstat's medians) is what lets
// the regression check recover statistical significance between commits; medians
// alone could only support untrustworthy raw-delta comparisons.

const (
	// historyCapDefault bounds repo growth: only the most recent N snapshots are
	// retained, the oldest pruned on each push (issue #51 retention policy).
	historyCapDefault = 100
	historyExt        = ".txt"
)

// HistoryEntry is one retained commit snapshot: its provenance (parsed from the
// filename) and the per-benchmark ns/op sample sets parsed from its contents.
type HistoryEntry struct {
	File    string                  // basename, e.g. 2026-06-16T19-23-20Z_a894dc3.txt
	Stamp   string                  // timestamp portion of the name (dashes for colons)
	Commit  string                  // short SHA
	Samples map[sampleKey][]float64 // conforming benchmark cell -> ns/op samples
}

// historyFilename builds the deterministic, sortable snapshot name for a given
// generation timestamp and commit. Colons in the ISO timestamp are replaced with
// dashes so the name is valid on every filesystem (Windows included); empty
// fields fall back to stable placeholders so a name is always well-formed.
func historyFilename(date, commit string) string {
	stamp := strings.ReplaceAll(date, ":", "-")
	if stamp == "" {
		stamp = "unknown"
	}
	if commit == "" {
		commit = "nocommit"
	}
	return stamp + "_" + commit + historyExt
}

// parseHistoryName splits a snapshot basename back into its timestamp and commit
// parts. The commit is whatever follows the final underscore; everything before
// it is the timestamp (which itself contains no underscore).
func parseHistoryName(file string) (stamp, commit string) {
	base := strings.TrimSuffix(file, historyExt)
	if i := strings.LastIndex(base, "_"); i >= 0 {
		return base[:i], base[i+1:]
	}
	return base, ""
}

// AddHistoryEntry writes one raw-bench snapshot into the history directory under
// its deterministic name, then prunes the oldest entries beyond cap. It returns
// the basename written and the basenames pruned (for logging). Re-running with
// the same timestamp and commit overwrites in place, so the step is idempotent.
func AddHistoryEntry(dir, raw, date, commit string, cap int) (added string, pruned []string, err error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", nil, err
	}
	added = historyFilename(date, commit)
	if err := os.WriteFile(filepath.Join(dir, added), []byte(raw), 0o644); err != nil {
		return "", nil, err
	}
	pruned, err = pruneHistory(dir, cap)
	if err != nil {
		return added, nil, err
	}
	return added, pruned, nil
}

// pruneHistory keeps only the newest cap snapshots, deleting the oldest. A cap
// below 1 disables pruning (keep everything). Returns the basenames removed.
func pruneHistory(dir string, cap int) ([]string, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*"+historyExt))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths) // chronological: oldest first
	if cap < 1 || len(paths) <= cap {
		return nil, nil
	}
	var pruned []string
	for _, p := range paths[:len(paths)-cap] {
		if err := os.Remove(p); err != nil {
			return pruned, err
		}
		pruned = append(pruned, filepath.Base(p))
	}
	return pruned, nil
}

// LoadHistory reads every snapshot in dir, oldest first, into parsed entries. A
// missing directory yields no entries (the trend section is simply omitted until
// the first push populates it). Files with no conforming benchmarks are skipped
// rather than failing the whole render.
func LoadHistory(dir string) ([]HistoryEntry, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*"+historyExt))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths) // oldest -> newest
	var entries []HistoryEntry
	for _, p := range paths {
		raw, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		samples := parseRawBench(string(raw))
		if len(samples) == 0 {
			continue
		}
		base := filepath.Base(p)
		stamp, commit := parseHistoryName(base)
		entries = append(entries, HistoryEntry{
			File: base, Stamp: stamp, Commit: commit, Samples: samples,
		})
	}
	return entries, nil
}

// parseRawBench extracts ns/op samples from raw `go test -bench` output, keyed by
// the same standardized (pkg, impl, op, size) cell the rest of the tool uses. It
// is the multi-sample counterpart to LoadDataset (which reads benchstat's
// median-only CSV): every repeat of a benchmark contributes one sample, so the
// significance test has a real distribution to work with.
//
// Lines that aren't `pkg:` config or a conforming "Benchmark<Impl>_<Op>/size_<N>"
// result (geomean rows, cross-impl comparison benchmarks, PASS/ok footers) are
// ignored, mirroring LoadDataset's allowlist.
func parseRawBench(s string) map[sampleKey][]float64 {
	out := map[sampleKey][]float64{}
	pkg := ""
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimRight(line, "\r")
		if rest, ok := strings.CutPrefix(line, "pkg:"); ok {
			pkg = strings.TrimSpace(rest)
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 || !strings.HasPrefix(fields[0], "Benchmark") {
			continue
		}
		m := nameRe.FindStringSubmatch(strings.TrimPrefix(fields[0], "Benchmark"))
		if m == nil {
			continue
		}
		ns, ok := nsOpValue(fields)
		if !ok {
			continue
		}
		size, _ := strconv.Atoi(m[3])
		key := sampleKey{path.Base(pkg), m[1], m[2], size}
		out[key] = append(out[key], ns)
	}
	return out
}

// nsOpValue finds the ns/op measurement on a benchmark result line: the numeric
// token immediately preceding the "ns/op" unit. Returns false if absent or
// unparseable (e.g. a benchmark that reported no time).
func nsOpValue(fields []string) (float64, bool) {
	for i := 1; i < len(fields); i++ {
		if fields[i] == "ns/op" {
			if v, err := strconv.ParseFloat(fields[i-1], 64); err == nil {
				return v, true
			}
		}
	}
	return 0, false
}
