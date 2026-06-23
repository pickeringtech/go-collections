// Package streaming provides online, bounded-memory algorithms for unbounded
// data streams — structures you feed one element at a time and query as you go,
// without ever holding the whole stream in memory. It composes with the rest of
// go-collections (heaps, channels, stats) rather than duplicating them.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/streaming"
//
//	top := streaming.NewTopKOrdered[int](3)
//	for _, v := range []int{5, 1, 9, 3, 7, 2, 8} {
//		top.Add(v)
//	}
//	result := top.Result()
//	// Result: [9 8 7] — the three largest, highest first, in O(k) memory
//
// # Why Use Streaming Algorithms?
//
// The obvious way to find the top-k of a stream is to buffer everything, sort,
// and slice — O(n) memory and O(n log n) time, impossible when the stream is
// unbounded:
//
//	// Native approach — buffers the entire stream:
//	all := make([]int, 0)
//	for v := range stream {
//		all = append(all, v)
//	}
//	sort.Sort(sort.Reverse(sort.IntSlice(all)))
//	top := all[:3]
//
//	// With this package — bounded O(k) memory, O(log k) per element:
//	top := streaming.NewTopKOrdered[int](3)
//	for v := range stream {
//		top.Add(v)
//	}
//	result := top.Result()
//
// # Exact vs Approximate
//
// TopK is exact: Result is precisely the k highest-ranked elements seen, never
// an estimate. Reservoir and WeightedReservoir are exact in their sampling
// probabilities — they are randomized, but the distribution they draw from is
// the precise one specified, not an approximation. The online aggregates
// (RunningMean, RunningVariance, EWMA, RunningMinMax) are exact too: each agrees
// with the corresponding batch stats reduction over the same data. (Approximate
// sketches — frequency, cardinality — are a separate concern and live
// elsewhere.)
//
// # Bounded vs Unbounded Memory
//
// Every algorithm here runs in memory that does not grow with the stream:
// O(k) for the samplers and TopK, O(1) for the online aggregates and the
// bootstrap. The bootstrap resamplers are the one batch-style exception — they
// take a finite slice and return slices of the same size — but they hold no
// per-stream state.
//
// # Determinism and Seeding
//
// The randomized operations — the sampling types and the bootstrap functions —
// take an explicit *math/rand/v2.Rand as their last parameter and are
// deterministic by default: passing nil uses a fixed seed (equivalent to
// NewRand(0)), so the same input always yields the same result. Supply your own
// generator — NewRand(seed) is provided for convenience — to vary or pin the
// sequence. The generator is non-cryptographic; it is meant for
// reproducibility, not security.
//
// # Bootstrap Resampling
//
// Bootstrap draws a single resample of a slice — len(input) elements sampled
// uniformly with replacement — and BootstrapN draws several independent
// resamples from one seed. They are the building block of the statistical
// bootstrap: recompute a statistic (compose a stats reduction) over many
// resamples to estimate its sampling distribution. Both never mutate the input
// and return non-nil empty results for nil/empty input or a non-positive count.
//
// # Online Aggregates
//
// These structures fold a stream one element at a time and answer with the
// (value, ok) contract, where ok reports whether enough data has been seen:
//
//   - RunningMean accumulates the arithmetic mean incrementally (no growing
//     running total). Result is ok once at least one value is added.
//   - RunningVariance runs Welford's algorithm online, exposing the mean,
//     SampleVariance (ok for n >= 2, Bessel-corrected) and PopulationVariance
//     (ok for n >= 1). The recurrence and these contracts intentionally match
//     stats.SampleVariance and stats.PopulationVariance exactly — it is a
//     streaming reimplementation rather than a reuse, because stats' helper is a
//     batch-over-slice function that cannot be fed incrementally without first
//     buffering the whole stream (which would defeat bounded memory). Nothing in
//     stats is exported or modified to share it.
//   - EWMA maintains an exponentially weighted moving average; alpha is clamped
//     into (0, 1] and the average is primed on the first Add so it is not biased
//     toward zero.
//   - RunningMinMax tracks the smallest and largest element of an ordered stream,
//     mirroring stats.MinMax.
//
// # Thread Safety
//
// The types in this package are single-threaded: an Add must not race with
// another Add or a Result. Guard a shared instance with your own lock, or give
// each goroutine its own and merge results.
//
// # Available Algorithms
//
// TopK[T]:
//   - Streaming top-k by an arbitrary LessFunc, backed by a size-bounded
//     min-heap (collections/heaps). Add is O(log k); Result is O(k log k).
//
// Reservoir[T]:
//   - Uniform fixed-size sampling over an unbounded stream (Vitter's Algorithm
//     R). Every element seen is retained with equal probability k/n. Add is
//     O(1) in O(k) memory.
//
// WeightedReservoir[T]:
//   - Weighted fixed-size sampling without replacement (Efraimidis & Spirakis
//     A-Res), backed by a size-bounded min-heap (collections/heaps). An
//     element's retention probability grows with its weight. Add is O(log k) in
//     O(k) memory.
//
// Bootstrap / BootstrapN:
//   - Resample a slice with replacement, for the statistical bootstrap. Seeded
//     for reproducibility; never mutates the input.
//
// RunningMean, RunningVariance, EWMA, RunningMinMax[T]:
//   - Online aggregates that fold a stream in O(1) memory and answer with the
//     (value, ok) contract. RunningVariance matches stats.SampleVariance and
//     stats.PopulationVariance; RunningMinMax mirrors stats.MinMax.
package streaming
