package stats

import (
	"math"
	"math/bits"

	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/internal/kahan"
)

// Mean returns the arithmetic mean (average) of input together with an ok flag.
// Empty or nil input yields (0, false): the mean of no data is undefined, so it
// is reported via the (result, ok) idiom rather than as a silent zero.
//
// The values are summed in float64 with Kahan compensated summation, so large
// inputs do not lose precision to naive round-off, and integer inputs do not
// overflow the way an exact-in-T Sum would.
//
// The mean of finite values is always finite: if the sum itself would overflow
// float64 (for example Mean({1e308, 1e308})), the values are rescaled by a power
// of two before summing, so the result is the correct finite mean rather than
// ±Inf. The result is also clamped to [min(input), max(input)], the interval a
// mean always lies within, so round-off cannot push it outside the data: the
// mean of identical values is exactly that value (Mean({0.1, 0.1, 0.1}) is 0.1,
// not 0.10000000000000002).
//
// ok is false (and the result is 0) when any value is non-finite (NaN or ±Inf),
// because the resulting mean would be undefined.
func Mean[T constraints.Numeric](input []T) (float64, bool) {
	if len(input) == 0 {
		return 0, false
	}
	// The mean of finite values is always finite, and a non-finite value always
	// makes it NaN or ±Inf, so a non-finite mean is exactly the rejection case.
	m := compensatedMean(input)
	if nonFinite(m) {
		return 0, false
	}
	return m, true
}

// compensatedMean is the package's single arithmetic-mean computation, shared
// by Mean, the central moments behind Skewness/Kurtosis, and Standardize: a
// compensated sum divided by the count, recomputed by overflowSafeMean when
// the sum overflows, and clamped to [min(input), max(input)]. input must be
// non-empty. Non-finite values propagate (a NaN gives NaN, a lone ±Inf gives
// ±Inf), which lets the propagating callers use it unchanged and Mean reject
// them.
func compensatedMean[T constraints.Numeric](input []T) float64 {
	var total kahan.Sum
	lo, hi := math.Inf(1), math.Inf(-1)
	for _, v := range input {
		x := float64(v)
		total.Add(x)
		lo = min(lo, x) // min/max propagate NaN into the clamp below
		hi = max(hi, x)
	}
	m := total.Total() / float64(len(input))
	if math.IsInf(total.Total(), 0) {
		// For finite values an infinite total is overflow, not data; for an
		// infinite value the rescaled sum is the same ±Inf.
		m = overflowSafeMean(input)
	}
	return min(max(m, lo), hi)
}

// overflowSafeMean recomputes the mean of finite input whose plain sum
// overflowed. Every value is first scaled by 2^-k, where 2^k > len(input), so
// no partial sum can exceed math.MaxFloat64; the mean of the scaled values is
// then scaled back by 2^k. Scaling by a power of two is exact, so only values
// too small to matter next to an overflowing total lose any bits.
func overflowSafeMean[T constraints.Numeric](input []T) float64 {
	k := bits.Len(uint(len(input)))
	scale := math.Ldexp(1, -k)
	var total kahan.Sum
	for _, v := range input {
		total.Add(float64(v) * scale)
	}
	return math.Ldexp(total.Total()/float64(len(input)), k)
}

// WeightedMean returns the weighted arithmetic mean of values, i.e.
// Σ(weightᵢ·valueᵢ) / Σ(weightᵢ), together with an ok flag.
//
// Both sums are accumulated with Kahan compensated summation for stability on
// large inputs. The weighted mean of finite values is always finite: if a
// product weightᵢ·valueᵢ or either sum would overflow float64, the values and
// the weights are each rescaled by a power of two so that every term stays
// small, and the result is scaled back — so WeightedMean({1e200, 1},
// {1e200, 1}) is ≈1e200 rather than ±Inf or NaN. Like Mean, the result is
// clamped to the range of the values that carry a positive weight, so the
// weighted mean of identical values is exactly that value.
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

	var weighted, total kahan.Sum
	lo, hi := math.Inf(1), math.Inf(-1)
	for i, v := range values {
		x := float64(v)
		w := float64(weights[i])
		if w < 0 || nonFinite(x) || nonFinite(w) {
			return 0, false
		}
		weighted.Add(w * x)
		total.Add(w)
		if w > 0 {
			lo = min(lo, x)
			hi = max(hi, x)
		}
	}

	if total.Total() == 0 {
		return 0, false
	}
	mean := weighted.Total() / total.Total()
	if nonFinite(weighted.Total()) || math.IsInf(total.Total(), 1) {
		// Every input is finite, so a non-finite sum is overflow (a product
		// overflowing to ±Inf, or opposite overflows cancelling to NaN).
		mean = overflowSafeWeightedMean(values, weights)
	}
	// A positive total means at least one weight is positive, so lo/hi are set.
	return min(max(mean, lo), hi), true
}

