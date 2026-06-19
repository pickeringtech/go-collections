package preprocessing_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func TestFixedWidthBinner(t *testing.T) {
	// Range [0, 8] into 4 bins of width 2: cut points at 2, 4, 6.
	binner := preprocessing.NewFixedWidthBinner(4).Fit([]float64{0, 8})
	if !floatSlicesClose(binner.Edges(), []float64{2, 4, 6}) {
		t.Fatalf("Edges() = %v, want [2 4 6]", binner.Edges())
	}

	got, ok := binner.Transform([]float64{0, 1, 2, 3, 5, 7, 8})
	want := []int{0, 0, 1, 1, 2, 3, 3}
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("Transform = (%v, %v), want (%v, true)", got, ok, want)
	}
}

func TestFixedWidthBinnerOutOfRangeClampsToEnds(t *testing.T) {
	binner := preprocessing.NewFixedWidthBinner(4).Fit([]float64{0, 8})
	got, _ := binner.Transform([]float64{-5, 100})
	if !reflect.DeepEqual(got, []int{0, 3}) {
		t.Fatalf("Transform out-of-range = %v, want [0 3]", got)
	}
}

func TestFixedWidthBinnerDegenerateRange(t *testing.T) {
	// All-equal data collapses to a single bin (index 0).
	binner := preprocessing.NewFixedWidthBinner(4).Fit([]float64{7, 7, 7})
	got, ok := binner.Transform([]float64{7, 7})
	if !ok || !reflect.DeepEqual(got, []int{0, 0}) {
		t.Fatalf("Transform = (%v, %v), want ([0 0], true)", got, ok)
	}
}

func TestFixedWidthBinnerInvalidBinCount(t *testing.T) {
	if _, ok := preprocessing.NewFixedWidthBinner(0).Fit([]float64{1, 2}).Transform([]float64{1}); ok {
		t.Fatalf("nBins=0 reported fitted")
	}
}

func TestFixedWidthBinnerUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewFixedWidthBinner(3).Transform([]float64{1}); ok {
		t.Fatalf("unfitted Transform reported ok")
	}
}

func TestQuantileBinner(t *testing.T) {
	// Two equal-population bins of [1,2,3,4] split at the median (2.5).
	binner := preprocessing.NewQuantileBinner(2).Fit([]float64{1, 2, 3, 4})
	if !floatSlicesClose(binner.Edges(), []float64{2.5}) {
		t.Fatalf("Edges() = %v, want [2.5]", binner.Edges())
	}

	got, ok := binner.Transform([]float64{1, 2, 3, 4})
	want := []int{0, 0, 1, 1}
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("Transform = (%v, %v), want (%v, true)", got, ok, want)
	}
}

func TestQuantileBinnerSingleBin(t *testing.T) {
	binner := preprocessing.NewQuantileBinner(1).Fit([]float64{1, 2, 3})
	got, ok := binner.Transform([]float64{1, 2, 3})
	if !ok || !reflect.DeepEqual(got, []int{0, 0, 0}) {
		t.Fatalf("Transform = (%v, %v), want ([0 0 0], true)", got, ok)
	}
}

func TestQuantileBinnerRejectsNonFinite(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1)} {
		if _, ok := preprocessing.NewQuantileBinner(2).Fit([]float64{1, 2, bad}).Transform([]float64{1}); ok {
			t.Fatalf("Fit accepted non-finite %v", bad)
		}
	}
}

func TestQuantileBinnerEmptyLeavesUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewQuantileBinner(2).Fit(nil).Transform([]float64{1}); ok {
		t.Fatalf("empty Fit reported fitted")
	}
}

func TestBinnersDoNotMutateInput(t *testing.T) {
	input := []float64{4, 1, 3, 2}
	original := []float64{4, 1, 3, 2}
	preprocessing.NewFixedWidthBinner(3).Fit(input).Transform(input)
	preprocessing.NewQuantileBinner(3).Fit(input).Transform(input)
	if !floatSlicesClose(input, original) {
		t.Fatalf("input mutated to %v, want %v", input, original)
	}
}

// The fitted binners satisfy the Transformer contract.
var (
	_ preprocessing.Transformer[float64, []int] = (*preprocessing.FixedWidthBinner)(nil)
	_ preprocessing.Transformer[float64, []int] = (*preprocessing.QuantileBinner)(nil)
)
