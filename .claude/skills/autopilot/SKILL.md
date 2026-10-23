---
name: autopilot
description: Autonomous, goal-seeking controller — pursue a stated goal by running a phase, reading its STRUCTURED results, and re-planning the next phase FROM those results, looping until a goal predicate, convergence, human gate, or budget backstop fires. The pipeline is chosen by an LLM planner each iteration, not baked in — this is the generalized, more powerful sibling of /full-work (which flies a fixed issue→PR path). Persists across turns (ScheduleWakeup) and pauses at human gates (AskUserQuestion). Use when the user says "/autopilot <goal>", "pursue this goal autonomously", "keep working until <X> is true", or wants an open-ended objective driven to done with result-driven re-planning. NOT /loop (which re-fires a fixed prompt) and NOT a raw Workflow (which can't pause for humans or persist across turns).
---

# /autopilot — result-driven re-planning controller

Pursue a **stated goal** autonomously: run a phase, read its *structured*
results, and decide the **next** phase **from those results** — looping until
the goal is met or a stop policy fires. Unlike `/full-work`, the pipeline is not
hardcoded; an LLM **planner picks each phase per iteration**. Unlike `/loop`, it
does not re-fire a fixed prompt — every iteration is chosen from where the work
actually is.

The controller does little work itself. It runs a thin **model-driven
re-planning loop** (the judgement) wrapped around **deterministic fan-out** (each
phase executes via `Workflow`, or by composing the existing delivery skills). It
sequences phases, folds their results into accumulated state, checks the stop
policy, and holds course against drift.

**Context hygiene is the whole point** (same discipline as `/full-work`): a
phase's diff, logs, and search dumps stay inside that phase's fork or `Workflow`
run; only its *compact structured artifact* returns to this loop. If you find
yourself reading a full diff or log in this context, push that work into a
subagent/`Workflow`.

## Where this sits — the three-layer gap it closes

| Layer | Loops over | Re-plans from results? | Across turns? | Human gates? |
|-------|-----------|------------------------|---------------|--------------|
| `Workflow` | fan-out **within one run** | no — returns once, can't re-plan a *different* next phase | no | no |
| `/loop` | a **fixed prompt/command** on an interval | no — runs the same thing each fire | yes | no |
| `/full-work` | a real controller loop | **pipeline hardcoded** to issue→PR | within a run | yes |
| **`/autopilot`** | **phases chosen per iteration** | **yes — the whole point** | **yes** (`ScheduleWakeup`) | **yes** (`AskUserQuestion`) |

The missing piece is the bottom row: a controller that reads its own structured
results and re-plans, **persists** that re-planning across turns, and **pauses**
at human gates — none of which a single `Workflow` run can do.

## Why a skill (not a raw Workflow script)

The fan-out *inside* each phase is exactly what `Workflow` is for, and
`/autopilot` uses it there. But the **outer** loop must be a skill because it
needs two things a single `Workflow` run cannot provide:

- **Persist a long-running re-plan across turns** — via `ScheduleWakeup`, so the
  drive survives context summarization and idle waits on external state.
- **Pause at human gates** — via `AskUserQuestion`, when the planner flags
  ambiguous/risky/oversized direction.

And the core decision — *which phase next, given these results* — is
**model-driven**, which the `Workflow` tool deliberately steers away from
(`Workflow` is for deterministic control flow). So the judgement lives in the
skill; the determinism lives in the `Workflow` calls it dispatches.

## Invocation & arguments

```
/autopilot <goal> [--max-iter N] [--converge-k K] [--budget T] [--no-gates] [--plan-model M]
```

- `<goal>` — **required**, free-text objective. Stored **verbatim and
  immutable** in drive state; it is the anchor every CRITIC check measures
  against. A vague goal is the single biggest risk (see CRITIC) — if the goal is
  not falsifiable ("improve the code"), **stop at a gate and ask for a goal
  predicate** before the first iteration.
- `--max-iter N` — hard iteration cap (backstop). **Default 12.**
- `--converge-k K` — stop after K consecutive iterations that FOLD adds nothing
  new. **Default 2.**
- `--budget T` — token target for the whole drive (passed to `Workflow` as its
  budget). Optional; the `--max-iter` cap is the always-present backstop.
- `--no-gates` — suppress *non-essential* gates only. **Never** suppresses the
  ambiguous/risky/oversized human gate or a destructive-action confirm.
- `--plan-model M` — model for the judgement steps (ASSESS / CRITIC / PLAN).
  **Default: the strongest available (Opus).** Phase bodies route their own
  models per the dispatch table.

## Drive state — the accumulated object (persisted across turns)

A single compact object threads through every iteration **and survives across
turns**. Keep it small — artifacts and counters, never diffs or logs. When the
loop pauses (`ScheduleWakeup`, or a gate), this object *is* the resume point;
restate it compactly so the next turn continues without re-reading transcripts.

```json
{
  "goal": "<original goal, verbatim — immutable>",
  "goal_predicate": "<explicit success test: tests green / criteria met / queue empty>",
  "iteration": 3,
  "max_iter": 12,
  "converge_k": 2,
  "no_new_streak": 0,
  "seen": ["<stable key of every item ever folded — the CUMULATIVE set>"],
  "phases": [
    {"i": 1, "goal": "phase goal", "shape": "parallel|pipeline", "artifact": {"…": "compact result"}}
  ],
  "last_results": {"…": "structured output of the most recent DISPATCH"},
  "budget_spent": 0,
  "needs_human": [],
  "status": "running|gated|converged|goal-met|capped|stopped"
}
```

- **`goal` is written once and never edited.** If you feel the urge to rewrite
  it because the work has moved on, that *is* the drift CRITIC exists to catch —
  escalate, don't rewrite.
- **`seen` is cumulative**, not per-round. This is load-bearing (see Convergence).

## Per-iteration state machine

Run the steps below in order each iteration. Three are **model-driven**
judgement (ASSESS / CRITIC / PLAN — use `--plan-model`); DISPATCH is
**deterministic** fan-out; FOLD is **plain merge code**, not an agent.

```
1.   ASSESS   → agent reads drive state + last_results → structured "where are we"
1.5  CRITIC   → cheap drift check: does the proposed direction still serve the ORIGINAL goal?
2.   CHECK    → evaluate the composite stop policy (below) on accumulated state; maybe stop/gate
3.   PLAN     → agent proposes the next phase: {goal, work-list, fan-out shape} + maybe a gate flag
3.5  GATE     → if PLAN flagged a gate, PAUSE (AskUserQuestion) BEFORE dispatching; resume on answer
4.   DISPATCH → Workflow runs that phase deterministically (parallel/pipeline)
5.   FOLD     → merge results into drive state; dedup against the CUMULATIVE seen-set
6.   loop     → ScheduleWakeup or continue inline
```

CHECK (step 2) and GATE (step 3.5) are **two distinct gate points, not one**.
CHECK evaluates the stop policy against *accumulated* state carried in from prior
iterations, **before** spending a model call on PLAN. GATE handles the *fresh*
risk PLAN just surfaced about the phase it proposed — it must fire **before**
DISPATCH, so a risky/destructive phase is never executed in the same iteration
it was flagged.

### 1. ASSESS — *(judgement, `--plan-model`)*
A subagent reads the drive state (goal, prior phase artifacts, `last_results`)
and returns a compact structured read of **where we are relative to the goal**:
what's done, what's outstanding, what the last results imply. No new work is
planned here — ASSESS only *describes the situation*. Returns a small object;
the controller never reads raw last-phase output, only this summary.

### 1.5 CRITIC — *(judgement, `--plan-model`) — the drift mitigation*
The chief risk of any result-driven loop is **rationalized drift**: the loop
convinces itself it's progressing while wandering off-goal. CRITIC is the
mitigation. A *cheap* subagent compares ASSESS's read (and the direction it
implies) against the **immutable original `goal`** and answers:

