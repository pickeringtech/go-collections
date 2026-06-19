# countmin

A generic [Count-Min sketch](https://en.wikipedia.org/wiki/Count%E2%80%93min_sketch):
approximate frequency counts over a stream, in space far smaller than an exact
counter table.

```go
import "github.com/pickeringtech/go-collections/collections/sketches/countmin"

s, err := countmin.New[string](0.001, 0.01) // overshoot ≤ 0.001·N, 99% confidence
s.Add("/index.html")
s.AddCount("/api/v1", 42)
s.Estimate("/index.html") // approximate count, never an under-count
```

## The guarantee

`Estimate` **never under-reports**. It returns the true count plus a one-sided
error of at most `epsilon·N` (with `N` the total count added) with probability at
least `1-delta`. Over-counting a rare item is possible; missing a frequent one is
not — ideal for heavy-hitter detection.

## Accuracy vs memory

`New(epsilon, delta)` sizes the counter table:

| Quantity            | Formula               |
| ------------------- | --------------------- |
| columns `w`         | `ceil(e/epsilon)`     |
| rows `d`            | `ceil(ln(1/delta))`   |

The table holds `w·d` 64-bit counters **regardless of stream cardinality** —
memory depends only on the target accuracy. Counters saturate rather than
overflow.

## Hashing

Seeded and deterministic; each row's column comes from two base hashes via
Kirsch–Mitzenmacher double hashing. `WithSeed` varies the seed; `WithHasher`
plugs in a custom hash.

## Mergeability

`Merge` sums two sketches' counters element-wise, combining per-shard or
per-worker sketches into one. Both must share dimensions and seed.

## Thread safety

`Sketch` is **not** safe for concurrent use. `ConcurrentSketch` (via
`NewConcurrent`) adds a read-write mutex: `Estimate` reads, mutations write. To
merge a `ConcurrentSketch` into another, pass its `Snapshot()`.

## Sibling sketches

- [`bloom`](../bloom) — approximate set membership.
- [`hll`](../hll) — approximate distinct-element cardinality.
