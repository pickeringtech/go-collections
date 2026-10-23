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

Only the **root library module** is mutated. `examples/` and `tools/*` are
separate modules the root suite never runs, so `.gremlins.yaml` excludes them
(`unleash.exclude-files`) — counting their mutants as "not covered" used to
roughly halve the reported mutator coverage without saying anything about the
library's tests.

## Two metrics

- **Efficacy** = `KILLED / (KILLED + LIVED)` — of the mutants the tests *could*
  catch, how many they did. This is the assertion-quality number.
- **Mutator coverage** = `(KILLED + LIVED) / all mutants` — how much of the
  mutated code the tests even exercise (the mutation analogue of line coverage).

## How it runs in CI

Mutation testing re-runs the affected package's tests *once per mutant* (~1,700
mutants for this module), so it is far too slow to gate every PR. Per #32 the cost is
bounded two ways (`.github/workflows/mutation.yml`):

| Job      | Trigger                                  | Scope                          |
|----------|------------------------------------------|--------------------------------|
| `full`   | Weekly (Mon 04:00 UTC) + `workflow_dispatch` | Whole module                |
| `scoped` | Per merge (the merge queue, `merge_group`) | Only lines the merge changes (`gremlins --diff`) |

Both publish a machine-readable `mutation-report.json` artifact. The `scoped`
job no-ops when a merge touches no non-test Go source. It runs once per merge in
the merge queue rather than on every PR push (issue #150): installing gremlins
fetches a newer Go toolchain, so paying that per push was waste — once per merge
keeps the signal at a fraction of the cost.

### Two guards every run needs

Both jobs set `GOFLAGS="-count=1 -exec=…/.github/scripts/memcap-exec.sh"`:

- **`-count=1`** — gremlins derives each mutant's timeout from how long its
  initial coverage pass took. setup-go restores the build cache, *including Go's
  test-result cache*, so without `-count=1` the coverage pass is a cache hit that
  "takes" ~0 s and every mutant is reported `TIMED OUT` — a meaningless report.
- **`-exec=memcap-exec.sh`** — runs each *test binary* under a per-process
  address-space cap (`ulimit -v`, 3 GiB by default, `MEMCAP_KB` to tune). Some
  mutants turn a bounded loop into an unbounded, allocating one (linked-list
  traversal steps, slice-growth bounds). A few of those in parallel exhaust the
  runner's RAM and the kernel takes the *whole runner* down ("The runner has
  received a shutdown signal") — which is why every weekly full run failed
  between June and October 2026. Capped, the runaway test dies alone with
  `fatal error: out of memory` and gremlins records the mutant as caught. The
  cap goes on the test binaries only (via `go test -exec`), never on gremlins or
  the compiler — capping gremlins itself starves it of thread stacks.

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

## Baseline

First complete full-module run (gremlins v0.6.0, October 2026, with the guards
above — earlier scheduled runs never finished):

- **1,679 mutants**: 1,460 killed, 202 lived, 9 timed out, 8 not covered.
- **Mutator coverage ≈ 99.5%** — nearly all mutated code is exercised.
- **Efficacy ≈ 87.9%** — the surviving mutants cluster in the probabilistic
  sketches (`tdigest`, `countmin`, `hll`), the linked lists and
  `slices/retrieval.go`; they are the triage queue below.

Record the floor here (just under the triaged baseline) before flipping on
gating.

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

# Always pass the same guards CI uses (see "Two guards every run needs"):
export GOFLAGS="-count=1 -exec=$PWD/.github/scripts/memcap-exec.sh"

gremlins unleash .                 # whole module (minutes)
gremlins unleash ./collections     # one package subtree (much faster)
gremlins unleash --dry-run .       # list mutants + mutator coverage, no test runs
gremlins unleash --diff main .     # only lines changed vs main (what CI's PR job does)
```

`gremlins unleash .` auto-loads `.gremlins.yaml` from the repo root.
