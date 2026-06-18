package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

// Example_quickStart is the runnable twin of the package godoc overview. Keep
// the two in sync: `go test` compiles and output-checks this, which is what
// guarantees the documented stats API actually exists and behaves as shown.
func Example_quickStart() {
	data := []float64{10, 20, 30}

	// Exact-in-T total and float64 average — each with an ok flag.
	total, _ := stats.Sum(data)
	mean, _ := stats.Mean(data)

	prices := []float64{10, 20, 30}
	weights := []float64{1, 2, 3}

	wm, _ := stats.WeightedMean(prices, weights)
	gm, _ := stats.GeometricMean([]float64{1, 10, 100})
	hm, _ := stats.HarmonicMean([]float64{1, 2, 4})

	fmt.Printf("sum=%.0f mean=%.0f weighted=%.4f geometric=%.4f harmonic=%.4f",
		total, mean, wm, gm, hm)
	// Output: sum=60 mean=20 weighted=23.3333 geometric=10.0000 harmonic=1.7143
}

func ExampleSum() {
	total, ok := stats.Sum([]int{1, 2, 3, 4, 5})
	fmt.Printf("%d %v", total, ok)
	// Output: 15 true
}

func ExampleSum_empty() {
	// The empty sum is undefined under the (result, ok) idiom.
	total, ok := stats.Sum([]int{})
	fmt.Printf("%d %v", total, ok)
	// Output: 0 false
}

func ExampleMean() {
	mean, ok := stats.Mean([]float64{1, 2, 3, 4, 5})
	fmt.Printf("%.1f %v", mean, ok)
	// Output: 3.0 true
}

func ExampleWeightedMean() {
	// Grades worth different numbers of credits.
	grades := []float64{90, 80, 70}
	credits := []float64{3, 4, 2}

	mean, ok := stats.WeightedMean(grades, credits)
	fmt.Printf("%.2f %v", mean, ok)
	// Output: 81.11 true
}

func ExampleWeightedMean_invalid() {
	// Mismatched lengths cannot be summarised, so ok is false.
	_, ok := stats.WeightedMean([]float64{1, 2, 3}, []float64{1, 2})
	fmt.Println(ok)
	// Output: false
}

func ExampleGeometricMean() {
	// Average growth factor across three periods (10%, 50%, 30% growth).
	factors := []float64{1.1, 1.5, 1.3}

	mean, _ := stats.GeometricMean(factors)
	fmt.Printf("%.4f", mean)
	// Output: 1.2897
}

func ExampleHarmonicMean() {
	// Average speed over equal distances driven at 30 and 60 mph.
	speeds := []float64{30, 60}

	mean, _ := stats.HarmonicMean(speeds)
	fmt.Printf("%.2f", mean)
	// Output: 40.00
}
