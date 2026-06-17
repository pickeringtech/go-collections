package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// rawBench is a structurally faithful slice of `go test -bench` output: a config
// preamble, two repeats of a conforming benchmark, a non-conforming cross-impl
// benchmark (no size sub-benchmark) that must be skipped, a line with no ns/op,
// and the PASS/ok footer.
const rawBench = `goos: linux
goarch: amd64
pkg: github.com/pickeringtech/go-collections/collections/dicts
cpu: TestCPU
BenchmarkHash_Get/size_1000-32    	1000000	    4.44 ns/op	   0 B/op	   0 allocs/op
BenchmarkHash_Get/size_1000-32    	1000000	    4.80 ns/op	   0 B/op	   0 allocs/op
BenchmarkComparison_Get/Hash-32   	 500000	   12.00 ns/op	   0 B/op	   0 allocs/op
BenchmarkBroken_NoUnit-32         	 500000
pkg: github.com/pickeringtech/go-collections/collections/sets
BenchmarkHash_Contains/size_1000-32	2000000	    5.23 ns/op	   0 B/op	   0 allocs/op
PASS
ok  	github.com/pickeringtech/go-collections/collections/sets	1.234s
`

func TestParseRawBench(t *testing.T) {
	samples := parseRawBench(rawBench)

	getKey := sampleKey{"dicts", "Hash", "Get", 1000}
	got := samples[getKey]
	if want := []float64{4.44, 4.80}; !reflect.DeepEqual(got, want) {
		t.Errorf("Hash_Get samples = %v, want %v", got, want)
	}
	if _, ok := samples[sampleKey{"sets", "Hash", "Contains", 1000}]; !ok {
		t.Error("missing sets Hash Contains 1000")
	}
	// The cross-impl comparison (no size_) and the unit-less line are skipped.
	if len(samples) != 2 {
		t.Errorf("got %d keys, want 2: %v", len(samples), samples)
	}
}

func TestNsOpValueMissing(t *testing.T) {
	if _, ok := nsOpValue([]string{"BenchmarkX-8", "500000"}); ok {
		t.Error("expected false when no ns/op token present")
	}
	if v, ok := nsOpValue([]string{"BenchmarkX-8", "5", "12.5", "ns/op"}); !ok || v != 12.5 {
		t.Errorf("nsOpValue = %v,%v want 12.5,true", v, ok)
	}
	// A non-numeric value before ns/op is rejected.
	if _, ok := nsOpValue([]string{"BenchmarkX-8", "5", "NaNish", "ns/op"}); ok {
		t.Error("expected false for unparseable ns/op value")
	}
}

func TestHistoryFilenameRoundTrip(t *testing.T) {
	name := historyFilename("2026-06-16T19:23:20Z", "a894dc3")
	if name != "2026-06-16T19-23-20Z_a894dc3.txt" {
		t.Errorf("filename = %q", name)
	}
	stamp, commit := parseHistoryName(name)
	if stamp != "2026-06-16T19-23-20Z" || commit != "a894dc3" {
		t.Errorf("parsed stamp=%q commit=%q", stamp, commit)
	}

	// Empty fields fall back to placeholders; a name without an underscore parses
	// to an empty commit.
	if got := historyFilename("", ""); got != "unknown_nocommit.txt" {
		t.Errorf("empty fallback = %q", got)
	}
	if s, c := parseHistoryName("loose.txt"); s != "loose" || c != "" {
		t.Errorf("underscore-less parse = %q,%q", s, c)
	}
}

func TestAddPruneAndLoadHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "history")

	// Add four snapshots with a cap of 2 — the two oldest must be pruned.
	stamps := []string{
		"2026-06-13T00:00:00Z", "2026-06-14T00:00:00Z",
		"2026-06-15T00:00:00Z", "2026-06-16T00:00:00Z",
	}
	shas := []string{"aaa1111", "bbb2222", "ccc3333", "ddd4444"}
	for i := range stamps {
		_, pruned, err := AddHistoryEntry(dir, rawBench, stamps[i], shas[i], 2)
		if err != nil {
			t.Fatalf("AddHistoryEntry: %v", err)
		}
		// Pruning only kicks in once the cap is exceeded (after the 3rd add).
		if i < 2 && len(pruned) != 0 {
			t.Errorf("add %d: unexpected prune %v", i, pruned)
		}
		if i >= 2 && len(pruned) != 1 {
			t.Errorf("add %d: pruned %v, want exactly 1", i, pruned)
		}
	}

	entries, err := LoadHistory(dir)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("retained %d entries, want 2 (cap)", len(entries))
	}
	// Oldest-first ordering, and only the two most recent survive.
	if entries[0].Commit != "ccc3333" || entries[1].Commit != "ddd4444" {
		t.Errorf("retained commits = %q,%q want ccc3333,ddd4444", entries[0].Commit, entries[1].Commit)
	}
	if len(entries[1].Samples) == 0 {
		t.Error("loaded entry has no parsed samples")
	}
}

func TestPruneDisabledAndNoOp(t *testing.T) {
	dir := t.TempDir()
	// cap < 1 disables pruning entirely.
	if _, _, err := AddHistoryEntry(dir, rawBench, "2026-06-16T00:00:00Z", "aaa", 0); err != nil {
		t.Fatalf("add: %v", err)
	}
	if _, _, err := AddHistoryEntry(dir, rawBench, "2026-06-17T00:00:00Z", "bbb", 0); err != nil {
		t.Fatalf("add: %v", err)
	}
	pruned, err := pruneHistory(dir, 0)
	if err != nil || pruned != nil {
		t.Errorf("pruneHistory(cap=0) = %v,%v want nil,nil", pruned, err)
	}
	// cap >= count is a no-op.
	if pruned, err := pruneHistory(dir, 10); err != nil || pruned != nil {
		t.Errorf("pruneHistory(cap=10) = %v,%v want nil,nil", pruned, err)
	}
}

func TestLoadHistoryEmptyAndSkips(t *testing.T) {
	// A missing directory yields no entries and no error.
	entries, err := LoadHistory(filepath.Join(t.TempDir(), "does-not-exist"))
	if err != nil || len(entries) != 0 {
		t.Errorf("missing dir: %v entries, err %v", len(entries), err)
	}

	// A file with no conforming benchmarks is skipped rather than failing.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "2026-06-16T00-00-00Z_junk.txt"), []byte("not a bench file\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if entries, err := LoadHistory(dir); err != nil || len(entries) != 0 {
		t.Errorf("junk file: %v entries, err %v", len(entries), err)
	}
}

func TestRunHistoryValidatesInput(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "bench.txt")
	if err := os.WriteFile(good, []byte(rawBench), 0o644); err != nil {
		t.Fatal(err)
	}
	store := filepath.Join(dir, "history")
	if err := runHistory([]string{"-in", good, "-dir", store, "-commit", "abc1234", "-date", "2026-06-16T00:00:00Z"}); err != nil {
		t.Fatalf("runHistory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(store, "2026-06-16T00-00-00Z_abc1234.txt")); err != nil {
		t.Errorf("expected snapshot written: %v", err)
	}

	// Missing -in.
	if err := runHistory([]string{"-dir", store}); err == nil {
		t.Error("expected error when -in is missing")
	}
	// Non-existent input file.
	if err := runHistory([]string{"-in", filepath.Join(dir, "nope.txt")}); err == nil {
		t.Error("expected error for unreadable input")
	}
	// Input with no conforming benchmarks.
	junk := filepath.Join(dir, "junk.txt")
	if err := os.WriteFile(junk, []byte("nothing here\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runHistory([]string{"-in", junk, "-dir", store}); err == nil {
		t.Error("expected error for non-conforming input")
	}
}