- Does continuing in this direction still serve the original goal? (yes/no + why)
- Is the loop substituting an *easier nearby* goal for the stated one?
- Has "done" quietly been redefined since iteration 1?

On **drift = yes** → do **not** proceed to PLAN as-is. Either correct course
back toward the goal, or — if the goal itself now looks wrong/ambiguous —
**escalate to a human gate**. Without CRITIC, the goal predicate is the *only*
thing pulling the loop back, and a vague predicate won't. Keep CRITIC cheap
(small model is fine for the check itself); its job is a tripwire, not a full
re-plan.

### 2. CHECK — *(controller, deterministic)*
Evaluate the **composite stop policy** (next section) in precedence order. If
any level fires, stop or gate **before** spending a model call on PLAN/DISPATCH.

### 3. PLAN — *(judgement, `--plan-model`)*
A subagent proposes the **single next phase**: its sub-goal, the **work-list**
(the concrete items to fan out over), and the **fan-out shape**
(`parallel` barrier vs `pipeline`, and roughly how many agents). PLAN also
**raises a gate flag** if the next phase is ambiguous, risky (destructive/outward
action), or oversized/needs-decomposition. That flag is handled by the **GATE
step immediately below — before DISPATCH**, in the *same* iteration; it does not
wait for the next loop. PLAN proposes **one** phase, not the whole remaining
pipeline — the next phase after that is re-decided from *its* results.

