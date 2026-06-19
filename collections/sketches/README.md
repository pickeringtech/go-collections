# collections/sketches

Probabilistic data-sketching structures for approximate set operations. Sketches trade exact answers for very low memory and sub-linear time — ideal when exact set sizes or intersections are prohibitively large.

## MinHash

MinHash estimates the Jaccard similarity between two sets without storing either set explicitly. It maintains a compact signature of `uint64` minimums — one per hash function — and compares signatures element-wise to approximate the Jaccard coefficient.

```go
import (
    "github.com/pickeringtech/go-collections/collections/sketches"
)

// nil rng → deterministic default seed (reproducible across runs).
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

Passing `nil` as the `rng` selects a fixed deterministic default. Two sketches with the same `numHashes` and the same `rng` (or both `nil`) produce the same permutation family and can be compared via `EstimatedJaccard`. Sketches built from different `rng` instances are **not** comparable and `EstimatedJaccard` returns `ok == false`.

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
