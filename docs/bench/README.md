# Benchmark datasets

This directory holds the per-environment benchmark datasets that feed the
generated [`BENCHMARKS.md`](../../BENCHMARKS.md) report, the
[`bench.svg`](../bench.svg) chart, and the README performance preview. Each file
is a [`benchstat`](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) CSV with a
small `# benchreport-meta:` provenance preamble, written by
[`tools/benchreport`](../../tools/benchreport).

The report surfaces two environments so a controlled baseline is never confused
with noisy shared-runner numbers:

| File | Environment | Tier | Refreshed by |
|------|-------------|------|--------------|
| `reference.csv` | A fixed, controlled machine (the project's reference box) | **primary** — drives the headline table + chart | a maintainer running `make bench-report` locally |
| `ci.csv` | The shared GitHub-hosted runner | **secondary** — indicative only | the `main`-only `bench-report` CI job, on every push to `main` |

> ⚠️ CI numbers come from a shared, noisy runner — trust them only for orders of
> magnitude and relative comparisons. The reference machine is the baseline.

## Refreshing

```bash
# Reference (primary) — run on the controlled machine and commit the result:
make bench-report

# Re-render the combined report/chart/README from the committed datasets only
# (no benchmarking):
make bench-render
```

The CI job refreshes `ci.csv` automatically; it never touches `reference.csv`.
`benchreport render` combines whichever datasets are present, so updating one
environment never disturbs the other. Re-rendering with unchanged data is a
no-op diff.
