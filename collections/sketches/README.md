# collections/sketches

Probabilistic data-sketching structures — compact summaries that trade exact
answers for very low memory and sub-linear time, useful when exact computation
over large sets or streams is prohibitive.

This package holds **MinHash** (set similarity, below). The **streaming
sketches** each live in their own sub-package:

| Package                  | Question it answers                       | Memory driver        |
| ------------------------ | ----------------------------------------- | -------------------- |
| [`bloom`](./bloom)       | "Have I seen this before?" (membership)   | items × target rate  |
| [`countmin`](./countmin) | "How often have I seen this?" (frequency) | target accuracy only |
| [`hll`](./hll)           | "How many distinct things?" (cardinality) | precision only       |

> A streaming-quantiles sketch (t-digest) is co-designed with the `stats`
> quantile work and lives with that package, not here.

The streaming sketches share a common design: bounded, configurable accuracy
with documented error bounds; seeded, pluggable hashing (`WithSeed`/`WithHasher`)
so behaviour is reproducible; `Merge` for parallel and distributed aggregation;
and a delegating `sync.RWMutex`-guarded `Concurrent` variant alongside a plain
type that is **not** safe for concurrent use.

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

## MinHash

MinHash estimates the Jaccard similarity between two sets without storing either
set explicitly. It maintains a compact signature of `uint64` minimums — one per
hash function — and compares signatures element-wise to approximate the Jaccard
coefficient.

```go
import (
    "github.com/pickeringtech/go-collections/collections/sketches"
)

// nil rng → deterministic default seed (comparable within the same process run).
a := sketches.NewMinHash[string](128, nil)
b := sketches.NewMinHash[string](128, nil)

for _, word := range []string{"the", "quick", "brown", "fox"} {
    a.Add(word)
}
for _, word := range []string{"the", "lazy", "brown", "dog"} {
    b.Add(word)
}

est, ok := sketches.EstimatedJaccard(a, b)
// ok == true; est ≈ Jaccard({the,quick,brown,fox}, {the,lazy,brown,dog}) = 2/6
```

### Accuracy

Accuracy improves with `numHashes`:
- 128 hashes → ~±7% error (95% confidence)
- 256 hashes → ~±5% error
- 512 hashes → ~±3% error

### Seeding and Reproducibility

Passing `nil` as the `rng` selects a fixed deterministic default. Two sketches with the same `numHashes` and the same `rng` (or both `nil`) produce the same permutation family within a process run and can be compared via `EstimatedJaccard`. Sketches built from different `rng` instances are **not** comparable and `EstimatedJaccard` returns `ok == false`.

Note: element hashing uses `maphash.MakeSeed()`, which is randomized per process. Sketches are therefore only comparable within the same program run — cross-process or serialized comparison requires a stable hashing strategy not yet provided by this package.

### Goroutine Safety

`MinHash` is **NOT goroutine-safe**. Do not call `Add` or `Signature` concurrently on the same instance. A `ConcurrentMinHash` variant is planned for a later issue.

### API

| Symbol | Description |
|--------|-------------|
| `NewMinHash[T](numHashes int, rng *rand.Rand) *MinHash[T]` | Create a new empty sketch |
| `(*MinHash[T]).Add(element T)` | Add an element to the sketch |
| `(*MinHash[T]).Signature() []uint64` | Return a copy of the current signature |
| `EstimatedJaccard[T](a, b *MinHash[T]) (float64, bool)` | Estimate Jaccard similarity |

For exact set similarity see [`ml/similarity`](../../ml/similarity/).
