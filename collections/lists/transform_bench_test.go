package lists_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// These benchmarks backfill the missing Benchmark leg for the free functions in
// lists/transform.go (issue #52 sweep). They reuse the seq helper and listSizes
// ladder from list_bench_test.go and follow the same Benchmark +
// b.Run("size_%d") + ReportAllocs shape used across the list suite.

func BenchmarkMap(b *testing.B) {
	for _, size := range listSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			l := lists.NewArray(seq(size)...)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = lists.Map(l, func(n int) int { return n * 2 })
			}
		})
	}
}

func BenchmarkFlatMap(b *testing.B) {
	for _, size := range listSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			l := lists.NewArray(seq(size)...)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = lists.FlatMap(l, func(n int) lists.List[int] {
					return lists.NewArray(n, n)
				})
			}
		})
	}
}

func BenchmarkReduce(b *testing.B) {
	for _, size := range listSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			l := lists.NewArray(seq(size)...)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = lists.Reduce(l, 0, func(acc, n int) int { return acc + n })
			}
		})
	}
}
