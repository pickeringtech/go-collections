package clustering

import (
	"math"

	"github.com/pickeringtech/go-collections/stats"
)

// DistanceFunc computes the distance between two equal-length coordinate
// vectors. It backs the pluggable metric in SilhouetteScoreWith and
// SilhouetteSamplesWith.
type DistanceFunc func(a, b []float64) float64

// EuclideanDistance is the default metric — the straight-line (L2) distance
// between two points, sqrt(Σ(aᵢ − bᵢ)²). It delegates to the module's canonical
// stats.EuclideanDistance. validate guarantees equal-length, finite coordinate
// rows before this runs, so the ok flag is discarded: the only remaining input
// that makes it false is a zero-dimensional row, whose distance is 0 anyway —
// the same value the previous local implementation returned.
func EuclideanDistance(a, b []float64) float64 {
	d, _ := stats.EuclideanDistance(a, b)
	return d
}

// validate checks the shared preconditions for the silhouette functions and
// returns the number of distinct clusters. The silhouette coefficient requires
// between 2 and n−1 clusters: with one cluster there is nothing to separate,
// and with every point in its own cluster there is no within-cluster structure.
func validate(points [][]float64, labels []int) (clusters int, ok bool) {
	n := len(points)
	if n < 2 || len(labels) != n {
		return 0, false
	}

	dim := len(points[0])
	sizes := make(map[int]int)
	for i := range points {
		if len(points[i]) != dim {
			return 0, false
		}
		for _, v := range points[i] {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return 0, false
			}
		}
		sizes[labels[i]]++
	}

	k := len(sizes)
	if k < 2 || k >= n {
		return 0, false
	}
	return k, true
}

// clusterDistances sums, for point i, the distances to every other point
// grouped by the other point's cluster label (i itself excluded). ok is false
// when dist returns a non-finite or negative distance — or when summing many
// large distances overflows to ±Inf — since either would push a coefficient
// outside the documented [−1, 1] range.
func clusterDistances(points [][]float64, labels []int, i int, dist DistanceFunc) (sums map[int]float64, counts map[int]int, ok bool) {
	sums = make(map[int]float64)
	counts = make(map[int]int)
	for j := range points {
		if i == j {
			continue
		}
		d := dist(points[i], points[j])
		if math.IsNaN(d) || math.IsInf(d, 0) || d < 0 {
			return nil, nil, false
		}
		sums[labels[j]] += d
		if math.IsInf(sums[labels[j]], 0) {
			// The running total overflowed; a mean built from it would
			// yield Inf/Inf == NaN downstream, voiding the range guarantee.
			return nil, nil, false
		}
		counts[labels[j]]++
	}
	return sums, counts, true
}

// samples computes the per-sample silhouette coefficients under dist, assuming
// validate has already passed. ok is false when dist misbehaves (see
// clusterDistances).
func samples(points [][]float64, labels []int, dist DistanceFunc) ([]float64, bool) {
	s := make([]float64, len(points))
	for i := range points {
		// Mean distance from point i to each cluster (excluding i itself).
		sums, counts, ok := clusterDistances(points, labels, i, dist)
		if !ok {
			return nil, false
		}

		own := labels[i]
		ownCount := counts[own] // size of i's cluster minus i
		if ownCount == 0 {
			// A lone point in its cluster has no defined cohesion; by
			// convention its silhouette is 0.
			s[i] = 0
			continue
		}

		a := sums[own] / float64(ownCount) // mean intra-cluster distance
		b := math.Inf(1)                   // nearest other cluster's mean distance
		for id, cnt := range counts {
			if id == own {
				continue
			}
			mean := sums[id] / float64(cnt)
			if mean < b {
				b = mean
			}
		}

		denom := math.Max(a, b)
		if denom == 0 {
			s[i] = 0
		} else {
			s[i] = (b - a) / denom
		}
	}
	return s, true
}

// SilhouetteSamplesWith returns the silhouette coefficient of every sample
// under the supplied distance metric, together with an ok flag. Each
// coefficient lies in [−1, 1]: near 1 the point sits comfortably in its
// cluster, near 0 it lies between two clusters, and negative it is probably
// mis-assigned.
//
// ok is false (and the result is nil) when the inputs cannot be summarised:
// fewer than two points, len(labels) != len(points), ragged coordinate rows,
// any non-finite coordinate, or a cluster count outside [2, n−1]. It is also
// false when dist is nil or returns a non-finite or negative distance, since
// such a metric would void the [−1, 1] guarantee.
func SilhouetteSamplesWith(points [][]float64, labels []int, dist DistanceFunc) ([]float64, bool) {
	_, ok := validate(points, labels)
	if !ok || dist == nil {
		return nil, false
	}
	return samples(points, labels, dist)
}

// SilhouetteSamples is SilhouetteSamplesWith using EuclideanDistance.
func SilhouetteSamples(points [][]float64, labels []int) ([]float64, bool) {
	return SilhouetteSamplesWith(points, labels, EuclideanDistance)
}

// SilhouetteScoreWith returns the mean silhouette coefficient over all samples
// under the supplied distance metric, together with an ok flag — a single
// number summarising how well-separated the clustering is, in [−1, 1]. It
// rejects the same inputs as SilhouetteSamplesWith.
func SilhouetteScoreWith(points [][]float64, labels []int, dist DistanceFunc) (float64, bool) {
	s, ok := SilhouetteSamplesWith(points, labels, dist)
	if !ok {
		return 0, false
	}
	return stats.Mean(s)
}

// SilhouetteScore is SilhouetteScoreWith using EuclideanDistance.
func SilhouetteScore(points [][]float64, labels []int) (float64, bool) {
	return SilhouetteScoreWith(points, labels, EuclideanDistance)
}