### 3.5 GATE — *(controller, deterministic) — honour PLAN's flag before dispatching*
If PLAN raised a gate flag, **stop here and ask** (`AskUserQuestion`) **before**
DISPATCH runs anything — present the proposed phase and why it was flagged, get
a decision, fold it into drive state, then continue (proceed, revise the phase,
or abandon it). This is the control-flow guarantee that a risky/destructive or
ambiguous phase is **never executed in the iteration it was flagged**. Honour it
even under `--no-gates` — that flag suppresses only non-essential gates, never a
risk/ambiguity/destructive one (see the stop policy). If PLAN raised no flag,
GATE is a no-op and DISPATCH proceeds.

### 4. DISPATCH — *(deterministic fan-out via `Workflow`)*
Execute the planned phase. Two routes:

- **Compose an existing delivery skill** when the phase maps to one
  (`/design`, `/code`, `/test`, `/code-review`, `/verify`, `/pr` →
  `/pr-watch`) — never reimplement them. Route models per the table below.
- **Run an ad-hoc `Workflow`** for novel fan-out (a discovery sweep, a
  multi-lens review, a migration over a work-list) — pass the work-list and
  shape from PLAN, and the `--budget` so its token spend counts against the
  drive. Each agent returns a **schema-validated** compact result.

DISPATCH returns one compact **structured** artifact (not transcripts) into
`last_results`.

### 5. FOLD — *(plain merge code — NOT an agent)*
Merge `last_results` into the drive state: append the phase artifact, advance
`budget_spent`, and **dedup new items against the cumulative `seen` set**. Count
how many items were genuinely *new*; if zero, increment `no_new_streak`,
otherwise reset it to 0 and add the new items' stable keys to `seen`.

### 6. loop
If a terminal/gate condition holds, stop/ask. Otherwise continue: inline if work
is ready now, or `ScheduleWakeup` if waiting on external state (CI, a deploy, a
queue) — see Cross-turn persistence.

## Composite stop policy (precedence matters — they mean different things)

Evaluate top-down in CHECK. The order is not arbitrary: each level means
something the one below it does not, and a lower level must **never** override a
higher one.

1. **Human gate (highest).** PLAN flagged the direction ambiguous / risky /
   oversized, or any phase returned `needs_human`. **Pause, ask, resume.**
   *Never* overridden by budget — running out of tokens is not a license to make
   a risky call unsupervised. Use `AskUserQuestion`; present the relevant
   artifact compactly. Suppressed by `--no-gates` **only** for non-essential
   gates, never for risk/ambiguity/destructive actions.
2. **Goal predicate.** The explicit success test in `goal_predicate` is met
   (tests green / acceptance criteria satisfied / work-queue empty). A **clean
   win → stop and report.** This is the *intended* exit.
3. **Convergence.** `no_new_streak >= converge_k` — K consecutive iterations
   that added nothing new. The loop has run dry; stop. (See below — this is the
   step most prone to a silent bug.)
4. **Token budget (backstop).** `budget.remaining()` is exhausted, or
   `iteration >= max_iter`. This should **rarely** be the one that fires — if it
   is, levels 1–3 underperformed (usually a goal predicate too vague to ever
   register as "met"). Report it as a backstop stop, not a success.

## Convergence & dedup — the #1 failure mode

Dedup in FOLD **must** be against the **cumulative `seen` set** (every item ever
folded), **not** the last round's items. This is the single most common reason
these loops never terminate:

> If you dedup against only the previous round, an item that was found, judged,
> and **rejected** looks "new" again next round → it gets re-planned, re-judged,
> re-rejected, forever. The streak never grows, convergence never fires, and the
> loop spins until the budget backstop kills it (reported as a failure, not a
> clean win).

So: every item that is ever folded — accepted **or** rejected — gets its stable
key added to `seen`. "New" means "key not in `seen`." Convergence counts
iterations that produced **zero** keys not already in `seen`. Use a *stable*
key (content/identity hash), not object identity or body text that bots may
edit in place.

## Cross-turn persistence

