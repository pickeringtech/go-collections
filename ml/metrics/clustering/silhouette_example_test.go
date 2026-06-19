package clustering_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/metrics/clustering"
)

// Example_quickStart is the runnable twin of the package godoc overview. Keep
// the two in sync: `go test` compiles and output-checks this, which guarantees
// the documented clustering API actually exists and behaves as shown.
func Example_quickStart() {
	points := [][]float64{{0, 0}, {0.5, 0}, {10, 0}, {10.5, 0}}
	labels := []int{0, 0, 1, 1}

	score, _ := clustering.SilhouetteScore(points, labels)

	fmt.Printf("score=%.4f", score)
	// Output: score=0.9500
}

func ExampleSilhouetteScore_misassigned() {
	// Interleaving the two true groups scores negative — the clustering is poor.
	points := [][]float64{{0, 0}, {10, 0}, {0.5, 0}, {10.5, 0}}
	labels := []int{0, 1, 1, 0}

	score, ok := clustering.SilhouetteScore(points, labels)
	fmt.Printf("%.4f %v", score, ok)
	// Output: -0.4737 true
}
