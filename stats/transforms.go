package stats

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
)

// Normalize rescales input to the range [0, 1] using min-max scaling:
// each element becomes (x − min) / (max − min). It returns a new []float64 the
// same length as input, together with an ok flag; input is not modified.
//
// It returns ok == false only for empty input, where the transform is
// undefined. When every element is identical (and finite) the range is
// degenerate (max == min, a divide-by-zero); rather than produce NaNs, Normalize
// maps every element to 0 — the low end of the target range — and returns
// ok == true.
//
// Non-finite inputs (NaN/Inf) propagate per IEEE-754: a NaN element yields a NaN
// at its position, and a non-finite element that becomes the min or max makes
// the span non-finite, so finite positions rescale through it (typically to 0 or
// a non-finite value) while the non-finite position itself stays non-finite.
// The all-zeros guarantee for a degenerate range therefore holds only for finite
// input; clean non-finite values beforehand if that is not what you want.
func Normalize[T constraints.Numeric](input []T) ([]float64, bool) {
	if len(input) == 0 {
		return nil, false
	}

	lo := float64(input[0])
	hi := lo
	for _, v := range input {
		f := float64(v)
		if f < lo {
			lo = f
		}
		if f > hi {
			hi = f
		}
	}

	out := make([]float64, len(input))
	span := hi - lo
	if span == 0 {
		// Degenerate range: every value is identical, so they all map to 0.
		return out, true
	}
	for i, v := range input {
		out[i] = (float64(v) - lo) / span
	}
	return out, true
}

// Standardize rescales input to zero mean and unit variance using the z-score
// (x − mean) / stddev, where stddev is the population standard deviation. It
// returns a new []float64 the same length as input, together with an ok flag;
// input is not modified.
//
// The mean and standard deviation are computed in a single numerically-stable
// Welford pass (see accumulate). It returns ok == false only for empty input.
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

	m := accumulate(input, input)
	mean := m.meanX
	stdDev := math.Sqrt(m.m2X / float64(m.n)) // population standard deviation

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
// is O(len(input)) regardless of window size. A consequence of the running sum
// is that a non-finite input value (NaN/Inf) propagates to its own window and
// to every subsequent window; clean non-finite values beforehand if you need
// strict per-window locality.
func MovingAverage[T constraints.Numeric](input []T, window int) ([]float64, bool) {
	if window < 1 || window > len(input) {
		return nil, false
	}

	out := make([]float64, len(input)-window+1)
	w := float64(window)

	var sum float64
	for i := 0; i < window; i++ {
		sum += float64(input[i])
	}
	out[0] = sum / w

	for i := window; i < len(input); i++ {
		// Slide the window: add the entering element, drop the leaving one.
		sum += float64(input[i]) - float64(input[i-window])
		out[i-window+1] = sum / w
	}
	return out, true
}
