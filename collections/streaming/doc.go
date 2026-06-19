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
// an estimate. (Approximate sketches — frequency, cardinality — are a separate
// concern and live elsewhere.)
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
package streaming
