package sets_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// These benchmarks backfill the missing Benchmark leg for the free functions in
// sets/transform.go (issue #52 sweep). They reuse the seq helper and setSizes
// ladder from set_bench_test.go and follow the same Benchmark + b.Run("size_%d")
// + ReportAllocs shape, so ns/op, B/op and allocs/op are all captured.

func BenchmarkMap(b *testing.B) {
	for _, size := range setSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s := sets.NewHash(seq(size)...)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = sets.Map(s, func(n int) int { return n * 2 })
			}
		})
	}
}

func BenchmarkReduce(b *testing.B) {
	for _, size := range setSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s := sets.NewHash(seq(size)...)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = sets.Reduce(s, 0, func(acc, n int) int { return acc + n })
			}
		})
	}
}
