package stats_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func BenchmarkProduct(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Product(values)
	}
}

func BenchmarkMedian(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Median(values)
	}
}

func BenchmarkMode(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Mode(values)
	}
}

func BenchmarkMinMax(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = stats.MinMax(values)
	}
}

func BenchmarkArgMin(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.ArgMin(values)
	}
}

func BenchmarkArgMax(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.ArgMax(values)
	}
}

func BenchmarkRange(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Range(values)
	}
}

func BenchmarkCumulativeSum(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stats.CumulativeSum(values)
	}
}

func BenchmarkClampAll(b *testing.B) {
	values := benchInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stats.ClampAll(values, 25, 75)
	}
}
