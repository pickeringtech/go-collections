package similarity_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
	"github.com/pickeringtech/go-collections/ml/similarity"
)

func benchFloats(n int) []float64 {
	v := make([]float64, n)
	for i := range v {
		v[i] = float64(i%100) + 1
	}
	return v
}

func BenchmarkDotProduct(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v []float64
	}{
		{name: "3 elements", a: benchFloats(3), v: benchFloats(3)},
		{name: "10 elements", a: benchFloats(10), v: benchFloats(10)},
		{name: "100 elements", a: benchFloats(100), v: benchFloats(100)},
		{name: "1_000 elements", a: benchFloats(1_000), v: benchFloats(1_000)},
		{name: "10_000 elements", a: benchFloats(10_000), v: benchFloats(10_000)},
		{name: "100_000 elements", a: benchFloats(100_000), v: benchFloats(100_000)},
		{name: "1_000_000 elements", a: benchFloats(1_000_000), v: benchFloats(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = similarity.DotProduct(bm.a, bm.v)
			}
		})
	}
}

func BenchmarkCosineSimilarity(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v []float64
	}{
		{name: "3 elements", a: benchFloats(3), v: benchFloats(3)},
		{name: "10 elements", a: benchFloats(10), v: benchFloats(10)},
		{name: "100 elements", a: benchFloats(100), v: benchFloats(100)},
		{name: "1_000 elements", a: benchFloats(1_000), v: benchFloats(1_000)},
		{name: "10_000 elements", a: benchFloats(10_000), v: benchFloats(10_000)},
		{name: "100_000 elements", a: benchFloats(100_000), v: benchFloats(100_000)},
		{name: "1_000_000 elements", a: benchFloats(1_000_000), v: benchFloats(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = similarity.CosineSimilarity(bm.a, bm.v)
			}
		})
	}
}

func benchSet(n int) sets.Set[int] {
	elems := make([]int, n)
	for i := range elems {
		elems[i] = i
	}
	return sets.NewHash(elems...)
}

func BenchmarkJaccard(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v sets.Set[int]
	}{
		{name: "3 elements", a: benchSet(3), v: benchSet(3)},
		{name: "10 elements", a: benchSet(10), v: benchSet(10)},
		{name: "100 elements", a: benchSet(100), v: benchSet(100)},
		{name: "1_000 elements", a: benchSet(1_000), v: benchSet(1_000)},
		{name: "10_000 elements", a: benchSet(10_000), v: benchSet(10_000)},
		{name: "100_000 elements", a: benchSet(100_000), v: benchSet(100_000)},
		{name: "1_000_000 elements", a: benchSet(1_000_000), v: benchSet(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = similarity.Jaccard(bm.a, bm.v)
			}
		})
	}
}

func BenchmarkDice(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v sets.Set[int]
	}{
		{name: "3 elements", a: benchSet(3), v: benchSet(3)},
		{name: "10 elements", a: benchSet(10), v: benchSet(10)},
		{name: "100 elements", a: benchSet(100), v: benchSet(100)},
		{name: "1_000 elements", a: benchSet(1_000), v: benchSet(1_000)},
		{name: "10_000 elements", a: benchSet(10_000), v: benchSet(10_000)},
		{name: "100_000 elements", a: benchSet(100_000), v: benchSet(100_000)},
		{name: "1_000_000 elements", a: benchSet(1_000_000), v: benchSet(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = similarity.Dice(bm.a, bm.v)
			}
		})
	}
}

func BenchmarkOverlap(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v sets.Set[int]
	}{
		{name: "3 elements", a: benchSet(3), v: benchSet(3)},
		{name: "10 elements", a: benchSet(10), v: benchSet(10)},
		{name: "100 elements", a: benchSet(100), v: benchSet(100)},
		{name: "1_000 elements", a: benchSet(1_000), v: benchSet(1_000)},
		{name: "10_000 elements", a: benchSet(10_000), v: benchSet(10_000)},
		{name: "100_000 elements", a: benchSet(100_000), v: benchSet(100_000)},
		{name: "1_000_000 elements", a: benchSet(1_000_000), v: benchSet(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = similarity.Overlap(bm.a, bm.v)
			}
		})
	}
}
