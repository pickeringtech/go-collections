package dicts_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// These benchmarks backfill the missing Benchmark leg for the free functions in
// dicts/transform.go (issue #52 sweep). They follow the same Benchmark +
// b.Run("size_%d") + ReportAllocs shape established by hash_bench_test.go.

// transformBenchSizes is the element-count matrix for the transform benchmarks,
// matching the read-style sizes used by the rest of the dict suite.
var transformBenchSizes = []int{10, 100, 1000, 10000}

// benchHash builds a Hash dict of `size` entries keyed 0..size-1 with the value
// mirroring the key, giving the transform benchmarks deterministic input.
func benchHash(size int) dicts.Dict[int, int] {
	pairs := make([]dicts.Pair[int, int], size)
	for i := 0; i < size; i++ {
		pairs[i] = dicts.Pair[int, int]{Key: i, Value: i}
	}
	return dicts.NewHash(pairs...)
}

func BenchmarkMap(b *testing.B) {
	for _, size := range transformBenchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			d := benchHash(size)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = dicts.Map(d, func(k, v int) (int, int) { return k, v * 2 })
			}
		})
	}
}

func BenchmarkReduce(b *testing.B) {
	for _, size := range transformBenchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			d := benchHash(size)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = dicts.Reduce(d, 0, func(acc, k, v int) int { return acc + v })
			}
		})
	}
}
