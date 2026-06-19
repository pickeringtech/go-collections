package preprocessing_test

import (
	"fmt"
	"math"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

// Example_quickStart mirrors the Quick Start in the package doc, so the
// documented API is compiled and verified.
func Example_quickStart() {
	// Scale features to zero mean / unit variance using train parameters.
	scaler := preprocessing.NewStandardScaler().Fit([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	z, _ := scaler.Transform([]float64{5}) // 5 is the train mean

	// One-hot encode categories in a stable, sorted column order.
	enc := preprocessing.NewOneHotEncoder[string]().Fit([]string{"b", "a", "a", "c"})
	rows, _ := enc.Transform([]string{"a"}) // columns a,b,c

	// Fill missing values with the train mean.
	imp := preprocessing.NewMeanImputer(nil).Fit([]float64{1, 2, 3})
	filled, _ := imp.Transform([]float64{math.NaN(), 2})

	// Reproducible train/test split.
	rng := preprocessing.NewRand(42)
	train, test, _ := preprocessing.TrainTestSplit([]int{0, 1, 2, 3, 4}, 0.4, rng)

	fmt.Println(z)
	fmt.Println(rows)
	fmt.Println(filled)
	fmt.Println(len(train), len(test))
	// Output:
	// [0]
	// [[1 0 0]]
	// [2 2]
	// 3 2
}
