---
name: full-work
description: Drive a GitHub issue through the whole delivery lifecycle — understand → design → implement → test → verify → review → PR — by composing focused skills/agents, each run in its own context fork so the parent stays lean. Mechanical phases use cheaper models; judgement phases use the strongest. Pauses at human gates (design approval, ambiguous/oversized scope). Use when the user says "/full-work <issue#>", "take issue N to a PR", or wants a routine issue delivered end-to-end.
---

# /full-work — issue → PR orchestrator

Take a GitHub issue from **read** to **opened PR** with minimal hand-holding, by
running each phase as a **forked subagent** and passing a **small structured
handoff** between phases — never raw transcripts. The orchestrator does little
work itself: it sequences phases, routes each to an appropriate model, enforces
human gates, and threads the compact artifact from one phase into the next.

**Context hygiene is the whole point.** Each phase's diff, test output, and
review chatter stays inside that phase's fork; only its *conclusion* (a few-line
JSON artifact) returns here. If you find yourself reading a full diff or test
log in this context, you've broken the design — push that work into a subagent.

## Invocation

```
/full-work <issue#> [--design-model M] [--implement-model M] [--no-design-gate] [--max-fix-loops N]
```

- `<issue#>` — required GitHub issue number.
- Model overrides and gate flags are optional (see Config).

## Why a skill (not a Workflow script)

This is built as a **skill that forks subagents via the Agent tool**, matching
how the rest of this repo's tooling composes (`/pr` → `/pr-watch`). A
`Workflow` script would also fit the deterministic fan-out, but workflows aren't
a checked-in pattern here and require explicit opt-in per run; a skill is
directly invokable as `/full-work N` and reuses the existing skills as-is. If
this pipeline later needs heavy parallel fan-out, promote the phase loop to a
`Workflow` and keep this skill as the thin entry point.

## Config — model routing & gates

Defaults below; override per-run via the flags above. Route mechanical phases to
cheaper/faster models and reserve the strongest model for judgement.

Every phase is its **own slash command** and is usable standalone — the
orchestrator just sequences them and threads the handoff. The command does the
work; the agent/model column is *how* the orchestrator forks it.

| Phase | Command | Forked as | Default model | Fork |
|-------|---------|-----------|---------------|------|
| 1. Understand | [`/understand`](../understand/SKILL.md) | `Explore` agent | Sonnet | yes |
| 2. Design | [`/design`](../design/SKILL.md) | `Plan` agent / `agent-os:shape-spec` | **Opus** | yes |
| 3. Implement | [`/code`](../code/SKILL.md) | subagent, **worktree** | Sonnet | yes — **worktree** |
| 4. Test | [`/test`](../test/SKILL.md) | subagent | Haiku | yes |
| 5. Verify | `/verify` | subagent | Sonnet | yes |
| 6. Review/cleanup | `/code-review`, `/simplify` | subagent | **Opus** | yes |
| 7. PR | `/pr`, `/pr-watch` | subagent | Sonnet | yes |

When you spawn a phase, invoke its command inside a forked subagent, passing
`model:` per this table (or the override) and `subagent_type` where one fits
(`Explore` for phase 1, `Plan` for phase 2). The implement phase passes
`isolation: "worktree"`. Each command returns its own compact artifact; the
orchestrator merges that into the shared handoff below.

## Handoff artifact

A single compact object threads through all phases. Each phase reads it and
returns an updated copy — keep it small (no diffs, no logs):

```json
{
  "issue": 130,
  "scope": "restated scope + acceptance criteria",
  "plan": ["step 1", "…"],
  "files": ["path/a.go", "path/a_test.go"],
  "risks": ["…"],
  "branch": "code/<slug>",
  "worktree": "../wt-<slug>",
  "test": "PASS|FAIL summary",
  "verify": "confirmed|failed + note",
  "review": "clean|fixes-applied summary",
  "pr": "url",
  "needs_human": []
}
```

## Phases

Run in order. Each is a forked subagent; feed it the current handoff, get back
the updated handoff. Honour the gates between phases.

