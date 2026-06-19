# sketches

Probabilistic data structures — compact **sketches** that answer questions about
massive streams in bounded memory, trading a little accuracy for a lot of space.

| Package                | Question it answers                       | Memory driver        |
| ---------------------- | ----------------------------------------- | -------------------- |
| [`bloom`](./bloom)     | "Have I seen this before?" (membership)   | items × target rate  |
| [`countmin`](./countmin) | "How often have I seen this?" (frequency) | target accuracy only |
| [`hll`](./hll)         | "How many distinct things?" (cardinality) | precision only       |

> A streaming-quantiles sketch (t-digest) is co-designed with the `stats`
> quantile work and lives with that package, not here.

## Shared design

Every sketch in this family follows the same conventions:

- **Bounded, configurable accuracy.** Construction takes the error target
  (false-positive rate, `epsilon`/`delta`, or precision); the documented error
  bound and the memory cost follow from it. Constructors return an `error`
  (rooted at each package's `ErrInvalidConfig`) rather than panic.
- **Seeded, pluggable hashing.** Hashing is deterministic from a stored seed, so
  behaviour is reproducible. `WithSeed` varies it; `WithHasher` overrides it for
  custom key types or cross-process reproducibility.
- **Mergeable.** `Merge` combines two compatible sketches into one covering the
  union of their inputs — the basis for parallel and distributed aggregation.
  Compatibility (matching configuration and seed) is checked.
- **A concurrent variant.** Each plain type owns all the logic; the matching
  `Concurrent` type wraps it with a `sync.RWMutex` and delegates. The plain
  types are **not** safe for concurrent use.

## Quick comparison

```go
// Membership — no false negatives, ~1% false positives.
bf, _ := bloom.New[string](1_000_000, 0.01)
bf.Add("alice"); bf.Contains("alice") // true

// Frequency — never under-counts.
cm, _ := countmin.New[string](0.001, 0.01)
cm.Add("/index"); cm.Estimate("/index") // ~1

// Cardinality — ~0.8% error in ~16 KB.
hl, _ := hll.New[string]()
hl.Add("visitor-1"); hl.Count() // ~1
```
