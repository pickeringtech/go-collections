package regression_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/ml/metrics/regression"
)

const epsilon = 1e-9

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

func TestMeanSquaredError(t *testing.T) {
	tests := []struct {
		name   string
		yTrue  []float64
		yPred  []float64
		want   float64
		wantOK bool
	}{
		{
			name:   "typical predictions",
			yTrue:  []float64{3, -0.5, 2, 7},
			yPred:  []float64{2.5, 0, 2, 8},
			want:   0.375,
			wantOK: true,
		},
		{
			name:   "perfect fit is zero error",
			yTrue:  []float64{1, 2, 3},
			yPred:  []float64{1, 2, 3},
			want:   0,
			wantOK: true,
		},
		{
			name:   "single pair",
			yTrue:  []float64{4},
			yPred:  []float64{1},
			want:   9,
			wantOK: true,
		},
		{
			name:   "empty is undefined",
			yTrue:  []float64{},
			yPred:  []float64{},
			wantOK: false,
		},
		{
			name:   "length mismatch is undefined",
			yTrue:  []float64{1, 2, 3},
			yPred:  []float64{1, 2},
			wantOK: false,
		},
		{
			name:   "non-finite is rejected",
			yTrue:  []float64{1, math.Inf(1)},
			yPred:  []float64{1, 2},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := regression.MeanSquaredError(tt.yTrue, tt.yPred)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && !approxEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRootMeanSquaredError(t *testing.T) {
	got, ok := regression.RootMeanSquaredError([]float64{3, -0.5, 2, 7}, []float64{2.5, 0, 2, 8})
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if !approxEqual(got, math.Sqrt(0.375)) {
		t.Errorf("got %v, want %v", got, math.Sqrt(0.375))
	}

	_, ok = regression.RootMeanSquaredError([]float64{1}, []float64{1, 2})
	if ok {
		t.Error("length mismatch: ok = true, want false")
	}
}

func TestMeanAbsoluteError(t *testing.T) {
	tests := []struct {
		name   string
		yTrue  []float64
		yPred  []float64
		want   float64
		wantOK bool
	}{
		{
			name:   "typical predictions",
			yTrue:  []float64{3, -0.5, 2, 7},
			yPred:  []float64{2.5, 0, 2, 8},
			want:   0.5,
			wantOK: true,
		},
		{
			name:   "outlier weighted linearly",
			yTrue:  []float64{0, 0, 0, 0},
			yPred:  []float64{0, 0, 0, 8},
			want:   2,
			wantOK: true,
		},
		{
			name:   "empty is undefined",
			yTrue:  []float64{},
			yPred:  []float64{},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := regression.MeanAbsoluteError(tt.yTrue, tt.yPred)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && !approxEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeanAbsolutePercentageError(t *testing.T) {
	// |（100-90)/100| + |（200-210)/200| + |（300-330)/300| = 0.1 + 0.05 + 0.1 = 0.25, /3
	got, ok := regression.MeanAbsolutePercentageError(
		[]float64{100, 200, 300},
		[]float64{90, 210, 330},
	)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if !approxEqual(got, 0.25/3) {
		t.Errorf("got %v, want %v", got, 0.25/3)
	}

	// A zero true value makes the relative error undefined.
	_, ok = regression.MeanAbsolutePercentageError([]float64{0, 1}, []float64{1, 1})
	if ok {
		t.Error("zero true value: ok = true, want false")
	}

	// Empty and unequal-length inputs are undefined.
	if _, ok := regression.MeanAbsolutePercentageError([]float64{}, []float64{}); ok {
		t.Error("empty: ok = true, want false")
	}
	if _, ok := regression.MeanAbsolutePercentageError([]float64{1, 2}, []float64{1}); ok {
		t.Error("length mismatch: ok = true, want false")
	}
}

func TestRSquared(t *testing.T) {
	tests := []struct {
		name   string
		yTrue  []float64
		yPred  []float64
		want   float64
		wantOK bool
	}{
		{
			name:   "good fit",
			yTrue:  []float64{3, -0.5, 2, 7},
			yPred:  []float64{2.5, 0, 2, 8},
			want:   0.9486081370449679,
			wantOK: true,
		},
		{
			name:   "perfect fit is one",
			yTrue:  []float64{1, 2, 3, 4},
			yPred:  []float64{1, 2, 3, 4},
			want:   1,
			wantOK: true,
		},
		{
			name:   "predicting the mean scores zero",
			yTrue:  []float64{1, 2, 3, 4},
			yPred:  []float64{2.5, 2.5, 2.5, 2.5},
			want:   0,
			wantOK: true,
		},
		{
			name:   "worse than the mean is negative",
			yTrue:  []float64{1, 2, 3},
			yPred:  []float64{3, 2, 1},
			want:   -3,
			wantOK: true,
		},
		{
			name:   "constant target has zero variance and is undefined",
			yTrue:  []float64{5, 5, 5},
			yPred:  []float64{5, 4, 6},
			wantOK: false,
		},
		{
			name:   "empty is undefined",
			yTrue:  []float64{},
			yPred:  []float64{},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := regression.RSquared(tt.yTrue, tt.yPred)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && !approxEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegerInput(t *testing.T) {
	// The metrics are generic over constraints.Numeric, so integer series work.
	got, ok := regression.MeanSquaredError([]int{1, 2, 3}, []int{1, 2, 5})
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if !approxEqual(got, 4.0/3) {
		t.Errorf("got %v, want %v", got, 4.0/3)
	}
}

func TestNoMutation(t *testing.T) {
	yTrue := []float64{3, -0.5, 2, 7}
	yPred := []float64{2.5, 0, 2, 8}
	trueCopy := append([]float64(nil), yTrue...)
	predCopy := append([]float64(nil), yPred...)

	_, _ = regression.MeanSquaredError(yTrue, yPred)
	_, _ = regression.RSquared(yTrue, yPred)

	for i := range yTrue {
		if yTrue[i] != trueCopy[i] || yPred[i] != predCopy[i] {
			t.Fatal("input slice was mutated")
		}
	}
}
