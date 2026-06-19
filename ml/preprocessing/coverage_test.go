package preprocessing_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

// This file rounds out coverage of the FitTransform conveniences, accessors,
// seed wrappers and the remaining edge branches not exercised by the
// behaviour-focused tests in the other files.

func TestRobustScalerEmptyLeavesUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewRobustScaler().Fit(nil).Transform([]float64{1}); ok {
		t.Fatalf("empty Fit reported fitted")
	}
}

func TestImputerFitTransform(t *testing.T) {
	mean, ok := preprocessing.NewMeanImputer(nil).FitTransform([]float64{1, 2, 3})
	if !ok || !floatSlicesClose(mean, []float64{1, 2, 3}) {
		t.Fatalf("MeanImputer.FitTransform = (%v, %v)", mean, ok)
	}
	median, ok := preprocessing.NewMedianImputer(nil).FitTransform([]float64{1, 2, 3})
	if !ok || !floatSlicesClose(median, []float64{1, 2, 3}) {
		t.Fatalf("MedianImputer.FitTransform = (%v, %v)", median, ok)
	}

	isMissing := func(v string) bool { return v == "" }
	mode, ok := preprocessing.NewModeImputer(isMissing).FitTransform([]string{"a", "a", "b"})
	if !ok || !reflect.DeepEqual(mode, []string{"a", "a", "b"}) {
		t.Fatalf("ModeImputer.FitTransform = (%v, %v)", mode, ok)
	}
	constant, ok := preprocessing.NewConstantImputer("z", isMissing).FitTransform([]string{"a", ""})
	if !ok || !reflect.DeepEqual(constant, []string{"a", "z"}) {
		t.Fatalf("ConstantImputer.FitTransform = (%v, %v)", constant, ok)
	}
}

func TestImputerAccessorsAndNoOpFit(t *testing.T) {
	median := preprocessing.NewMedianImputer(nil).Fit([]float64{1, 2, 3, 4, 100})
	if !floatsClose(median.Fill(), 3) {
		t.Fatalf("MedianImputer.Fill() = %v, want 3", median.Fill())
	}
	// ConstantImputer.Fit is a no-op returning the receiver.
	constant := preprocessing.NewConstantImputer(9, nil).Fit([]int{1, 2})
	if constant.Fill() != 9 {
		t.Fatalf("ConstantImputer.Fill() = %v, want 9", constant.Fill())
	}
}

func TestMedianImputerAllMissingLeavesUnfitted(t *testing.T) {
	nan := math.NaN()
	if _, ok := preprocessing.NewMedianImputer(nil).Fit([]float64{nan, nan}).Transform([]float64{1}); ok {
		t.Fatalf("all-missing Fit reported fitted")
	}
}

func TestModeImputerAllMissingLeavesUnfitted(t *testing.T) {
	isMissing := func(v string) bool { return v == "" }
	if _, ok := preprocessing.NewModeImputer(isMissing).Fit([]string{"", ""}).Transform([]string{"x"}); ok {
		t.Fatalf("all-missing Fit reported fitted")
	}
}

func TestModeImputerUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewModeImputer[string](nil).Transform([]string{"x"}); ok {
		t.Fatalf("unfitted Transform reported ok")
	}
}

func TestEncoderFitTransform(t *testing.T) {
	oneHot, ok := preprocessing.NewOneHotEncoder[string]().FitTransform([]string{"b", "a"})
	want := [][]float64{{0, 1}, {1, 0}}
	if !ok || !reflect.DeepEqual(oneHot, want) {
		t.Fatalf("OneHot.FitTransform = (%v, %v), want (%v, true)", oneHot, ok, want)
	}
	label, ok := preprocessing.NewLabelEncoder[string]().FitTransform([]string{"b", "a"})
	if !ok || !reflect.DeepEqual(label, []int{1, 0}) {
		t.Fatalf("Label.FitTransform = (%v, %v), want ([1 0], true)", label, ok)
	}
	ordinal, ok := preprocessing.NewOrdinalEncoder[string]().FitTransform([]string{"b", "a"})
	if !ok || !reflect.DeepEqual(ordinal, []int{1, 0}) {
		t.Fatalf("Ordinal.FitTransform = (%v, %v), want ([1 0], true)", ordinal, ok)
	}
}

func TestLabelEncoderCategories(t *testing.T) {
	enc := preprocessing.NewLabelEncoder[string]().Fit([]string{"b", "a", "c"})
	if !reflect.DeepEqual(enc.Categories(), []string{"a", "b", "c"}) {
		t.Fatalf("Categories() = %v, want [a b c]", enc.Categories())
	}
}

func TestLabelEncoderInverseUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewLabelEncoder[string]().InverseTransform([]int{0}); ok {
		t.Fatalf("unfitted InverseTransform reported ok")
	}
}

func TestOrdinalEncoderExplicitFitIsNoOp(t *testing.T) {
	// Fit on an explicitly-ordered encoder keeps the caller's order.
	enc := preprocessing.NewOrdinalEncoder("low", "medium", "high").Fit([]string{"high", "low"})
	got, ok := enc.Transform([]string{"medium", "high"})
	if !ok || !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("Transform = (%v, %v), want ([1 2], true)", got, ok)
	}
}

func TestOrdinalEncoderUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewOrdinalEncoder[string]().Transform([]string{"x"}); ok {
		t.Fatalf("unfitted Transform reported ok")
	}
}

func TestOrdinalEncoderEmptyFitLeavesUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewOrdinalEncoder[string]().Fit(nil).Transform([]string{"x"}); ok {
		t.Fatalf("empty Fit reported fitted")
	}
}

func TestTargetEncoderEmptyLeavesUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewTargetEncoder[string]().Fit(nil, nil).Transform([]string{"a"}); ok {
		t.Fatalf("empty Fit reported fitted")
	}
}

