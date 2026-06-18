package stats

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
)

// welford runs Welford's online algorithm over input, returning the element
// count and M2 — the sum of squared deviations from the running mean.
//
// We deliberately accumulate M2 incrementally rather than computing the naive
// Σx² − (Σx)²/n. That textbook formula subtracts two large, nearly-equal
// numbers and loses catastrophic precision (even going negative) on large or
// near-constant magnitudes. Welford only ever works with deltas from the
// running mean, so it stays accurate regardless of the inputs' scale.
//
// NaN/Inf propagate naturally: any non-finite element poisons the running mean
// and therefore M2, so callers surface a non-finite statistic rather than a
// silently-wrong finite one. This matches the variance/covariance/correlation
// family's package-documented propagate policy.
func welford[T constraints.Numeric](input []T) (n int, m2 float64) {
	var mean float64
	for i, x := range input {
		fx := float64(x)
		delta := fx - mean
		mean += delta / float64(i+1)
		delta2 := fx - mean
		m2 += delta * delta2
	}
	return len(input), m2
}

// PopulationVariance computes the population variance (the average squared
// deviation from the mean, no Bessel's correction) of input using Welford's
// numerically-stable online algorithm.
//
// It returns ok == false for empty input, where variance is undefined; a
// single element yields a defined variance of 0. Non-finite inputs (NaN/Inf)
// propagate to a non-finite result with ok == true — consistent with the
// variance/covariance/correlation family — rather than being silently dropped.
// Use SampleVariance when input is a sample drawn from a larger population
// rather than the whole population.
func PopulationVariance[T constraints.Numeric](input []T) (float64, bool) {
	n, m2 := welford(input)
	if n < 1 {
		return 0, false
	}
	return m2 / float64(n), true
}

// SampleVariance computes the sample variance of input — the sum of squared
// deviations divided by n−1 (Bessel's correction) — using Welford's
// numerically-stable online algorithm.
//
// Sample variance is undefined for fewer than two elements, so it returns
// ok == false for empty or single-element input. Non-finite inputs (NaN/Inf)
// propagate to a non-finite result with ok == true. Use PopulationVariance
// when input is the entire population rather than a sample of it.
func SampleVariance[T constraints.Numeric](input []T) (float64, bool) {
	n, m2 := welford(input)
	if n < 2 {
		return 0, false
	}
	return m2 / float64(n-1), true
}

// PopulationStdDev computes the population standard deviation — the square root
// of the population variance — of input.
//
// It mirrors PopulationVariance's contract: ok == false for empty input, 0 for
// a single element, and non-finite results propagate with ok == true.
func PopulationStdDev[T constraints.Numeric](input []T) (float64, bool) {
	variance, ok := PopulationVariance(input)
	if !ok {
		return 0, false
	}
	return math.Sqrt(variance), true
}

// SampleStdDev computes the sample standard deviation — the square root of the
// sample variance — of input.
//
// It mirrors SampleVariance's contract: ok == false for fewer than two
// elements, and non-finite results propagate with ok == true.
func SampleStdDev[T constraints.Numeric](input []T) (float64, bool) {
	variance, ok := SampleVariance(input)
	if !ok {
		return 0, false
	}
	return math.Sqrt(variance), true
}
