// Package distance is documented in doc.go.
package distance

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/ml/similarity"
	"github.com/pickeringtech/go-collections/stats"
)

// Euclidean computes the L2 (straight-line) distance between two numeric
// vectors: √(Σ (aᵢ - bᵢ)²). It is the most common notion of physical distance
// in n-dimensional space.
//
// It delegates to stats.EuclideanDistance, the module's canonical vector
// geometry, whose scaled sum-of-squares is both overflow-safe and
// Kahan-precise. It returns ok == false for empty input or vectors of differing
// lengths — the distance is undefined in those cases. Non-finite inputs
// (NaN/Inf) propagate to a non-finite result with ok == true, following the
// package's NaN/Inf policy.
func Euclidean[T constraints.Numeric](a, b []T) (float64, bool) {
	return stats.EuclideanDistance(a, b)
}

// Manhattan computes the L1 (taxicab) distance between two numeric vectors:
// Σ |aᵢ - bᵢ|. It measures the total variation across all dimensions, treating
// each independently.
//
// It returns ok == false for empty input or vectors of differing lengths.
// Non-finite inputs (NaN/Inf) propagate to a non-finite result with ok == true.
func Manhattan[T constraints.Numeric](a, b []T) (float64, bool) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	var dist float64
	for i := range a {
		diff := float64(a[i]) - float64(b[i])
		dist += math.Abs(diff)
	}
	return dist, true
}

// Minkowski computes the Lp distance between two numeric vectors:
// (Σ |aᵢ - bᵢ|ᵖ)^(1/p). It generalises both Manhattan (p=1) and Euclidean
// (p=2). For p → ∞ the result approaches the Chebyshev distance.
//
// It returns ok == false for:
//   - empty input or vectors of differing lengths;
//   - p < 1 — values below 1 do not satisfy the triangle inequality and
//     therefore do not define a valid distance metric.
//
// Non-finite inputs (NaN/Inf) propagate to a non-finite result with ok == true.
func Minkowski[T constraints.Numeric](a, b []T, p float64) (float64, bool) {
	if p < 1 {
		return 0, false
	}
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	var sum float64
	for i := range a {
		diff := math.Abs(float64(a[i]) - float64(b[i]))
		sum += math.Pow(diff, p)
	}
	return math.Pow(sum, 1/p), true
}

// CosineDistance computes the cosine distance between two numeric vectors:
// 1 - CosineSimilarity(a, b). The result is in [0, 2]: 0 for identical
// direction, 1 for orthogonal vectors, and 2 for anti-parallel vectors.
//
// It delegates to ml/similarity.CosineSimilarity and inherits its ok == false
// conditions: empty input, mismatched lengths, or either vector having zero
// magnitude. Non-finite inputs propagate per the package's NaN/Inf policy.
func CosineDistance[T constraints.Numeric](a, b []T) (float64, bool) {
	cos, ok := similarity.CosineSimilarity(a, b)
	if !ok {
		return 0, false
	}
	return 1 - cos, true
}
