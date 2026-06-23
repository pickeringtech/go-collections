# tdigest

A [t-digest](https://github.com/tdunning/t-digest): a streaming, mergeable sketch
that estimates quantiles (and the inverse CDF) of a stream of `float64` values in
memory bounded by a compression parameter, independent of how many values it has
seen.

```go
import "github.com/pickeringtech/go-collections/collections/sketches/tdigest"

d, _ := tdigest.New() // DefaultCompression (100)
for _, v := range latencies {
	d.Add(v)
}
p99, ok := d.Percentile(99) // estimated 99th-percentile latency
```

It is the approximate, bounded-memory counterpart of the exact
[`stats`](../../../stats) quantile functions: where `stats.Quantile` sorts the
whole sample, a `Digest` keeps only a small set of weighted centroids.

## The guarantee

- **Bounded memory.** The retained centroid count is roughly proportional to the
  compression parameter and **independent of the stream length** — a long-running
  stream never grows the digest.
- **Exact extremes.** `Quantile(0)` and `Quantile(1)` return the exact observed
  minimum and maximum; a single distinct value is returned exactly for any
  quantile.
- **Tail-accurate.** The scale function makes error **smallest at the tails**
  (q near 0 or 1, e.g. p99/p999) and largest in the middle — the opposite of a
  fixed histogram, and exactly what latency monitoring wants.

## Accuracy vs memory

`New()` uses `DefaultCompression` (100); `WithCompression(c)` trades memory for
accuracy.

| Compression `c` | Approx. centroids | Typical use                          |
| --------------- | ----------------- | ------------------------------------ |
| 20              | ~tens             | rough estimates, very small footprint |
| 100 (default)   | ~hundreds         | good balance for monitoring          |
| 500             | ~thousands        | tight tails for SLO reporting        |

Error shrinks as `c` grows and is concentrated away from the tails: at the
default compression, extreme percentiles are typically estimated to within a
fraction of a percent of the true value, while mid-distribution quantiles carry a
larger relative error.

## Determinism and order-dependence

A `Digest` uses **no randomness**, so a fixed sequence of operations always
produces the same result. The retained centroids — and therefore the estimates —
can depend on the **order** of `Add` and `Merge` calls: different orderings of the
same data yield close but not bit-identical quantiles.

## Mergeability

`Merge` folds one digest's centroids into another, so per-shard or per-worker
digests combine into one covering the union — the basis for parallel and
distributed quantile aggregation. Both digests must share the same compression;
a mismatch (or a nil argument) returns an error wrapping `ErrInvalidConfig` and
leaves the receiver unchanged.

## Thread safety

`Digest` is **not** safe for concurrent use. `ConcurrentDigest` (via
`NewConcurrent`) wraps it with a read-write mutex. The quantile queries take the
write lock because the digest compresses lazily on read. To merge a
`ConcurrentDigest` into another, pass its `Snapshot()`.

## Float64, not generic

Unlike the comparable-typed sibling sketches, `Digest` is deliberately **not
generic**: it operates on `float64` only, matching the value domain of the stats
quantile functions it approximates.

## Sibling sketches

- [`bloom`](../bloom) — approximate set membership.
- [`countmin`](../countmin) — approximate frequency counts.
- [`hll`](../hll) — approximate distinct-element cardinality.
