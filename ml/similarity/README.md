# ml/similarity

Similarity metrics for vectors and sets — higher values mean items are more alike.
This package complements [`ml/distance`](../distance/), which uses the inverse convention.

## Vector Similarity

| Function | Description | Returns |
|----------|-------------|---------|
| `DotProduct[T](a, b []T)` | Inner product Σ aᵢ·bᵢ | `(float64, bool)` |
| `CosineSimilarity[T](a, b []T)` | Cosine of the angle in [−1, 1] | `(float64, bool)` |

Both functions accept any `constraints.Numeric` element type and return `ok == false` for empty or mismatched-length inputs. `CosineSimilarity` also returns `ok == false` for zero-magnitude vectors. Non-finite inputs (NaN/Inf) propagate with `ok == true`.

```go
a := []float64{1, 2, 3}
b := []float64{2, 4, 6}

cos, ok := similarity.CosineSimilarity(a, b) // 1.0, true — same direction
```

## Set Similarity

| Function | Formula | Returns |
|----------|---------|---------|
| `Jaccard[T](a, b Set[T])` | \|A∩B\| / \|A∪B\| | `(float64, bool)` |
| `Dice[T](a, b Set[T])` | 2\|A∩B\| / (\|A\|+\|B\|) | `(float64, bool)` |
| `Overlap[T](a, b Set[T])` | \|A∩B\| / min(\|A\|,\|B\|) | `(float64, bool)` |

All three compose `Intersection` and `Length` from the `collections/sets.Set[T]` interface — no set algebra is reimplemented. Like the vector similarities they return `ok == false` when the coefficient is undefined (a zero denominator): both sets empty for `Jaccard` and `Dice`, either set empty for `Overlap`. That keeps an undefined result distinct from the computed `0` of two disjoint sets.

```go
s1 := sets.NewHash("a", "b", "c", "d")
s2 := sets.NewHash("b", "c", "d", "e")

j, ok := similarity.Jaccard(s1, s2)  // 3/5 = 0.6, true
d, ok := similarity.Dice(s1, s2)     // 6/8 = 0.75, true
o, ok := similarity.Overlap(s1, s2)  // 3/4 = 0.75, true

_, ok = similarity.Jaccard(sets.NewHash[string](), sets.NewHash[string]()) // false: 0/0
```

For probabilistic set-similarity over large datasets consider [`collections/sketches.MinHash`](../../collections/sketches/).
