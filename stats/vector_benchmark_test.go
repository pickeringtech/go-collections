package stats_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func BenchmarkDot(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v []float64
	}{
		{name: "3 elements", a: benchInput(3), v: benchInput(3)},
		{name: "10 elements", a: benchInput(10), v: benchInput(10)},
		{name: "100 elements", a: benchInput(100), v: benchInput(100)},
		{name: "1_000 elements", a: benchInput(1_000), v: benchInput(1_000)},
		{name: "10_000 elements", a: benchInput(10_000), v: benchInput(10_000)},
		{name: "100_000 elements", a: benchInput(100_000), v: benchInput(100_000)},
		{name: "1_000_000 elements", a: benchInput(1_000_000), v: benchInput(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = stats.Dot(bm.a, bm.v)
			}
		})
	}
}

func BenchmarkNorm(b *testing.B) {
	benchmarks := []struct {
		name  string
		input []float64
	}{
		{name: "3 elements", input: benchInput(3)},
		{name: "10 elements", input: benchInput(10)},
		{name: "100 elements", input: benchInput(100)},
		{name: "1_000 elements", input: benchInput(1_000)},
		{name: "10_000 elements", input: benchInput(10_000)},
		{name: "100_000 elements", input: benchInput(100_000)},
		{name: "1_000_000 elements", input: benchInput(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = stats.Norm(bm.input)
			}
		})
	}
}
