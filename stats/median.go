package stats

import (
	"sort"

	"github.com/pickeringtech/go-collections/constraints"
)

// Median returns the middle value of input as a float64: for an odd-length
// sample it is the single middle element, and for an even-length sample it is
// the arithmetic mean of the two middle elements (hence float64, since that
// mean need not be representable in T). The caller's slice is never mutated; a
// sorted copy is taken internally.
//
// The second return value is false — and the result 0 — when input is empty or
// contains a non-finite value (NaN or ±Inf): NaN has no defined position in a
// sorted order, so the median of a contaminated sample is undefined by policy
// rather than silently wrong.
func Median[T constraints.Numeric](input []T) (float64, bool) {
	sorted, ok := sortedFloats(input)
	if !ok {
		return 0, false
	}
	n := len(sorted)
	mid := n / 2
	if n%2 == 1 {
		return sorted[mid], true
	}
	return (sorted[mid-1] + sorted[mid]) / 2, true
}

// sortedFloats copies input into an ascending-sorted float64 slice without
// touching the caller's slice. It returns ok=false when input is empty or
// contains a non-finite value (NaN or ±Inf).
func sortedFloats[T constraints.Numeric](input []T) ([]float64, bool) {
	if len(input) == 0 {
		return nil, false
	}
	out := make([]float64, len(input))
	for i, v := range input {
		if nonFiniteT(v) {
			return nil, false
		}
		out[i] = float64(v)
	}
	sort.Float64s(out)
	return out, true
}
