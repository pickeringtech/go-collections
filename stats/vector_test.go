package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestDot(t *testing.T) {
	got, ok := stats.Dot([]float64{1, 2, 3}, []float64{4, 5, 6})
	if !ok || !approxEqual(got, 32) {
		t.Fatalf("Dot = %v, %v; want 32, true", got, ok)
	}

	t.Run("rejects empty and mismatched lengths", func(t *testing.T) {
		if _, ok := stats.Dot([]float64{}, []float64{}); ok {
			t.Errorf("empty reported ok")
		}
		if _, ok := stats.Dot([]float64{1, 2}, []float64{1}); ok {
			t.Errorf("mismatched lengths reported ok")
		}
	})

	t.Run("non-finite propagates", func(t *testing.T) {
		got, ok := stats.Dot([]float64{1, math.NaN()}, []float64{1, 1})
		if !ok || !math.IsNaN(got) {
			t.Fatalf("Dot = %v, %v; want NaN, true", got, ok)
		}
	})
}

func TestNorm(t *testing.T) {
	got, ok := stats.Norm([]float64{3, 4})
	if !ok || !approxEqual(got, 5) {
		t.Fatalf("Norm = %v, %v; want 5, true", got, ok)
	}
	if _, ok := stats.Norm([]int{}); ok {
		t.Errorf("empty reported ok")
	}
}

func TestEuclideanDistance(t *testing.T) {
	got, ok := stats.EuclideanDistance([]float64{0, 0}, []float64{3, 4})
	if !ok || !approxEqual(got, 5) {
		t.Fatalf("EuclideanDistance = %v, %v; want 5, true", got, ok)
	}

	t.Run("identical points are distance zero", func(t *testing.T) {
		d, ok := stats.EuclideanDistance([]int{1, 2, 3}, []int{1, 2, 3})
		if !ok || d != 0 {
			t.Fatalf("EuclideanDistance = %v, %v; want 0, true", d, ok)
		}
	})

	t.Run("rejects empty and mismatched lengths", func(t *testing.T) {
		if _, ok := stats.EuclideanDistance([]float64{}, []float64{}); ok {
			t.Errorf("empty reported ok")
		}
		if _, ok := stats.EuclideanDistance([]float64{1}, []float64{1, 2}); ok {
			t.Errorf("mismatched lengths reported ok")
		}
	})
}

func TestCosineSimilarity(t *testing.T) {
	cases := map[string]struct {
		a, b []float64
		want float64
	}{
		"identical direction": {[]float64{1, 1}, []float64{2, 2}, 1},
		"orthogonal":          {[]float64{1, 0}, []float64{0, 1}, 0},
		"opposite":            {[]float64{1, 0}, []float64{-1, 0}, -1},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, ok := stats.CosineSimilarity(tc.a, tc.b)
			if !ok || !approxEqual(got, tc.want) {
				t.Fatalf("CosineSimilarity = %v, %v; want %v, true", got, ok, tc.want)
			}
		})
	}

	t.Run("rejects undefined inputs", func(t *testing.T) {
		if _, ok := stats.CosineSimilarity([]float64{}, []float64{}); ok {
			t.Errorf("empty reported ok")
		}
		if _, ok := stats.CosineSimilarity([]float64{1, 2}, []float64{1}); ok {
			t.Errorf("mismatched lengths reported ok")
		}
		if _, ok := stats.CosineSimilarity([]float64{0, 0}, []float64{1, 2}); ok {
			t.Errorf("zero vector reported ok")
		}
	})

	t.Run("non-finite propagates", func(t *testing.T) {
		got, ok := stats.CosineSimilarity([]float64{1, math.Inf(1)}, []float64{1, 1})
		if !ok || (!math.IsNaN(got) && !math.IsInf(got, 0)) {
			t.Fatalf("CosineSimilarity = %v, %v; want non-finite, true", got, ok)
		}
	})
}
