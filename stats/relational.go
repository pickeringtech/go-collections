package stats

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
)

// moments holds the running results of one pass of Welford's online algorithm
// over a pair of equal-length series: the element count, both running means,
// both sums of squared deviations from the running mean (M2X, M2Y) and the
// co-moment C (the running sum of products of paired deviations).
//
// Covariance, Correlation and Standardize are all thin arithmetic over these
// accumulated moments. We compute them in a single stable sweep rather than the
// textbook Σxy − ΣxΣy/n form, which subtracts two large, nearly-equal numbers
// and bleeds precision (it can even go negative) on large or near-constant
// magnitudes. Welford only ever works with deltas from the running mean, so it
// stays accurate regardless of the inputs' scale.
type moments struct {
	n            int
	meanX, meanY float64
	m2X, m2Y, c  float64
}

// accumulate runs Welford's online algorithm over the paired series x and y,
// which the caller guarantees to be equal length. Every deviation is taken from
// the running mean, so the accumulation is numerically stable on large or
// near-constant magnitudes. Non-finite inputs (NaN/Inf) propagate into the
// moments rather than being silently dropped — the package's NaN/Inf policy.
//
// For single-series statistics (variance, standardization) pass the same slice
// as both x and y; the X moments then describe that series.
func accumulate[T constraints.Numeric](x, y []T) moments {
	var m moments
	for i := range x {
		fx := float64(x[i])
		fy := float64(y[i])
		m.n++
		n := float64(m.n)

		// dx/dy use the running means *before* this element is folded in; the
		// means are then updated, and the moment terms use the post-update mean
		// for the second factor. This is the standard online update for the
		// co-moment and for each series' M2.
		dx := fx - m.meanX
		dy := fy - m.meanY
		m.meanX += dx / n
		m.meanY += dy / n
		m.c += dx * (fy - m.meanY)
		m.m2X += dx * (fx - m.meanX)
		m.m2Y += dy * (fy - m.meanY)
	}
	return m
}

// PopulationCovariance computes the population covariance of x and y — the
// average product of paired deviations from their means, with no Bessel's
// correction — using Welford's numerically-stable online algorithm.
//
// It returns ok == false when the inputs cannot be summarised: when they are
// empty or of differing lengths. A single pair yields a defined covariance of
// 0. Non-finite inputs (NaN/Inf) propagate to a non-finite result with
// ok == true. Use SampleCovariance when x and y are samples drawn from a larger
// population rather than the whole population.
func PopulationCovariance[T constraints.Numeric](x, y []T) (float64, bool) {
	if len(x) != len(y) || len(x) == 0 {
		return 0, false
	}
	m := accumulate(x, y)
	return m.c / float64(m.n), true
}

// SampleCovariance computes the sample covariance of x and y — the sum of
// products of paired deviations divided by n−1 (Bessel's correction) — using
// Welford's numerically-stable online algorithm.
//
// Sample covariance is undefined for fewer than two pairs, so it returns
// ok == false for empty, single-pair or differing-length inputs. Non-finite
// inputs (NaN/Inf) propagate to a non-finite result with ok == true. Use
// PopulationCovariance when x and y are the entire population.
func SampleCovariance[T constraints.Numeric](x, y []T) (float64, bool) {
	if len(x) != len(y) || len(x) < 2 {
		return 0, false
	}
	m := accumulate(x, y)
	return m.c / float64(m.n-1), true
}

// Correlation computes the Pearson product-moment correlation coefficient of x
// and y — their covariance divided by the product of their standard deviations,
// a scale-free measure of linear association in [−1, 1].
//
// Because the n (or n−1) factors cancel in the ratio, sample and population
// conventions yield the same coefficient, so a single function suffices. It is
// computed from one stable Welford pass as C / √(M2X · M2Y).
//
// It returns ok == false when the coefficient is undefined: when the inputs are
// empty, of differing lengths, have fewer than two pairs, or when either series
// is constant (zero variance, so there is no linear relationship to measure).
// Non-finite inputs (NaN/Inf) propagate to a non-finite result with ok == true.
func Correlation[T constraints.Numeric](x, y []T) (float64, bool) {
	if len(x) != len(y) || len(x) < 2 {
		return 0, false
	}
	m := accumulate(x, y)
	denom := math.Sqrt(m.m2X * m.m2Y)
	if denom == 0 {
		// A constant series has zero variance, so the correlation is 0/0 —
		// undefined rather than zero. NaN/Inf inputs make denom non-finite
		// (not zero) and so fall through to propagate as documented.
		return 0, false
	}
	return m.c / denom, true
}
