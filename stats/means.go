package stats

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
)

// WeightedMean returns the weighted arithmetic mean of values, i.e.
// Σ(weightᵢ·valueᵢ) / Σ(weightᵢ), together with an ok flag.
//
// Both sums are accumulated with Kahan compensated summation for stability on
// large inputs.
//
// ok is false (and the result is 0) when the input cannot be summarised:
//   - values is empty, or len(values) != len(weights);
//   - any weight is negative (a weight is, by definition, non-negative);
//   - the weights sum to zero (the mean would divide by zero);
//   - any value or weight is non-finite (NaN or ±Inf).
//
// Weights need not sum to one; they are normalised by their own total.
func WeightedMean[T constraints.Numeric](values, weights []T) (float64, bool) {
	if len(values) == 0 || len(values) != len(weights) {
		return 0, false
	}

	var weighted, total kahan
	for i, v := range values {
		x := float64(v)
		w := float64(weights[i])
		if w < 0 || nonFinite(x) || nonFinite(w) {
			return 0, false
		}
		weighted.add(w * x)
		total.add(w)
	}

	if total.sum == 0 {
		return 0, false
	}
	return weighted.sum / total.sum, true
}

// GeometricMean returns the geometric mean of values — the nth root of their
// product — together with an ok flag. It is the right average for ratios,
// growth rates and other multiplicative quantities.
//
// The mean is computed in log space (exp of the mean of the natural logs)
// rather than as a direct product, so it neither overflows nor loses precision
// on large inputs; the log sum uses Kahan compensated summation.
//
// ok is false (and the result is 0) when the input cannot be summarised:
//   - values is empty;
//   - any value is non-positive (≤ 0), for which the geometric mean is
//     undefined;
//   - any value is non-finite (NaN or +Inf).
func GeometricMean[T constraints.Numeric](values []T) (float64, bool) {
	if len(values) == 0 {
		return 0, false
	}

	var logSum kahan
	for _, v := range values {
		x := float64(v)
		if !(x > 0) || math.IsInf(x, 1) {
			return 0, false
		}
		logSum.add(math.Log(x))
	}

	return math.Exp(logSum.sum / float64(len(values))), true
}

// HarmonicMean returns the harmonic mean of values — n divided by the sum of
// their reciprocals — together with an ok flag. It is the right average for
// rates defined per unit (speeds, prices-per-share), where the arithmetic mean
// would overweight large values.
//
// The reciprocal sum uses Kahan compensated summation.
//
// ok is false (and the result is 0) when the input cannot be summarised:
//   - values is empty;
//   - any value is non-positive (≤ 0); zero divides by zero, and mixing signs
//     makes the harmonic mean ill-defined, so all values must be positive;
//   - any value is non-finite (NaN or +Inf).
func HarmonicMean[T constraints.Numeric](values []T) (float64, bool) {
	if len(values) == 0 {
		return 0, false
	}

	var recipSum kahan
	for _, v := range values {
		x := float64(v)
		if !(x > 0) || math.IsInf(x, 1) {
			return 0, false
		}
		recipSum.add(1 / x)
	}

	// Every value is positive and finite, so each reciprocal is positive and
	// the sum is strictly greater than zero — no divide-by-zero guard needed.
	return float64(len(values)) / recipSum.sum, true
}
