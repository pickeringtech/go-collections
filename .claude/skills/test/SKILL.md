---
name: test
description: Phase 4 of the delivery pipeline — run the Go test suite (root + nested modules, race, vet, gofmt, and the coverage floor) and return a compact pass/fail verdict, not raw logs. Use when the user says "/test", "run the tests", "check the suite", or before pushing. Diagnoses failures briefly; fixing them is /code's job.
---

# /test — run the suite and summarize pass/fail

Phase 4 of the [`/full-work`](../full-work/SKILL.md) pipeline, usable on its
own. Run the project's tests the way CI does and return a **compact verdict** —
the parent context never sees raw test logs. This skill **runs and reports**; it
does not fix failures (that loops back to [`/code`](../code/SKILL.md)).

## Invocation

```
/test            # default: full suite, CI parity
/test quick      # fast inner-loop: root suite only
/test <pkg>      # scope to a package path, e.g. /test ./slices
```

## What it does

A root `go test ./...` does **not** descend into the nested modules
(`examples/`, `tools/benchreport/`), so use the Makefile targets that test every
module the way CI does.

**Full (default) — CI parity:**
```bash
make test                  # root + every nested module, -shuffle=on, race
go vet ./...
gofmt -l .                 # must print nothing
make cover                 # root suite with -race -shuffle + coverage floor (COVERAGE_MIN)
```

**Quick:**
```bash
go test ./...              # root suite + Examples only
go vet ./...
```

Run from a worktree when invoked after [`/code`](../code/SKILL.md) so you test
the change in isolation. Capture failing output to diagnose, but **do not** echo
full logs upward — distil to a one-line cause per failure.

## Output — compact verdict

```json
{
  "result": "PASS | FAIL",
  "ran": ["make test", "go vet", "gofmt -l", "make cover"],
  "failures": [
    {"target": "slices", "cause": "TestDedupe: got [1 1 2], want [1 2]"}
  ],
  "coverage": "92.3% (floor 90%) | n/a",
  "gofmt_dirty": []
}
```

On `PASS`, `failures` is empty. When run inside `/full-work`, a `FAIL` verdict
triggers the bounded fix loop back to `/code` with these `failures` attached;
the parent never ingests the logs.

## Guardrails

- Report honestly — if a gate was too slow to run here (e.g. full `make fuzz`),
  say so in `ran`/`notes` rather than implying it passed.
- Diagnose, don't fix — never weaken a test or assertion to make the suite go
  green; surface the failure so `/code` addresses the cause.
