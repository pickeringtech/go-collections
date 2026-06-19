package ranking_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/metrics/ranking"
)

func benchRanking(n int) (rel, scores []float64) {
	rel = make([]float64, n)
	scores = make([]float64, n)
	for i := range rel {
		rel[i] = float64(i % 4)    // relevance grades 0..3
		scores[i] = float64(n - i) // already descending
	}
	return rel, scores
}

func BenchmarkNDCG(b *testing.B) {
	rel, scores := benchRanking(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ranking.NDCG(rel, scores, 0)
	}
}

func BenchmarkAveragePrecision(b *testing.B) {
	ranked := make([]bool, 1000)
	for i := range ranked {
		ranked[i] = i%3 == 0
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ranking.AveragePrecision(ranked)
	}
}
