package preprocessing_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func ExampleStandardScaler() {
	// Fit learns mean 5 and population stddev 2 from the training data; the test
	// data is then standardized with those same parameters.
	scaler := preprocessing.NewStandardScaler().Fit([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	got, ok := scaler.Transform([]float64{5, 7})
	fmt.Println(got, ok)
	// Output: [0 1] true
}

func ExampleMinMaxScaler() {
	scaler := preprocessing.NewMinMaxScaler().Fit([]float64{2, 4, 6, 8, 10})
	got, ok := scaler.Transform([]float64{2, 4, 6, 8, 10})
	fmt.Println(got, ok)
	// Output: [0 0.25 0.5 0.75 1] true
}

func ExampleRobustScaler() {
	// Median 5, IQR 4, so r = (x − 5) / 4 — robust to the outlier-free tails.
	scaler := preprocessing.NewRobustScaler().Fit([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	got, ok := scaler.Transform([]float64{5, 9, 1})
	fmt.Println(got, ok)
	// Output: [0 1 -1] true
}

func ExampleStandardScaler_unfitted() {
	// Transforming before Fit reports ok == false rather than panicking.
	_, ok := preprocessing.NewStandardScaler().Transform([]float64{1, 2, 3})
	fmt.Println(ok)
	// Output: false
}
