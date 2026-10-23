package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

// The canonical dataset (mean 5): central moments m2=4, m3=5.25, m4=44.5, so
// skewness = 5.25/4^1.5 = 0.65625 and excess kurtosis = 44.5/16 - 3 = -0.21875.
var shapeData = []float64{2, 4, 4, 4, 5, 5, 7, 9}

func TestSkewness(t *testing.T) {
	t.Run("canonical dataset", func(t *testing.T) {
		got, ok := stats.Skewness(shapeData)
		if !ok || !approxEqual(got, 0.65625) {
			t.Fatalf("Skewness = %v, %v; want 0.65625, true", got, ok)
		}
	})

	t.Run("symmetric sample is zero", func(t *testing.T) {
		got, ok := stats.Skewness([]float64{1, 2, 3, 4, 5})
		if !ok || !approxEqual(got, 0) {
			t.Fatalf("Skewness = %v, %v; want 0, true", got, ok)
		}
	})

	t.Run("rejects undefined inputs", func(t *testing.T) {
		cases := map[string][]float64{
			"empty":    nil,
			"single":   {42},
			"constant": {7, 7, 7},
		}
		for name, in := range cases {
			t.Run(name, func(t *testing.T) {
				got, ok := stats.Skewness(in)
				if ok {
					t.Errorf("ok = true (%v), want false", got)
				}
			})
		}
	})

	t.Run("non-finite propagates", func(t *testing.T) {
		got, ok := stats.Skewness([]float64{1, 2, math.NaN()})
		if !ok || !math.IsNaN(got) {
			t.Fatalf("Skewness = %v, %v; want NaN, true", got, ok)
		}
	})
}

// TestShapeRejectsInexactConstantData guards the constant-data check against
// round-off: values like 0.1 are not exactly representable, so a sum/n mean
// drifts an ulp away from them and leaves a tiny non-zero m2 that would turn
// 0/0 into a spurious finite skewness or kurtosis.
func TestShapeRejectsInexactConstantData(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
	}{
		{name: "three tenths", input: []float64{0.1, 0.1, 0.1}},
		{name: "many thirds", input: []float64{1.0 / 3, 1.0 / 3, 1.0 / 3, 1.0 / 3, 1.0 / 3, 1.0 / 3, 1.0 / 3}},
		{name: "large inexact constant", input: []float64{1e16 + 2, 1e16 + 2, 1e16 + 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skew, ok := stats.Skewness(tt.input)
			if ok {
				t.Errorf("Skewness(%v) = %v, true; want ok == false", tt.input, skew)
			}
			kurt, ok := stats.Kurtosis(tt.input)
			if ok {
				t.Errorf("Kurtosis(%v) = %v, true; want ok == false", tt.input, kurt)
			}
		})
	}
}

// TestShapeConstantNonFinitePropagates checks that the constant-data rejection
// does not swallow non-finite input: identical infinities still propagate a
// non-finite result with ok == true, per the variance family's policy.
func TestShapeConstantNonFinitePropagates(t *testing.T) {
	input := []float64{math.Inf(1), math.Inf(1)}
	skew, ok := stats.Skewness(input)
	if !ok || !math.IsNaN(skew) {
		t.Errorf("Skewness = %v, %v; want NaN, true", skew, ok)
	}
	kurt, ok := stats.Kurtosis(input)
	if !ok || !math.IsNaN(kurt) {
		t.Errorf("Kurtosis = %v, %v; want NaN, true", kurt, ok)
	}
}

func TestKurtosis(t *testing.T) {
	t.Run("canonical dataset", func(t *testing.T) {
		got, ok := stats.Kurtosis(shapeData)
		if !ok || !approxEqual(got, -0.21875) {
			t.Fatalf("Kurtosis = %v, %v; want -0.21875, true", got, ok)
		}
	})

	t.Run("rejects undefined inputs", func(t *testing.T) {
		cases := map[string][]float64{
			"empty":    nil,
			"single":   {42},
			"constant": {3, 3, 3, 3},
		}
		for name, in := range cases {
			t.Run(name, func(t *testing.T) {
				got, ok := stats.Kurtosis(in)
				if ok {
					t.Errorf("ok = true (%v), want false", got)
				}
			})
		}
	})

	t.Run("non-finite propagates", func(t *testing.T) {
		got, ok := stats.Kurtosis([]float64{1, 2, 3, math.Inf(1)})
		if !ok || !math.IsInf(got, 0) && !math.IsNaN(got) {
			t.Fatalf("Kurtosis = %v, %v; want non-finite, true", got, ok)
		}
	})
}
