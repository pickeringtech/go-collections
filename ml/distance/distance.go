// Package distance is documented in doc.go.
package distance

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/ml/similarity"
)

// Euclidean computes the L2 (straight-line) distance between two numeric
// vectors: √(Σ (aᵢ - bᵢ)²). It is the most common notion of physical distance
// in n-dimensional space.
//
// It returns ok == false for empty input or vectors of differing lengths —
// the distance is undefined in those cases. Non-finite inputs (NaN/Inf)
// propagate to a non-finite result with ok == true, following the package's
// NaN/Inf policy.
func Euclidean[T constraints.Numeric](a, b []T) (float64, bool) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	var dist float64
	for i := range a {
		diff := float64(a[i]) - float64(b[i])
		dist = math.Hypot(dist, diff)
	}
	return dist, true
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
// (p=2). As p → ∞ the metric approaches the Chebyshev (max-coordinate)
// distance, but p itself must be finite — pass increasing finite values to
// approximate that limit.
//
// It returns ok == false for:
//   - empty input or vectors of differing lengths;
//   - p that is not a finite value ≥ 1 (p < 1, NaN, or ±Inf) — values below 1
//     do not satisfy the triangle inequality, and a non-finite p is not a
//     valid metric parameter.
//
// Large finite p is handled stably: the largest absolute difference is factored
// out before exponentiation, so the result neither overflows to +Inf nor
// collapses to a single dimension as it would under a naive Σ diffᵖ.
//
// Non-finite inputs (NaN/Inf) propagate to a non-finite result with ok == true.
func Minkowski[T constraints.Numeric](a, b []T, p float64) (float64, bool) {
	if math.IsNaN(p) || math.IsInf(p, 1) || p < 1 {
		return 0, false
	}
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	// Factor out the largest absolute difference so math.Pow operates on values
	// in [0, 1], avoiding overflow/underflow for large p. A NaN difference
	// propagates per the package policy; an infinite max dominates the sum.
	var maxDiff float64
	for i := range a {
		diff := math.Abs(float64(a[i]) - float64(b[i]))
		if math.IsNaN(diff) {
			return math.NaN(), true
		}
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	if math.IsInf(maxDiff, 1) {
		return math.Inf(1), true
	}
	if maxDiff == 0 {
		return 0, true
	}
	var sum float64
	for i := range a {
		diff := math.Abs(float64(a[i]) - float64(b[i]))
		sum += math.Pow(diff/maxDiff, p)
	}
	return maxDiff * math.Pow(sum, 1/p), true
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
