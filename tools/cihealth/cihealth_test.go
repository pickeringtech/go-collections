package main

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	ts, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("bad test timestamp %q: %v", s, err)
	}
	return ts
}

// rec builds a Record at `daysAgo` days before the fixed `now` used in tests.
func rec(id int64, conclusion string, ts time.Time) Record {
	return Record{RunID: id, SHA: "abcdef0", Conclusion: conclusion, Timestamp: ts}
}

func TestCount_WindowAndConclusionFilter(t *testing.T) {
	now := mustTime(t, "2026-06-18T00:00:00Z")
	ago := func(d int) time.Time { return now.Add(-time.Duration(d) * 24 * time.Hour) }

	records := []Record{
		rec(1, "success", ago(1)),         // in all windows, healthy
		rec(2, "failure", ago(2)),         // in all windows, counted-not-healthy
		rec(3, "timed_out", ago(3)),       // counted
		rec(4, "startup_failure", ago(4)), // counted
		rec(5, "cancelled", ago(1)),       // EXCLUDED from denominator
		rec(6, "skipped", ago(1)),         // EXCLUDED
		rec(7, "neutral", ago(1)),         // EXCLUDED (unfamiliar/neutral)
		rec(8, "success", ago(20)),        // only in 30d/90d
		rec(9, "success", ago(60)),        // only in 90d
		rec(10, "failure", ago(200)),      // outside every window
	}

	tests := []struct {
		days        int
		wantHealthy int
		wantCounted int
	}{
		{7, 1, 4},  // recs 1-4 counted (5,6,7 excluded), 1 healthy
		{30, 2, 5}, // + rec 8 (success)
		{90, 3, 6}, // + rec 9 (success); rec 10 too old
	}
	for _, tc := range tests {
		got := Count(records, now, tc.days)
		if got.Healthy != tc.wantHealthy || got.Counted != tc.wantCounted {
			t.Errorf("Count(%dd) = {healthy:%d counted:%d}, want {healthy:%d counted:%d}",
				tc.days, got.Healthy, got.Counted, tc.wantHealthy, tc.wantCounted)
		}
	}
}

func TestRender_MessageColorAndGreying(t *testing.T) {
	tests := []struct {
		name        string
		tally       Tally
		minSample   int
		wantMessage string
		wantColor   string
	}{
		{"empty", Tally{0, 0}, 3, "no runs", "lightgrey"},
		{"perfect", Tally{30, 30}, 3, "100% (30/30)", "brightgreen"},
		{"brightgreen boundary", Tally{19, 20}, 3, "95% (19/20)", "brightgreen"},
		{"yellow", Tally{18, 20}, 3, "90% (18/20)", "yellow"},
		{"yellow boundary", Tally{17, 20}, 3, "85% (17/20)", "yellow"},
		{"orange", Tally{16, 20}, 3, "80% (16/20)", "orange"},
		{"orange boundary", Tally{14, 20}, 3, "70% (14/20)", "orange"},
		{"red", Tally{13, 20}, 3, "65% (13/20)", "red"},
		{"low sample greyed", Tally{1, 2}, 3, "50% (1/2)", "lightgrey"},
		{"at min sample colored", Tally{3, 3}, 3, "100% (3/3)", "brightgreen"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := Render("main health (7d)", tc.tally, tc.minSample)
			if b.SchemaVersion != 1 {
				t.Errorf("SchemaVersion = %d, want 1", b.SchemaVersion)
			}
			if b.Label != "main health (7d)" {
				t.Errorf("Label = %q, want %q", b.Label, "main health (7d)")
			}
			if b.Message != tc.wantMessage {
				t.Errorf("Message = %q, want %q", b.Message, tc.wantMessage)
			}
			if b.Color != tc.wantColor {
				t.Errorf("Color = %q, want %q", b.Color, tc.wantColor)
			}
		})
	}
}

func TestMerge_DedupAndFetchWins(t *testing.T) {
	now := mustTime(t, "2026-06-18T00:00:00Z")
	stored := []Record{
		rec(1, "success", now.Add(-48*time.Hour)),
		rec(2, "failure", now.Add(-24*time.Hour)), // will be re-run to success
	}
	fetched := []Record{
		rec(2, "success", now.Add(-24*time.Hour)), // re-run: fetch wins
		rec(3, "success", now.Add(-1*time.Hour)),  // brand new
	}
	got := Merge(stored, fetched)
	if len(got) != 3 {
		t.Fatalf("Merge len = %d, want 3 (deduped by run id)", len(got))
	}
	// Sorted oldest first; run 2 must carry the fetched conclusion.
	byID := map[int64]Record{}
	for _, r := range got {
		byID[r.RunID] = r
	}
	if byID[2].Conclusion != "success" {
		t.Errorf("run 2 conclusion = %q, want %q (fetch should win)", byID[2].Conclusion, "success")
	}
	// Confirm ascending order.
	for i := 1; i < len(got); i++ {
		if got[i].Timestamp.Before(got[i-1].Timestamp) {
			t.Errorf("Merge not sorted ascending at %d", i)
		}
	}
}