### 1. Understand — *([`/understand`](../understand/SKILL.md), Explore, Sonnet)*
Run `/understand <issue#>` to read the issue + linked issues/PRs, locate the
relevant code, and **restate scope + acceptance criteria**. Merge its `scope`,
`acceptance`, and `touches` into the handoff. **Scope guard:** a
`too_big`/`needs_clarification` verdict stops at a human gate.

### 2. Design — *([`/design`](../design/SKILL.md), Plan / agent-os:shape-spec, Opus)*
Run `/design` on the scope artifact to produce an implementation plan: ordered
steps, the exact `files` to touch, `risks`, and any **decisions that need a
human**. This is the judgement phase — use the strongest model. Merge `plan`,
`files`, `risks` into the handoff.

> **GATE — design approval (default ON).** Present the plan and acceptance
> criteria to the user and wait for approval before implementing. Skip only if
> `--no-design-gate` was passed. **Always** stop here regardless of the flag if
> the design surfaced ambiguous scope, an oversized/"needs decomposition"
> verdict, or an API/contract change.

### 3. Implement — *(`/code`, Sonnet, worktree)*
Hand the `{plan, files, scope, acceptance}` slice of the handoff to
[`/code`](../code/SKILL.md), spawned with `isolation: "worktree"`. It writes the
code + tests to the repo standards and returns a compact summary. Thread its
`branch`, `worktree`, and `files_changed` back into the handoff. If `/code`
returns `needs_human`, escalate to a gate.

### 4. Test — *([`/test`](../test/SKILL.md), Haiku)*
Run `/test` in the implementation worktree (CI parity: `make test`, `go vet`,
`gofmt -l`, `make cover`). It returns a compact `result: PASS|FAIL` verdict with
one-line failure causes — no logs in the parent. Merge into `test`. On `FAIL` →
Fix loop.

### 5. Verify — *(`/verify`, Sonnet)*
Run `/verify` to confirm the change actually does what the
**issue** asked (not just that tests pass). Return `verify`. On failure that
implies the *design* was wrong → stop at a human gate rather than looping.

### 6. Review / cleanup — *(`/code-review`, `/simplify`, Opus)*
Self-review the diff with `/code-review` and apply boy-scout cleanups with
`/simplify`. This is judgement work — strongest model. Return a `review`
summary. If review surfaces real bugs → Fix loop.

### 7. PR — *(`/pr` → `/pr-watch`, Sonnet)*
From the implementation branch, run [`/pr`](../pr/SKILL.md) to commit, push, and
open the PR (commit message follows repo convention: `pkg: summary (closes
#<issue>)`), then it hands off to [`/pr-watch`](../pr-watch/SKILL.md) to drive
to green. Return `pr` (the URL).

## Fix loop (bounded)

If **Test** or **Review** finds a fixable problem, loop back to **Implement**
(`/code`) with the failure summary appended to the handoff. Cap at
`--max-fix-loops` (**default 2**). On exhausting the cap, **stop and report** —
do not thrash. A **Verify** failure that implies the design was wrong skips the
loop and goes straight to a human gate.

## Human gates — where to stop and ask

- **After Design** (default), and always on ambiguous scope, oversized/"needs
  decomposition" verdict, or API/contract change.
- **Any phase returning `needs_human`** (e.g. `/code` hit a standard conflict).
- **Verify failure** implying the design was wrong.
- **Fix-loop cap** exhausted.

Use the `AskUserQuestion` tool at gates; present the relevant artifact (plan,
failure summary) compactly so the user can decide without re-reading context.

## Guardrails

- Reuse the existing skills/agents — never reimplement `/code-review`,
  `/verify`, `/pr`, `/pr-watch`.
- Keep the parent context lean: phases return artifacts, not transcripts.
- Implementation is worktree-isolated; the main tree is never edited directly.
- Follow `agent-os/standards/` throughout (the implement phase enforces this).
- Don't add Claude/AI attribution to commits or the PR unless asked.
