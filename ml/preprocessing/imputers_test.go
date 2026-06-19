package preprocessing_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func TestMeanImputer(t *testing.T) {
	nan := math.NaN()
	// Train mean over non-missing {1,2,3} is 2.
	imp := preprocessing.NewMeanImputer(nil).Fit([]float64{1, 2, 3, nan})
	got, ok := imp.Transform([]float64{nan, 5, nan})
	if !ok || !floatSlicesClose(got, []float64{2, 5, 2}) {
		t.Fatalf("Transform = (%v, %v), want ([2 5 2], true)", got, ok)
	}
	if !floatsClose(imp.Fill(), 2) {
		t.Fatalf("Fill() = %v, want 2", imp.Fill())
	}
}

func TestMeanImputerCustomMissing(t *testing.T) {
	// Treat -1 as the sentinel for missing.
	isMissing := func(v float64) bool { return v == -1 }
	got, ok := preprocessing.NewMeanImputer(isMissing).Fit([]float64{2, 4, -1}).Transform([]float64{-1, 10})
	if !ok || !floatSlicesClose(got, []float64{3, 10}) {
		t.Fatalf("Transform = (%v, %v), want ([3 10], true)", got, ok)
	}
}

func TestMeanImputerAllMissingLeavesUnfitted(t *testing.T) {
	nan := math.NaN()
	if _, ok := preprocessing.NewMeanImputer(nil).Fit([]float64{nan, nan}).Transform([]float64{1}); ok {
		t.Fatalf("Fit on all-missing reported fitted")
	}
}

func TestMeanImputerUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewMeanImputer(nil).Transform([]float64{1}); ok {
		t.Fatalf("unfitted Transform reported ok")
	}
}

func TestMedianImputer(t *testing.T) {
	nan := math.NaN()
	// Median over non-missing {1,2,3,4,100} is 3.
	got, ok := preprocessing.NewMedianImputer(nil).Fit([]float64{1, 2, 3, 4, 100, nan}).Transform([]float64{nan, 7})
	if !ok || !floatSlicesClose(got, []float64{3, 7}) {
		t.Fatalf("Transform = (%v, %v), want ([3 7], true)", got, ok)
	}
}

func TestModeImputer(t *testing.T) {
	// "" marks missing; the modal non-missing category is "a".
	isMissing := func(v string) bool { return v == "" }
	imp := preprocessing.NewModeImputer(isMissing).Fit([]string{"a", "a", "b", ""})
	got, ok := imp.Transform([]string{"", "b", ""})
	want := []string{"a", "b", "a"}
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("Transform = (%v, %v), want (%v, true)", got, ok, want)
	}
	if imp.Fill() != "a" {
		t.Fatalf("Fill() = %q, want \"a\"", imp.Fill())
	}
}

func TestModeImputerTieFirstSeenWins(t *testing.T) {
	// "b" and "a" tie; "b" appears first, so it is the fill.
	isMissing := func(v string) bool { return v == "" }
	imp := preprocessing.NewModeImputer(isMissing).Fit([]string{"b", "a", "b", "a"})
	if imp.Fill() != "b" {
		t.Fatalf("Fill() = %q, want \"b\"", imp.Fill())
	}
}

func TestConstantImputer(t *testing.T) {
	// Needs no Fit; ready to Transform immediately.
	isMissing := func(v string) bool { return v == "?" }
	got, ok := preprocessing.NewConstantImputer("UNKNOWN", isMissing).Transform([]string{"x", "?", "y"})
	want := []string{"x", "UNKNOWN", "y"}
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("Transform = (%v, %v), want (%v, true)", got, ok, want)
	}
}

func TestConstantImputerNilMissingIsIdentity(t *testing.T) {
	got, ok := preprocessing.NewConstantImputer(0, nil).Transform([]int{1, 2, 3})
	if !ok || !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Transform = (%v, %v), want ([1 2 3], true)", got, ok)
	}
}

func TestImputersZeroValueOutput(t *testing.T) {
	got, ok := preprocessing.NewMeanImputer(nil).Fit([]float64{1, 2}).Transform(nil)
	if !ok || got == nil || len(got) != 0 {
		t.Fatalf("Transform(nil) = (%v, %v), want (non-nil empty, true)", got, ok)
	}
}

func TestImputersDoNotMutateInput(t *testing.T) {
	nan := math.NaN()
	input := []float64{nan, 2, 3}
	imp := preprocessing.NewMeanImputer(nil).Fit([]float64{2, 4})
	imp.Transform(input)
	if !math.IsNaN(input[0]) || input[1] != 2 || input[2] != 3 {
		t.Fatalf("input mutated to %v", input)
	}
}

func TestIsNaN(t *testing.T) {
	if !preprocessing.IsNaN(math.NaN()) || preprocessing.IsNaN(1) {
		t.Fatalf("IsNaN behaved unexpectedly")
	}
}

// The fitted imputers satisfy the Transformer contract.
var (
	_ preprocessing.Transformer[float64, []float64] = (*preprocessing.MeanImputer)(nil)
	_ preprocessing.Transformer[float64, []float64] = (*preprocessing.MedianImputer)(nil)
	_ preprocessing.Transformer[string, []string]   = (*preprocessing.ModeImputer[string])(nil)
	_ preprocessing.Transformer[int, []int]         = (*preprocessing.ConstantImputer[int])(nil)
)
