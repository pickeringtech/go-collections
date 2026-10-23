package stats

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/internal/kahan"
)

// Normalize rescales input to the range [0, 1] using min-max scaling:
// each element becomes (x − min) / (max − min). It returns a new []float64 the
// same length as input, together with an ok flag; input is not modified.
//
// It returns ok == false only for empty input, where the transform is
// undefined. When every element is identical (and finite) the range is
// degenerate (max == min, a divide-by-zero); rather than produce NaNs, Normalize
// maps every element to 0 — the low end of the target range — and returns
// ok == true. A finite range whose width max − min overflows float64 (such as
// {−math.MaxFloat64, math.MaxFloat64}) is rescaled in half-units, which is exact
// at that magnitude, so it still maps onto [0, 1].
//
// Non-finite inputs (NaN/Inf) propagate per IEEE-754. A NaN element has no
// place in the ordering, so min and max are taken over the non-NaN elements and
// a NaN yields NaN at its own position only; every other position rescales as
// if it were absent (an all-NaN input yields all NaN). An infinite element that
// becomes the min or max makes the span infinite, so finite positions rescale
// through it (typically to 0) while the infinite position itself becomes NaN.
// The all-zeros guarantee for a degenerate range therefore holds only for
// finite input; clean non-finite values beforehand if that is not what you
// want.
func Normalize[T constraints.Numeric](input []T) ([]float64, bool) {
	if len(input) == 0 {
		return nil, false
	}

	lo, hi := nonNaNRange(input)
	// hi − lo can only overflow for finite values of opposite sign near
	// ±MaxFloat64 (or with an infinite bound, where halving changes nothing);
	// halving such values is exact, so rescale in half-units there.
	scale := 1.0
	if math.IsInf(hi-lo, 1) {
		scale = 0.5
	}
	span := hi*scale - lo*scale

	out := make([]float64, len(input))
	for i, v := range input {
		x := float64(v)
		if span == 0 {
			// Degenerate range: every non-NaN value is lo and maps to 0, and
			// x − x keeps a NaN at its own position.
			out[i] = x - x
			continue
		}
		out[i] = (x*scale - lo*scale) / span
	}
	return out, true
}

// nonNaNRange returns the smallest and largest non-NaN values of input, or
// (NaN, NaN) when every value is NaN. Skipping NaN keeps a single NaN from
// seeding (and so poisoning) both bounds.
func nonNaNRange[T constraints.Numeric](input []T) (lo, hi float64) {
	lo, hi = math.NaN(), math.NaN()
	for _, v := range input {
		f := float64(v)
		switch {
		case f != f:
			// NaN is unordered; leave the bounds to the other values.
		case lo != lo:
			lo, hi = f, f // the first non-NaN value seeds both bounds
		case f < lo:
			lo = f
		case f > hi:
			hi = f
		}
	}
	return lo, hi
}

// Standardize rescales input to zero mean and unit variance using the z-score
// (x − mean) / stddev, where stddev is the population standard deviation. It
// returns a new []float64 the same length as input, together with an ok flag;
// input is not modified.
//
// The mean is the package's compensated mean (the one Mean returns for finite
// data, accurate to within an ulp even over millions of large values) and the
// standard deviation comes from a numerically-stable Welford pass. It returns
// ok == false only for empty input.
// When finite data has zero spread (a single element, or all elements identical)
// every value sits exactly at the mean — zero standard deviations away — so
// Standardize returns all zeros with ok == true rather than dividing by zero.
// (Non-finite identical inputs instead propagate, per the policy below.)
//
// Non-finite inputs (NaN/Inf) propagate: a non-finite element poisons the mean
// and standard deviation, so the whole result becomes non-finite.
func Standardize[T constraints.Numeric](input []T) ([]float64, bool) {
	if len(input) == 0 {
		return nil, false
	}

	mean := compensatedMean(input)
	n, m2 := welford(input)
	stdDev := math.Sqrt(m2 / float64(n)) // population standard deviation

	out := make([]float64, len(input))
	if stdDev == 0 {
		// Zero spread: every element is the mean, i.e. zero std devs away.
		return out, true
	}
	for i, v := range input {
		out[i] = (float64(v) - mean) / stdDev
	}
	return out, true
}

// MovingAverage computes the rolling mean of input over a sliding window of the
// given size, returning a new []float64 together with an ok flag; input is not
// modified.
//
// Only full windows are emitted (the "valid" convention): the result has length
// len(input) − window + 1, where result[i] is the mean of input[i : i+window].
// Partial leading windows are deliberately not produced — every output value is
// the mean of exactly window elements, so none is a statistically weaker
// average over fewer points.
//
// Edge handling is explicit:
//   - window < 1 is invalid and returns ok == false.
//   - window > len(input) cannot form a single full window and returns
//     ok == false (this also covers empty input).
//   - window == len(input) yields a single element: the mean of the whole input.
//
// The mean is computed with an incremental running sum, so the whole transform
// is O(len(input)) regardless of window size. The running sum is compensated
// (Kahan summation in Neumaier's form), so the round-off from adding and later
// removing a large value is carried in the compensation term rather than lost:
// MovingAverage({1e17, 1, 2, 3, 4}, 2) is exactly [5e16 1.5 2.5 3.5], where a
// naive running sum would give [5e16 0 1 2]. A consequence of the running sum is that a non-finite input value (NaN/Inf)
// propagates to its own window and to every subsequent window; clean
// non-finite values beforehand if you need strict per-window locality.
func MovingAverage[T constraints.Numeric](input []T, window int) ([]float64, bool) {
	if window < 1 || window > len(input) {
		return nil, false
	}

	out := make([]float64, len(input)-window+1)
	w := float64(window)

	var sum kahan.Sum
	for i := 0; i < window; i++ {
		sum.Add(float64(input[i]))
	}
	out[0] = sum.Total() / w

	for i := window; i < len(input); i++ {
		// Slide the window: add the entering element, drop the leaving one.
		// Both go through the compensated sum separately — folding them into
		// one difference first would round away the small entering value.
		sum.Add(float64(input[i]))
		sum.Add(-float64(input[i-window]))
		out[i-window+1] = sum.Total() / w
	}
	return out, true
}
