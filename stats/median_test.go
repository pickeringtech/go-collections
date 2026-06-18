package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestMedian(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		want   float64
		wantOK bool
	}{
		{name: "odd length picks middle", input: []int{3, 1, 2}, want: 2, wantOK: true},
		{name: "even length averages middle two", input: []int{1, 2, 3, 4}, want: 2.5, wantOK: true},
		{name: "unsorted input", input: []int{5, 1, 4, 2, 3}, want: 3, wantOK: true},
		{name: "single value", input: []int{42}, want: 42, wantOK: true},
		{name: "even with averaging to fraction", input: []int{1, 2}, want: 1.5, wantOK: true},
		{name: "empty input", input: []int{}, wantOK: false},
		{name: "nil input", input: nil, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.Median(tt.input)
			if ok != tt.wantOK || (ok && !approxEqual(got, tt.want)) {
				t.Fatalf("Median(%v) = (%v, %v), want (%v, %v)", tt.input, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestMedianDoesNotMutateInput(t *testing.T) {
	input := []int{3, 1, 2}
	_, _ = stats.Median(input)
	for i, v := range []int{3, 1, 2} {
		if input[i] != v {
			t.Fatalf("Median mutated input to %v", input)
		}
	}
}

func TestMedianRejectsNonFinite(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, ok := stats.Median([]float64{1, bad, 3}); ok {
			t.Fatalf("Median with %v should be ok=false", bad)
		}
	}
}
