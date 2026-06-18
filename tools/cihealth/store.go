package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// The persisted tally store (issue #209) is a single CSV, one row per completed
// `CI` run on `main`, that outlives GitHub's 90-day Actions retention. The live
// `actions/runs` API undercounts the 90-day window the moment the oldest runs
// age out of retention; appending each run we observe to a committed store keeps
// the quarter badge honest. It mirrors the docs/bench/history/ trend store
// (issue #51) in spirit — a small committed dataset the badge job refreshes —
// but is row-based rather than file-per-commit because each run is one datum.
//
//	docs/ci-health/history.csv
//	run_id,sha,conclusion,timestamp
//	123456789,a894dc3,success,2026-06-16T19:23:20Z
//
// Rows are keyed by run_id (the stable GitHub identifier), so a re-run that
// changes a run's conclusion overwrites in place rather than double-counting.

// Record is one completed CI run on main: its GitHub run id, head commit, final
// conclusion, and the moment it started (the window-membership timestamp).
type Record struct {
	RunID      int64
	SHA        string
	Conclusion string
	Timestamp  time.Time
}

// LoadStore reads the tally CSV into records. A missing file is not an error —
// it yields no records, which is the correct bootstrap state before the first
// run populates the store.
func LoadStore(path string) ([]Record, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 4
	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading store %s: %w", path, err)
	}

	var records []Record
	for i, row := range rows {
		// Skip the header row if present (first row, non-numeric run id).
		if i == 0 && row[0] == "run_id" {
			continue
		}
		id, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("store %s line %d: bad run_id %q: %w", path, i+1, row[0], err)
		}
		ts, err := time.Parse(time.RFC3339, row[3])
		if err != nil {
			return nil, fmt.Errorf("store %s line %d: bad timestamp %q: %w", path, i+1, row[3], err)
		}
		records = append(records, Record{RunID: id, SHA: row[1], Conclusion: row[2], Timestamp: ts})
	}
	return records, nil
}

// SaveStore writes records to the tally CSV, oldest first, creating any missing
// parent directories. Output is deterministic (sorted by timestamp then run id,
// timestamps normalised to UTC RFC-3339) so an unchanged dataset re-serialises
// byte-identically and the badge job's "commit if changed" stays a no-op.
func SaveStore(path string, records []Record) error {
	sortRecords(records)
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return writeStore(f, records)
}

func writeStore(w io.Writer, records []Record) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"run_id", "sha", "conclusion", "timestamp"}); err != nil {
		return err
	}
	for _, r := range records {
		row := []string{
			strconv.FormatInt(r.RunID, 10),
			r.SHA,
			r.Conclusion,
			r.Timestamp.UTC().Format(time.RFC3339),
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// Merge unions the existing store with freshly-fetched runs, keyed by run id so
// a run observed in both wins from the fetch (a re-run's updated conclusion
// replaces the stored one) and is never counted twice. Runs that have aged out
// of the live API but still live in the store are retained — that retention is
// the whole reason the store exists.
func Merge(stored, fetched []Record) []Record {
	byID := make(map[int64]Record, len(stored)+len(fetched))
	for _, r := range stored {
		byID[r.RunID] = r
	}
	for _, r := range fetched {
		byID[r.RunID] = r
	}
	out := make([]Record, 0, len(byID))
	for _, r := range byID {
		out = append(out, r)
	}
	sortRecords(out)
	return out
}

// Prune drops records older than retentionDays before now, bounding store growth
// to a little past the widest (90-day) window. A retentionDays below 1 disables
// pruning (keep everything).
func Prune(records []Record, now time.Time, retentionDays int) []Record {
	if retentionDays < 1 {
		return records
	}
	cutoff := now.Add(-time.Duration(retentionDays) * 24 * time.Hour)
	kept := records[:0]
	for _, r := range records {
		if r.Timestamp.Before(cutoff) {
			continue
		}
		kept = append(kept, r)
	}
	return kept
}

func sortRecords(records []Record) {
	sort.Slice(records, func(i, j int) bool {
		if records[i].Timestamp.Equal(records[j].Timestamp) {
			return records[i].RunID < records[j].RunID
		}
		return records[i].Timestamp.Before(records[j].Timestamp)
	})
}
