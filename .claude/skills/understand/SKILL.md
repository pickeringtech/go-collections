---
name: understand
description: Phase 1 of the delivery pipeline — read a GitHub issue (plus its linked issues/PRs), locate the relevant code, and restate scope + acceptance criteria as a compact artifact. Flags issues too big/vague for an autonomous run. Use when the user says "/understand <issue#>", "scope out issue N", or wants an issue triaged before planning. Reads only — writes no code.
---

# /understand — read an issue and restate scope

Phase 1 of the [`/full-work`](../full-work/SKILL.md) pipeline, usable on its
own. Turn a GitHub issue into a **compact, agreed statement of scope** — what's
being asked, what "done" looks like, and which code it touches — so the design
phase starts from a clean brief. **Read-only**: locate and summarize; write no
code.

## Invocation

```
/understand <issue#>
```

## What it does

Spawn an `Explore` agent (read-only fan-out) so the parent context stays lean,
and have it:

1. **Read the issue and its thread** — including linked issues/PRs:
   ```bash
   gh issue view <issue#> --comments
   ```
   Follow `#`-references and any "closes/depends on" links.
2. **Locate the relevant code** — the packages, files, and existing patterns the
   change will touch. Note neighbouring conventions to follow.
3. **Restate scope + acceptance criteria** in your own words — the concrete
   behaviour to deliver and how it will be judged done.
4. **Scope guard** — if the issue is too big or too vague for one autonomous
   run, say so explicitly with a `too_big` / `needs_decomposition` verdict and a
   one-line reason, rather than forcing a plan.

## Output — compact scope artifact

Return this (no code, no file dumps):

```json
{
  "issue": 130,
  "title": "…",
  "scope": "plain-English restatement of what to build",
  "acceptance": ["criterion 1", "criterion 2"],
  "touches": ["pkg/file.go", "pkg/…"],
  "links": ["#118", "#125"],
  "verdict": "ok | needs_clarification | too_big",
  "notes": "open questions or ambiguities, if any"
}
```

When run inside `/full-work`, this artifact is the handoff into the **Design**
phase. A `needs_clarification` or `too_big` verdict should stop at a human gate.

## Guardrails

- Read-only — no edits, no branches.
- Don't invent acceptance criteria the issue doesn't support; surface gaps as
  `notes` / `needs_clarification` instead of guessing.
