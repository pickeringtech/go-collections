---
name: pr
description: Take the current working tree to an open PR — commit if dirty, branch off main when needed, push, create the PR via gh, then hand off to /pr-watch to babysit it to green. Use when the user says "open a PR", "/pr", "raise a PR", or wants their local changes turned into a reviewed, merge-ready pull request.
---

# /pr — open a pull request from the current work

Take whatever is in the working tree to an **open PR**, then hand off to
[`/pr-watch`](../pr-watch/SKILL.md) to drive it to green. Run the steps in
order; skip any that are already satisfied.

## Preconditions

- `gh` is authenticated (`gh auth status`). If not, stop and tell the user to
  run `! gh auth login`.
- Inside a git repo with an `origin` remote.

## Step 1 — Ensure a topic branch (never commit to `main`)

```bash
git rev-parse --abbrev-ref HEAD
```

If the current branch is `main` (or `master`), create a topic branch first.
Derive a short kebab-case `<summary>` from the change (the diff/intended work),
and prefix it with the current contributor's handle so branches stay namespaced
per author rather than hard-coding any one person's prefix:

```bash
PREFIX=$(gh api user -q .login 2>/dev/null || git config user.name | tr 'A-Z ' 'a-z-')
git switch -c "${PREFIX}/<summary>"
```

If already on a topic branch, stay on it.

## Step 2 — Commit if the tree is dirty

```bash
git status --porcelain
```

- **Clean** → skip to Step 3.
- **Dirty** → stage and commit. Synthesize the message from the actual diff,
  not the filenames:

```bash
git add -A
git diff --cached --stat        # understand what changed
git commit -m "<subject derived from the diff>"
```

Follow the repo's commit conventions (look at `git log` for tone/format). Do
**not** add Claude/AI attribution unless the user asks.

## Step 3 — Push and set upstream

```bash
git push -u origin HEAD          # normal push; sets upstream if missing
```

Never force-push here — this is a fresh topic branch. (Force-push is only ever
considered later, during an explicit rebase, and only behind the user-confirm
guardrail.)

## Step 4 — Create the PR

First check one doesn't already exist for this branch:

```bash
gh pr view --json number,url -q '.number' 2>/dev/null
```

If none, create it. Synthesize the title and body from the commits and diff —
title is a concise summary; body covers **what changed and why**, plus a short
test/verification note. Use the repo's standards (e.g. the testing trio) as the
bar for "done" when describing coverage.

```bash
gh pr create --base main --head "$(git branch --show-current)" \
  --title "<summary>" \
  --body  "<what changed and why; how it was verified>"
```

Capture the PR number from the output.

## Step 5 — Hand off to the responder

Invoke [`/pr-watch`](../pr-watch/SKILL.md) on the freshly opened PR so `/pr`
alone takes the user from local changes → open PR → babysat to green:

> Now running `/pr-watch <pr-number>` to drive the PR to green.

Default iteration cap applies (15) unless the user specified otherwise.

## Guardrails

- Branch off `main` rather than committing to it — always.
- Normal pushes only. Anything that rewrites pushed history (force-push) needs
  explicit user confirmation.
- If the diff is large or ambiguous and you can't write an honest title/body,
  ask the user rather than guessing.
