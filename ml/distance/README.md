# ml/distance

Distance metrics for vectors and sequences — lower values mean items are more alike.
This package complements [`ml/similarity`](../similarity/), which uses the inverse convention.

All functions are **DISTANCE** metrics (lower = closer, 0 = identical).

## Continuous Vector Distances

Accept any `constraints.Numeric` element type and return `(float64, bool)`.
Returns `ok == false` for empty or mismatched-length inputs.
Non-finite inputs (NaN/Inf) propagate with `ok == true`.

| Function | Formula | Notes |
|----------|---------|-------|
| `Euclidean[T](a, b []T)` | √(Σ (aᵢ−bᵢ)²) | L2 straight-line distance |
| `Manhattan[T](a, b []T)` | Σ \|aᵢ−bᵢ\| | L1 taxicab distance |
| `Minkowski[T](a, b []T, p)` | (Σ \|aᵢ−bᵢ\|ᵖ)^(1/p) | Generalises L1 and L2; `p<1` → `ok=false` |
| `CosineDistance[T](a, b []T)` | 1 − CosineSimilarity(a, b) | In [0, 2]; zero-vector → `ok=false` |

```go
a := []float64{0, 0}
b := []float64{3, 4}

e, _ := distance.Euclidean(a, b)  // 5.0
m, _ := distance.Manhattan(a, b)  // 7.0
```

## Discrete Distances

| Function | Returns | Notes |
|----------|---------|-------|
| `Hamming[T](a, b []T)` | `(int, bool)` | Positional differences; any `comparable` T; `ok=false` on length mismatch |
| `Levenshtein(a, b string)` | `int` | Min edit operations (insert/delete/substitute); operates over runes |

```go
h, _ := distance.Hamming([]string{"a","b","c"}, []string{"a","x","c"}) // 1

lev := distance.Levenshtein("kitten", "sitting") // 3
```
