package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

// floatSlicesClose compares two []float64 element-wise with floatsClose (NaN/Inf
// aware), also requiring equal length.
func floatSlicesClose(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !floatsClose(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		want  []float64
		ok    bool
	}{
		{name: "evenly spaced", input: []float64{2, 4, 6, 8, 10}, want: []float64{0, 0.25, 0.5, 0.75, 1}, ok: true},
		{name: "already in range", input: []float64{0, 0.5, 1}, want: []float64{0, 0.5, 1}, ok: true},
		{name: "negative values", input: []float64{-10, 0, 10}, want: []float64{0, 0.5, 1}, ok: true},
		{name: "unsorted, min and max not at ends", input: []float64{5, 1, 9, 3}, want: []float64{0.5, 0, 1, 0.25}, ok: true},
		{name: "constant maps to zero", input: []float64{5, 5, 5}, want: []float64{0, 0, 0}, ok: true},
		{name: "single element", input: []float64{7}, want: []float64{0}, ok: true},
		{name: "empty is undefined", input: []float64{}, want: nil, ok: false},
		{name: "nil is undefined", input: nil, want: nil, ok: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stats.Normalize(tc.input)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !floatSlicesClose(got, tc.want) {
				t.Fatalf("normalized = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStandardize(t *testing.T) {
	// canonical has mean 5 and population stddev 2, so z = (x − 5) / 2.
	canonical := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	tests := []struct {
		name  string
		input []float64
		want  []float64
		ok    bool
	}{
		{name: "canonical dataset", input: canonical, want: []float64{-1.5, -0.5, -0.5, -0.5, 0, 0, 1, 2}, ok: true},
		{name: "constant maps to zero", input: []float64{3, 3, 3}, want: []float64{0, 0, 0}, ok: true},
		{name: "single element", input: []float64{7}, want: []float64{0}, ok: true},
		{name: "empty is undefined", input: []float64{}, want: nil, ok: false},
		{name: "nil is undefined", input: nil, want: nil, ok: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stats.Standardize(tc.input)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !floatSlicesClose(got, tc.want) {
				t.Fatalf("standardized = %v, want %v", got, tc.want)
			}
		})
	}
}

// Standardized data has, by construction, a population mean of 0 and a
// population standard deviation of 1.
func TestStandardizeYieldsUnitVariance(t *testing.T) {
	z, ok := stats.Standardize([]float64{3, 1, 4, 1, 5, 9, 2, 6})
	if !ok {
		t.Fatalf("ok = false, want true")
	}

	var sum, sumSq float64
	for _, v := range z {
		sum += v
		sumSq += v * v
	}
	n := float64(len(z))
	mean := sum / n
	popVar := sumSq/n - mean*mean
	if !floatsClose(mean, 0) {
		t.Fatalf("mean of standardized data = %v, want 0", mean)
	}
	if !floatsClose(popVar, 1) {
		t.Fatalf("population variance of standardized data = %v, want 1", popVar)
	}
}

func TestMovingAverage(t *testing.T) {
	tests := []struct {
		name   string
		input  []float64
		window int
		want   []float64
		ok     bool
	}{
		{name: "window of two", input: []float64{1, 2, 3, 4, 5}, window: 2, want: []float64{1.5, 2.5, 3.5, 4.5}, ok: true},
		{name: "window of three", input: []float64{1, 2, 3, 4, 5}, window: 3, want: []float64{2, 3, 4}, ok: true},
		{name: "window of one is identity", input: []float64{1, 2, 3}, window: 1, want: []float64{1, 2, 3}, ok: true},
		{name: "window spanning whole input", input: []float64{2, 4, 6}, window: 3, want: []float64{4}, ok: true},
		{name: "window larger than input is undefined", input: []float64{1, 2, 3}, window: 4, want: nil, ok: false},
		{name: "zero window is invalid", input: []float64{1, 2, 3}, window: 0, want: nil, ok: false},
		{name: "negative window is invalid", input: []float64{1, 2, 3}, window: -1, want: nil, ok: false},
		{name: "empty input is undefined", input: []float64{}, window: 1, want: nil, ok: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stats.MovingAverage(tc.input, tc.window)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !floatSlicesClose(got, tc.want) {
				t.Fatalf("moving average = %v, want %v", got, tc.want)
			}
		})
	}
}

// The transforms are generic over constraints.Numeric, so integer slices work
// without conversion.
func TestTransformsWithIntegers(t *testing.T) {
	if got, ok := stats.Normalize([]int{0, 5, 10}); !ok || !floatSlicesClose(got, []float64{0, 0.5, 1}) {
		t.Fatalf("Normalize = %v, %v; want [0 0.5 1], true", got, ok)
	}
	if got, ok := stats.MovingAverage([]int{2, 4, 6, 8}, 2); !ok || !floatSlicesClose(got, []float64{3, 5, 7}) {
		t.Fatalf("MovingAverage = %v, %v; want [3 5 7], true", got, ok)
	}
}

// Transforms return a fresh slice and never mutate the caller's input.
func TestTransformsDoNotMutateInput(t *testing.T) {
	input := []float64{1, 2, 3, 4}
	original := []float64{1, 2, 3, 4}

	stats.Normalize(input)
	stats.Standardize(input)
	stats.MovingAverage(input, 2)

	if !floatSlicesClose(input, original) {
		t.Fatalf("input mutated to %v, want %v", input, original)
	}
}

// A NaN element propagates into the rescaled output rather than being dropped.
func TestNormalizeNaNPropagates(t *testing.T) {
	got, ok := stats.Normalize([]float64{1, 2, math.NaN(), 4})
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	if !math.IsNaN(got[2]) {
		t.Fatalf("got[2] = %v, want NaN", got[2])
	}
}
