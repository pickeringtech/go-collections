package tdigest_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/tdigest"
)

// sketchSizes is the shared element-count ladder; see the bloom suite for the
// naming rationale (Benchmark<Impl>_<Op> / size_%d) the report generator needs.
var sketchSizes = []int{10, 100, 1000, 10000}

func BenchmarkDigest_Add(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			d, _ := tdigest.New()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				d.Add(float64(i % size))
			}
		})
	}
}

func BenchmarkDigest_Quantile(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			d, _ := tdigest.New()
			for i := 0; i < size; i++ {
				d.Add(float64(i))
			}
			d.Quantile(0.5) // force a compress before timing
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = d.Quantile(0.99)
			}
		})
	}
}

func BenchmarkDigest_Merge(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			src, _ := tdigest.New()
			for i := 0; i < size; i++ {
				src.Add(float64(i))
			}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				dst, _ := tdigest.New()
				_ = dst.Merge(src)
			}
		})
	}
}

func BenchmarkConcurrentDigest_Add(b *testing.B) {
	for _, size := range sketchSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			c, _ := tdigest.NewConcurrent()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				c.Add(float64(i % size))
			}
		})
	}
}
