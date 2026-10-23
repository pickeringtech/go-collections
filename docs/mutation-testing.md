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
| "Keep the extremum" updates on equality: `stats/accumulate.go:23` (`maxAbs`), `stats/aggregate.go:52,55` (`Range`), `stats/minmax.go:21,24` (`MinMax`), `stats/clamp.go:14,17` (`Clamp`), `stats/mode.go:41`, `ml/distance/distance.go:79` (Minkowski max difference), `ml/metrics/clustering/silhouette.go:113` (nearest other cluster) | CONDITIONALS_BOUNDARY | On equality the mutant assigns (or returns) a value equal to the one already held. For floats the only distinguishable case is the sign of a zero (−0 vs +0), which compare equal and whose sign these functions do not specify. | None — equivalent |
| Guards whose boundary branch is a no-op: `slices/chunk.go:23`, `slices/zip.go:22,39`, `slices/transformreverse.go:8` (the middle element swaps with itself), `slices/modify.go:56` (an empty fill range is a plain copy either way), `channels/window.go:135` (at `step == len(buf)` both branches empty the buffer with `skip` 0), `concurrency/bgworklimiter.go:55`, `concurrency/blockingworklimiter.go:15`, `concurrency/parallel.go:92`, `ml/metrics/ranking/ranking.go:18` | CONDITIONALS_BOUNDARY | At the boundary the guarded branch sets the variable to the value it already has (`limit` 1, `policy` StopOnError, `k` n, `end`/`n` the length). | None — equivalent |
| `slices/retrieval.go:162,166,170,173,176` (`SubSlice` range guards) | CONDITIONALS_BOUNDARY ×7 | Each boundary case the mutant moves between the early `return []T{}` and the clamped `input[from:to:to]` yields an empty, non-nil, zero-capacity slice either way, or clamps an index to the value it already has. | None — equivalent |
| `stats/histogram.go:69` (the `(x−min)/width` bin estimate), `:99` (`pos > 0`), `:102` (the downward walk) | ARITHMETIC_BASE ×4, INVERT_NEGATIVES, CONDITIONALS_BOUNDARY ×3, CONDITIONALS_NEGATION | `locateBin` uses the estimate only as a starting index, then steps onto the bin whose reported edges hold `x`, so every start in `[0, last]` lands on the same bin. `pos` is never negative (`x ≥ min`) and `bins[0].Min == min`, so the walk never leaves the slice; only the number of steps changes. `FuzzHistogram` passes under these mutants. | None — affects only speed |
| `stats/vector.go:58` (`total >= minNormal`) | CONDITIONALS_BOUNDARY | Chooses between the fast and the scaled `Norm` at a sum of squares of exactly 2⁻¹⁰²², where both compute the same norm. | None — equivalent |
| Capacity hints: `ml/preprocessing/binners.go:104,166` (`nBins-1`), `ml/preprocessing/splits.go:71,184` (`len(input)-testN`) | ARITHMETIC_BASE, INVERT_NEGATIVES (the same `-`→`+` rewrite) | Only the capacity preallocated for a slice that `append` then fills changes; its contents and length do not. | None — equivalent |
| Sorts that differ only on equal keys: `ml/preprocessing/encoders.go:37` (`sortedUnique` sorts already-deduplicated values), `ml/metrics/classification/roc.go:86` (`>`→`>=`), `roc.go:120` | CONDITIONALS_BOUNDARY | The comparator changes only how equal keys are ordered, and none of these depend on it: there are no equal keys after deduplication, tied scores collapse to one ROC point, and they share one averaged rank. | None — equivalent |
| Tie-breaks in `StratifiedSplit`: `ml/preprocessing/splits.go:258` (largest-remainder order `>`→`>=`), `:314:51` (`pick` `>`→`>=`) | CONDITIONALS_BOUNDARY | Changes only which of several equally entitled labels wins a tie. The winner is still chosen by the seeded priority, so every documented property (exact test size, largest-remainder shares, both-sides coverage, the same split for the same seed) still holds; only which seed gives which split changes. A brute-force search over every layout of up to four labels of up to 12 rows found no seed-independent case that tells them apart. | None — pinning it would freeze one seed's output |
| `ml/metrics/ranking/ranking.go:121` (`hits++`) | INCREMENT_DECREMENT | `sum` and `hits` both flip sign, so `sum/hits` is unchanged and the `hits == 0` check is unaffected. | None — equivalent |
| `ml/metrics/classification/classification.go:295` (`binaryStats` `s.support++`) | INCREMENT_DECREMENT | Dead store: the binary metrics read only `tp`/`fp`/`fn`, and nothing reads `support` on this path. | None — could be removed as dead code; not a bug |
| `ml/metrics/clustering/silhouette.go:46` (`sizes[labels[i]]++`) | INCREMENT_DECREMENT | Only `len(sizes)`, the number of distinct labels, is read, and `--` creates the key too. | None — equivalent |
| `ml/metrics/clustering/silhouette.go:31` (`n < 2`) | CONDITIONALS_BOUNDARY | With `n == 2` the cluster-count check (2 ≤ k ≤ n−1 = 1) rejects the input anyway. | None — equivalent |
| `ml/distance/discrete.go:43` (Levenshtein's swap so that `b` is the shorter input) | CONDITIONALS_BOUNDARY, CONDITIONALS_NEGATION | The swap only chooses which sequence becomes the DP row. Edit distance is symmetric, so the result is the same; under the negation, memory is O(max) instead of O(min). | None — affects only memory |
| `ml/preprocessing/scalers.go:54` (`sum / n` in the non-finite fallback) | ARITHMETIC_BASE | Runs only when the data holds a NaN or ±Inf, so `sum` is non-finite, and dividing or multiplying it by `n ≥ 1` gives the same value. | None — equivalent |

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
