package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// runJSON mirrors the object the ci-health-badges workflow emits per run via
// `gh api --jq`. The field names are the jq projection's, not GitHub's raw
// schema, so the workflow and this tool agree on a small stable contract:
//
//	gh api .../runs --jq '.workflow_runs[]
//	  | {id: .id, sha: .head_sha, conclusion: .conclusion, timestamp: .created_at}'
type runJSON struct {
	ID         int64  `json:"id"`
	SHA        string `json:"sha"`
	Conclusion string `json:"conclusion"`
	Timestamp  string `json:"timestamp"`
}

// ParseRuns reads a stream of NDJSON (or otherwise whitespace-separated) run
// objects into records. A run with no conclusion yet (still in progress) or no
// timestamp is skipped rather than failing the whole refresh — the workflow asks
// for completed runs, but a malformed straggler must not wedge the badge job.
// The SHA is shortened to 7 chars to match the store's existing short-SHA rows
// (and the bench history naming) regardless of what width the API returns.
func ParseRuns(r io.Reader) ([]Record, error) {
	dec := json.NewDecoder(r)
	var records []Record
	for {
		var rj runJSON
		err := dec.Decode(&rj)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decoding run JSON: %w", err)
		}
		if rj.Conclusion == "" || rj.Timestamp == "" {
			continue
		}
		ts, err := time.Parse(time.RFC3339, rj.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("run %d: bad timestamp %q: %w", rj.ID, rj.Timestamp, err)
		}
		records = append(records, Record{
			RunID:      rj.ID,
			SHA:        shortSHA(rj.SHA),
			Conclusion: rj.Conclusion,
			Timestamp:  ts.UTC(),
		})
	}
	return records, nil
}

func shortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}
