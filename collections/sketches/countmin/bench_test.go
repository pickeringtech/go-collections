package countmin_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/countmin"
)

// sketchSizes is the shared element-count ladder; see the bloom suite for the
// naming rationale (Benchmark<Impl>_<Op> / size_%d) the report generator needs.
var sketchSizes = []int{10, 100, 1000, 10000}

func BenchmarkSketch_Add(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s, _ := countmin.New[int](0.001, 0.001)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				s.Add(i % size)
			}
		})
	}
}

func BenchmarkSketch_Estimate(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s, _ := countmin.New[int](0.001, 0.001)
			for i := 0; i < size; i++ {
				s.Add(i)
			}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = s.Estimate(i % size)
			}
		})
	}
}

func BenchmarkConcurrentSketch_Add(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			s, _ := countmin.NewConcurrent[int](0.001, 0.001)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				s.Add(i % size)
			}
		})
	}
}
