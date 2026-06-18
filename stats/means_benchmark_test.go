package stats_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func benchInput(n int) []float64 {
	values := make([]float64, n)
	for i := range values {
		values[i] = float64(i%100) + 1 // 1..100, all positive
	}
	return values
}

func BenchmarkSum(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Sum(values)
	}
}

func BenchmarkMean(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Mean(values)
	}
}

func BenchmarkWeightedMean(b *testing.B) {
	values := benchInput(1000)
	weights := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.WeightedMean(values, weights)
	}
}

func BenchmarkGeometricMean(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.GeometricMean(values)
	}
}

func BenchmarkHarmonicMean(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.HarmonicMean(values)
	}
}
