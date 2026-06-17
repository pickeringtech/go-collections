---
name: design
description: Phase 2 of the delivery pipeline — turn a scope statement (or an issue) into an implementation plan: ordered steps, the exact files to touch, risks, and any decisions that need a human. The judgement phase; use the strongest model. Use when the user says "/design <issue#>", "plan out this change", or hands over a scope artifact. Produces a plan — writes no code.
---

# /design — produce an implementation plan

Phase 2 of the [`/full-work`](../full-work/SKILL.md) pipeline, usable on its
own. Take a scope statement (the artifact from [`/understand`](../understand/SKILL.md))
or a raw issue and produce a **concrete implementation plan** the
[`/code`](../code/SKILL.md) phase can execute. This is the **judgement** step —
run it on the strongest model. It **plans**; it writes no code.

## Invocation

```
/design <issue#>        # will scope first if no scope artifact is supplied
/design                 # consume a scope artifact handed in from /understand
```

## What it does

Use the `Plan` agent (or `agent-os:shape-spec` for larger specs) to produce a
plan that **respects the repo standards** in
[`agent-os/standards/`](../../../agent-os/standards/index.yml) — so the plan
already accounts for the testing trio, `InPlace` suffix, no-`if init`, parameter
ordering, package docs, and the boy-scout rule. Specifically:

1. **Ordered steps** — the smallest sensible sequence, test-first where it fits.
2. **Exact file list** — every file to add or change (implementation **and**
   its `_test.go` / `_example_test.go` / fuzz files).
3. **Risks** — edge cases, concurrency concerns, API-compatibility impact.
4. **Decisions needing a human** — anything ambiguous, any public API/contract
   change, or a "this is too big, decompose it" verdict.

## Output — compact plan artifact

Return this (no code):

```json
{
  "issue": 130,
  "plan": ["step 1 …", "step 2 …"],
  "files": ["slices/dedupe.go", "slices/dedupe_test.go", "slices/dedupe_example_test.go"],
  "risks": ["…"],
  "acceptance": ["carried through from scope"],
  "needs_human": ["decision 1 …"],
  "verdict": "ready | needs_approval | too_big"
}
```

When run inside `/full-work`, this artifact is the handoff into **Implement**
(`/code`), and the orchestrator pauses at the design-approval gate before
proceeding.

## Gate

Always surface `needs_human` decisions and **stop for approval** when the plan
involves a public API/contract change, ambiguous scope, or a `too_big` verdict —
don't slide into implementation on those.

## Guardrails

- Plan only — no edits, no branches.
- Keep the plan executable and standards-aware; don't hand `/code` a plan that
  would force a standard to be broken.
