package stats_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

// approxEqual compares two float64s with a small tolerance. Interpolated
// quantiles (e.g. 4.6) are not always bit-exact, so result comparisons use this
// rather than ==.
func approxEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}

func ExampleQuantile() {
	data := []float64{1, 2, 3, 4, 5}

	med, _ := stats.Quantile(data, 0.5)
	p90, _ := stats.Percentile(data, 90)

	fmt.Printf("median: %v, p90: %v", med, p90)
	// Output: median: 3, p90: 4.6
}

func ExampleQuartiles() {
	data := []float64{1, 2, 3, 4, 5}

	qs, _ := stats.Quartiles(data)
	iqr, _ := stats.IQR(data)

	fmt.Printf("Q1: %v, Q2: %v, Q3: %v, IQR: %v", qs.Q1, qs.Q2, qs.Q3, iqr)
	// Output: Q1: 2, Q2: 3, Q3: 4, IQR: 2
}

func TestQuantile(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		q     float64
		want  float64
		wOK   bool
	}{
		{name: "min", input: []float64{1, 2, 3, 4, 5}, q: 0, want: 1, wOK: true},
		{name: "max", input: []float64{1, 2, 3, 4, 5}, q: 1, want: 5, wOK: true},
		{name: "median odd", input: []float64{1, 2, 3, 4, 5}, q: 0.5, want: 3, wOK: true},
		{name: "median even", input: []float64{1, 2, 3, 4}, q: 0.5, want: 2.5, wOK: true},
		{name: "type-7 interpolated", input: []float64{1, 2, 3, 4, 5}, q: 0.9, want: 4.6, wOK: true},
		{name: "unsorted input", input: []float64{5, 3, 1, 4, 2}, q: 0.25, want: 2, wOK: true},
		{name: "single element", input: []float64{42}, q: 0.5, want: 42, wOK: true},
		{name: "empty", input: []float64{}, q: 0.5, want: 0, wOK: false},
		{name: "nil", input: nil, q: 0.5, want: 0, wOK: false},
		{name: "q below range", input: []float64{1, 2, 3}, q: -0.1, want: 0, wOK: false},
		{name: "q above range", input: []float64{1, 2, 3}, q: 1.1, want: 0, wOK: false},
		{name: "q is NaN", input: []float64{1, 2, 3}, q: math.NaN(), want: 0, wOK: false},
		{name: "NaN poisons", input: []float64{1, 2, math.NaN(), 4}, q: 0.5, want: 0, wOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.Quantile(tt.input, tt.q)
			if !approxEqual(got, tt.want) || ok != tt.wOK {
				t.Errorf("Quantile(%v, %v) = (%v, %v), want (%v, %v)", tt.input, tt.q, got, ok, tt.want, tt.wOK)
			}
		})
	}
}

func TestQuantileIntegerInput(t *testing.T) {
	got, ok := stats.Quantile([]int{10, 20, 30, 40}, 0.5)
	if got != 25 || !ok {
		t.Errorf("Quantile([10 20 30 40], 0.5) = (%v, %v), want (25, true)", got, ok)
	}
}

func TestQuantileWith(t *testing.T) {
	// data sorts to [1 2 3 4 5]; q=0.9 -> continuous rank 3.6 (between 4 and 5).
	input := []float64{1, 2, 3, 4, 5}
	tests := []struct {
		name   string
		method stats.InterpolationMethod
		want   float64
	}{
		{name: "linear", method: stats.Linear, want: 4.6},
		{name: "lower", method: stats.Lower, want: 4},
		{name: "higher", method: stats.Higher, want: 5},
		{name: "nearest rounds up", method: stats.Nearest, want: 5},
		{name: "midpoint", method: stats.Midpoint, want: 4.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.QuantileWith(input, 0.9, tt.method)
			if !approxEqual(got, tt.want) || !ok {
				t.Errorf("QuantileWith(%v, 0.9, %v) = (%v, %v), want (%v, true)", input, tt.name, got, ok, tt.want)
			}
		})
	}
}

