package sets_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
	"github.com/pickeringtech/go-collections/slices"
)

// treeSetBenchSizes is the standard scaling ladder for benchmarks.
var treeSetBenchSizes = []int{3, 10, 100, 1_000, 10_000, 100_000, 1_000_000}

func buildBenchTreeSet(size int) *sets.TreeSet[int] {
	return sets.NewTreeSet(slices.Generate(size, slices.NumericIdentityGenerator[int])...)
}

func BenchmarkTreeSet_Floor(b *testing.B) {
	for _, size := range treeSetBenchSizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			s := buildBenchTreeSet(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = s.Floor(i % size)
			}
		})
	}
}

func BenchmarkTreeSet_All(b *testing.B) {
	for _, size := range treeSetBenchSizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			s := buildBenchTreeSet(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for range s.All() {
				}
			}
		})
	}
}

func BenchmarkTreeSet_Range(b *testing.B) {
	for _, size := range treeSetBenchSizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			s := buildBenchTreeSet(size)
			lo, hi := size/4, size*3/4
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = s.Range(lo, hi)
			}
		})
	}
}

func BenchmarkTreeSet_Intersection(b *testing.B) {
	for _, size := range treeSetBenchSizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			a := buildBenchTreeSet(size)
			other := buildBenchTreeSet(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = a.Intersection(other)
			}
		})
	}
}
