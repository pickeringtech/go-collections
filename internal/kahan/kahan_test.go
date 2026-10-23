package kahan_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/internal/kahan"
)

func sumOf(values []float64) float64 {
	var total kahan.Sum
	for _, x := range values {
		total.Add(x)
	}
	return total.Total()
}

func TestSum(t *testing.T) {
	inf := math.Inf(1)
	tests := []struct {
		name  string
		input []float64
		want  float64
	}{
		{name: "nil input is an empty sum", input: nil, want: 0},
		{name: "empty input is an empty sum", input: []float64{}, want: 0},
		{name: "small integers are exact", input: []float64{1, 2, 3}, want: 6},
		{name: "recovers round-off a naive sum loses", input: []float64{1, 1e100, 1, -1e100}, want: 2},
		{name: "a term larger than the total keeps the absorbed bits", input: []float64{1e17, 1, 2, -1e17}, want: 3},
		{name: "positive infinity stays infinite", input: []float64{inf, 1}, want: inf},
		{name: "negative infinity stays infinite", input: []float64{-inf, 1, 2}, want: -inf},
		{name: "overflow stays infinite", input: []float64{math.MaxFloat64, math.MaxFloat64, 1}, want: inf},
		{name: "opposite infinities are NaN", input: []float64{inf, -inf}, want: math.NaN()},
		{name: "NaN propagates", input: []float64{1, math.NaN(), 2}, want: math.NaN()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sumOf(tt.input)
			same := got == tt.want || (math.IsNaN(got) && math.IsNaN(tt.want))
			if !same {
				t.Errorf("Total() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumIsMoreAccurateThanNaive(t *testing.T) {
	// 0.1 is not exactly representable, so a million naive additions drift far
	// from the correctly rounded total; the compensated sum stays within an ulp.
	const n = 1_000_000
	var naive float64
	var total kahan.Sum
	for i := 0; i < n; i++ {
		naive += 0.1
		total.Add(0.1)
	}
	want := 100_000.0
	got := total.Total()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("Total() = %v, want %v", got, want)
	}
	if math.Abs(naive-want) <= math.Abs(got-want) {
		t.Errorf("naive sum %v was not less accurate than compensated %v", naive, got)
	}
}
