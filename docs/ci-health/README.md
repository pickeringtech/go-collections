# Main-branch CI health badges

The data behind the **main health (7d / 30d / 90d)** badges in the project
[`README`](../../README.md) (issue #209) — the percentage of healthy `CI` runs on
`main` over rolling windows, so "is `main` reliably green?" is visible at a
glance and a regression in build health is trended rather than anecdotal.

## What's here

| File | Purpose |
| --- | --- |
| `history.csv` | the persisted tally — one row per counted run (the source of truth) |
| `badge-7d.json` / `badge-30d.json` / `badge-90d.json` | [shields.io endpoint](https://shields.io/badges/endpoint-badge) JSON the README points at |

`history.csv` is `run_id,sha,conclusion,timestamp`, oldest first, one row per
completed `CI` run on `main`:

```
run_id,sha,conclusion,timestamp
27757102375,101b1ae,success,2026-06-18T11:44:19Z
```

Rows are keyed by `run_id`, so a re-run that changes a run's conclusion updates
in place rather than double-counting.

## Definitions

- **A "build"** is a run of the `CI` workflow (`.github/workflows/ci.yml`) whose
  `head_branch == main` and `event == push` — the post-merge runs the merge
  queue lands on `main`. The merge-queue *validation* runs sit on the queue ref
  (`gh-readonly-queue/...`), not `main`, so they're excluded; and the
  `bench-report` refresh commits carry `[skip ci]`, so they never start a run.
- **Healthy** = `conclusion == success`. The denominator counts
  `success + failure + timed_out + startup_failure`; `cancelled` / `skipped` /
  `action_required` / `stale` / `neutral` are **excluded** — they're queue churn
  or human intervention, not signal about main's code health.
- **Percentage** = healthy ÷ counted over each window. The badge message carries
  `n/d` so a small denominator isn't mistaken for a strong signal, and a window
  with too few counted runs is greyed rather than colour-coded.
- **Colour** (tune once real numbers settle): `≥95%` brightgreen, `≥85%` yellow,
  `≥70%` orange, else red.

## Why a persisted tally (not a live query)

GitHub's Actions run history is subject to a retention limit (default 90 days),
so the **90-day window sits right on the boundary** — a live
`actions/runs?branch=main` query undercounts the quarter the moment the oldest
runs age out. Appending every observed run to this committed store and computing
all three windows from it keeps the quarter badge honest after retention prunes
the API. It mirrors the [`docs/bench/history/`](../bench/history/README.md)
trend store (issue #51), but is row-based because each run is one datum.

**Retention:** rows older than ~100 days (a margin past the widest window) are
pruned on each refresh to bound growth (`-retention-days` in `tools/cihealth`).

## How it's maintained

The scheduled [`ci-health-badges`](../../.github/workflows/ci-health-badges.yml)
workflow recomputes everything **daily** (and on demand). It fetches the run
history via `gh api`, runs `make ci-health-report`, and commits any change here
to `main` via the `bench-report` write deploy key (a merge-queue ruleset bypass
actor, issue #199) with `[skip ci]` — so the refresh never trips the queue and
its own commit never becomes a counted build. The push is best-effort and never
gates.

Refresh locally the same way the workflow does:

```sh
gh api --paginate \
  '/repos/pickeringtech/go-collections/actions/workflows/ci.yml/runs?branch=main&event=push&status=completed&per_page=100' \
  --jq '.workflow_runs[] | {id:.id, sha:.head_sha, conclusion:.conclusion, timestamp:.created_at}' \
  > build/ci-runs.ndjson
make ci-health-report CI_RUNS=build/ci-runs.ndjson
```

The generator (`tools/cihealth`) is a pure function of its inputs, so re-running
with unchanged data produces byte-identical output.
