package sets_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// This file standardises the set benchmark suite so the generated BENCHMARKS.md
// report (issue #50) has consistent, comparable coverage across implementations.
// It mirrors the operation × implementation × size matrix established by
// collections/dicts/hash_bench_test.go: every benchmark is named
// Benchmark<Impl>_<Op> and sub-benchmarks across the shared size set via
// b.Run("size_%d"), with b.ReportAllocs()/b.ResetTimer() so ns/op, B/op and
// allocs/op are all captured. The benchreport generator parses exactly this
// shape, so keep new set benchmarks to the same naming and structure.

// setSizes is the representative element-count matrix shared by every set
// benchmark. Read-style benchmarks (Contains/ForEach) build the set once and
// reuse it; the mutating Add benchmark also builds once and undoes its single
// add each iteration (see benchAdd), so it carries no per-iteration rebuild and
// can run the full ladder without the CI blow-up of issue #112.
var setSizes = []int{10, 100, 1000, 10000}

// setCtor builds a MutableSet of the implementation under test, pre-filled with
// elements 0..n-1. Each implementation supplies one of these so the
// per-operation helpers below stay implementation-agnostic.
type setCtor func(elements ...int) sets.MutableSet[int]

func seq(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

// benchContains measures membership testing on a set holding `size` elements.
func benchContains(b *testing.B, ctor setCtor) {
	for _, size := range setSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s := ctor(seq(size)...)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = s.Contains(i % size)
			}
		})
	}
}

// benchAdd measures adding a single new element to a set already holding `size`
// elements. The set is built once and each iteration adds an element then
// removes it again, measuring an add+remove round-trip at a steady size.
// Rebuilding the whole set under StopTimer every iteration instead made
// wall-time ≈ b.N × O(size), unbounded by -benchtime — the hour-plus CI blow-up
// of issue #112. The cheap O(1) inverse is timed rather than excluded with a
// per-iteration b.StopTimer(), which reads memstats under -benchmem and would
// re-introduce the blow-up at the ns scale these ops run at.
func benchAdd(b *testing.B, ctor setCtor) {
	for _, size := range setSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s := ctor(seq(size)...)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				s.AddInPlace(size)
				s.RemoveInPlace(size)
			}
		})
	}
}

// benchForEach measures a full iteration over a set of `size` elements.
func benchForEach(b *testing.B, ctor setCtor) {
	for _, size := range setSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s := ctor(seq(size)...)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				count := 0
				s.ForEach(func(v int) { count++ })
			}
		})
	}
}

// Implementation constructors. Each adapts a concrete constructor to the shared
// MutableSet[int] interface so the helpers above can drive any implementation.
func newHash(e ...int) sets.MutableSet[int]             { return sets.NewHash(e...) }
func newConcurrentHash(e ...int) sets.MutableSet[int]   { return sets.NewConcurrentHash(e...) }
func newConcurrentHashRW(e ...int) sets.MutableSet[int] { return sets.NewConcurrentHashRW(e...) }

func BenchmarkHash_Contains(b *testing.B) { benchContains(b, newHash) }
func BenchmarkHash_Add(b *testing.B)      { benchAdd(b, newHash) }
func BenchmarkHash_ForEach(b *testing.B)  { benchForEach(b, newHash) }

func BenchmarkConcurrentHash_Contains(b *testing.B) { benchContains(b, newConcurrentHash) }
func BenchmarkConcurrentHash_Add(b *testing.B)      { benchAdd(b, newConcurrentHash) }
func BenchmarkConcurrentHash_ForEach(b *testing.B)  { benchForEach(b, newConcurrentHash) }

func BenchmarkConcurrentHashRW_Contains(b *testing.B) { benchContains(b, newConcurrentHashRW) }
func BenchmarkConcurrentHashRW_Add(b *testing.B)      { benchAdd(b, newConcurrentHashRW) }
func BenchmarkConcurrentHashRW_ForEach(b *testing.B)  { benchForEach(b, newConcurrentHashRW) }
