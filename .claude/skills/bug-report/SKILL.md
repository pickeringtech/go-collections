---
name: bug-report
description: Author a house-style bug issue and file it. Locate the offending code with an Explore agent and cite exact file:line, confirm the bug where feasible (build the failing snippet / reason the edge case) and label it BUILD-CONFIRMED vs SUSPECTED honestly, then open it via gh issue create. Use when the user says "/bug-report <description>", "file a bug for …", or wants a defect turned into a clean issue. Sits upstream of /understand — its output should be consumable by /understand <new#> immediately.
---

# /bug-report — author and file a bug issue

The step *before* [`/understand`](../understand/SKILL.md): turn a defect
description into a **well-formed GitHub issue** that matches the house style, so
the delivery pipeline can pick it up directly (`/bug-report` → `/understand
<new#>` → `/full-work`). The body mirrors the
[`.github/ISSUE_TEMPLATE/bug_report.md`](../../../.github/ISSUE_TEMPLATE/bug_report.md)
form so issues filed by skill and by the web UI look the same.

## Invocation

```
/bug-report [description of the defect]
```

## What it does

Keep the parent context lean: spawn an `Explore` agent (read-only fan-out, like
`/understand`) to locate the code, then file the issue from this context.

1. **Duplicate guard** — before anything else, search open issues and surface
   likely duplicates rather than blindly filing:
   ```bash
   gh issue list --search "<keywords>" --state open
   ```
   If a strong match exists, stop and ask the user whether to add to it instead.
2. **Locate the offending code** — spawn an `Explore` agent to find the exact
   package/file and the lines at fault. Require it to return concrete
   `path/to/file.go:NN` citations and the offending snippet, not a summary.
3. **Confirm where feasible** — build the failing snippet, run the case, or
   reason through the edge case. In the body, label the finding honestly as
   **BUILD-CONFIRMED** (reproduced/compiled) or **SUSPECTED** (reasoned, not
   run). Never overstate — this discipline is the point of the skill.
4. **Pick labels from the live set** — query, don't hardcode:
   ```bash
   gh label list
   ```
   Apply `bug`, plus `github_actions` (CI/workflow defects) or `documentation`
   (doc defects) when they fit. Add `feedback` if the bug came from a review
   round.
5. **Emit the house-style body** and file it (see below).

## Issue body — house style

Match the existing bug issues (#224–#240) and the template:

```markdown
## Context
<one or two sentences: what's wrong and where it came from>

## Evidence
<BUILD-CONFIRMED or SUSPECTED, then the citation + fence>

```go
// path/to/file.go:NN
<offending code>
```

## Expected vs actual
- **Expected:** …
- **Actual:** …

## Fix
<proposed fix; if there's a real design choice, list options as "decide one of …">

## Related
<#issues / PRs / agent-os/standards/… that bear on this>
```

## Filing

Open the issue and return its URL:

```bash
gh issue create --title "<pkg>: <concise defect>" --label bug[,…] --body "<body>"
```

Title convention: `<pkg>: <concise defect>` (e.g. `stats: Median panics on
empty input`).

## Output

Return a compact confirmation, not a transcript:

```json
{
  "url": "https://github.com/.../issues/NNN",
  "title": "stats: …",
  "labels": ["bug"],
  "confirmation": "build-confirmed | suspected",
  "evidence": ["pkg/file.go:NN"],
  "duplicates_checked": true
}
```

## Guardrails

- **Verification honesty** — never label a finding BUILD-CONFIRMED unless you
  actually reproduced or compiled it.
- **Real citations** — every claim of a defect cites `file:line`; no vague "the
  X function is wrong".
- **Don't blind-file** — always run the duplicate guard first.
- **Don't add Claude/AI attribution** to the issue body.
