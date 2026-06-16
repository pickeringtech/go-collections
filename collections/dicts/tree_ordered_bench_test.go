package dicts_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
	"github.com/pickeringtech/go-collections/slices"
)

// orderedBenchSizes is the standard scaling ladder for benchmarks.
var orderedBenchSizes = []int{3, 10, 100, 1_000, 10_000, 100_000, 1_000_000}

func buildBenchTree(size int) *dicts.Tree[int, int] {
	keys := slices.Generate(size, slices.NumericIdentityGenerator[int])
	pairs := make([]dicts.Pair[int, int], size)
	for i, k := range keys {
		pairs[i] = dicts.Pair[int, int]{Key: k, Value: k}
	}
	return dicts.NewTree(pairs...)
}

func BenchmarkTree_Floor(b *testing.B) {
	for _, size := range orderedBenchSizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			tree := buildBenchTree(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = tree.Floor(i % size)
			}
		})
	}
}

func BenchmarkTree_Ceiling(b *testing.B) {
	for _, size := range orderedBenchSizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			tree := buildBenchTree(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = tree.Ceiling(i % size)
			}
		})
	}
}

func BenchmarkTree_All(b *testing.B) {
	for _, size := range orderedBenchSizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			tree := buildBenchTree(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for range tree.All() {
				}
			}
		})
	}
}

func BenchmarkTree_Range(b *testing.B) {
	for _, size := range orderedBenchSizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			tree := buildBenchTree(size)
			lo, hi := size/4, size*3/4
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = tree.Range(lo, hi)
			}
		})
	}
}
