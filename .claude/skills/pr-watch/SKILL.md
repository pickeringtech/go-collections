---
name: pr-watch
description: Autonomously drive an open PR to merge-ready by reacting to whichever fires first — a CI failure, a new review comment (Copilot/bot/human), or a base-branch merge conflict. On green it keeps the branch current with base and arms auto-merge so the PR lands before it can go stale. Loops up to N iterations (default 15), dedups handled comments, detects stuck loops, and stops on success/cap/human-needed. Use when the user says "/pr-watch", "/pr-babysit", "babysit the PR", "watch the PR", or "get this PR to green".
---

# /pr-watch — drive a PR to merge-ready (OR-race responder)

Watch an open PR and react **autonomously** to whichever of these signals
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

**Before finalizing, wait for automated review.** Arming auto-merge is **not** a
race signal — it is the *last* step, gated on the automated reviewer (Copilot)
having finished reviewing the current head commit **and** all its comments being
handled. CI going green is necessary but not sufficient: Copilot's review is
asynchronous and is **not** a required status check, so auto-merge (which only
waits for required checks) will otherwise merge the PR before Copilot reviews or
before its comments are addressed. See "Gate before finalizing" below.

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

## Each iteration — poll all signals, then act on the first present

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

### Poll: automated review status (is Copilot still reviewing?)

```bash
# Pending if the automated reviewer is still in the requested-reviewers list.
gh pr view "$PR" --json reviewRequests -q '.reviewRequests[].login'
# Has it reviewed at all yet?
gh pr view "$PR" --json reviews \
  -q '[.reviews[] | select(.author.login=="copilot-pull-request-reviewer")] | length'
```
- **Review pending** if `copilot-pull-request-reviewer` is in `reviewRequests`,
  **or** it has been requested but has not yet submitted a review for the current
  head commit. GitHub re-requests Copilot on each new push (when auto-review is
  enabled), so this re-arms after every commit you push.
- **Not applicable** if no automated reviewer is configured/requested at all —
  then there is nothing to wait for; don't block.
- This is a **finalization gate**, not an OR-race signal: it only governs whether
  it's safe to arm auto-merge, never preempts handling CI/comments/conflicts.

### Act on the first present signal (priority order)

```
if new comments:           handle each → commit/push if code changed → continue
elif merge conflict:       fetch base, merge, resolve, commit, push   → continue
elif behind base:          update branch from base, push             → continue
elif ci failed:            diagnose, fix, commit, push → iteration++ → continue
elif ci still running:     ScheduleWakeup 30–60s, re-poll
elif review pending:       ScheduleWakeup 30–60s, re-poll   # give Copilot time to finish
else:  # green + no unhandled comments + automated review settled
        update branch from base (stay current) → if it conflicts, run the
            conflict playbook and continue
        enable auto-merge so the PR lands the moment it's ready → DONE(auto-merge armed)
```

Reaching green is **not** where we stop and walk away — that's exactly when a PR
starts silently going stale as other PRs merge into `main`. Instead, make the
PR *land itself*: bring it up to date with base and arm auto-merge, so there is
no idle window in which an unwatched conflict can accumulate.

