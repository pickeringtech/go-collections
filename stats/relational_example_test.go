package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

func ExamplePopulationCovariance() {
	heights := []float64{1, 2, 3, 4, 5}
	weights := []float64{2, 4, 6, 8, 10}

	cov, ok := stats.PopulationCovariance(heights, weights)
	fmt.Printf("%.1f %v", cov, ok)
	// Output: 4.0 true
}

func ExampleCorrelation() {
	// Two perfectly, positively linearly related series.
	xs := []float64{1, 2, 3, 4, 5}
	ys := []float64{2, 4, 6, 8, 10}

	r, ok := stats.Correlation(xs, ys)
	fmt.Printf("%.1f %v", r, ok)
	// Output: 1.0 true
}

func ExampleCorrelation_constant() {
	// A constant series has no variance, so the correlation is undefined.
	_, ok := stats.Correlation([]float64{1, 2, 3}, []float64{4, 4, 4})
	fmt.Println(ok)
	// Output: false
}
