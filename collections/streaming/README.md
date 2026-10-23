# collections/streaming

Online algorithms for unbounded data streams. You feed each structure one
element at a time and query it as you go, so the stream never has to fit in
memory.

```go
import "github.com/pickeringtech/go-collections/collections/streaming"
```

## Choose an algorithm

| Type / function | Question it answers | State held | Cost per element |
| --- | --- | --- | --- |
| `TopK[T]` (`NewTopK`, `NewTopKOrdered`) | "What are the k highest-ranked elements so far?" (exact) | O(k) | `Add` O(log k); `Result` O(k log k) |
| `Reservoir[T]` | "Give me a uniform random sample of k elements" (Vitter's Algorithm R) | O(k) | `Add` O(1) |
| `WeightedReservoir[T]` | "Give me k elements, heavier ones more likely" (Efraimidis–Spirakis A-Res, without replacement) | O(k) | `Add` O(log k) |
| `Bootstrap` / `BootstrapN` | "Resample this slice with replacement" (for the statistical bootstrap) | none between calls | O(n) per resample |
| `RunningMean` | "What is the mean so far?" | O(1) | O(1) |
| `RunningVariance` | "What are the mean and the sample/population variance so far?" (Welford) | O(1) | O(1) |
| `EWMA` | "What is the exponentially weighted moving average?" | O(1) | O(1) |
| `RunningMinMax[T]` | "What are the smallest and largest elements so far?" | O(1) | O(1) |

State never grows with the stream. `Bootstrap`/`BootstrapN` are the one
batch-style exception: they take a finite slice and return resamples of the
same size, holding nothing between calls.

`TopK` is exact. The samplers are randomized, but they draw from exactly the
distribution they specify. The online aggregates agree with their batch
counterparts in `stats`, including the `(value, ok)` contracts:
`RunningMean` matches `stats.Mean` (up to floating-point rounding),
`RunningVariance` matches `stats.SampleVariance`/`stats.PopulationVariance`,
and `RunningMinMax` matches `stats.MinMax`. `EWMA` has no batch counterpart.

## Online aggregates and `(value, ok)`

The aggregates report `ok == false` until they have seen enough data to give a
defined answer. For example, a sample variance needs at least two values, and
`RunningMean` stays `ok == false` after any NaN or ±Inf input, as
`stats.Mean` does.

```go
mean := streaming.NewRunningMean()
for _, v := range []float64{2, 4, 6, 8} {
	mean.Add(v)
}

avg, ok := mean.Result()
fmt.Println(avg, ok, mean.Count())
// Output: 5 true 4
```

`EWMA`'s alpha is clamped into (0, 1]: a non-positive or NaN alpha becomes 0.5.
The average is primed by the first `Add`, so it is not biased towards zero.

## Reproducible sampling

The randomized operations (`Reservoir`, `WeightedReservoir`, `Bootstrap`,
`BootstrapN`) take a `*math/rand/v2.Rand` as their last parameter. Passing
`nil` uses a fixed default seed, equivalent to `NewRand(0)`, so the same input
always gives the same sample. Pass `NewRand(seed)` to choose the sequence. The
generator is non-cryptographic: it is for reproducibility, not security.

```go
// A fixed seed makes the sample reproducible; pass nil for the default.
r := streaming.NewReservoir[int](3, streaming.NewRand(42))
for v := 0; v < 100; v++ {
	r.Add(v)
}

sample := r.Result()
sort.Ints(sample) // sample order is unspecified; sort for a stable print
fmt.Println(sample)
// Output: [9 15 91]
```

`Result` returns a copy, and the order of a sample is **not** specified. Sort
it if you need a stable order, as above. `WeightedReservoir` needs strictly
positive, finite weights. It ignores an element whose weight is zero,
negative, NaN or ±Inf. Only the ratios between weights matter, so any positive
finite scale works.

## Thread safety

These types are **single-threaded**: an `Add` must not race with another `Add`
or with a `Result`. Guard a shared instance with your own lock, or give each
goroutine its own instance and combine the results.

## Further documentation

- [Package overview](./doc.go)
- Executable examples:
  [aggregates](./aggregate_example_test.go),
  [reservoirs](./reservoir_example_test.go),
  [top-k](./topk_example_test.go),
  [bootstrap](./bootstrap_example_test.go)
- [API reference](https://pkg.go.dev/github.com/pickeringtech/go-collections/collections/streaming)
- Approximate frequency, cardinality and quantile sketches live in
  [`collections/sketches`](../sketches/README.md).
