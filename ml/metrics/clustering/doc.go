// Package clustering scores the quality of a clustering without reference to
// ground-truth labels — currently the silhouette coefficient, as pure
// functions over a slice of coordinate vectors and their cluster assignments.
//
// It is part of the ml/metrics family (see the ml umbrella package). Where the
// classification package needs true labels, silhouette is an internal measure:
// it judges a clustering purely by how compact and well-separated the clusters
// are.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/ml/metrics/clustering"
//
//	points := [][]float64{{0, 0}, {0.5, 0}, {10, 0}, {10.5, 0}}
//	labels := []int{0, 0, 1, 1}
//
//	score, _ := clustering.SilhouetteScore(points, labels)   // 0.9500 — tight, separated
//	samples, _ := clustering.SilhouetteSamples(points, labels)
//
//	_ = score
//	_ = samples
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
//
// # Silhouette
//
// For each sample the silhouette coefficient is (b − a) / max(a, b), where a is
// the mean distance to the other points in its own cluster and b is the mean
// distance to the points of the nearest other cluster. It lies in [−1, 1]: near
// 1 the point is well inside its cluster, near 0 it sits on a boundary, and
// negative it is probably assigned to the wrong cluster. SilhouetteScore is the
// mean over all samples; SilhouetteSamples returns the per-point values.
//
// SilhouetteScoreWith and SilhouetteSamplesWith accept a DistanceFunc, so any
// metric can replace the default EuclideanDistance.
//
// # Conventions
//
// The functions return (result, ok) in the library's idiom rather than
// panicking or returning an error. ok is false — and the result the zero value
// — when the inputs cannot be summarised: fewer than two points, a label slice
// of the wrong length, ragged coordinate rows, any non-finite coordinate, or a
// cluster count outside [2, n−1] (silhouette is undefined for a single cluster
// or for one-point-per-cluster). For the SilhouetteScoreWith and
// SilhouetteSamplesWith variants ok is also false when the supplied
// DistanceFunc is nil or returns a non-finite or negative distance, which would
// otherwise void the [−1, 1] guarantee. A lone point in its cluster is given a
// silhouette of 0. Inputs are never mutated, and the mean over samples routes
// through stats.Mean.
package clustering
