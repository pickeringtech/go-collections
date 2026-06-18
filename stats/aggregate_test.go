package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestProduct(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		want   int
		wantOK bool
	}{
		{name: "several values", input: []int{2, 3, 4}, want: 24, wantOK: true},
		{name: "single value", input: []int{7}, want: 7, wantOK: true},
		{name: "contains zero", input: []int{2, 0, 5}, want: 0, wantOK: true},
		{name: "negatives multiply", input: []int{-2, 3, -4}, want: 24, wantOK: true},
		{name: "empty input", input: []int{}, wantOK: false},
		{name: "nil input", input: nil, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.Product(tt.input)
			if ok != tt.wantOK || (ok && got != tt.want) {
				t.Fatalf("Product(%v) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestProductFloat(t *testing.T) {
	got, ok := stats.Product([]float64{0.5, 4, 2})
	if !ok || got != 4 {
		t.Fatalf("Product = (%v, %v), want (4, true)", got, ok)
	}
}

func TestProductRejectsNonFinite(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, ok := stats.Product([]float64{1, bad, 3}); ok {
			t.Fatalf("Product with %v should be ok=false", bad)
		}
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		want   int
		wantOK bool
	}{
		{name: "spread", input: []int{3, 1, 4, 1, 5}, want: 4, wantOK: true},
		{name: "single value has zero range", input: []int{9}, want: 0, wantOK: true},
		{name: "all equal", input: []int{2, 2, 2}, want: 0, wantOK: true},
		{name: "negatives", input: []int{-5, 0, 5}, want: 10, wantOK: true},
		{name: "empty input", input: []int{}, wantOK: false},
		{name: "nil input", input: nil, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.Range(tt.input)
			if ok != tt.wantOK || (ok && got != tt.want) {
				t.Fatalf("Range(%v) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestRangeRejectsNonFinite(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, ok := stats.Range([]float64{1, bad, 3}); ok {
			t.Fatalf("Range with %v should be ok=false", bad)
		}
	}
}

func TestCumulativeSum(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{name: "prefix sums", input: []int{3, 1, 4, 1, 5}, want: []int{3, 4, 8, 9, 14}},
		{name: "single value", input: []int{7}, want: []int{7}},
		{name: "negatives", input: []int{1, -1, 2}, want: []int{1, 0, 2}},
		{name: "empty input", input: []int{}, want: []int{}},
		{name: "nil input", input: nil, want: []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stats.CumulativeSum(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("CumulativeSum(%v) = %v, want %v", tt.input, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("CumulativeSum(%v) = %v, want %v", tt.input, got, tt.want)
				}
			}
		})
	}
}

func TestCumulativeSumDoesNotMutateInput(t *testing.T) {
	input := []int{1, 2, 3}
	_ = stats.CumulativeSum(input)
	for i, v := range []int{1, 2, 3} {
		if input[i] != v {
			t.Fatalf("CumulativeSum mutated input to %v", input)
		}
	}
}
