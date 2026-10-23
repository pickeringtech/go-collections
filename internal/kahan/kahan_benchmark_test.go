package kahan_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/internal/kahan"
)

func benchValues(n int) []float64 {
	values := make([]float64, n)
	for i := range values {
		values[i] = float64(i%100) + 0.1
	}
	return values
}

func BenchmarkSum(b *testing.B) {
	benchmarks := []struct {
		name   string
		values []float64
	}{
		{name: "3 elements", values: benchValues(3)},
		{name: "10 elements", values: benchValues(10)},
		{name: "100 elements", values: benchValues(100)},
		{name: "1_000 elements", values: benchValues(1_000)},
		{name: "10_000 elements", values: benchValues(10_000)},
		{name: "100_000 elements", values: benchValues(100_000)},
		{name: "1_000_000 elements", values: benchValues(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var total kahan.Sum
				for _, x := range bm.values {
					total.Add(x)
				}
				_ = total.Total()
			}
		})
	}
}
