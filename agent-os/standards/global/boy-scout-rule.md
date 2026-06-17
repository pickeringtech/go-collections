# Boy Scout Rule

**Always leave the code cleaner than you found it.** When you touch a file to add a feature or fix a bug, fix the small things you notice on the way past — don't walk by a problem just because it isn't the task.

The codebase improves continuously through many small, opportunistic corrections rather than rare big-bang cleanups.

## What to clean up in passing

When editing a file, opportunistically bring what you touch into line with the existing standards:

- A method missing its test trio → add it ([[coverage-requirements]]).
- An `if init; cond` you have to read → split it ([[no-if-init-statement]]).
- A stale or fictional doc example → correct it ([[package-doc]]).
- A clarity win with no behaviour change → take it ([[readability-and-performance]]).
- A dead branch, a typo in a comment, an inconsistent name next to the line you're already changing → fix it.

Prefer the cleanup that the *next* reader of this file would thank you for.

## Guardrails — clean up without derailing the change

The rule is "leave it cleaner," not "rewrite it." Keep the boy-scout work proportionate and reviewable:

- **Stay near what you touched.** Tidy the function, type, or file you're already in. Don't go hunting three packages away.
- **Keep cleanups separable.** Put unrelated tidy-ups in their own commit (following the repo's `area: summary` convention) so the feature/fix diff stays readable and revertible. Never bury a behaviour change inside a "cleanup."
- **Never lower a bar.** A cleanup must not drop coverage, weaken a test oracle, or break an existing standard. Cleaner means *more* aligned with the standards, never less.
- **Don't expand the public API on a whim.** Adding a `Fast` variant, a new constructor, or a facade alias is a deliberate decision, not a drive-by ([[readability-and-performance]]).

## When it's too big to fix now

If you spot a real problem that's larger than a quick fix — a deadlock risk, a missing fuzz oracle, an API contract that needs a decision — **don't silently leave it, and don't balloon the current change to fix it.** Capture it: open an issue (or add a `TODO` with a tracking reference) so the campsite gets cleaned later by someone who can scope it properly. Surfacing the problem is itself leaving the code better than you found it.
