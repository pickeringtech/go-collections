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

// setSizes is the representative element-count matrix for read-style benchmarks
// (Contains/ForEach), which build the set once and reuse it.
var setSizes = []int{10, 100, 1000, 10000}

// mutateSizes caps the matrix for benchmarks that rebuild the set under
// StopTimer every iteration (Add). The rebuild is wall-clock the framework
// can't amortise, so the 10k cell would dominate CI time for little signal —
// mirroring the existing dicts BenchmarkHash_Put, which stops at 1000.
var mutateSizes = []int{10, 100, 1000}

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
// elements. The set is rebuilt under StopTimer each iteration so the timed work
// is a single AddInPlace, not unbounded growth across b.N.
func benchAdd(b *testing.B, ctor setCtor) {
	for _, size := range mutateSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				s := ctor(seq(size)...)
				b.StartTimer()

				s.AddInPlace(size + i)
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
