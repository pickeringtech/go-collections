package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

func ExamplePercentileOfScore() {
	// A score of 80 sits at or above 60% of the cohort's marks.
	marks := []int{55, 62, 71, 80, 95}

	rank, ok := stats.PercentileOfScore(marks, 80)
	fmt.Printf("%.0fth percentile %v", rank, ok)
	// Output: 80th percentile true
}
