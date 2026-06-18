package stats

import "github.com/pickeringtech/go-collections/constraints"

// Median returns the middle value of input as a float64: for an odd-length
// sample it is the single middle element, and for an even-length sample it is
// the arithmetic mean of the two middle elements (hence float64, since that
// mean need not be representable in T). The caller's slice is never mutated; a
// sorted copy is taken internally.
//
// Median is exactly Quantile(input, 0.5) and shares its contract: the second
// return value is false — and the result 0 — when input is empty or contains a
// non-finite value (NaN or ±Inf), for which the order statistic is undefined.
func Median[T constraints.Numeric](input []T) (float64, bool) {
	return Quantile(input, 0.5)
}
