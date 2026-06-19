package stats

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
)

// centralMoments returns the element count and the 2nd, 3rd and 4th central
// moments of input — m2 = (1/n)Σ(x−mean)², m3 = (1/n)Σ(x−mean)³ and
// m4 = (1/n)Σ(x−mean)⁴ — the building blocks of skewness and kurtosis.
//
// They are computed in two passes: a Kahan-compensated pass for the mean, then
// a second pass accumulating the powered deviations (also with Kahan). Taking
// every deviation from the final mean keeps the moments accurate, unlike the
// naive power-sum forms (Σx³, Σx⁴ …) which subtract large nearly-equal numbers
// and bleed precision on large or near-constant magnitudes. Non-finite inputs
// (NaN/Inf) propagate into the moments rather than being silently dropped,
// matching the variance family's policy.
func centralMoments[T constraints.Numeric](input []T) (n int, m2, m3, m4 float64) {
	count := len(input)
	if count == 0 {
		return 0, 0, 0, 0
	}

	var meanSum kahan
	for _, v := range input {
		meanSum.add(float64(v))
	}
	mean := meanSum.sum / float64(count)

	var s2, s3, s4 kahan
	for _, v := range input {
		d := float64(v) - mean
		d2 := d * d
		s2.add(d2)
		s3.add(d2 * d)
		s4.add(d2 * d2)
	}
	fn := float64(count)
	return count, s2.sum / fn, s3.sum / fn, s4.sum / fn
}

// Skewness returns the population skewness of input — the Fisher–Pearson
// coefficient m3 / m2^1.5, where m2 and m3 are the second and third central
// moments. It measures the asymmetry of the distribution: zero for a symmetric
// sample, positive when the longer tail is to the right, negative when it is to
// the left.
//
// The moments are accumulated stably (see centralMoments). The result is
// float64 with an ok flag.
//
// ok is false (and the result 0) when skewness is undefined: when input is
// empty, or when every element is identical (zero variance, so m2^1.5 is 0 and
// the ratio is 0/0). A single element is the degenerate constant case and is
// likewise rejected. Non-finite inputs (NaN/Inf) propagate to a non-finite
// result with ok == true, consistent with the variance family.
func Skewness[T constraints.Numeric](input []T) (float64, bool) {
	n, m2, m3, _ := centralMoments(input)
	if n == 0 || m2 == 0 {
		// Empty has no moments; a constant sample has zero spread, so the
		// normalising m2^1.5 is zero and skewness is undefined. NaN/Inf make m2
		// non-finite (not zero), so they fall through and propagate.
		return 0, false
	}
	return m3 / math.Pow(m2, 1.5), true
}

// Kurtosis returns the excess kurtosis of input — m4 / m2² − 3, where m2 and m4
// are the second and fourth central moments. The −3 normalises a normal
// distribution to zero, so positive values indicate heavier tails (more
// outlier-prone) and negative values lighter tails than the normal.
//
// The moments are accumulated stably (see centralMoments). The result is
// float64 with an ok flag.
//
// ok is false (and the result 0) when kurtosis is undefined: when input is
// empty, or when every element is identical (zero variance, so m2² is 0 and the
// ratio is 0/0). A single element is the degenerate constant case and is
// likewise rejected. Non-finite inputs (NaN/Inf) propagate to a non-finite
// result with ok == true, consistent with the variance family.
func Kurtosis[T constraints.Numeric](input []T) (float64, bool) {
	n, m2, _, m4 := centralMoments(input)
	if n == 0 || m2 == 0 {
		return 0, false
	}
	return m4/(m2*m2) - 3, true
}
