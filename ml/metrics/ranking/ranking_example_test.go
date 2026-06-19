package ranking_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/metrics/ranking"
)

// Example_quickStart is the runnable twin of the package godoc overview. Keep
// the two in sync: `go test` compiles and output-checks this, which guarantees
// the documented ranking API actually exists and behaves as shown.
func Example_quickStart() {
	trueRelevance := []float64{3, 2, 3, 0, 1, 2}
	scores := []float64{6, 5, 4, 3, 2, 1}
	ndcg, _ := ranking.NDCG(trueRelevance, scores, 0)

	ap, _ := ranking.AveragePrecision([]bool{true, false, true, false, true})

	fmt.Printf("ndcg=%.4f ap=%.4f", ndcg, ap)
	// Output: ndcg=0.9608 ap=0.7556
}

func ExampleNDCG_cutoff() {
	// Only the top 3 ranks count toward NDCG@3.
	trueRelevance := []float64{3, 2, 3, 0, 1, 2}
	scores := []float64{6, 5, 4, 3, 2, 1}
	ndcg, ok := ranking.NDCG(trueRelevance, scores, 3)
	fmt.Printf("%.4f %v", ndcg, ok)
	// Output: 0.9778 true
}

func ExampleMeanAveragePrecision() {
	queries := [][]bool{
		{true, false, true, false, true},
		{false, true, true, false},
	}
	mapScore, ok := ranking.MeanAveragePrecision(queries)
	fmt.Printf("%.4f %v", mapScore, ok)
	// Output: 0.6694 true
}
