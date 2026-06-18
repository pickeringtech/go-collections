package main

import (
	"fmt"
	"math"
	"time"
)

// Health definition (issue #209), kept deliberately explicit so the badge means
// the same thing as CI changes underneath it:
//
//   - counted (the denominator) is every run that is signal about main's code
//     health: success + failure + timed_out + startup_failure.
//   - healthy (the numerator) is success alone.
//   - cancelled / skipped / action_required / stale / neutral are EXCLUDED — they
//     reflect queue churn or human intervention, not whether main built green.
//
// A run with any other (or empty/in-progress) conclusion is excluded too, so an
// unfamiliar future conclusion is ignored rather than silently miscounted.
var countedConclusions = map[string]bool{
	"success":         true,
	"failure":         true,
	"timed_out":       true,
	"startup_failure": true,
}

// Window is one rolling-window badge: how many days back it looks, the file it
// renders to, and the badge label.
type Window struct {
	Days  int
	File  string
	Label string
}

// windows are the three badges the README shows. Ordered widest-last to match
// the README badge row (7d, 30d, 90d).
var windows = []Window{
	{Days: 7, File: "badge-7d.json", Label: "main health (7d)"},
	{Days: 30, File: "badge-30d.json", Label: "main health (30d)"},
	{Days: 90, File: "badge-90d.json", Label: "main health (90d)"},
}

// Tally is the outcome of counting one window: healthy successes over the runs
// that count as health signal.
type Tally struct {
	Healthy int
	Counted int
}

// Count tallies the records that fall within the last `days` before now and
// carry a counted conclusion.
func Count(records []Record, now time.Time, days int) Tally {
	cutoff := now.Add(-time.Duration(days) * 24 * time.Hour)
	var t Tally
	for _, r := range records {
		if r.Timestamp.Before(cutoff) {
			continue
		}
		if !countedConclusions[r.Conclusion] {
			continue
		}
		t.Counted++
		if r.Conclusion == "success" {
			t.Healthy++
		}
	}
	return t
}

// Badge is the shields.io endpoint schema
// (https://shields.io/badges/endpoint-badge). shields can't compute "% of last N
// runs" itself, so we emit the numbers and let it render them.
type Badge struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
	Color         string `json:"color"`
}

// Render turns a window's tally into its badge. The message always carries `n/d`
// so a small denominator can't be mistaken for a strong signal, and a window
// with fewer than minSample counted runs is greyed (its percentage is too noisy
// to colour-code — a single failure in a near-empty week shouldn't read as a
// crisis). An empty window renders an explicit "no runs" rather than 0%.
func Render(label string, t Tally, minSample int) Badge {
	b := Badge{SchemaVersion: 1, Label: label}
	if t.Counted == 0 {
		b.Message = "no runs"
		b.Color = "lightgrey"
		return b
	}
	ratio := float64(t.Healthy) / float64(t.Counted)
	pct := int(math.Round(ratio * 100))
	b.Message = fmt.Sprintf("%d%% (%d/%d)", pct, t.Healthy, t.Counted)
	if t.Counted < minSample {
		b.Color = "lightgrey"
	} else {
		b.Color = colorFor(ratio * 100)
	}
	return b
}

// colorFor maps a health percentage to a shields colour. Thresholds per issue
// #209 — tune once real numbers exist.
func colorFor(pct float64) string {
	switch {
	case pct >= 95:
		return "brightgreen"
	case pct >= 85:
		return "yellow"
	case pct >= 70:
		return "orange"
	default:
		return "red"
	}
}