func TestPrune_DropsAgedRows(t *testing.T) {
	now := mustTime(t, "2026-06-18T00:00:00Z")
	records := []Record{
		rec(1, "success", now.Add(-99*24*time.Hour)),  // kept (< 100d)
		rec(2, "success", now.Add(-101*24*time.Hour)), // dropped (> 100d)
		rec(3, "success", now.Add(-1*24*time.Hour)),   // kept
	}
	got := Prune(records, now, 100)
	if len(got) != 2 {
		t.Fatalf("Prune len = %d, want 2", len(got))
	}
	for _, r := range got {
		if r.RunID == 2 {
			t.Errorf("run 2 should have been pruned")
		}
	}

	// retention < 1 disables pruning.
	if len(Prune(records, now, 0)) != 3 {
		t.Errorf("Prune with retention 0 should keep all rows")
	}
}

func TestParseRuns_NDJSON(t *testing.T) {
	in := strings.Join([]string{
		`{"id":101,"sha":"a894dc3aaaaaaaa","conclusion":"success","timestamp":"2026-06-16T19:23:20Z"}`,
		`{"id":102,"sha":"b133a28","conclusion":"failure","timestamp":"2026-06-17T10:00:00Z"}`,
		`{"id":103,"sha":"ccc","conclusion":"","timestamp":"2026-06-17T11:00:00Z"}`, // in progress: skipped
		``,
	}, "\n")
	got, err := ParseRuns(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ParseRuns: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("ParseRuns len = %d, want 2 (in-progress skipped)", len(got))
	}
	if got[0].SHA != "a894dc3" {
		t.Errorf("SHA not shortened: %q", got[0].SHA)
	}
	if got[0].RunID != 101 || got[0].Conclusion != "success" {
		t.Errorf("unexpected first record: %+v", got[0])
	}
}

func TestStoreRoundTrip(t *testing.T) {
	now := mustTime(t, "2026-06-18T00:00:00Z")
	path := filepath.Join(t.TempDir(), "history.csv")
	want := []Record{
		rec(2, "failure", now.Add(-24*time.Hour)),
		rec(1, "success", now.Add(-48*time.Hour)),
	}
	if err := SaveStore(path, want); err != nil {
		t.Fatalf("SaveStore: %v", err)
	}
	got, err := LoadStore(path)
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	// SaveStore sorts ascending: run 1 (older) then run 2.
	if len(got) != 2 || got[0].RunID != 1 || got[1].RunID != 2 {
		t.Fatalf("round-trip order wrong: %+v", got)
	}
	wantTS := now.Add(-48 * time.Hour) // run 1's timestamp
	if !got[0].Timestamp.Equal(wantTS) {
		t.Errorf("timestamp not preserved: got %v want %v", got[0].Timestamp, wantTS)
	}
}

func TestLoadStore_MissingFileIsEmpty(t *testing.T) {
	got, err := LoadStore(filepath.Join(t.TempDir(), "does-not-exist.csv"))
	if err != nil {
		t.Fatalf("LoadStore on missing file should not error: %v", err)
	}
	if got != nil {
		t.Errorf("LoadStore on missing file = %v, want nil", got)
	}
}

// TestRunEndToEnd exercises the whole pipeline: NDJSON in, store + badges out,
// then a second idempotent run that must change nothing.
func TestRunEndToEnd(t *testing.T) {
	dir := t.TempDir()
	store := filepath.Join(dir, "history.csv")
	runs := filepath.Join(dir, "runs.ndjson")

	ndjson := strings.Join([]string{
		`{"id":1,"sha":"aaaaaaa","conclusion":"success","timestamp":"2026-06-17T00:00:00Z"}`,
		`{"id":2,"sha":"bbbbbbb","conclusion":"failure","timestamp":"2026-06-17T01:00:00Z"}`,
		`{"id":3,"sha":"ccccccc","conclusion":"cancelled","timestamp":"2026-06-17T02:00:00Z"}`,
	}, "\n")
	if err := writeFileForTest(runs, ndjson); err != nil {
		t.Fatal(err)
	}

	args := []string{
		"-runs", runs,
		"-store", store,
		"-out", dir,
		"-now", "2026-06-18T00:00:00Z",
		"-min-sample", "1",
	}
	if err := run(args); err != nil {
		t.Fatalf("run: %v", err)
	}

	// 1 success of 2 counted (cancelled excluded) → 50%.
	badge := readBadgeForTest(t, filepath.Join(dir, "badge-7d.json"))
	if badge.Message != "50% (1/2)" {
		t.Errorf("7d message = %q, want %q", badge.Message, "50% (1/2)")
	}
	if badge.SchemaVersion != 1 {
		t.Errorf("schemaVersion = %d, want 1", badge.SchemaVersion)
	}

	// Idempotency: a second run with no new runs leaves the store byte-identical.
	before := readFileForTest(t, store)
	if err := writeFileForTest(runs, ""); err != nil {
		t.Fatal(err)
	}
	if err := run(args); err != nil {
		t.Fatalf("second run: %v", err)
	}
	after := readFileForTest(t, store)
	if before != after {
		t.Errorf("store changed on a no-op refresh:\nbefore=%q\nafter=%q", before, after)
	}

	// All three windows produced a file.
	for _, w := range windows {
		if _, err := readFileForTestErr(filepath.Join(dir, w.File)); err != nil {
			t.Errorf("missing badge file %s: %v", w.File, err)
		}
	}

	// Sanity: the stored records survive a reload with the cancelled run retained.
	got, err := LoadStore(store)
	if err != nil {
		t.Fatal(err)
	}
	wantIDs := []int64{1, 2, 3}
	var gotIDs []int64
	for _, r := range got {
		gotIDs = append(gotIDs, r.RunID)
	}
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Errorf("stored run ids = %v, want %v", gotIDs, wantIDs)
	}
}
