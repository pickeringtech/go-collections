package distance_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/distance"
)

func benchVec(n int) []float64 {
	v := make([]float64, n)
	for i := range v {
		v[i] = float64(i%100) + 1
	}
	return v
}

func BenchmarkEuclidean(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v []float64
	}{
		{name: "3 elements", a: benchVec(3), v: benchVec(3)},
		{name: "10 elements", a: benchVec(10), v: benchVec(10)},
		{name: "100 elements", a: benchVec(100), v: benchVec(100)},
		{name: "1_000 elements", a: benchVec(1_000), v: benchVec(1_000)},
		{name: "10_000 elements", a: benchVec(10_000), v: benchVec(10_000)},
		{name: "100_000 elements", a: benchVec(100_000), v: benchVec(100_000)},
		{name: "1_000_000 elements", a: benchVec(1_000_000), v: benchVec(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = distance.Euclidean(bm.a, bm.v)
			}
		})
	}
}

func BenchmarkManhattan(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v []float64
	}{
		{name: "3 elements", a: benchVec(3), v: benchVec(3)},
		{name: "10 elements", a: benchVec(10), v: benchVec(10)},
		{name: "100 elements", a: benchVec(100), v: benchVec(100)},
		{name: "1_000 elements", a: benchVec(1_000), v: benchVec(1_000)},
		{name: "10_000 elements", a: benchVec(10_000), v: benchVec(10_000)},
		{name: "100_000 elements", a: benchVec(100_000), v: benchVec(100_000)},
		{name: "1_000_000 elements", a: benchVec(1_000_000), v: benchVec(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = distance.Manhattan(bm.a, bm.v)
			}
		})
	}
}

func BenchmarkMinkowski(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v []float64
	}{
		{name: "3 elements", a: benchVec(3), v: benchVec(3)},
		{name: "10 elements", a: benchVec(10), v: benchVec(10)},
		{name: "100 elements", a: benchVec(100), v: benchVec(100)},
		{name: "1_000 elements", a: benchVec(1_000), v: benchVec(1_000)},
		{name: "10_000 elements", a: benchVec(10_000), v: benchVec(10_000)},
		{name: "100_000 elements", a: benchVec(100_000), v: benchVec(100_000)},
		{name: "1_000_000 elements", a: benchVec(1_000_000), v: benchVec(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = distance.Minkowski(bm.a, bm.v, 3)
			}
		})
	}
}

func BenchmarkCosineDistance(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v []float64
	}{
		{name: "3 elements", a: benchVec(3), v: benchVec(3)},
		{name: "10 elements", a: benchVec(10), v: benchVec(10)},
		{name: "100 elements", a: benchVec(100), v: benchVec(100)},
		{name: "1_000 elements", a: benchVec(1_000), v: benchVec(1_000)},
		{name: "10_000 elements", a: benchVec(10_000), v: benchVec(10_000)},
		{name: "100_000 elements", a: benchVec(100_000), v: benchVec(100_000)},
		{name: "1_000_000 elements", a: benchVec(1_000_000), v: benchVec(1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = distance.CosineDistance(bm.a, bm.v)
			}
		})
	}
}

func BenchmarkHamming(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v []int
	}{
		{name: "3 elements", a: make([]int, 3), v: make([]int, 3)},
		{name: "10 elements", a: make([]int, 10), v: make([]int, 10)},
		{name: "100 elements", a: make([]int, 100), v: make([]int, 100)},
		{name: "1_000 elements", a: make([]int, 1_000), v: make([]int, 1_000)},
		{name: "10_000 elements", a: make([]int, 10_000), v: make([]int, 10_000)},
		{name: "100_000 elements", a: make([]int, 100_000), v: make([]int, 100_000)},
		{name: "1_000_000 elements", a: make([]int, 1_000_000), v: make([]int, 1_000_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = distance.Hamming(bm.a, bm.v)
			}
		})
	}
}

func BenchmarkLevenshtein(b *testing.B) {
	benchmarks := []struct {
		name string
		a, v string
	}{
		{name: "3 chars", a: "cat", v: "bat"},
		{name: "10 chars", a: "abcdefghij", v: "jihgfedcba"},
		{name: "100 chars", a: repeatStr("abcdefghij", 10), v: repeatStr("jihgfedcba", 10)},
		{name: "1_000 chars", a: repeatStr("abcdefghij", 100), v: repeatStr("jihgfedcba", 100)},
		{name: "10_000 chars", a: repeatStr("abcdefghij", 1_000), v: repeatStr("jihgfedcba", 1_000)},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = distance.Levenshtein(bm.a, bm.v)
			}
		})
	}
}

func repeatStr(s string, n int) string {
	result := make([]byte, len(s)*n)
	for i := 0; i < n; i++ {
		copy(result[i*len(s):], s)
	}
	return string(result)
}
