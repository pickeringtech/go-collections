package clustering_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/metrics/clustering"
)

func benchClusters(n int) (points [][]float64, labels []int) {
	points = make([][]float64, n)
	labels = make([]int, n)
	for i := range points {
		cluster := i % 3
		// Three Gaussian-ish blobs spread along the x axis.
		points[i] = []float64{float64(cluster*10) + float64(i%5)*0.1, float64(i % 7)}
		labels[i] = cluster
	}
	return points, labels
}

func BenchmarkSilhouetteScore(b *testing.B) {
	points, labels := benchClusters(300)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = clustering.SilhouetteScore(points, labels)
	}
}
