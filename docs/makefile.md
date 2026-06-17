# Makefile reference

The `Makefile` is the set of **developer entry points** for this repository:
the test suite, a local mirror of the blocking CI gates, and the benchmark
report pipeline. Its guiding principle is **lockstep with CI** — `make ci` runs
the same checks the PR `CI Gate` aggregates, and the main-only benchmark CI job
is a thin wrapper over `make bench-report`, so local and CI behaviour can't
drift.

This document explains every target and every variable. For the contributor
workflow around them, see [CONTRIBUTING.md](../CONTRIBUTING.md); for the
benchmark results themselves, see [BENCHMARKS.md](../BENCHMARKS.md).

> `make help` prints the same target list (sourced from the `##` comments in the
> Makefile), so the in-terminal summary and this document stay in agreement.

## Quick start

```bash
make help     # list every target with a one-line purpose
make test     # full test suite — root module + every nested module
make ci        # every blocking CI gate at once — green here predicts a green PR
make bench     # run the benchmarks once and print the numbers
```

## Targets at a glance

| Target | Purpose | Mirrors CI job |
|--------|---------|----------------|
| `help` | List available targets with their purpose | — |
| `test` | Full suite: root module + every nested module | (composite) |
| `test-root` | Root library suite with `-race -shuffle=on` | part of `test` |
| `test-nested` | Every nested module's suite (`examples/`, `tools/benchreport`, …) | `examples-e2e`, `benchreport-test` |
| `ci` | Run every blocking gate locally in one shot | `CI Gate` (all blocking jobs) |
| `hygiene` | Compile, `go.mod` tidy, zero-dependency, module integrity | `build` |
| `cover` | Root suite with `-race -shuffle`, then enforce the coverage floor | `test` |
| `lint` | `gofmt` check, `go vet`, golangci-lint (pinned) | `lint` |
| `security` | `govulncheck` + `gosec` (pinned) | `security` |
| `cross-arch` | 386/arm64/s390x build+vet, plus run the 386 tests | `cross-arch` |
| `fuzz` | Count-based smoke run of every `FuzzXxx` target | `fuzz` |
| `bench` | Run the collections benchmarks once (no report) | — (report-only `benchmarks`) |
| `bench-report` | Benchmark this environment, capture its dataset, regenerate the report | `bench-report` (main-only) |
| `bench-render` | Re-render `BENCHMARKS.md`, `docs/bench.svg`, README preview from committed datasets | part of `bench-report` |

## Testing

### `make test`

The full suite a contributor should run before pushing: it depends on both
`test-root` and `test-nested`, so it covers the root library **and** every
nested Go module.

### `make test-root`

The root library only, with the race detector and randomised test order:

```bash
go test -race -shuffle=on ./...
```

### `make test-nested`

