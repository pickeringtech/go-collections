package bloom_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/bloom"
)

// sketchSizes is the representative element-count ladder shared across the
// sketch benchmarks, mirroring the set/dict suites so the BENCHMARKS.md report
// stays comparable. Benchmarks are named Benchmark<Impl>_<Op> with size_%d
// sub-benchmarks and ReportAllocs, matching the report generator's expected
// shape.
var sketchSizes = []int{10, 100, 1000, 10000}

func BenchmarkFilter_Add(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			f, _ := bloom.New[int](size, 0.01)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				f.Add(i % size)
			}
		})
	}
}

func BenchmarkFilter_Contains(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			f, _ := bloom.New[int](size, 0.01)
			for i := 0; i < size; i++ {
				f.Add(i)
			}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = f.Contains(i % size)
			}
		})
	}
}

func BenchmarkConcurrentFilter_Add(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			f, _ := bloom.NewConcurrent[int](size, 0.01)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				f.Add(i % size)
			}
		})
	}
}

func BenchmarkConcurrentFilter_Contains(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			f, _ := bloom.NewConcurrent[int](size, 0.01)
			for i := 0; i < size; i++ {
				f.Add(i)
			}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = f.Contains(i % size)
			}
		})
	}
}