**But green is also not where we finalize.** The `review pending` branch above is
what gives Copilot its chance: when CI is green but Copilot hasn't finished
reviewing the current head commit, **wait** — do not arm auto-merge yet. When
Copilot then posts comments, they arrive via the `new comments` arm and get
handled (which pushes a fix, re-triggering Copilot's review, so we wait again).
Only when CI is green, the automated review has settled, and no unhandled
comments remain do we update + arm auto-merge.

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

### Gate before finalizing — let the automated review finish

Before doing anything below, confirm the automated review has settled (see the
"Poll: automated review status" step). **While Copilot is still requested /
hasn't reviewed the current head commit, do not arm auto-merge** — `ScheduleWakeup`
and re-poll. When its comments arrive, handle them via the new-comment playbook
(which pushes a fix and re-triggers the review, so you wait again). Proceed only
once: CI green **and** automated review settled **and** no unhandled comments.

Why this matters: Copilot's review is asynchronous and is **not** a required
status check, so auto-merge (which waits only for required checks) will merge the
PR out from under an in-flight review — leaving its comments stranded on a merged
PR, exactly what we want to avoid.

### Green — keep current and auto-merge (close the staleness window)

Once the gate above is satisfied, **don't just terminate** — a PR that merely
*is* mergeable now will conflict later as other PRs land. Make it land itself:

```bash
# 1. Bring the branch up to date with base (GitHub-side merge of base → head).
gh pr update-branch "$PR"          # no-op if already current
# If update-branch reports a conflict, fall through to the merge-conflict
# playbook above (resolve locally, push), then retry.
#
# GOTCHA: if the update would touch a workflow file (.github/workflows/*) and
# the gh token lacks the `workflow` OAuth scope, this fails with
# "refusing to allow an OAuth App to create or update workflow ... without
# workflow scope". `git push` over SSH is NOT scope-gated, so fall back to a
# local merge instead of giving up (uses <baseRefName> from Setup):
#   git fetch origin <baseRefName> && git merge origin/<baseRefName> && git push
# (resolve any conflicts on their merits, as in the conflict playbook).

# 2. Arm auto-merge so the PR merges the instant it is green + up to date.
#    Pick the merge method GitHub reports as the repo default
#    (gh repo view --json viewerDefaultMergeMethod -q .viewerDefaultMergeMethod);
#    --squash here is a sensible fallback when no default is detected.
gh pr merge "$PR" --auto --squash
```

- **Auto-merge is now in-scope** for this skill (the user opted into it): arming
  it is allowed without a per-PR confirmation. It does **not** merge immediately
  — GitHub merges only once required checks pass and the branch is mergeable.
- **Preconditions:** auto-merge must be enabled in repo settings (and typically
  a protected base with required checks). If `gh pr merge --auto` errors with
  auto-merge-not-allowed, report it once and fall back to terminating on
  green+mergeable (the old behaviour) rather than looping.
- **Proactive update on every pass:** also run `gh pr update-branch` whenever the
  PR is `BEHIND` (not only at green), so it never drifts far from base.
- **`workflow`-scope fallback:** `gh pr update-branch` (and any `gh` write to a
  `.github/workflows/*` file) needs the `workflow` OAuth scope; without it the
  call is refused. Don't treat that as "can't update" — `git push` over SSH is
  not scope-gated, so update via local `git merge origin/<base>` + push instead.
- **Residual gap (be honest):** if the PR can't auto-complete (e.g. a required
  human review is missing) it can still go stale after the watcher exits. For
  full coverage of "conflicts that appear long after I've moved on," a scheduled
  re-arm is the complete fix; auto-merge shrinks the window but doesn't remove it
  entirely.

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
- **Success** — CI green, **automated review settled**, no unhandled comments,
  branch up to date with base, and the PR is **set to land**: either auto-merge
  armed, or — when auto-merge is not permitted (see the Green playbook's
  precondition) — green + mergeable after one last update. Do not exit on "green +
  mergeable" *without* first letting the review finish, handling its comments, and
  updating the branch — those are the gaps that merge too early or let later
  conflicts go unhandled.
- **Cap** — `iteration == max`.
- **Stuck** — repeated identical failure with no progress.
- **Human-needed** — ambiguous conflict, a change needing a product decision, or
  a required outward action behind a guardrail.

Surface a short summary: what was handled, current CI/mergeability state, and
exactly what (if anything) the human needs to decide.

## Guardrails

- Autonomous in scope: code fixes, adding tests, threaded replies, push-safe
  merges, normal pushes, **updating the branch from base (`gh pr update-branch`)
  and arming auto-merge (`gh pr merge --auto`)**.
- **Confirm-first** (do not do autonomously): force-push, rebasing pushed
  history, resolving/closing conversations, an **immediate** `gh pr merge` (i.e.
  merging *now* rather than arming auto-merge to land when ready), or closing the
  PR.
- Never weaken logic or assertions to pass a coverage gate — write tests.
