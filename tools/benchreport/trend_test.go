package main

import (
	"strings"
	"testing"
)

var hashGet = sampleKey{"dicts", "Hash", "Get", 1000}

// mkEntry builds a history entry whose listed cells each carry eight identical
// ns/op samples — enough for the Mann–Whitney test to find a constant shift
// highly significant, which keeps the regression tests deterministic.
func mkEntry(commit, stamp string, cells map[sampleKey]float64) HistoryEntry {
	s := map[sampleKey][]float64{}
	for k, v := range cells {
		s[k] = []float64{v, v, v, v, v, v, v, v}
	}
	return HistoryEntry{Commit: commit, Stamp: stamp, Samples: s}
}

func TestRenderTrendSectionEmpty(t *testing.T) {
	if got := RenderTrendSection(nil); got != "" {
		t.Errorf("empty history should render nothing, got %q", got)
	}
}

func TestRenderTrendSectionContent(t *testing.T) {
	entries := []HistoryEntry{
		mkEntry("aaa1111", "2026-06-14T00-00-00Z", map[sampleKey]float64{hashGet: 4.4}),
		mkEntry("bbb2222", "2026-06-15T00-00-00Z", map[sampleKey]float64{hashGet: 4.5}),
		mkEntry("ccc3333", "2026-06-16T00-00-00Z", map[sampleKey]float64{hashGet: 4.6}),
	}
	out := RenderTrendSection(entries)
	for _, want := range []string{
		"## Trend (recent `main` commits)",
		"| Operation |",
		"Dict — Hash.Get",
		"ccc3333", // newest first → leads the table header
		"Commits, newest first:",
		"`ccc3333` (2026-06-16)",
		"### Regression check (report-only)",
		"### Compare any two commits manually",
		"benchstat docs/bench/history/",
		"—", // a headline with no samples renders as an em dash
	} {
		if !strings.Contains(out, want) {
			t.Errorf("trend section missing %q\n%s", want, out)
		}
	}
	// Newest commit must appear before the oldest in the header row.
	if strings.Index(out, "ccc3333") > strings.Index(out, "aaa1111") {
		t.Error("trend table is not newest-first")
	}
}

func TestRenderTrendSectionColumnCap(t *testing.T) {
	var entries []HistoryEntry
	for i := 0; i < trendMaxColumns+5; i++ {
		entries = append(entries, mkEntry(
			string(rune('a'+i))+"000000",
			"2026-06-16T00-00-00Z",
			map[sampleKey]float64{hashGet: 4.0},
		))
	}
	out := RenderTrendSection(entries)
	// The narrative names how many of how many are shown.
	if !strings.Contains(out, "last 12 of 17 retained") {
		t.Errorf("expected column cap narrative, got:\n%s", out)
	}
}

func TestRegressionInsufficientHistory(t *testing.T) {
	entries := []HistoryEntry{
		mkEntry("aaa", "2026-06-15T00-00-00Z", map[sampleKey]float64{hashGet: 4.0}),
		mkEntry("bbb", "2026-06-16T00-00-00Z", map[sampleKey]float64{hashGet: 9.0}),
	}
	out := RenderRegressionAlert(entries)
	if !strings.Contains(out, "Need at least 3 retained commits") {
		t.Errorf("expected insufficient-history message, got:\n%s", out)
	}
}

func TestRegressionNoneFlagged(t *testing.T) {
	entries := []HistoryEntry{
		mkEntry("aaa", "d1", map[sampleKey]float64{hashGet: 4.0}),
		mkEntry("bbb", "d2", map[sampleKey]float64{hashGet: 4.0}),
		mkEntry("ccc", "d3", map[sampleKey]float64{hashGet: 4.0}),
	}
	out := RenderRegressionAlert(entries)
	if !strings.Contains(out, "✅ No headline operation regressed") {
		t.Errorf("expected clean verdict, got:\n%s", out)
	}
}

