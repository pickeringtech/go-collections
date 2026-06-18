package stats

import (
	"math"
	"sort"

	"github.com/pickeringtech/go-collections/constraints"
)

// InterpolationMethod selects how a quantile is resolved when the desired rank
// falls between two adjacent samples. The names and behaviour mirror the
// "interpolation" options of numpy.percentile so results are directly
// comparable.
type InterpolationMethod int

const (
	// Linear interpolates linearly between the two bracketing samples. This is
	// the default used by Quantile/Percentile/Quartiles/IQR. It corresponds to
	// "type 7" in Hyndman & Fan (1996) and to the default of numpy.percentile,
	// which is what most users expect.
	Linear InterpolationMethod = iota
	// Lower takes the sample at the lower of the two bracketing ranks.
	Lower
	// Higher takes the sample at the higher of the two bracketing ranks.
	Higher
	// Nearest takes the sample whose rank is nearest the desired rank (ties
	// round to the higher rank).
	Nearest
	// Midpoint takes the arithmetic mean of the two bracketing samples.
	Midpoint
)

// Quantile returns the q-quantile of input for q in [0, 1] using Linear
// ("type 7") interpolation — the numpy.percentile default. For example q=0
// is the minimum, q=0.5 the median and q=1 the maximum.
//
// The result is float64 and the second return value reports whether it is
// meaningful. ok is false (and the result 0) when input is empty, when q lies
// outside [0, 1], or when input contains a NaN — NaN poisons ordering, so the
// quantile of a NaN-contaminated sample is undefined by policy rather than
// silently wrong.
//
// The caller's slice is never mutated; input is copied before sorting.
//
// Median is Quantile(input, 0.5).
func Quantile[T constraints.Numeric](input []T, q float64) (float64, bool) {
	return QuantileWith(input, q, Linear)
}

// QuantileWith is Quantile with an explicit InterpolationMethod. See Quantile
// for the empty/range/NaN contract, which is identical.
func QuantileWith[T constraints.Numeric](input []T, q float64, method InterpolationMethod) (float64, bool) {
	if q < 0 || q > 1 {
		return 0, false
	}
	sorted, ok := sortedCopy(input)
	if !ok {
		return 0, false
	}
	return quantileSorted(sorted, q, method), true
}

// Percentile returns the p-th percentile of input for p in [0, 100]. It is
// exactly Quantile(input, p/100); see Quantile for the empty/range/NaN
// contract.
func Percentile[T constraints.Numeric](input []T, p float64) (float64, bool) {
	return PercentileWith(input, p, Linear)
}

// PercentileWith is Percentile with an explicit InterpolationMethod.
func PercentileWith[T constraints.Numeric](input []T, p float64, method InterpolationMethod) (float64, bool) {
	if p < 0 || p > 100 {
		return 0, false
	}
	return QuantileWith(input, p/100, method)
}

// QuartileSet holds the three quartiles of a sample: Q1 (the 25th percentile),
// Q2 (the median) and Q3 (the 75th percentile).
type QuartileSet struct {
	Q1 float64
	Q2 float64
	Q3 float64
}

// Quartiles returns the Q1/Q2/Q3 quartiles of input using Linear interpolation.
// The contract for the bool matches Quantile (false for empty or
// NaN-contaminated input). input is sorted once and never mutated.
func Quartiles[T constraints.Numeric](input []T) (QuartileSet, bool) {
	sorted, ok := sortedCopy(input)
	if !ok {
		return QuartileSet{}, false
	}
	return QuartileSet{
		Q1: quantileSorted(sorted, 0.25, Linear),
		Q2: quantileSorted(sorted, 0.5, Linear),
		Q3: quantileSorted(sorted, 0.75, Linear),
	}, true
}

// IQR returns the interquartile range (Q3 - Q1) of input using Linear
// interpolation. The contract for the bool matches Quantile.
func IQR[T constraints.Numeric](input []T) (float64, bool) {
	q, ok := Quartiles(input)
	if !ok {
		return 0, false
	}
	return q.Q3 - q.Q1, true
}

// sortedCopy copies input into a float64 slice and sorts it ascending without
// touching the caller's slice. It returns ok=false when input is empty or
// contains a NaN.
func sortedCopy[T constraints.Numeric](input []T) ([]float64, bool) {
	if len(input) == 0 {
		return nil, false
	}
	out := make([]float64, len(input))
	for i, v := range input {
		f := float64(v)
		if math.IsNaN(f) {
			return nil, false
		}
		out[i] = f
	}
	sort.Float64s(out)
	return out, true
}

// quantileSorted resolves the q-quantile of an already-sorted, non-empty,
// NaN-free slice. This is the single home of the interpolation maths shared by
// every public entry point.
func quantileSorted(sorted []float64, q float64, method InterpolationMethod) float64 {
	n := len(sorted)
	if n == 1 {
		return sorted[0]
	}
	// "type 7" continuous rank, 0-indexed.
	rank := q * float64(n-1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	if lo == hi {
		return sorted[lo]
	}
	frac := rank - float64(lo)
	switch method {
	case Lower:
		return sorted[lo]
	case Higher:
		return sorted[hi]
	case Nearest:
		if frac < 0.5 {
			return sorted[lo]
		}
		return sorted[hi]
	case Midpoint:
		return sorted[lo] + (sorted[hi]-sorted[lo])/2
	default: // Linear
		return sorted[lo] + frac*(sorted[hi]-sorted[lo])
	}
}
