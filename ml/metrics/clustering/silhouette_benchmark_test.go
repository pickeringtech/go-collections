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

// TestSilhouetteScore_AllocatesOnlyTheResult guards the per-point cluster maps
// staying off the heap: they are created inside the per-point loop and never
// escape, so a whole score costs a single allocation (the per-sample slice). A
// refactor that lets them escape costs two heap allocations per point, which
// the merge-queue benchmark gate flagged as a 1200x allocation regression.
func TestSilhouetteScore_AllocatesOnlyTheResult(t *testing.T) {
	points, labels := benchClusters(60)
	got := testing.AllocsPerRun(20, func() {
		_, _ = clustering.SilhouetteScore(points, labels)
	})
	if got > 1 {
		t.Errorf("SilhouetteScore allocations per call = %v, want 1", got)
	}
}
