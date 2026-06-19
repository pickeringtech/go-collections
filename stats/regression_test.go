package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestLinearRegression(t *testing.T) {
	const eps = 1e-9

	t.Run("fits a noisy line", func(t *testing.T) {
		fit, ok := stats.LinearRegression([]float64{1, 2, 3, 4, 5}, []float64{2, 4, 5, 4, 6})
		if !ok {
			t.Fatalf("ok = false, want true")
		}
		if math.Abs(fit.Slope-0.8) > eps {
			t.Errorf("Slope = %v, want 0.8", fit.Slope)
		}
		if math.Abs(fit.Intercept-1.8) > eps {
			t.Errorf("Intercept = %v, want 1.8", fit.Intercept)
		}
		if math.Abs(fit.R2-0.7272727272727273) > eps {
			t.Errorf("R2 = %v, want 0.727...", fit.R2)
		}
	})

	t.Run("perfect line has R2 of 1", func(t *testing.T) {
		// y = 3x - 1 exactly.
		fit, ok := stats.LinearRegression([]int{0, 1, 2, 3}, []int{-1, 2, 5, 8})
		if !ok {
			t.Fatalf("ok = false, want true")
		}
		if math.Abs(fit.Slope-3) > eps || math.Abs(fit.Intercept+1) > eps {
			t.Errorf("fit = %+v, want slope 3 intercept -1", fit)
		}
		if math.Abs(fit.R2-1) > eps {
			t.Errorf("R2 = %v, want 1", fit.R2)
		}
	})

	t.Run("Predict and residuals", func(t *testing.T) {
		fit, _ := stats.LinearRegression([]float64{1, 2, 3, 4, 5}, []float64{2, 4, 5, 4, 6})
		// Fitted value at x=3 is 0.8*3 + 1.8 = 4.2; residual against observed 5.
		if got := fit.Predict(3); math.Abs(got-4.2) > eps {
			t.Errorf("Predict(3) = %v, want 4.2", got)
		}
		if residual := 5.0 - fit.Predict(3); math.Abs(residual-0.8) > eps {
			t.Errorf("residual = %v, want 0.8", residual)
		}
	})

	t.Run("rejects undefined inputs", func(t *testing.T) {
		tests := map[string]struct {
			x, y []float64
		}{
			"empty":           {nil, nil},
			"single point":    {[]float64{1}, []float64{2}},
			"length mismatch": {[]float64{1, 2, 3}, []float64{1, 2}},
			"constant x":      {[]float64{4, 4, 4}, []float64{1, 2, 3}},
			"constant y":      {[]float64{1, 2, 3}, []float64{7, 7, 7}},
		}
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				if fit, ok := stats.LinearRegression(tc.x, tc.y); ok {
					t.Errorf("ok = true (fit %+v), want false", fit)
				}
			})
		}
	})

	t.Run("non-finite propagates", func(t *testing.T) {
		fit, ok := stats.LinearRegression([]float64{1, 2, math.NaN()}, []float64{1, 2, 3})
		if !ok {
			t.Fatalf("ok = false, want true (NaN propagates)")
		}
		if !math.IsNaN(fit.Slope) {
			t.Errorf("Slope = %v, want NaN", fit.Slope)
		}
	})
}
