package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

func ExampleLinearRegression() {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 5, 4, 6}

	fit, ok := stats.LinearRegression(x, y)
	fmt.Printf("y = %.1fx + %.1f (R²=%.3f) %v", fit.Slope, fit.Intercept, fit.R2, ok)
	// Output: y = 0.8x + 1.8 (R²=0.727) true
}

func ExampleLineFit_Predict() {
	fit, _ := stats.LinearRegression([]float64{1, 2, 3, 4, 5}, []float64{2, 4, 5, 4, 6})

	// Predict the fitted value at x=6, and the residual of an observed y=7.
	predicted := fit.Predict(6)
	residual := 7 - predicted
	fmt.Printf("predicted=%.1f residual=%.1f", predicted, residual)
	// Output: predicted=6.6 residual=0.4
}

func ExampleLinearRegression_undefined() {
	// A constant x has no spread to fit a slope against, so ok is false.
	fit, ok := stats.LinearRegression([]float64{3, 3, 3}, []float64{1, 2, 3})
	fmt.Printf("%+v %v", fit, ok)
	// Output: {Slope:0 Intercept:0 R2:0} false
}