func TestRegressionFlaggedAcrossTwoComparisons(t *testing.T) {
	// Hash.Get gets steadily, significantly slower across BOTH comparisons.
	entries := []HistoryEntry{
		mkEntry("aaa1111", "d1", map[sampleKey]float64{hashGet: 4.0}),
		mkEntry("bbb2222", "d2", map[sampleKey]float64{hashGet: 8.0}),
		mkEntry("ccc3333", "d3", map[sampleKey]float64{hashGet: 16.0}),
	}
	out := RenderRegressionAlert(entries)
	if !strings.Contains(out, "⚠️") || !strings.Contains(out, "Dict — Hash.Get") {
		t.Errorf("expected Hash.Get flagged, got:\n%s", out)
	}
	if !strings.Contains(out, "+100%") { // 8 → 16 in the latest comparison
		t.Errorf("expected +100%% magnitude, got:\n%s", out)
	}
	if !strings.Contains(out, "`aaa1111` → `bbb2222` → `ccc3333`") {
		t.Errorf("expected the three-commit lineage, got:\n%s", out)
	}
}

func TestRegressionRequiresTwoConsecutive(t *testing.T) {
	// Slower only in the latest comparison (stable across the prior one) →
	// not flagged, because one noisy run must not trip the alarm.
	entries := []HistoryEntry{
		mkEntry("aaa", "d1", map[sampleKey]float64{hashGet: 4.0}),
		mkEntry("bbb", "d2", map[sampleKey]float64{hashGet: 4.0}),
		mkEntry("ccc", "d3", map[sampleKey]float64{hashGet: 16.0}),
	}
	out := RenderRegressionAlert(entries)
	if !strings.Contains(out, "✅ No headline operation regressed") {
		t.Errorf("single-comparison spike should not flag, got:\n%s", out)
	}
}

func TestRegressionPriorOnlyNotFlagged(t *testing.T) {
	// Regressed in the prior comparison but recovered in the latest one → not
	// flagged (the slowdown is no longer current).
	entries := []HistoryEntry{
		mkEntry("aaa", "d1", map[sampleKey]float64{hashGet: 4.0}),
		mkEntry("bbb", "d2", map[sampleKey]float64{hashGet: 16.0}),
		mkEntry("ccc", "d3", map[sampleKey]float64{hashGet: 16.0}),
	}
	out := RenderRegressionAlert(entries)
	if !strings.Contains(out, "✅ No headline operation regressed") {
		t.Errorf("recovered regression should not flag, got:\n%s", out)
	}
}

func TestRegressedWithMediansGates(t *testing.T) {
	old4 := mkEntry("o", "d", map[sampleKey]float64{hashGet: 4.0})
	new5 := mkEntry("n", "d", map[sampleKey]float64{hashGet: 5.0})    // +25%, significant
	newTiny := mkEntry("n", "d", map[sampleKey]float64{hashGet: 4.2}) // +5%, below threshold
	empty := HistoryEntry{Samples: map[sampleKey][]float64{}}

	if ok, _, _ := regressedWithMedians(old4, new5, hashGet); !ok {
		t.Error("expected a significant +25% slowdown to regress")
	}
	if ok, _, _ := regressedWithMedians(old4, newTiny, hashGet); ok {
		t.Error("a +5% slowdown is below the 10% threshold")
	}
	if ok, _, _ := regressedWithMedians(old4, empty, hashGet); ok {
		t.Error("missing samples on one side must not regress")
	}
	if ok, _, _ := regressedWithMedians(empty, new5, hashGet); ok {
		t.Error("missing samples on the old side must not regress")
	}

	// Median > threshold but statistically insignificant (wide overlap, n=2) is
	// NOT a regression — significance, not raw delta, is the gate.
	noisyOld := HistoryEntry{Samples: map[sampleKey][]float64{hashGet: {10, 100}}}
	noisyNew := HistoryEntry{Samples: map[sampleKey][]float64{hashGet: {12, 120}}}
	if ok, _, _ := regressedWithMedians(noisyOld, noisyNew, hashGet); ok {
		t.Error("a +20% but insignificant change must not regress")
	}

	// A zero (or negative) old median can't yield a percentage — guarded out.
	zeroOld := HistoryEntry{Samples: map[sampleKey][]float64{hashGet: {0, 0}}}
	if ok, _, _ := regressedWithMedians(zeroOld, new5, hashGet); ok {
		t.Error("zero old median must not regress")
	}
}
