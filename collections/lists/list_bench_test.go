package lists_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// This file standardises the list benchmark suite so the generated BENCHMARKS.md
// report (issue #50) has consistent, comparable coverage across implementations.
// It mirrors the operation × implementation × size matrix established by
// collections/dicts/hash_bench_test.go: every benchmark is named
// Benchmark<Impl>_<Op> and sub-benchmarks across the shared size set via
// b.Run("size_%d"), with b.ReportAllocs()/b.ResetTimer() so ns/op, B/op and
// allocs/op are all captured. The benchreport generator parses exactly this
// shape, so keep new list benchmarks to the same naming and structure.

// listSizes is the representative element-count matrix for read-style
// benchmarks (Get/Filter/ForEach), which build the list once and reuse it.
var listSizes = []int{10, 100, 1000, 10000}

// mutateSizes caps the matrix for benchmarks that rebuild the list under
// StopTimer every iteration (Push). The rebuild is wall-clock the framework
// can't amortise, so the 10k cell would dominate CI time for little signal —
// mirroring the existing dicts BenchmarkHash_Put, which stops at 1000.
var mutateSizes = []int{10, 100, 1000}

// listCtor builds a MutableList of the implementation under test, pre-filled
// with elements 0..n-1. Each implementation supplies one of these so the
// per-operation helpers below stay implementation-agnostic.
type listCtor func(elements ...int) lists.MutableList[int]

func seq(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

// benchGet measures positional access on a list already holding `size` elements.
func benchGet(b *testing.B, ctor listCtor) {
	for _, size := range listSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			l := ctor(seq(size)...)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, _ = l.Get(i%size, 0)
			}
		})
	}
}

// benchPush measures appending a single element to a list already holding
// `size` elements. The list is rebuilt under StopTimer each iteration so the
// timed work is a single PushInPlace, not unbounded growth across b.N.
func benchPush(b *testing.B, ctor listCtor) {
	for _, size := range mutateSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				l := ctor(seq(size)...)
				b.StartTimer()

				l.PushInPlace(size + i)
			}
		})
	}
}

// benchFilter measures filtering a list of `size` elements into a new slice.
func benchFilter(b *testing.B, ctor listCtor) {
	for _, size := range listSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			l := ctor(seq(size)...)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = l.Filter(func(v int) bool { return v%2 == 0 })
			}
		})
	}
}

// benchForEach measures a full iteration over a list of `size` elements.
func benchForEach(b *testing.B, ctor listCtor) {
	for _, size := range listSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			l := ctor(seq(size)...)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				count := 0
				l.ForEach(func(v int) { count++ })
			}
		})
	}
}

// Implementation constructors. Each adapts a concrete constructor to the shared
// MutableList[int] interface so the helpers above can drive any implementation.
func newArray(e ...int) lists.MutableList[int]        { return lists.NewArray(e...) }
func newLinked(e ...int) lists.MutableList[int]       { return lists.NewLinked(e...) }
func newDoublyLinked(e ...int) lists.MutableList[int] { return lists.NewDoublyLinked(e...) }

func BenchmarkArray_Get(b *testing.B)     { benchGet(b, newArray) }
func BenchmarkArray_Push(b *testing.B)    { benchPush(b, newArray) }
func BenchmarkArray_Filter(b *testing.B)  { benchFilter(b, newArray) }
func BenchmarkArray_ForEach(b *testing.B) { benchForEach(b, newArray) }

func BenchmarkLinked_Get(b *testing.B)     { benchGet(b, newLinked) }
func BenchmarkLinked_Push(b *testing.B)    { benchPush(b, newLinked) }
func BenchmarkLinked_Filter(b *testing.B)  { benchFilter(b, newLinked) }
func BenchmarkLinked_ForEach(b *testing.B) { benchForEach(b, newLinked) }

func BenchmarkDoublyLinked_Get(b *testing.B)     { benchGet(b, newDoublyLinked) }
func BenchmarkDoublyLinked_Push(b *testing.B)    { benchPush(b, newDoublyLinked) }
func BenchmarkDoublyLinked_Filter(b *testing.B)  { benchFilter(b, newDoublyLinked) }
func BenchmarkDoublyLinked_ForEach(b *testing.B) { benchForEach(b, newDoublyLinked) }
