package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

func ExampleCosineSimilarity() {
	// Two feature vectors pointing in a similar direction score close to 1.
	a := []float64{1, 2, 3}
	b := []float64{2, 4, 6}

	sim, ok := stats.CosineSimilarity(a, b)
	fmt.Printf("%.1f %v", sim, ok)
	// Output: 1.0 true
}

func ExampleEuclideanDistance() {
	d, ok := stats.EuclideanDistance([]float64{0, 0}, []float64{3, 4})
	fmt.Printf("%.1f %v", d, ok)
	// Output: 5.0 true
}
