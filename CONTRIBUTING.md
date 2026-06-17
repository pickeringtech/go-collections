# Contributing to go-collections

Thanks for contributing — whether you're a human or an agent. This guide gets
you from a clean checkout to a PR that passes CI on the first try.

The substance of *how* we build lives in two places, and this guide deliberately
**summarizes and links** to them rather than copying them (so it can't drift):

- [`agent-os/standards/`](agent-os/standards/) — the design, concurrency,
  functional, documentation and testing conventions (index:
  [`agent-os/standards/index.yml`](agent-os/standards/index.yml)).
- [`.github/workflows/ci.yml`](.github/workflows/ci.yml) — exactly what gates a
  PR, with the rationale inline.

If anything here disagrees with those files, **those files win** — please fix
this guide.

## Project principles

go-collections ends per-repo reinvention with one comprehensive, reliable
collections library. Four pillars (see
[`agent-os/product/mission.md`](agent-os/product/mission.md)):

- **Zero runtime dependencies** — pure Go, nothing external to pull in. CI fails
  if a `require` appears in `go.mod`.
- **Type-safe generics** — full generics, no `interface{}` casts.
- **Consistent, thoughtful API design** — uniform patterns across every
  collection, so you learn it once and apply it everywhere.
- **Real verification** — every public function ships an Example, a Test and a
  Benchmark; collection and transformation types add a fuzz target. Coverage
  floor is **100%**.

## Local setup

You need the Go toolchain pinned in [`go.mod`](go.mod) (currently **1.24**) — no
other dependencies.

```bash
go build ./...        # compile every package
go test ./...         # run the suite (Examples are verified here too)
make test             # the suite with -race -shuffle=on (what CI runs)
make help             # list the developer entry points
```

## Design conventions

These are the house rules that keep the API uniform. Read the linked standard
before adding or changing a collection — each is one short page.

- **Capability-interface composition** — small role interfaces (`Indexable`,
  `Iterable`, `Filterable`…) composed into aggregate `Dict`/`Set`/`List`, with
  compile-time `var _ Iface = &Type{}` conformance guards.
  → [`standards/collections/interface-composition.md`](agent-os/standards/collections/interface-composition.md),
  [`interface-guards.md`](agent-os/standards/collections/interface-guards.md)
- **Immutable/mutable dual hierarchy + `InPlace` suffix** — the immutable base
  returns a *new* collection; `Mutable*` embeds it and adds `InPlace` forms that
  mutate and return nothing.
  → [`mutable-immutable-hierarchy.md`](agent-os/standards/collections/mutable-immutable-hierarchy.md),
  [`inplace-suffix.md`](agent-os/standards/collections/inplace-suffix.md)
- **Concurrent variants** — every collection gets `Concurrent*` (`sync.Mutex`)
  and `ConcurrentRW*` (`sync.RWMutex`) variants. Lock + `defer` unlock are the
  first lines of every method; reads use `RLock`. An immutable op on a
  concurrent type returns the **same concurrent type**; `InPlace` stays void.
  Callback-taking methods (`ForEach`, `Filter`, `AllMatch`, `Find`…) are the
  exception: they snapshot under the lock, release it, then invoke the callback
  outside the lock so it can safely re-enter the collection.
  → [`standards/concurrency/`](agent-os/standards/concurrency/)
- **Functional standards** — higher-order params use named `XxxFunc` type
  aliases (`FilterFunc`, `MapFunc`…); package-level funcs never mutate input and
  build a new result; parameter order is collection-first, transforming-fn-last;
  `nil`/empty input returns an initialized **non-nil** empty result, never
  `nil`.
  → [`standards/functional/`](agent-os/standards/functional/)
- **`no if init; cond` house style** — never use the `if init; cond` form;
  declare the variable on its own line, then a plain `if`.
  → [`standards/global/no-if-init-statement.md`](agent-os/standards/global/no-if-init-statement.md)
- **Readable-default, prove-then-optimize** — write clear production code; add a
  `Fast`-suffixed variant only when a benchmark proves it's worth the
  complexity.
  → [`standards/global/readability-and-performance.md`](agent-os/standards/global/readability-and-performance.md)

## Testing requirements

Full detail in [`standards/testing/`](agent-os/standards/testing/). The
essentials:

- **The trio for every public function** — a runnable **Example** with
  `// Output:`, a table-driven **Test**, and a scaling-ladder **Benchmark**
  (sub-benchmarks across 3 → 1,000,000 elements). Private functions get a Test
  only.
  → [`coverage-requirements.md`](agent-os/standards/testing/coverage-requirements.md),
  [`table-driven-tests.md`](agent-os/standards/testing/table-driven-tests.md),
  [`benchmark-scaling.md`](agent-os/standards/testing/benchmark-scaling.md)
- **A `FuzzXxx` target** for collection types and slice/map/channel
  transformation functions, in `fuzz_test.go`. Prefer invariant / differential /
  no-panic assertions over exact outputs.
- **Black-box tests** — tests live in an external `_test` package (e.g.
  `package slices_test`); exercise the public API only.
- **Standard library only** — `reflect.DeepEqual` + `t.Errorf`. No testify or
  other assertion libraries.
- **Always cover** `nil` input and empty input (both → empty non-nil output),
  and run everything with **`-race`** (`make test`).
- **Coverage floor is 100%** — every new statement ships with a covering test.
  Refactor away truly-unreachable defensive branches rather than lowering the
  floor.

## Documentation

- Document **every exported** `func`, `type`, `var` and `const`; the comment
  starts with the symbol name and explains *why* it exists, in plain language.
- Every package ships a rich [`doc.go`](agent-os/standards/documentation/package-doc.md)
  with a summary, a `# Quick Start`, and a native-vs-this `# Why Use…` /
  `# When to Use X vs Y` section.
- **No ancillary emojis** in package docs or code comments (see #34).

## CI — what gates a PR

Branch protection requires exactly **one** check: the **`CI Gate`** aggregator.
It `needs:` every blocking job, so a single stable context survives matrix and
job-name changes (#41). The full policy and rationale live at the top of
[`.github/workflows/ci.yml`](.github/workflows/ci.yml).

**Blocking** (these must pass — reproduce each locally before pushing):

| CI job | What it checks | Reproduce locally |
|---|---|---|
| Build & module hygiene | compiles; `go.mod` tidy; **zero deps**; module integrity | `go build ./...` · `go mod tidy` (no diff) · `go mod verify` |
| Test (race + coverage) | suite on Linux/macOS/Windows × Go 1.24 (+ Go 1.23 Linux); 100% floor | `make test` · `go test -race -shuffle=on -coverprofile=coverage.out ./... && go tool cover -func=coverage.out \| tail -1` |
| Lint, format & complexity | `gofmt`, `go vet`, golangci-lint (staticcheck, revive, cyclop, gocognit…) | `gofmt -l .` · `go vet ./...` · `golangci-lint run` |
| Security | known-vuln scan + security lint | `govulncheck ./...` · `gosec ./...` |
| Examples E2E | the separate `examples/` module builds and matches golden stdout | `cd examples && go test -shuffle=on ./...` |

**Report-only** (surface findings/warnings, never block a merge): Go tip,
cross-arch (386/arm64/s390x), fuzz smoke, the benchmark base-vs-head benchstat
table, and API-compatibility (`gorelease`, report-only pre-1.0). **main-only:**
`bench-report` regenerates `BENCHMARKS.md`, `docs/bench.svg` and the README
preview and bot-commits them with `[skip ci]` — don't hand-edit those.

## PR workflow

- **Branch off `main`; never commit directly to `main`.** Keep each PR focused
  on one change.
- **Conventional commits**, matching `git log` — a `scope: subject` summary
  (e.g. `collections: add LRU cache`, `ci: …`, `docs(agent-os): …`,
  `test: …`), and reference the issue in the subject when one applies
  (e.g. `(closes #56)`).
- **No AI attribution** in commits or PR descriptions unless the maintainer
  explicitly asks for it.
- Open the PR against `main`. Once the `CI Gate` is green and any required
  review is in, auto-merge lands it; branch protection blocks the merge until
  the gate passes.

---

Out of scope for this issue, but welcome as follow-ups: a
`.github/PULL_REQUEST_TEMPLATE.md` and issue templates.
