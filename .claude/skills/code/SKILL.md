---
name: code
description: Implement a change — write the code and its tests to the repo standards in a git worktree, then return a compact summary (files changed, what was done, what was deferred) rather than the full diff. Takes a structured plan + file list (the handoff from a design phase) or a plain task instruction. Use when the user says "/code", "implement this", "write the code for", or when an orchestrator hands off an implementation plan. Distinct from /code-review — this one writes, that one reviews.
---

# /code — implement a change to spec, in an isolated worktree

Turn an implementation plan (or a direct task) into working code **and its
tests**, following this repo's standards, without polluting the main working
tree. Return a **compact summary**, not the raw diff — the caller (a human or
the `/full-work` orchestrator) keeps its context lean.

This skill **writes**; `/code-review` and `/simplify` review. Keep them
separate — do not self-review here beyond making the change correct and
standards-compliant.

## Input

One of:

- **Structured handoff** (from a design phase) — a small artifact like:
  ```json
  {
    "task": "Add slices.PartitionInPlace",
    "files": ["slices/partition.go", "slices/partition_test.go", "slices/partition_example_test.go"],
    "plan": ["step 1 …", "step 2 …"],
    "risks": ["…"],
    "acceptance": ["…"]
  }
  ```
- **Direct instruction** — a plain-English task (`/code add a Dedupe slice helper`).
  Derive the file list yourself from the repo layout and standards.

If the task is ambiguous, contradicts a standard, or implies a public
API/contract change that wasn't sanctioned, **stop and ask** rather than
guessing (see Gates).

## Step 1 — Isolate in a git worktree

Never write the implementation directly into the primary working tree — aborted
or parallel runs must not corrupt it.

- **When spawned as a subagent** with worktree isolation already provided
  (e.g. the orchestrator launched this with `isolation: "worktree"`), you are
  already isolated — skip creating another worktree and work in place.
- **When run standalone**, create one off the current branch:
  ```bash
  BRANCH=$(git rev-parse --abbrev-ref HEAD)
  SLUG=<short-kebab-summary-of-task>
  git worktree add -b "code/${SLUG}" "../wt-${SLUG}" "${BRANCH}"
  ```
  Do all edits and test runs inside `../wt-${SLUG}`. Leave the worktree in place
  on success so the caller can inspect/commit it; report its path in the summary.
  Only `git worktree remove` it if the run is aborted and produced nothing worth
  keeping — and confirm before discarding work you didn't create.

## Step 2 — Read the standards that apply

Before writing anything, load the relevant standards so the code matches the
house style. The index lives at
[`agent-os/standards/index.yml`](../../../agent-os/standards/index.yml); read the
specific files the task touches. The ones that bite most often:

- **Testing trio** (`testing/coverage-requirements`): every **public** function
  ships an **Example + Test + Benchmark**; private funcs get a Test. Collection
  and transformation types add a **fuzz** target. Table-driven tests
  (`testing/table-driven-tests`) using an args struct, `t.Run`, `reflect.DeepEqual`,
  stdlib only. Benchmarks ladder 3 → 1,000,000 via `b.Run`
  (`testing/benchmark-scaling`).
- **`InPlace` suffix** (`collections/inplace-suffix`): mutating methods end in
  `InPlace`; unsuffixed methods return a new collection. Package-level funcs
  never mutate input (`functional/non-mutating`) and return a non-nil empty
  result for nil/empty input (`functional/zero-value-output`).
- **No `if init; cond`** (`global/no-if-init-statement`): declare the variable
  on its own line, then a plain `if`. (Mirrors the user's standing preference.)
- **Parameter ordering** (`functional/parameter-ordering`): collection first,
  transforming fn last. Higher-order params use named `XxxFunc` aliases
  (`functional/func-type-aliases`).
- **Docs** (`documentation/package-doc`): every exported symbol has a doc
  comment starting with its name; every package has a `doc.go`.
- **Concurrency** standards if touching `concurrency/` (lock discipline,
  safe variants, immutable-return contract).
- **Boy-scout rule**: leave touched code cleaner than you found it — but keep
  cleanups small and in-scope; don't smuggle a refactor into a feature change.

## Step 3 — Implement test-first

Follow the repo's scientific TDD flow
([`.ai-agent-guidelines-go.md`](../../../.ai-agent-guidelines-go.md)):

1. Write the focused failing test(s) first — smallest testable unit, edge cases
   included.
2. Write the simplest code that passes; no speculative generalisation.
3. Add Example, Benchmark, and (where applicable) Fuzz targets to complete the
   trio.
4. Document exported symbols; add/extend `doc.go` if a package gains public API.
5. `gofmt`/`go vet` as you go.

## Step 4 — Verify locally

Run the suite for every module you touched (a root `go test ./...` does **not**
descend into `examples/` or `tools/`):

```bash
go test ./...                 # root suite (Examples run here too)
make test                     # root + nested modules, -shuffle=on (what CI runs)
go vet ./... && gofmt -l .     # lint hygiene; gofmt -l should print nothing
```

If you added or changed exported behaviour, also run the relevant
`make` gate (`make cover`, `make fuzz`) when practical. Don't claim green you
didn't see — if a gate is too slow to run here, say so in the summary.

## Step 5 — Return a compact summary

Do **not** dump the diff. Return a small artifact the caller can act on:

```json
{
  "status": "done | blocked | partial",
  "worktree": "../wt-<slug>",          // or "in-place" when run as an isolated subagent
  "branch": "code/<slug>",
  "files_changed": ["slices/partition.go", "slices/partition_test.go", "..."],
  "implemented": "one-paragraph plain-English summary of what now works",
  "tests": "what was added (trio/fuzz) + local result, e.g. 'go test ./... PASS, make test PASS'",
  "deferred": ["anything intentionally left out + why"],
  "needs_human": ["decisions you hit a gate on, if any"]
}
```

## Gates — stop and ask instead of guessing

- **Ambiguous scope** the plan doesn't resolve.
- **Public API / contract change** that wasn't explicitly sanctioned by the plan
  or the user.
- **A standard would have to be broken** to satisfy the task.
- **Repeated test failure** that suggests the *plan itself* is wrong (not just a
  bug in your code) — report it up rather than thrashing.

When run inside `/full-work`, surface these via `needs_human` so the
orchestrator can escalate to a human gate.
