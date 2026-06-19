package preprocessing_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func TestOneHotEncoder(t *testing.T) {
	// Categories are learned in sorted order regardless of training row order.
	enc := preprocessing.NewOneHotEncoder[string]().Fit([]string{"b", "a", "a", "c"})
	if !reflect.DeepEqual(enc.Categories(), []string{"a", "b", "c"}) {
		t.Fatalf("Categories() = %v, want [a b c]", enc.Categories())
	}

	got, ok := enc.Transform([]string{"c", "a"})
	want := [][]float64{{0, 0, 1}, {1, 0, 0}}
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("Transform = (%v, %v), want (%v, true)", got, ok, want)
	}
}

func TestOneHotEncoderUnseenIsAllZero(t *testing.T) {
	enc := preprocessing.NewOneHotEncoder[string]().Fit([]string{"a", "b"})
	got, ok := enc.Transform([]string{"z"})
	want := [][]float64{{0, 0}}
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("Transform = (%v, %v), want (%v, true)", got, ok, want)
	}
}

func TestOneHotEncoderUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewOneHotEncoder[int]().Transform([]int{1}); ok {
		t.Fatalf("unfitted Transform reported ok")
	}
}

func TestLabelEncoder(t *testing.T) {
	enc := preprocessing.NewLabelEncoder[string]().Fit([]string{"b", "a", "c"})
	got, ok := enc.Transform([]string{"c", "a", "b"})
	want := []int{2, 0, 1}
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("Transform = (%v, %v), want (%v, true)", got, ok, want)
	}
}

func TestLabelEncoderUnseenIsMinusOne(t *testing.T) {
	enc := preprocessing.NewLabelEncoder[string]().Fit([]string{"a", "b"})
	got, _ := enc.Transform([]string{"z"})
	if !reflect.DeepEqual(got, []int{-1}) {
		t.Fatalf("Transform unseen = %v, want [-1]", got)
	}
}

func TestLabelEncoderInverseRoundTrip(t *testing.T) {
	enc := preprocessing.NewLabelEncoder[string]().Fit([]string{"b", "a", "c"})
	codes, _ := enc.Transform([]string{"a", "c", "b"})
	back, ok := enc.InverseTransform(codes)
	want := []string{"a", "c", "b"}
	if !ok || !reflect.DeepEqual(back, want) {
		t.Fatalf("InverseTransform = (%v, %v), want (%v, true)", back, ok, want)
	}
}

func TestLabelEncoderInverseRejectsOutOfRange(t *testing.T) {
	enc := preprocessing.NewLabelEncoder[string]().Fit([]string{"a", "b"})
	if _, ok := enc.InverseTransform([]int{0, 5}); ok {
		t.Fatalf("InverseTransform accepted out-of-range code")
	}
	if _, ok := enc.InverseTransform([]int{-1}); ok {
		t.Fatalf("InverseTransform accepted -1")
	}
}

func TestOrdinalEncoderExplicitOrder(t *testing.T) {
	// Caller-defined order, not alphabetical: low < medium < high.
	enc := preprocessing.NewOrdinalEncoder("low", "medium", "high")
	got, ok := enc.Transform([]string{"high", "low", "medium"})
	want := []int{2, 0, 1}
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("Transform = (%v, %v), want (%v, true)", got, ok, want)
	}
}

func TestOrdinalEncoderExplicitDeduplicates(t *testing.T) {
	enc := preprocessing.NewOrdinalEncoder("low", "high", "low", "high")
	if !reflect.DeepEqual(enc.Categories(), []string{"low", "high"}) {
		t.Fatalf("Categories() = %v, want [low high]", enc.Categories())
	}
}

func TestOrdinalEncoderLearnsSortedWhenNoOrder(t *testing.T) {
	enc := preprocessing.NewOrdinalEncoder[string]().Fit([]string{"b", "a", "c"})
	got, ok := enc.Transform([]string{"a", "b", "c"})
	if !ok || !reflect.DeepEqual(got, []int{0, 1, 2}) {
		t.Fatalf("Transform = (%v, %v), want ([0 1 2], true)", got, ok)
	}
}

func TestOrdinalEncoderUnseenIsMinusOne(t *testing.T) {
	enc := preprocessing.NewOrdinalEncoder("a", "b")
	got, _ := enc.Transform([]string{"z"})
	if !reflect.DeepEqual(got, []int{-1}) {
		t.Fatalf("Transform unseen = %v, want [-1]", got)
	}
}

func TestTargetEncoder(t *testing.T) {
	// a -> mean(1,3)=2; b -> mean(10)=10; global mean = 14/3.
	enc := preprocessing.NewTargetEncoder[string]().Fit(
		[]string{"a", "a", "b"},
		[]float64{1, 3, 10},
	)
	got, ok := enc.Transform([]string{"a", "b", "unseen"})
	want := []float64{2, 10, 14.0 / 3.0}
	if !ok || !floatSlicesClose(got, want) {
		t.Fatalf("Transform = (%v, %v), want (%v, true)", got, ok, want)
	}
	if !floatsClose(enc.GlobalMean(), 14.0/3.0) {
		t.Fatalf("GlobalMean() = %v, want %v", enc.GlobalMean(), 14.0/3.0)
	}
}

func TestTargetEncoderMismatchedLengthsLeavesUnfitted(t *testing.T) {
	enc := preprocessing.NewTargetEncoder[string]().Fit([]string{"a", "b"}, []float64{1})
	if _, ok := enc.Transform([]string{"a"}); ok {
		t.Fatalf("Fit with mismatched lengths reported fitted")
	}
}

func TestEncodersDoNotMutateInput(t *testing.T) {
	input := []string{"b", "a", "c"}
	original := []string{"b", "a", "c"}
	preprocessing.NewOneHotEncoder[string]().Fit(input).Transform(input)
	preprocessing.NewLabelEncoder[string]().Fit(input).Transform(input)
	if !reflect.DeepEqual(input, original) {
		t.Fatalf("input mutated to %v", input)
	}
}

// The fitted encoders satisfy the Transformer contract.
var (
	_ preprocessing.Transformer[string, [][]float64] = (*preprocessing.OneHotEncoder[string])(nil)
	_ preprocessing.Transformer[string, []int]       = (*preprocessing.LabelEncoder[string])(nil)
	_ preprocessing.Transformer[string, []int]       = (*preprocessing.OrdinalEncoder[string])(nil)
	_ preprocessing.Transformer[string, []float64]   = (*preprocessing.TargetEncoder[string])(nil)
)
