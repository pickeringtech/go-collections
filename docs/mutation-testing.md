# Mutation testing

Line coverage tells you a line *ran*; it cannot tell you a test would have
*failed* if that line were wrong. A test can execute code and assert nothing
meaningful. Mutation testing closes that gap: it injects small faults
("mutants") into the source — flip `<` to `<=`, negate a condition, swap `+`
for `-` — and re-runs the suite. If a test fails, the mutant is **killed**
(good — the suite caught the regression). If every test still passes, the
mutant **lived** — a hole in the assertions.

This directly serves the project's "real verification" goal (#32): with ~99%
line coverage, mutation score is the honest answer to *"do our tests actually
protect us?"*

## Tool

[`gremlins`](https://gremlins.dev) (`github.com/go-gremlins/gremlins`), pinned
in [`.github/workflows/mutation.yml`](../.github/workflows/mutation.yml) via
`GREMLINS_VERSION`. What it mutates is pinned in
[`.gremlins.yaml`](../.gremlins.yaml) (the v0.6.0 default mutant set, listed
explicitly so an upstream default change can't silently move the score).

> gremlins v0.6.0 is built with Go ≥ 1.25, newer than this module's declared
> `go 1.24`. CI installs the tool with `GOTOOLCHAIN=auto` so the newer Go is
> fetched only to *build gremlins*; the project's own tests still run under the
> `go.mod` version.

## Two metrics

- **Efficacy** = `KILLED / (KILLED + LIVED)` — of the mutants the tests *could*
  catch, how many they did. This is the assertion-quality number.
- **Mutator coverage** = `(KILLED + LIVED) / all mutants` — how much of the
  mutated code the tests even exercise (the mutation analogue of line coverage).

## How it runs in CI

Mutation testing re-runs the whole suite *once per mutant* (~400 mutants for
this module), so it is far too slow to gate every PR. Per #32 the cost is
bounded two ways (`.github/workflows/mutation.yml`):

| Job      | Trigger                                  | Scope                          |
|----------|------------------------------------------|--------------------------------|
| `full`   | Weekly (Mon 04:00 UTC) + `workflow_dispatch` | Whole module                |
| `scoped` | Every pull request                       | Only lines the PR changed (`gremlins --diff`) |

Both publish a machine-readable `mutation-report.json` artifact. The `scoped`
job no-ops when a PR touches no non-test Go source.

## Gating policy (report-only → ratchet)

Both jobs are **report-only today** (`continue-on-error: true`): the first runs
need triage (below) before failures are meaningful. The score floors live in
the workflow as env vars — `MUTATION_EFFICACY_MIN` and `MUTATION_MCOVER_MIN`,
currently `0` — the single ratchet knob, mirroring how `COVERAGE_MIN` works in
`ci.yml`.

To start gating, in two deliberate steps:

1. **Set a floor.** Raise `MUTATION_EFFICACY_MIN` (and/or `MUTATION_MCOVER_MIN`)
   to just under the triaged baseline. gremlins then exits non-zero below it,
   but `continue-on-error` still keeps the job green — a grace period to confirm
   the score is stable run-to-run.
2. **Enforce.** Remove `continue-on-error` from the job so the floor blocks.

Thereafter ratchet the floor **upward only, never down** — exactly the
coverage-floor discipline. A drop in mutation score is a regression in test
quality and should fail.

## Baseline (first observations)

From the initial local run (gremlins v0.6.0):

- **Mutator coverage ≈ 98.5%** across the module (406 runnable mutants, 6 not
  covered) — nearly all mutated code is exercised.
- On the `collections` subtree, **efficacy ≈ 91%** with a small number of
  `LIVED` mutants to chase down (e.g. `collections/sets/hash.go` arithmetic
  mutants that no assertion distinguishes).

A full-module efficacy number is established by the first scheduled/dispatched
`full` run — read it off the run's summary or the `mutation-report` artifact,
then record the floor here before flipping on gating.

## Known-equivalent / non-actionable mutants

Some surviving or timed-out mutants are **not** test gaps and should be excluded
from triage rather than "fixed":

- **`TIMED OUT` in linked-list traversal** (`collections/lists/linked.go`).
  Mutating a loop boundary or pointer step turns a finite traversal into an
  infinite loop; the test never returns and gremlins times it out rather than
  observing a failure. These dominate the `lists` package's mutant count. They
  are inherent to mutating linked structures, not missing assertions. Prefer
  raising `unleash.timeout-coefficient` or excluding the specific construct over
  adding contrived tests; track decisions in the table below.

| File / construct | Mutant type | Why it's equivalent / non-actionable | Action |
|------------------|-------------|--------------------------------------|--------|
| _(fill in during triage of the first full run)_ | | | |

To permanently exclude a file from mutation, add it to `unleash.exclude-files`
in `.gremlins.yaml` (a filepath regexp) with a comment explaining why.

## Running it locally

```bash
# Install the pinned version (Go's toolchain auto-fetches the Go gremlins needs):
GOTOOLCHAIN=auto go install github.com/go-gremlins/gremlins/cmd/gremlins@v0.6.0

gremlins unleash .                 # whole module (slow — minutes)
gremlins unleash ./collections     # one package subtree (much faster)
gremlins unleash --dry-run .       # list mutants + mutator coverage, no test runs
gremlins unleash --diff main .     # only lines changed vs main (what CI's PR job does)
```

`gremlins unleash .` auto-loads `.gremlins.yaml` from the repo root.
