package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

func ExampleSampleVariance() {
	variance, ok := stats.SampleVariance([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	fmt.Printf("%.4f %v", variance, ok)
	// Output: 4.5714 true
}

func ExampleSampleVariance_undefined() {
	// Sample variance needs at least two elements, so ok is false.
	variance, ok := stats.SampleVariance([]int{42})
	fmt.Printf("%.1f %v", variance, ok)
	// Output: 0.0 false
}

func ExamplePopulationVariance() {
	variance, ok := stats.PopulationVariance([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	fmt.Printf("%.1f %v", variance, ok)
	// Output: 4.0 true
}

func ExampleSampleStdDev() {
	stddev, ok := stats.SampleStdDev([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	fmt.Printf("%.4f %v", stddev, ok)
	// Output: 2.1381 true
}

func ExamplePopulationStdDev() {
	stddev, ok := stats.PopulationStdDev([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	fmt.Printf("%.1f %v", stddev, ok)
	// Output: 2.0 true
}
