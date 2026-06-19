# bloom

A generic [Bloom filter](https://en.wikipedia.org/wiki/Bloom_filter): a compact,
probabilistic set that answers membership queries in space independent of element
size, at the cost of a tunable false-positive rate.

```go
import "github.com/pickeringtech/go-collections/collections/sketches/bloom"

f, err := bloom.New[string](1_000_000, 0.01) // 1M items, 1% false positives
f.Add("alice")
f.Contains("alice") // true
f.Contains("carol") // false — never added
```

## The guarantee

- **No false negatives.** If `Contains` returns `false`, the element was
  definitely never added.
- **Tunable false positives.** `Contains` may return `true` for an element that
  was never added, with probability close to the configured rate.
- **No removal, no enumeration.** A Bloom filter only ever accumulates; you
  cannot delete an element or list the contents.

## Accuracy vs memory

`New(n, p)` sizes the filter with the standard optimal formulas:

| Quantity            | Formula                       |
| ------------------- | ----------------------------- |
| bits `m`            | `ceil(-n·ln p / (ln 2)^2)`    |
| hash functions `k`  | `round((m/n)·ln 2)`           |

That works out to about `-1.44·log2(p)` bits **per element** — roughly 9.6 bits
each at `p = 0.01` — no matter how large each element is. Add more than `n`
elements and the false-positive rate climbs above the target;
`EstimatedFalsePositiveRate()` reports the rate at the current fill and
`ApproxCount()` estimates how many distinct elements have been added.

## Hashing

Hashing is seeded and deterministic. `k` bit positions come from two base hashes
via Kirsch–Mitzenmacher double hashing (two hashes per op, not `k`). Use
`WithSeed` to vary the seed and `WithHasher` to plug in a custom hash for exotic
key types or cross-process reproducibility.

## Mergeability

`Merge` ORs one filter's bits into another, so per-shard or per-worker filters
combine into one covering the union — the basis for parallel/distributed
aggregation. Both filters must share capacity and seed.

## Thread safety

`Filter` is **not** safe for concurrent use. `ConcurrentFilter` (via
`NewConcurrent`) wraps it with a read-write mutex: reads take a read lock, writes
the write lock. To merge a `ConcurrentFilter` into another, pass its `Snapshot()`.

## Sibling sketches

- [`countmin`](../countmin) — approximate frequency counts.
- [`hll`](../hll) — approximate distinct-element cardinality.
