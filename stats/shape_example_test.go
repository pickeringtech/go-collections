package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

func ExampleSkewness() {
	// A right-skewed sample: most values are small, one long tail to the right.
	skew, ok := stats.Skewness([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	fmt.Printf("%.5f %v", skew, ok)
	// Output: 0.65625 true
}

func ExampleKurtosis() {
	// Excess kurtosis: negative here means lighter tails than a normal curve.
	kurt, ok := stats.Kurtosis([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	fmt.Printf("%.5f %v", kurt, ok)
	// Output: -0.21875 true
}
