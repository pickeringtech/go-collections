package preprocessing_test

import (
	"fmt"
	"math"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func ExampleMeanImputer() {
	// NaN marks missing by default; the learned fill is the train mean (2).
	imp := preprocessing.NewMeanImputer(nil).Fit([]float64{1, 2, 3})
	got, ok := imp.Transform([]float64{math.NaN(), 5, math.NaN()})
	fmt.Println(got, ok)
	// Output: [2 5 2] true
}

func ExampleModeImputer() {
	// An empty string marks missing; the modal non-missing category is "a".
	isMissing := func(v string) bool { return v == "" }
	imp := preprocessing.NewModeImputer(isMissing).Fit([]string{"a", "a", "b"})
	got, ok := imp.Transform([]string{"", "b"})
	fmt.Println(got, ok)
	// Output: [a b] true
}

func ExampleConstantImputer() {
	// A ConstantImputer needs no training data.
	isMissing := func(v string) bool { return v == "?" }
	got, ok := preprocessing.NewConstantImputer("UNKNOWN", isMissing).Transform([]string{"x", "?", "y"})
	fmt.Println(got, ok)
	// Output: [x UNKNOWN y] true
}