func TestQuantileWithNearestRoundsDown(t *testing.T) {
	// q=0.6 on [1 2 3 4 5] -> continuous rank 2.4, frac 0.4 < 0.5, so Nearest
	// takes the lower sample (3 rather than 4).
	got, ok := stats.QuantileWith([]float64{1, 2, 3, 4, 5}, 0.6, stats.Nearest)
	if !approxEqual(got, 3) || !ok {
		t.Errorf("QuantileWith(..., 0.6, Nearest) = (%v, %v), want (3, true)", got, ok)
	}
}

func TestPercentile(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		p     float64
		want  float64
		wOK   bool
	}{
		{name: "p50 median", input: []float64{1, 2, 3, 4, 5}, p: 50, want: 3, wOK: true},
		{name: "p90", input: []float64{1, 2, 3, 4, 5}, p: 90, want: 4.6, wOK: true},
		{name: "p0", input: []float64{1, 2, 3, 4, 5}, p: 0, want: 1, wOK: true},
		{name: "p100", input: []float64{1, 2, 3, 4, 5}, p: 100, want: 5, wOK: true},
		{name: "p below range", input: []float64{1, 2, 3}, p: -1, want: 0, wOK: false},
		{name: "p above range", input: []float64{1, 2, 3}, p: 101, want: 0, wOK: false},
		{name: "p is NaN", input: []float64{1, 2, 3}, p: math.NaN(), want: 0, wOK: false},
		{name: "empty", input: nil, p: 50, want: 0, wOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.Percentile(tt.input, tt.p)
			if !approxEqual(got, tt.want) || ok != tt.wOK {
				t.Errorf("Percentile(%v, %v) = (%v, %v), want (%v, %v)", tt.input, tt.p, got, ok, tt.want, tt.wOK)
			}
		})
	}
}

func TestPercentileWith(t *testing.T) {
	got, ok := stats.PercentileWith([]float64{1, 2, 3, 4, 5}, 90, stats.Midpoint)
	if !approxEqual(got, 4.5) || !ok {
		t.Errorf("PercentileWith(..., 90, Midpoint) = (%v, %v), want (4.5, true)", got, ok)
	}
}

func TestQuartiles(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		want  stats.QuartileSet
		wOK   bool
	}{
		{name: "odd length", input: []float64{1, 2, 3, 4, 5}, want: stats.QuartileSet{Q1: 2, Q2: 3, Q3: 4}, wOK: true},
		{name: "even length", input: []float64{1, 2, 3, 4}, want: stats.QuartileSet{Q1: 1.75, Q2: 2.5, Q3: 3.25}, wOK: true},
		{name: "empty", input: nil, want: stats.QuartileSet{}, wOK: false},
		{name: "NaN poisons", input: []float64{1, math.NaN(), 3}, want: stats.QuartileSet{}, wOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.Quartiles(tt.input)
			if !reflect.DeepEqual(got, tt.want) || ok != tt.wOK {
				t.Errorf("Quartiles(%v) = (%+v, %v), want (%+v, %v)", tt.input, got, ok, tt.want, tt.wOK)
			}
		})
	}
}

func TestIQR(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		want  float64
		wOK   bool
	}{
		{name: "odd length", input: []float64{1, 2, 3, 4, 5}, want: 2, wOK: true},
		{name: "even length", input: []float64{1, 2, 3, 4}, want: 1.5, wOK: true},
		{name: "empty", input: nil, want: 0, wOK: false},
		{name: "NaN poisons", input: []float64{1, math.NaN(), 3}, want: 0, wOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.IQR(tt.input)
			if got != tt.want || ok != tt.wOK {
				t.Errorf("IQR(%v) = (%v, %v), want (%v, %v)", tt.input, got, ok, tt.want, tt.wOK)
			}
		})
	}
}

// TestNoMutation guards the ownership-isolation contract: the caller's slice
// must be untouched even though the implementation sorts internally.
func TestNoMutation(t *testing.T) {
	input := []float64{5, 3, 1, 4, 2}
	original := make([]float64, len(input))
	copy(original, input)

	stats.Quantile(input, 0.5)
	stats.Quartiles(input)
	stats.IQR(input)

	if !reflect.DeepEqual(input, original) {
		t.Errorf("input was mutated: got %v, want %v", input, original)
	}
}
