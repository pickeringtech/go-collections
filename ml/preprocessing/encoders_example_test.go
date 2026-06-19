package preprocessing_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func ExampleOneHotEncoder() {
	// Columns follow sorted category order: a, b, c.
	enc := preprocessing.NewOneHotEncoder[string]().Fit([]string{"b", "a", "a", "c"})
	got, ok := enc.Transform([]string{"a", "c"})
	fmt.Println(got, ok)
	// Output: [[1 0 0] [0 0 1]] true
}

func ExampleLabelEncoder() {
	enc := preprocessing.NewLabelEncoder[string]().Fit([]string{"b", "a", "c"})
	got, ok := enc.Transform([]string{"a", "b", "c"})
	fmt.Println(got, ok)
	// Output: [0 1 2] true
}

func ExampleOrdinalEncoder() {
	// A caller-defined order, not alphabetical.
	enc := preprocessing.NewOrdinalEncoder("low", "medium", "high")
	got, ok := enc.Transform([]string{"high", "low", "medium"})
	fmt.Println(got, ok)
	// Output: [2 0 1] true
}

func ExampleTargetEncoder() {
	// a -> mean(0, 2) = 1; b -> mean(10, 10) = 10.
	enc := preprocessing.NewTargetEncoder[string]().Fit(
		[]string{"a", "a", "b", "b"},
		[]float64{0, 2, 10, 10},
	)
	got, ok := enc.Transform([]string{"a", "b"})
	fmt.Println(got, ok)
	// Output: [1 10] true
}
