package stats

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
)

// Dot returns the dot product of the vectors a and b — Σ aᵢ·bᵢ — together with
// an ok flag. The products are summed with Kahan compensated summation, so long
// vectors do not lose precision to naive round-off.
//
// ok is false (and the result 0) when the product is undefined: when the
// vectors are empty or of differing lengths. Non-finite inputs (NaN/Inf)
// propagate to a non-finite result with ok == true, consistent with the
// transforms and relational statistics, rather than being silently dropped.
func Dot[T constraints.Numeric](a, b []T) (float64, bool) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	var sum kahan
	for i := range a {
		sum.add(float64(a[i]) * float64(b[i]))
	}
	return sum.sum, true
}

// Norm returns the Euclidean (L2) norm of the vector a — √(Σ aᵢ²), the length
// of the vector — together with an ok flag. The squares are summed with Kahan
// compensated summation.
//
// ok is false (and the result 0) only for empty input, where the norm is
// undefined. Non-finite inputs (NaN/Inf) propagate to a non-finite result with
// ok == true.
func Norm[T constraints.Numeric](a []T) (float64, bool) {
	if len(a) == 0 {
		return 0, false
	}
	var sum kahan
	for _, v := range a {
		f := float64(v)
		sum.add(f * f)
	}
	return math.Sqrt(sum.sum), true
}

// EuclideanDistance returns the straight-line distance between the points a and
// b — √(Σ (aᵢ−bᵢ)²) — together with an ok flag. The squared differences are
// summed with Kahan compensated summation.
//
// ok is false (and the result 0) when the distance is undefined: when the
// vectors are empty or of differing lengths. Non-finite inputs (NaN/Inf)
// propagate to a non-finite result with ok == true.
func EuclideanDistance[T constraints.Numeric](a, b []T) (float64, bool) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	var sum kahan
	for i := range a {
		d := float64(a[i]) - float64(b[i])
		sum.add(d * d)
	}
	return math.Sqrt(sum.sum), true
}

// CosineSimilarity returns the cosine of the angle between the vectors a and b —
// their dot product divided by the product of their norms — a scale-free measure
// of orientation in [−1, 1]: 1 when they point the same way, 0 when orthogonal,
// −1 when opposite. It is the standard similarity measure for embeddings and
// other high-dimensional feature vectors.
//
// ok is false (and the result 0) when the similarity is undefined: when the
// vectors are empty, of differing lengths, or when either vector is the zero
// vector (a zero norm, so the ratio is 0/0 — a zero vector has no orientation).
// Non-finite inputs (NaN/Inf) make a norm non-finite (not zero) and so fall
// through to propagate to a non-finite result with ok == true.
func CosineSimilarity[T constraints.Numeric](a, b []T) (float64, bool) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	var dot, sumA, sumB kahan
	for i := range a {
		fa := float64(a[i])
		fb := float64(b[i])
		dot.add(fa * fb)
		sumA.add(fa * fa)
		sumB.add(fb * fb)
	}
	denom := math.Sqrt(sumA.sum) * math.Sqrt(sumB.sum)
	if denom == 0 {
		// A zero vector has no orientation, so the cosine is 0/0 — undefined
		// rather than zero. NaN/Inf inputs make denom non-finite (not zero) and
		// so fall through to propagate as documented.
		return 0, false
	}
	return dot.sum / denom, true
}