When the loop must wait on something this context can't be notified about (CI
finishing, a deploy, a queue draining), don't block — **`ScheduleWakeup`** with
the drive state as the resume point, and re-enter this skill on wake. Pacing
follows the built-in `/loop` skill's dynamic mode (`/loop` ships with Claude
Code; it is not defined in this repo):

- **< 5 min (60–270s)** — actively polling external state (CI/deploy/queue);
  keeps the prompt cache warm.
- **20–30 min (1200–1800s)** — idle heartbeat when there's no specific signal to
  watch; don't burn cache polling for nothing.
- **Stop scheduling** once a terminal condition (goal-met / converged / capped /
  gated) is reached.

Harness-tracked work (a `Workflow` you launched, a background task) re-invokes
you on completion automatically — **don't** schedule a short poll just to check
on it; schedule only a long fallback in case it hangs.

## Phase bodies — compose, don't reimplement

DISPATCH should reuse what already works. Default model routing (mechanical →
cheaper/faster, judgement → strongest); override via `--plan-model` for the
controller's own ASSESS/CRITIC/PLAN steps.

| Phase kind | Compose | Default model |
|------------|---------|---------------|
| Scope / understand | [`/understand`](../understand/SKILL.md), `Explore` agent | Sonnet |
| Design / plan a sub-pipeline | [`/design`](../design/SKILL.md), `Plan` agent | **Opus** |
| Implement | [`/code`](../code/SKILL.md), **worktree** isolation | Sonnet |
| Test | [`/test`](../test/SKILL.md) | Haiku |
| Verify behaviour | `/verify` | Sonnet |
| Review / cleanup | `/code-review`, `/simplify` | **Opus** |
| Ship | [`/pr`](../pr/SKILL.md) → [`/pr-watch`](../pr-watch/SKILL.md) | Sonnet |
| Novel fan-out (sweep/audit/migrate) | ad-hoc `Workflow` script | per phase |
| ASSESS / CRITIC / PLAN (controller) | this skill's own subagents | **`--plan-model` (Opus)** |

When a sub-goal *is* the issue→PR lifecycle, the right move is often to DISPATCH
**`/full-work`** as a single phase rather than re-deriving its pipeline —
`/autopilot` is a layer *above* it, not a replacement for it.

## Human gates — where to stop and ask

- **Vague/unfalsifiable goal** at start — ask for a goal predicate before
  iterating.
- **PLAN flags** ambiguous scope, a destructive/outward action, or an
  oversized/"needs decomposition" phase.
- **CRITIC detects drift** that implies the goal itself is wrong or ambiguous.
- **Any phase returns `needs_human`** (e.g. `/code` hit a standards conflict, a
  conflict needs a product decision).

Use `AskUserQuestion`; present the relevant artifact (ASSESS read, the proposed
phase, the failure) compactly so the user decides without re-reading context.
On resume, fold the answer into drive state and continue from CHECK.

## Termination & reporting

Stop and report when any holds (mapped to the stop policy):

- **Goal met** — `goal_predicate` satisfied. The clean win. Report what was
  achieved.
- **Converged** — `no_new_streak >= converge_k`. Report what was done and why it
  ran dry (likely the goal is effectively complete, or the remaining work is
  out of this loop's reach).
- **Gated** — a human gate fired; report the decision needed and the resume
  state.
- **Capped** — `iteration >= max_iter` or budget exhausted. Report honestly as a
  **backstop** stop, not a success, and name the most likely cause (usually a
  goal predicate too vague to register as met).

Surface a short summary: the goal, iterations run, phases dispatched, what's
done vs outstanding, current status, and exactly what (if anything) the human
must decide.

## Guardrails

- **Compose, never reimplement** the existing skills/agents (`/design`,
  `/code`, `/test`, `/code-review`, `/verify`, `/pr`, `/pr-watch`,
  `/full-work`).
- **Keep the parent context lean** — phases return compact structured artifacts,
  not transcripts. Diffs/logs stay in the fork or `Workflow` run.
- **`goal` is immutable**; the `seen` set is cumulative. These two invariants
  are what keep the loop on-goal and terminating.
- **Implementation is worktree-isolated** (via `/code`); the main tree is never
  edited directly by a phase.
- **Confirm-first** on destructive/outward actions (force-push, merging *now*,
  closing PRs/issues, anything irreversible) — a budget backstop never licenses
  skipping this.
- **Honest reporting** — a budget/iteration-cap stop is not a success; say so.
  If a phase failed, surface it with the cause.
- Follow `agent-os/standards/` throughout (phase bodies enforce this). Don't add
  Claude/AI attribution to commits or PRs unless asked.
