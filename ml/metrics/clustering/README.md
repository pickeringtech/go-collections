# Clustering Metrics

The `clustering` package scores the quality of a clustering **without**
ground-truth labels — currently the silhouette coefficient. It is part of the
[`ml/metrics`](../) family. Unlike [`classification`](../classification), which
needs true labels, silhouette judges a clustering purely by how compact and
well-separated its clusters are.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/ml/metrics/clustering"

points := [][]float64{{0, 0}, {0.5, 0}, {10, 0}, {10.5, 0}}
labels := []int{0, 0, 1, 1}

score, _ := clustering.SilhouetteScore(points, labels)    // 0.9500 — tight, separated
samples, _ := clustering.SilhouetteSamples(points, labels)
```

## Silhouette

For each sample the coefficient is `(b − a) / max(a, b)`, where `a` is the mean
distance to the other points of its own cluster and `b` the mean distance to the
nearest other cluster. It lies in `[−1, 1]`: near 1 the point is well inside its
cluster, near 0 it is on a boundary, negative it is probably mis-assigned.

| Function                                  | Returns             | Notes                              |
| ----------------------------------------- | ------------------- | ---------------------------------- |
| `SilhouetteScore(points, labels)`         | `(float64, bool)`   | mean coefficient over all samples  |
| `SilhouetteSamples(points, labels)`       | `([]float64, bool)` | per-point coefficients             |
| `SilhouetteScoreWith(…, dist)`            | `(float64, bool)`   | custom `DistanceFunc`              |
| `SilhouetteSamplesWith(…, dist)`          | `([]float64, bool)` | custom `DistanceFunc`              |

The default metric is `EuclideanDistance`; pass any `DistanceFunc` to the
`…With` variants.

## Conventions

Returns follow the `(result, ok)` idiom; `ok` is `false` (result zero) when the
input cannot be summarised: fewer than two points, a label slice of the wrong
length, ragged coordinate rows, any non-finite coordinate, or a cluster count
outside `[2, n−1]` (silhouette is undefined for one cluster or for
one-point-per-cluster). For the `…With` variants `ok` is also `false` when the
supplied `DistanceFunc` is `nil` or returns a non-finite or negative distance,
which would otherwise void the `[−1, 1]` guarantee. A lone point in its cluster
is given a silhouette of 0.
Inputs are never mutated; the mean over samples routes through
[`stats.Mean`](../../../stats).
