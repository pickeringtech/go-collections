# tdigest

A [t-digest](https://github.com/tdunning/t-digest): a streaming, mergeable sketch
that estimates quantiles (and the inverse CDF) of a stream of `float64` values in
memory that scales with a compression parameter and grows only logarithmically
with how many values it has seen.

```go
import "github.com/pickeringtech/go-collections/collections/sketches/tdigest"

d, _ := tdigest.New() // DefaultCompression (100)
for _, v := range latencies {
	d.Add(v)
}
p99, ok := d.Percentile(99) // estimated 99th-percentile latency
```

It is the approximate, small-memory counterpart of the exact
[`stats`](../../../stats) quantile functions: where `stats.Quantile` sorts the
whole sample, a `Digest` keeps only a small set of weighted centroids.

## The guarantee

- **Small memory.** The retained centroid count is proportional to the
  compression parameter and grows only **logarithmically** with the stream
  length — each tenfold increase in the number of values adds about
  1.5×compression centroids — so even a very long stream keeps the digest small.
- **Exact extremes.** `Quantile(0)` and `Quantile(1)` return the exact observed
  minimum and maximum; a single distinct value is returned exactly for any
  quantile.
- **Tail-accurate.** The scale function makes error **smallest at the tails**
  (q near 0 or 1, e.g. p99/p999) and largest in the middle — the opposite of a
  fixed histogram, and exactly what latency monitoring wants.

## Accuracy vs memory

`New()` uses `DefaultCompression` (100); `WithCompression(c)` trades memory for
accuracy.

| Compression `c` | Centroids after 1e3 / 1e5 / 1e7 values | Typical use                           |
| --------------- | -------------------------------------- | ------------------------------------- |
| 20              | ~75 / ~140 / ~200                      | rough estimates, very small footprint |
| 100 (default)   | ~280 / ~590 / ~890                     | good balance for monitoring           |
| 500             | ~1000 / ~2400 / ~3900                  | tight tails for SLO reporting         |

(Measured on uniform data; up to 5·`c` further values sit in the merge buffer
between compressions.)

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
