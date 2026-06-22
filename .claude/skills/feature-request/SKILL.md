---
name: feature-request
description: Author a house-style enhancement issue and file it. Frame the use-case, sketch an API surface consistent with house idioms (the (value, ok) contract, generic-first, non-mutating standard), note alternatives/trade-offs, cross-reference related packages/issues/standards, then open it via gh issue create with the enhancement label. Flags RFC-sized directions that should be a discussion, not one issue. Use when the user says "/feature-request <description>", "propose a feature for …", or wants a capability turned into a clean issue. Sits upstream of /understand — its output feeds /full-work directly.
---

# /feature-request — author and file an enhancement issue

The step *before* [`/understand`](../understand/SKILL.md): turn a capability idea
into a **well-formed GitHub issue** framed by use-case and house idioms, clean
enough that `/understand <new#>` → `/full-work` can consume it immediately. The
body mirrors the
[`.github/ISSUE_TEMPLATE/feature_request.md`](../../../.github/ISSUE_TEMPLATE/feature_request.md)
form so skill- and web-filed issues look the same.

## Invocation

```
/feature-request [description of the capability]
```

## What it does

Spawn an `Explore` agent (read-only) to ground the proposal in the existing code
before filing from this context.

1. **Duplicate guard** — search open issues first and surface likely duplicates
   rather than blindly filing:
   ```bash
   gh issue list --search "<keywords>" --state open
   ```
2. **Ground it in the codebase** — have the `Explore` agent find the packages
   this composes with, the neighbouring idioms to match, and any relevant
   `agent-os/standards/`. The API sketch should fit what exists, not invent a
   parallel style.
3. **Scope guard** — if this is really an RFC-sized *direction* (a whole package
   or research lane, cf. #190 / #179) rather than one shippable change, say so:
   recommend it become a discussion/RFC, not a single issue, and stop to confirm
   with the user before filing.
4. **Pick labels from the live set** (query, don't hardcode): `gh label list` —
   apply `enhancement` (plus `feedback` if review-sourced).
5. **Emit the house-style body** and file it (see below).

## Issue body — house style

Match the existing enhancement issues and the template. Keep the API sketch
consistent with house idioms:

- the **`(value, ok)`** contract for partial/undefined results,
- **generic-first** (`[T comparable]` / numeric constraints) over type-specific,
- the **non-mutating standard** (package-level funcs allocate fresh, never write
  back into arguments — see `agent-os/standards/functional/non-mutating.md`).

```markdown
## Context / motivation
<the use-case: what problem this solves and who hits it>

## Proposed API
```go
<signature sketch in house idiom>
```

## Alternatives / trade-offs
<other shapes considered and why this one; note numeric/NaN-policy or
 thread-safety implications if relevant>

## Related
<existing packages this composes with, related issues/PRs, applicable
 agent-os/standards/; flag if RFC-sized>
```

## Filing

Open the issue and return its URL:

```bash
gh issue create --title "<pkg>: <concise capability>" --label enhancement[,…] --body "<body>"
```

Title convention: `<pkg>: <concise capability>` (e.g. `stats: add weighted
mean`). For a brand-new package, use the proposed package name.

## Output

Return a compact confirmation, not a transcript:

```json
{
  "url": "https://github.com/.../issues/NNN",
  "title": "stats: …",
  "labels": ["enhancement"],
  "verdict": "filed | rfc-sized (not filed)",
  "related": ["#NNN", "agent-os/standards/…"],
  "duplicates_checked": true
}
```

## Guardrails

- **House idioms** — the API sketch follows `(value, ok)`, generic-first, and
  non-mutating; don't propose a mutating or type-specific surface without saying
  why.
- **Flag RFC-sized scope** — don't force a research direction into one issue.
- **Don't blind-file** — run the duplicate guard first.
- **Don't add Claude/AI attribution** to the issue body.
