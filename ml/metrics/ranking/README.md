# Ranking Metrics

The `ranking` package scores ordered result lists — the information-retrieval
metrics NDCG (with its DCG building block) and mean average precision. It is
part of the [`ml/metrics`](../) family. Unlike
[`classification`](../classification), which scores an unordered set of
predictions, ranking cares about **order**: a relevant item near the top is
worth more than the same item further down.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/ml/metrics/ranking"

// Rank items by model score; judge against their true relevance grades.
trueRelevance := []float64{3, 2, 3, 0, 1, 2}
scores := []float64{6, 5, 4, 3, 2, 1}
ndcg, _ := ranking.NDCG(trueRelevance, scores, 0) // 0.9608

ap, _ := ranking.AveragePrecision([]bool{true, false, true, false, true}) // 0.7556
```

## Metrics

| Function                         | Returns           | Undefined (`ok == false`) when                         |
| -------------------------------- | ----------------- | ------------------------------------------------------ |
| `DCG(relevances, k)`             | `(float64, bool)` | empty, or any non-finite relevance                     |
| `NDCG(trueRelevance, scores, k)` | `(float64, bool)` | empty, differing lengths, non-finite term, zero ideal DCG |
| `AveragePrecision(ranked)`       | `(float64, bool)` | empty, or no relevant items                            |
| `MeanAveragePrecision(queries)`  | `(float64, bool)` | empty, or any query undefined                          |

- **DCG** is `Σ relᵢ / log₂(i+2)` over the first `k` items; later positions are
  discounted. A cutoff `k <= 0` (or larger than the list) scores the whole list.
- **NDCG** normalises a score-ranked DCG by the ideal DCG, giving `[0, 1]` (1 is
  a perfect ranking). Items are ordered by `scores` (highest first), ties keeping
  their input order.
- **AveragePrecision** averages precision-at-k over the relevant positions of a
  single ranked list. `ranked[i]` says whether rank `i+1` is relevant.
- **MeanAveragePrecision** averages `AveragePrecision` across queries.

## Conventions

Returns follow the `(result, ok)` idiom; `ok` is `false` (result zero) for input
that cannot be summarised. A degenerate query makes `MeanAveragePrecision` as a
whole undefined, so the mean is always over a well-defined set. Inputs are never
mutated; the mean over queries routes through [`stats.Mean`](../../../stats).
