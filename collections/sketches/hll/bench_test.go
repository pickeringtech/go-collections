package hll_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/hll"
)

// sketchSizes is the shared element-count ladder; see the bloom suite for the
// naming rationale (Benchmark<Impl>_<Op> / size_%d) the report generator needs.
var sketchSizes = []int{10, 100, 1000, 10000}

func BenchmarkSketch_Add(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s, _ := hll.New[int]()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				s.Add(i % size)
			}
		})
	}
}

func BenchmarkSketch_Count(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s, _ := hll.New[int]()
			for i := 0; i < size; i++ {
				s.Add(i)
			}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = s.Count()
			}
		})
	}
}

func BenchmarkConcurrentSketch_Add(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s, _ := hll.NewConcurrent[int]()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				s.Add(i % size)
			}
		})
	}
}
