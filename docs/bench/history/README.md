# Benchmark trend store

Long-term, per-commit benchmark history (issue #51) — the store behind the
**Trend** section and regression check in [`BENCHMARKS.md`](../../../BENCHMARKS.md#trend-recent-main-commits).

## What's here

One file per push to `main`, named `<timestamp>_<sha>.txt`, e.g.:

```
2026-06-16T19-23-20Z_a894dc3.txt
```

Each file is the **raw, multi-sample** `go test -bench -benchmem` output for that
commit (the same `build/bench.txt` the snapshot pipeline produces), with no
preamble. Two consequences follow from that format choice:

- **Significance is recoverable.** Keeping every sample (not just benchstat's
  median) is what lets the regression check apply a real statistical test
  between commits, rather than comparing untrustworthy single deltas.
- **benchstat reads them directly.** No preprocessing is needed for the manual
  recipe below; provenance lives entirely in the filename (timestamp + commit),
  which is also why a plain name sort is chronological.

## How it's maintained

Written only by the main-only `bench-report` CI job (via `make bench-report
BENCH_HISTORY=1`), so every snapshot comes from the **same** shared GitHub-hosted
runner and is therefore commit-to-commit comparable. A maintainer's local
`make bench-report` does **not** archive here (it would mix in a different
machine's numbers).

**Retention:** the newest 100 snapshots are kept; older ones are pruned on each
push to bound repo growth (`BENCH_HISTORY_CAP` in the `Makefile`).

## Compare any two commits manually

```sh
benchstat docs/bench/history/<older>_<sha>.txt docs/bench/history/<newer>_<sha>.txt
```

benchstat prints each benchmark's delta and `p`-value; trust a change only when
`p` is low (a non-`~` delta). The runners are shared and noisy, so the automated
check (and you) should weigh significance, not raw numbers.
