package preprocessing_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func TestFixedWidthBinner(t *testing.T) {
	tests := []struct {
		name      string
		nBins     int
		input     []float64
		wantEdges []float64
		want      []int
	}{
		{
			// Range [0, 8] into 4 bins of width 2: cut points at 2, 4, 6.
			name:      "four bins over a range of eight",
			nBins:     4,
			input:     []float64{0, 1, 2, 3, 5, 7, 8},
			wantEdges: []float64{2, 4, 6},
			want:      []int{0, 0, 1, 1, 2, 3, 3},
		},
		{
			// nBins = 1 is the smallest valid count: no cut points, one bin.
			name:      "a single bin holds every value",
			nBins:     1,
			input:     []float64{0, 1, 4, 8},
			wantEdges: []float64{},
			want:      []int{0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			binner := preprocessing.NewFixedWidthBinner(tt.nBins).Fit([]float64{0, 8})
			if !floatSlicesClose(binner.Edges(), tt.wantEdges) {
				t.Fatalf("Edges() = %v, want %v", binner.Edges(), tt.wantEdges)
			}
			got, ok := binner.Transform(tt.input)
			if !ok || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Transform = (%v, %v), want (%v, true)", got, ok, tt.want)
			}
		})
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
	_, ok := preprocessing.NewFixedWidthBinner(0).Fit([]float64{1, 2}).Transform([]float64{1})
	if ok {
		t.Fatalf("nBins=0 reported fitted")
	}
}

func TestFixedWidthBinnerUnfitted(t *testing.T) {
	_, ok := preprocessing.NewFixedWidthBinner(3).Transform([]float64{1})
	if ok {
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
		_, ok := preprocessing.NewQuantileBinner(2).Fit([]float64{1, 2, bad}).Transform([]float64{1})
		if ok {
			t.Fatalf("Fit accepted non-finite %v", bad)
		}
	}
}

func TestQuantileBinnerEmptyLeavesUnfitted(t *testing.T) {
	_, ok := preprocessing.NewQuantileBinner(2).Fit(nil).Transform([]float64{1})
	if ok {
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

// TestBinnersRejectNonFinite pins the binners to the package's reject policy
// for non-finite values: NaN has no order and so no bin (it used to land in
// the top bin), ±Inf has no finite bin either, and training data containing
// either cannot define a range or quantiles.
func TestBinnersRejectNonFinite(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name string
		fn   func() bool
	}{
		{name: "QuantileBinner transforming NaN", fn: func() bool {
			_, ok := preprocessing.NewQuantileBinner(2).Fit([]float64{1, 2, 3, 4}).Transform([]float64{1, nan})
			return ok
		}},
		{name: "QuantileBinner transforming +Inf", fn: func() bool {
			_, ok := preprocessing.NewQuantileBinner(2).Fit([]float64{1, 2, 3, 4}).Transform([]float64{inf})
			return ok
		}},
		{name: "FixedWidthBinner transforming NaN", fn: func() bool {
			_, ok := preprocessing.NewFixedWidthBinner(3).Fit([]float64{0, 3, 6}).Transform([]float64{nan})
			return ok
		}},
		{name: "FixedWidthBinner transforming -Inf", fn: func() bool {
			_, ok := preprocessing.NewFixedWidthBinner(3).Fit([]float64{0, 3, 6}).Transform([]float64{2, -inf})
			return ok
		}},
		{name: "FixedWidthBinner fitting NaN", fn: func() bool {
			_, ok := preprocessing.NewFixedWidthBinner(3).Fit([]float64{nan, 1, 2}).Transform([]float64{1})
			return ok
		}},
		{name: "FixedWidthBinner fitting +Inf", fn: func() bool {
			_, ok := preprocessing.NewFixedWidthBinner(3).Fit([]float64{1, inf}).Transform([]float64{1})
			return ok
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fn() {
				t.Errorf("ok = true, want false")
			}
		})
	}
}

// TestFixedWidthBinnerOverflowingRange checks a finite training range whose
// width overflows float64 still yields finite, evenly spaced cut points.
func TestFixedWidthBinnerOverflowingRange(t *testing.T) {
	maxF := math.MaxFloat64
	binner := preprocessing.NewFixedWidthBinner(2).Fit([]float64{-maxF, maxF})
	if !reflect.DeepEqual(binner.Edges(), []float64{0}) {
		t.Errorf("Edges() = %v, want [0]", binner.Edges())
	}
	four := preprocessing.NewFixedWidthBinner(4).Fit([]float64{-maxF, maxF})
	got, ok := four.Transform([]float64{-maxF, -1, 1, maxF})
	if !ok || !reflect.DeepEqual(got, []int{0, 1, 2, 3}) {
		t.Errorf("Transform() = %v, %v; want [0 1 2 3], true", got, ok)
	}
	// The half-unit edges must be scaled back up: with two bins the only edge
	// is 0 whatever the scaling, so four bins are needed to check it. The
	// upper edge rounds to within an ulp of maxF/2, hence the relative check.
	edges := four.Edges()
	wantEdges := []float64{-maxF / 2, 0, maxF / 2}
	if len(edges) != len(wantEdges) {
		t.Fatalf("Edges() = %v, want %v", edges, wantEdges)
	}
	for i, want := range wantEdges {
		if math.Abs(edges[i]-want) > 1e-15*math.Abs(want) {
			t.Errorf("Edges()[%d] = %v, want %v", i, edges[i], want)
		}
	}
}
