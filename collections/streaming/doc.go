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
// the precise one specified, not an approximation. (Approximate sketches —
// frequency, cardinality — are a separate concern and live elsewhere.)
//
// # Determinism and Seeding
//
// The sampling types take an explicit *math/rand/v2.Rand as their last
// constructor parameter and are deterministic by default: passing nil uses a
// fixed seed (equivalent to NewRand(0)), so the same stream always yields the
// same sample. Supply your own generator — NewRand(seed) is provided for
// convenience — to vary or pin the sequence. The generator is non-cryptographic;
// it is meant for reproducibility, not security.
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
package streaming
