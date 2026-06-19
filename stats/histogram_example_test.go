package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

func ExampleHistogram() {
	// Span 4 over 2 bins gives equal-width buckets [1,3) and [3,5]; the maximum
	// (5) is folded into the final, inclusive bin.
	bins, ok := stats.Histogram([]float64{1, 2, 3, 4, 5}, 2)
	fmt.Println(ok)
	for _, b := range bins {
		fmt.Printf("[%.0f,%.0f) count=%d\n", b.Min, b.Max, b.Count)
	}
	// Output:
	// true
	// [1,3) count=2
	// [3,5) count=3
}