func TestTargetEncoderNonFiniteTargetLeavesUnfitted(t *testing.T) {
	enc := preprocessing.NewTargetEncoder[string]().Fit([]string{"a", "b"}, []float64{1, math.NaN()})
	if _, ok := enc.Transform([]string{"a"}); ok {
		t.Fatalf("non-finite target reported fitted")
	}
}

func TestBinnerFitTransform(t *testing.T) {
	fixed, ok := preprocessing.NewFixedWidthBinner(2).FitTransform([]float64{0, 2})
	if !ok || !reflect.DeepEqual(fixed, []int{0, 1}) {
		t.Fatalf("FixedWidth.FitTransform = (%v, %v), want ([0 1], true)", fixed, ok)
	}
	quantile, ok := preprocessing.NewQuantileBinner(2).FitTransform([]float64{1, 2, 3, 4})
	if !ok || !reflect.DeepEqual(quantile, []int{0, 0, 1, 1}) {
		t.Fatalf("Quantile.FitTransform = (%v, %v), want ([0 0 1 1], true)", quantile, ok)
	}
}

func TestFixedWidthBinnerEmptyLeavesUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewFixedWidthBinner(3).Fit(nil).Transform([]float64{1}); ok {
		t.Fatalf("empty Fit reported fitted")
	}
}

func TestQuantileBinnerInvalidBinCount(t *testing.T) {
	if _, ok := preprocessing.NewQuantileBinner(0).Fit([]float64{1, 2}).Transform([]float64{1}); ok {
		t.Fatalf("nBins=0 reported fitted")
	}
}

func TestQuantileBinnerSingleBinEmptyLeavesUnfitted(t *testing.T) {
	// nBins == 1 takes the dedicated guard path; empty data must still reject.
	if _, ok := preprocessing.NewQuantileBinner(1).Fit(nil).Transform([]float64{1}); ok {
		t.Fatalf("empty single-bin Fit reported fitted")
	}
}

func TestSplitSeedWrappers(t *testing.T) {
	input := benchInts(12)
	folds, ok := preprocessing.KFoldSeed(input, 3, 4)
	if !ok || len(folds) != 3 {
		t.Fatalf("KFoldSeed = (%d, %v), want (3, true)", len(folds), ok)
	}
	labels := make([]int, 12)
	for i := range labels {
		labels[i] = i % 2
	}
	train, test, ok := preprocessing.StratifiedSplitSeed(input, labels, 0.5, 4)
	if !ok || len(train)+len(test) != 12 {
		t.Fatalf("StratifiedSplitSeed = (%d+%d, %v)", len(train), len(test), ok)
	}
}

func TestSplitsRejectNaNFraction(t *testing.T) {
	nan := math.NaN()
	if _, _, ok := preprocessing.TrainTestSplit(benchInts(5), nan, preprocessing.NewRand(1)); ok {
		t.Fatalf("TrainTestSplit accepted NaN testFrac")
	}
	labels := []int{0, 0, 1, 1, 1}
	if _, _, ok := preprocessing.StratifiedSplit(benchInts(5), labels, nan, preprocessing.NewRand(1)); ok {
		t.Fatalf("StratifiedSplit accepted NaN testFrac")
	}
}

// Float NaN cannot be a category (NaN != NaN). Encoders drop NaN at Fit time, so
// a NaN at Transform falls through to the unseen path rather than growing the
// category set or creating an unmatchable column.
func TestEncodersDropNaNCategories(t *testing.T) {
	nan := math.NaN()

	// OneHot/Label/learned-Ordinal go through sortedUnique: NaN is excluded from
	// the learned categories.
	oneHot := preprocessing.NewOneHotEncoder[float64]().Fit([]float64{1, nan, 2, nan, 1})
	if !floatSlicesClose(oneHot.Categories(), []float64{1, 2}) {
		t.Fatalf("OneHot.Categories() = %v, want [1 2]", oneHot.Categories())
	}
	rows, _ := oneHot.Transform([]float64{nan, 2})
	if !reflect.DeepEqual(rows, [][]float64{{0, 0}, {0, 1}}) {
		t.Fatalf("OneHot NaN transform = %v, want [[0 0] [0 1]]", rows)
	}

	// Explicit ordinal order drops NaN too.
	ordinal := preprocessing.NewOrdinalEncoder(1.0, nan, 2.0)
	if !floatSlicesClose(ordinal.Categories(), []float64{1, 2}) {
		t.Fatalf("Ordinal.Categories() = %v, want [1 2]", ordinal.Categories())
	}
	codes, _ := ordinal.Transform([]float64{nan, 1})
	if !reflect.DeepEqual(codes, []int{-1, 0}) {
		t.Fatalf("Ordinal NaN transform = %v, want [-1 0]", codes)
	}

	// TargetEncoder skips NaN categories; a NaN input maps to the global mean.
	target := preprocessing.NewTargetEncoder[float64]().Fit(
		[]float64{1, 1, nan},
		[]float64{2, 4, 99},
	)
	got, ok := target.Transform([]float64{1, nan})
	if !ok || !floatsClose(got[0], 3) || !floatsClose(got[1], target.GlobalMean()) {
		t.Fatalf("Target NaN transform = (%v, %v), want [3 globalMean]", got, ok)
	}
}

func TestSplitsNilRandUsesDefault(t *testing.T) {
	// A nil generator falls back to the deterministic default rather than
	// panicking, and stays a permutation.
	got := preprocessing.Shuffle(benchInts(8), nil)
	if !reflect.DeepEqual(sortedCopy(got), benchInts(8)) {
		t.Fatalf("Shuffle(_, nil) is not a permutation: %v", got)
	}
}
