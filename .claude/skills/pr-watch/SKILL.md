---
name: pr-watch
description: Autonomously drive an open PR to green by reacting to whichever fires first — a CI failure, a new review comment (Copilot/bot/human), or a base-branch merge conflict. Loops up to N iterations (default 15), dedups handled comments, detects stuck loops, and stops on success/cap/human-needed. Use when the user says "/pr-watch", "/pr-babysit", "babysit the PR", "watch the PR", or "get this PR to green".
---

# /pr-watch — drive a PR to green (OR-race responder)

Watch an open PR and react **autonomously** to whichever of three signals
occurs first, then loop:

1. **CI failure** — a required check fails (includes the coverage floor and,
   once configured, `codecov/*`).
2. **New review comment** — from Copilot, any other bot, or a human.
3. **Merge conflict / behind base** — the branch conflicts with or has fallen
   behind `main`.

**First-to-fire wins.** Across iterations, act on whichever signal *appears
first in time* — do not wait for CI to finish before handling a comment, and do
not ignore a CI failure while waiting on comments. **Tie-break:** if more than
one signal is already present in the *same* poll, handle them in the priority
order below (new comments → conflict → behind base → CI failure). That ordering
is only a deterministic tie-break for a single poll; it never overrides
first-to-fire across polls.

## Arguments & pacing

- Optional integer = max iterations (e.g. `/pr-watch 25`). **Default 15.**
- Optional PR number; otherwise use the PR for the current branch.
- This skill **composes the built-in `/loop` skill's dynamic (self-paced)
  pacing** (`/loop` ships with Claude Code; it is not defined in this repo) — drive the
  cadence with `ScheduleWakeup` rather than a blocking sleep. Use a short delay
  (**30–60s**) while CI is in progress so comments are picked up quickly; stop
  scheduling once a terminal condition is reached.

## Setup

```bash
PR=<number or: gh pr view --json number -q .number>
gh pr view "$PR" --json number,headRefName,baseRefName,url
```

Initialize run state (in-memory across iterations):
- `seen_comment_ids = {}` — comment IDs already handled.
- `iteration = 0`, `max = arg or 15`.
- `last_failure_signature = none` — for loop detection.

## Each iteration — poll all three signals, then act on the first present

### Poll: CI status (one call, all checks)

```bash
gh pr view "$PR" --json statusCheckRollup \
  -q '.statusCheckRollup[] | {name, status, conclusion}'
```
- **Failed** if any element has `status == COMPLETED` and `conclusion` in
  `{FAILURE, TIMED_OUT, STARTUP_FAILURE, ...}`.
- **Still running** if any element has `status != COMPLETED`.
- **Green** if all `COMPLETED` with `conclusion == SUCCESS` (or `NEUTRAL`).

### Poll: new review comments (any author)

```bash
OWNER_REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)
gh api "repos/$OWNER_REPO/pulls/$PR/comments"   # inline review comments
gh api "repos/$OWNER_REPO/issues/$PR/comments"  # PR-level comments
gh pr view "$PR" --json reviews -q '.reviews[] | {author:.author.login, state, body}'
```
Filter to comments whose `id` is **not** in `seen_comment_ids`. Copilot appears
as `Copilot` (inline) / `copilot-pull-request-reviewer` (review summary); bots
as `*[bot]`. **Treat every author as a trigger** — bot and human alike.

### Poll: mergeability (retry until resolved)

```bash
gh pr view "$PR" --json mergeable,mergeStateStatus,baseRefName,headRefName
```
- `mergeable == CONFLICTING` or `mergeStateStatus == DIRTY` → **conflict**.
- `mergeStateStatus == BEHIND` → **behind base**, update it.
- **Gotcha:** these are computed lazily and often return `UNKNOWN` right after a
  push. **Poll/retry until the value resolves** — never treat `UNKNOWN` as
  conflict-free.

### Act on the first present signal (priority order)

```
if new comments:        handle each → commit/push if code changed → continue
elif merge conflict:    fetch base, merge, resolve, commit, push   → continue
elif behind base:       update branch from base, push             → continue
elif ci failed:         diagnose, fix, commit, push → iteration++ → continue
elif ci still running:  ScheduleWakeup 30–60s, re-poll
else:                   green + mergeable + no unhandled comments → DONE
```

## Remediation playbooks

### New review comment
Read it. Address it with a code change and/or a threaded reply, then **mark the
`id` seen**. If code changed, commit + push (re-triggers CI). Reply so the
thread is clearly handled:

```bash
gh api "repos/$OWNER_REPO/pulls/$PR/comments" -F in_reply_to=<comment-id> -f body="<reply>"
```
Resolving a conversation thread is **GraphQL-only** (`resolveReviewThread`) and
counts as an outward/destructive action — **a threaded reply is sufficient for
v1; resolving follows the confirm-first guardrail.**

### CI failure
```bash
RUN=$(gh run list --branch <headRefName> --limit 1 --json databaseId -q '.[0].databaseId')
gh run view "$RUN" --log-failed     # only the failed steps
```
Diagnose from the log, fix the **code**, commit, push.

**Coverage failures are special — add tests, never weaken logic:**
- The **Test (race + coverage)** job enforces a local `COVERAGE_MIN` floor
  (currently 90%). A drop fails CI.
- Once [#14](https://github.com/pickeringtech/go-collections/issues/14) lands,
  `codecov/project` (total) and `codecov/patch` (changed lines) become real
  checks too. **Until then Codecov is informational only** (non-gating, no
  `codecov.yml`).
- Remediate by **adding/extending tests** per the repo testing standards
  (`agent-os/standards/testing/`): the Example + table-driven Test + Benchmark
  trio, nil→empty / empty→empty edge cases, and a fuzz target for collection &
  transformation types. Find uncovered lines from `go test -coverprofile` /
  `codecov/patch` annotations / the Codecov comment. Then commit + push.

### Merge conflict / behind base
```bash
git fetch origin <baseRefName>
git merge origin/<baseRefName>      # push-safe, non-rewriting — the default
# resolve conflicted files on their merits, then:
git add -A && git commit
git push
```
- **Prefer merge** (push-safe). Use **rebase** only when the repo explicitly
  wants linear history — it needs force-push, so it's behind the force-push
  guardrail.
- Resolve conflicts on their **merits** (understand both sides). If a conflict
  is genuinely ambiguous or risks losing intent, **stop and escalate** — don't
  guess.

## State, dedup & loop detection

- **Dedup comments by `id`** (stable), not body. Codecov-style bots **edit their
  comment in place**, so deduping by `id` alone would re-trigger on every edit —
  key such bot comments off `updated_at`/a content hash, or prefer the
  `codecov/*` **status checks** as the trigger and treat the comment as a data
  source.
- **Loop detection:** record a signature of each CI failure (check name + root
  cause). If the **same failure recurs two iterations running with no diff
  progress**, stop and escalate rather than looping.

## Termination

Stop and report when any holds:
- **Success** — CI green **and** mergeable **and** no unhandled comments.
- **Cap** — `iteration == max`.
- **Stuck** — repeated identical failure with no progress.
- **Human-needed** — ambiguous conflict, a change needing a product decision, or
  a required outward action behind a guardrail.

Surface a short summary: what was handled, current CI/mergeability state, and
exactly what (if anything) the human needs to decide.

## Guardrails

- Autonomous in scope: code fixes, adding tests, threaded replies, push-safe
  merges, normal pushes.
- **Confirm-first** (do not do autonomously): force-push, rebasing pushed
  history, resolving/closing conversations, closing or merging the PR.
- Never weaken logic or assertions to pass a coverage gate — write tests.