Runs the suite of every **nested** Go module — `examples/`,
`tools/benchreport/`, and any future ones. These are separate modules that the
root `go test ./...` never descends into (issue #79), so they need their own
pass. The module list is **discovered dynamically** (any nested `go.mod`), so
new modules are picked up automatically. It uses `-shuffle=on` without `-race`,
mirroring CI's examples job — those tests shell out via `go run`, so the race
detector can't see into the child process anyway. A failure in any module aborts
with a non-zero status.

## Local CI mirror

`make ci` runs the **same blocking gates** the PR `CI Gate` aggregates
(`.github/workflows/ci.yml` → the `ci-gate` job's `needs:` list), in
cheap-/common-failure-first order so a typical mistake (formatting, a failing
test) aborts before the slow cross-arch builds. A green `make ci` predicts a
green PR. Each gate is also a standalone target so you can reproduce a single
failing job.

| Target | What it does | CI job |
|--------|--------------|--------|
| `hygiene` | `go build ./...`; verifies `go.mod`/`go.sum` are tidy; enforces **zero runtime dependencies**; `go mod verify` | `build` |
| `cover` | Root suite with `-race -shuffle` + coverage profile, then fails if total coverage is below the floor | `test` |
| `lint` | Fails on non-`gofmt`-clean files; `go vet ./...`; golangci-lint at the pinned version | `lint` |
| `security` | `govulncheck` (known-vuln scan) + `gosec` (security lint), both pinned | `security` |
| `cross-arch` | Cross-compiles + vets for 386/arm64/s390x; runs the 386 tests on the amd64 host | `cross-arch` |
| `fuzz` | Runs every `FuzzXxx` target for a fixed iteration count (`-fuzztime=2000x`) | `fuzz` |

The report-only CI jobs (`test-tip`, `benchmarks`, `api-compat`) are
deliberately **not** mirrored — they never gate a merge.

> The lint/security tools run via `go run <module>@<version>` rather than
> expecting a binary on `PATH`: that pins the exact version with zero setup and
> never touches this repo's dependency-free `go.mod`. The first run builds the
> tool, then caches it.

## Benchmarks

### `make bench`

Runs the collections benchmarks once and prints the results — no report, no
files written. The quickest way to see numbers locally.

### `make bench-report`

The full pipeline: run the benchmarks → summarise with `benchstat` → capture the
result as a committed dataset (`docs/bench/<env>.csv`) → re-render the report.
This is what the main-only CI job wraps.

### `make bench-render`

Re-renders `BENCHMARKS.md`, `docs/bench.svg`, and the README preview region from
the **already-committed** datasets under `docs/bench/` (plus the trend store).
Cheap; runs no benchmarks. Useful after editing the report generator.

### The two benchmark environments

The report surfaces two datasets, each refreshed independently and committed
under `docs/bench/`:

- **`reference` (primary)** — a fixed, controlled machine; the trustworthy
  baseline that drives the README headline table and chart. Refreshed by a
  maintainer running `make bench-report` locally.
- **`ci` (secondary)** — the shared, noisy GitHub-hosted runner; indicative
  only. Refreshed by the main-only CI job, which overrides the bench variables.

### Refreshing the reference baseline (Framework Desktop)

The Makefile **defaults to the reference environment**, so on the reference
machine the refresh is simply:

```bash
make bench-report
```

With no overrides this benchmarks `./collections/...`, captures the result to
`docs/bench/reference.csv` labelled for the Framework Desktop, and regenerates
the report. It does **not** archive a trend snapshot (history archiving is
opt-in via `BENCH_HISTORY`, reserved for the consistent CI environment so the
trend store isn't polluted by ad-hoc local runs).

If you refresh the reference from a **different** machine, override the labels so
the report's provenance stays accurate:

```bash
make bench-report \
  BENCH_LABEL="Reference — <machine>" \
  BENCH_MACHINE="<cpu / memory / OS>"
```

For a fast local preview (not a real refresh), shrink the sampling:

```bash
make bench-report BENCH_TIME=10ms BENCH_COUNT=2
```

## Variables

Override any of these on the command line (`make <target> VAR=value`).

| Variable | Default | Purpose |
|----------|---------|---------|
| `BENCH_TIME` | `50ms` | Per-sample benchmark duration (`-benchtime`); mirrors CI |
| `BENCH_COUNT` | `8` | benchstat samples per benchmark (`-count`); mirrors CI |
| `BENCH_PKGS` | `./collections/...` | Which packages to benchmark (the standardized matrix) |
| `BENCH_ENV` | `reference` | Environment captured; also the output filename (`docs/bench/<env>.csv`) |
| `BENCH_TIER` | `primary` | Dataset tier (`primary` reference vs `secondary` ci) |
| `BENCH_LABEL` | `Reference — Framework Desktop` | Human label in the report |
| `BENCH_MACHINE` | `Framework Desktop · AMD Ryzen AI MAX+ 395 · …` | Machine provenance string |
| `BENCHSTAT` | `benchstat` (from `PATH`) | benchstat binary; CI installs a pinned version first |
| `BENCH_HISTORY` | empty | If non-empty, archive a trend snapshot (CI-only; opt-in) |
| `BENCH_HISTORY_CAP` | `100` | Max trend snapshots kept (oldest pruned) |
| `GOLANGCI_VERSION` | `v2.1.6` | golangci-lint version — **must track** the CI pin (#88) |
| `GOVULNCHECK_VERSION` | `v1.3.0` | govulncheck version — must track the CI pin |
| `GOSEC_VERSION` | `v2.27.1` | gosec version — must track the CI pin |
| `COVERAGE_MIN` | `100` | Coverage floor (%) enforced by `cover` |

> **Tool-version pins must track `.github/workflows/ci.yml`.** golangci-lint in
> particular changes its findings between releases, so local and CI parity
> depends on bumping both together (issue #88).

## Generated files

Build artifacts land under `build/` (git-ignored): `bench.txt` (raw output),
`bench.csv` (benchstat summary), the compiled `benchreport` binary,
`coverage.out`, and `bench-alert.md`. Committed outputs are the per-environment
datasets under `docs/bench/`, the trend store under `docs/bench/history/`, and
the rendered `BENCHMARKS.md` / `docs/bench.svg` / README preview region.