// overflowSafeWeightedMean recomputes WeightedMean for finite, validated input
// whose plain sums overflowed. The values are scaled by 2^-ex and the weights by
// 2^-ew, where 2^ex and 2^ew exceed the largest |value| and weight, so every
// scaled value and weight is below 1 in magnitude: the products and both sums
// then stay far below math.MaxFloat64. The weights' scale cancels in the ratio
// and the values' scale is undone at the end. Scaling by a power of two is
// exact, so only terms too small to matter next to the overflowing ones lose
// any bits. The total is non-zero because the largest weight scales to at
// least 0.5.
func overflowSafeWeightedMean[T constraints.Numeric](values, weights []T) float64 {
	_, ex := math.Frexp(maxAbs(values))
	_, ew := math.Frexp(maxAbs(weights))
	var weighted, total kahan.Sum
	for i, v := range values {
		w := math.Ldexp(float64(weights[i]), -ew)
		weighted.Add(w * math.Ldexp(float64(v), -ex))
		total.Add(w)
	}
	return math.Ldexp(weighted.Total()/total.Total(), ex)
}

// GeometricMean returns the geometric mean of values — the nth root of their
// product — together with an ok flag. It is the right average for ratios,
// growth rates and other multiplicative quantities.
//
// The mean is computed in log space rather than as a direct product, so it
// neither overflows nor loses precision on large inputs. Each value is split
// into its binary exponent and a fraction in [0.5, 1) (math.Frexp): the
// exponents are summed exactly as integers and only the fractions' logs are
// accumulated (with Kahan compensated summation), which keeps the result
// accurate across the whole float64 range, subnormals included. The result is
// finally clamped to [min(values), max(values)] — the interval a geometric mean
// always lies within — so log/exp round-off can never push it past the largest
// value: GeometricMean({math.MaxFloat64}) is math.MaxFloat64, not +Inf.
//
// ok is false (and the result is 0) when the input cannot be summarised:
//   - values is empty;
//   - any value is non-positive (≤ 0), for which the geometric mean is
//     undefined;
//   - any value is non-finite (NaN or +Inf).
func GeometricMean[T constraints.Numeric](values []T) (float64, bool) {
	lo, hi, ok := positiveRange(values)
	if !ok {
		return 0, false
	}

	var logFrac kahan.Sum
	var exponents int64
	for _, v := range values {
		frac, exp := math.Frexp(float64(v))
		logFrac.Add(math.Log(frac))
		exponents += int64(exp)
	}

	// mean(log x) = mean(log frac) + (exponents/n)·ln 2. Split exponents/n into
	// an integer part q, applied exactly with Ldexp, and a remainder r/n with
	// |r/n| < 1, so the argument to Exp stays within (−1.4, 0.7) and cannot
	// overflow or underflow.
	n := int64(len(values))
	q, r := exponents/n, exponents%n
	gm := math.Ldexp(math.Exp(logFrac.Total()/float64(n)+float64(r)/float64(n)*math.Ln2), int(q))
	return min(max(gm, lo), hi), true
}

// HarmonicMean returns the harmonic mean of values — n divided by the sum of
// their reciprocals — together with an ok flag. It is the right average for
// rates defined per unit (speeds, prices-per-share), where the arithmetic mean
// would overweight large values.
//
// The reciprocal sum uses Kahan compensated summation. The values are first
// rescaled by a power of two so the smallest lies in [0.5, 1): every scaled
// reciprocal is then at most 2, so the sum cannot overflow however tiny the
// values are (HarmonicMean({1e-310, 1e-310}) is 1e-310, not 0), and the
// reciprocals of huge values underflow only when they are negligible next to
// the smallest value's. Scaling by a power of two is exact, so ordinary inputs
// give the same result as the unscaled sum. As with GeometricMean, the result
// is clamped to [min(values), max(values)], the interval a harmonic mean always
// lies within.
//
// ok is false (and the result is 0) when the input cannot be summarised:
//   - values is empty;
//   - any value is non-positive (≤ 0); zero divides by zero, and mixing signs
//     makes the harmonic mean ill-defined, so all values must be positive;
//   - any value is non-finite (NaN or +Inf).
func HarmonicMean[T constraints.Numeric](values []T) (float64, bool) {
	lo, hi, ok := positiveRange(values)
	if !ok {
		return 0, false
	}

	_, e := math.Frexp(lo)
	var recipSum kahan.Sum
	for _, v := range values {
		recipSum.Add(1 / math.Ldexp(float64(v), -e))
	}

	// The smallest value scales into [0.5, 1), so its reciprocal alone makes
	// the sum exceed 1 — no divide-by-zero guard is needed.
	hm := math.Ldexp(float64(len(values))/recipSum.Total(), e)
	return min(max(hm, lo), hi), true
}

// positiveRange returns the smallest and largest of values, with ok == false
// when values is empty or any value is not a positive, finite number — the
// domain shared by the geometric and harmonic means.
func positiveRange[T constraints.Numeric](values []T) (lo, hi float64, ok bool) {
	if len(values) == 0 {
		return 0, 0, false
	}
	lo, hi = math.MaxFloat64, 0
	for _, v := range values {
		x := float64(v)
		if !(x > 0) || math.IsInf(x, 1) {
			return 0, 0, false
		}
		lo = min(lo, x)
		hi = max(hi, x)
	}
	return lo, hi, true
}
