---
name: Feature request
about: Propose a new API or capability, framed by use-case and house idioms
title: ""
labels: enhancement
assignees: ""
---

## Context / motivation

<!-- The use-case. What problem does this solve, and who hits it? -->

## Proposed API

<!--
Sketch the surface. Keep it consistent with house idioms:
- the `(value, ok)` contract for partial/undefined results
- generic-first (`[T comparable]` / `[T constraints.Numeric]`) over type-specific
- the non-mutating standard (package-level funcs allocate fresh; never write back
  into arguments) — see agent-os/standards/functional/non-mutating.md
-->

```go

```

## Alternatives / trade-offs

<!-- Other shapes you considered, and why this one. Note numeric/NaN-policy or
     thread-safety implications if relevant. -->

## Related

<!--
Existing packages this composes with, related issues/PRs, and any
agent-os/standards/ that apply. If this is really an RFC-sized direction
(cf. #190 / #179), say so — it may belong as an RFC rather than one issue.
-->
