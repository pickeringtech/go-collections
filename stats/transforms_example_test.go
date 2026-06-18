package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

func ExampleNormalize() {
	got, ok := stats.Normalize([]float64{2, 4, 6, 8, 10})
	fmt.Println(got, ok)
	// Output: [0 0.25 0.5 0.75 1] true
}

func ExampleStandardize() {
	// Mean 5, population standard deviation 2, so z = (x − 5) / 2.
	got, ok := stats.Standardize([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	fmt.Println(got, ok)
	// Output: [-1.5 -0.5 -0.5 -0.5 0 0 1 2] true
}

func ExampleMovingAverage() {
	// Only full windows are emitted, so the result is shorter than the input.
	got, ok := stats.MovingAverage([]float64{1, 2, 3, 4, 5}, 2)
	fmt.Println(got, ok)
	// Output: [1.5 2.5 3.5 4.5] true
}

func ExampleMovingAverage_windowTooLarge() {
	// A window larger than the input cannot form a single full window.
	_, ok := stats.MovingAverage([]float64{1, 2, 3}, 4)
	fmt.Println(ok)
	// Output: false
}
