# hll

A generic [HyperLogLog](https://en.wikipedia.org/wiki/HyperLogLog): estimate the
number of **distinct** elements in a stream using a fixed few kilobytes, even for
billions of distinct items.

```go
import "github.com/pickeringtech/go-collections/collections/sketches/hll"

s, err := hll.New[string]() // default precision 14 (~16 KB, ~0.81% error)
for _, visitor := range stream {
    s.Add(visitor)
}
s.Count() // estimated distinct visitors
```

## How it works

Each element is hashed; the top `p` bits pick one of `m = 2^p` registers and the
rest contribute the length of their leading zero-run. Each register keeps the
longest run it has seen, and the harmonic mean of the registers (bias-corrected)
estimates the cardinality. Re-adding an element never changes the estimate.

## Accuracy vs memory

Precision `p` (via `WithPrecision`, range `[4, 18]`) sets the trade-off:

| Precision `p` | Registers `m` | Memory   | Std. error |
| ------------- | ------------- | -------- | ---------- |
| 10            | 1024          | ~1 KB    | ~3.25%     |
| 14 (default)  | 16384         | ~16 KB   | ~0.81%     |
| 16            | 65536         | ~64 KB   | ~0.41%     |

Standard error is about `1.04/sqrt(m)`, and memory is fixed by `p` alone — it
never grows with the stream. `StandardError()` reports the expected relative
error. Small cardinalities use linear counting for better accuracy.

## Hashing

Seeded and deterministic. HyperLogLog reads both the high index bits and the
low-bit zero-run, so the default hasher applies a strong finalizer for good
distribution. `WithSeed` varies the seed; `WithHasher` plugs in a custom hash.

## Mergeability

`Merge` takes the register-wise maximum of two sketches, estimating the
cardinality of the **union** — and it is exact with respect to that union (the
merged registers equal those of one sketch fed both streams). Both must share
precision and seed.

## Thread safety

`Sketch` is **not** safe for concurrent use. `ConcurrentSketch` (via
`NewConcurrent`) adds a read-write mutex: `Count` reads, `Add`/`Merge`/`Clear`
write. To merge a `ConcurrentSketch` into another, pass its `Snapshot()`.

## Sibling sketches

- [`bloom`](../bloom) — approximate set membership.
- [`countmin`](../countmin) — approximate frequency counts.
