package preprocessing_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func ExampleFixedWidthBinner() {
	// Range [0, 8] into 4 equal-width bins (cut points 2, 4, 6).
	binner := preprocessing.NewFixedWidthBinner(4).Fit([]float64{0, 8})
	got, ok := binner.Transform([]float64{1, 3, 5, 7})
	fmt.Println(got, ok)
	// Output: [0 1 2 3] true
}

func ExampleQuantileBinner() {
	// Two equal-population bins split at the median (2.5).
	binner := preprocessing.NewQuantileBinner(2).Fit([]float64{1, 2, 3, 4})
	got, ok := binner.Transform([]float64{1, 2, 3, 4})
	fmt.Println(got, ok)
	// Output: [0 0 1 1] true
}
