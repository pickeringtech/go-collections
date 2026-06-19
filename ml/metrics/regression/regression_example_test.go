package regression_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/metrics/regression"
)

// Example_quickStart is the runnable twin of the package godoc overview. Keep
// the two in sync: `go test` compiles and output-checks this, which is what
// guarantees the documented regression API actually exists and behaves as shown.
func Example_quickStart() {
	yTrue := []float64{3, -0.5, 2, 7}
	yPred := []float64{2.5, 0, 2, 8}

	mse, _ := regression.MeanSquaredError(yTrue, yPred)
	rmse, _ := regression.RootMeanSquaredError(yTrue, yPred)
	mae, _ := regression.MeanAbsoluteError(yTrue, yPred)
	r2, _ := regression.RSquared(yTrue, yPred)

	fmt.Printf("mse=%.3f rmse=%.4f mae=%.1f r2=%.4f", mse, rmse, mae, r2)
	// Output: mse=0.375 rmse=0.6124 mae=0.5 r2=0.9486
}

func ExampleMeanSquaredError() {
	mse, ok := regression.MeanSquaredError([]float64{1, 2, 3}, []float64{1, 2, 5})
	fmt.Printf("%.4f %v", mse, ok)
	// Output: 1.3333 true
}

func ExampleMeanAbsolutePercentageError() {
	// Returned as a fraction: 0.0722 means the predictions are off by ~7.2%.
	mape, ok := regression.MeanAbsolutePercentageError(
		[]float64{100, 200, 300},
		[]float64{110, 190, 320},
	)
	fmt.Printf("%.4f %v", mape, ok)
	// Output: 0.0722 true
}

func ExampleRSquared() {
	// A constant target has zero variance, so R² is undefined.
	_, ok := regression.RSquared([]float64{5, 5, 5}, []float64{5, 4, 6})
	fmt.Printf("%v", ok)
	// Output: false
}
