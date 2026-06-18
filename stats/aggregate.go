package stats

import "github.com/pickeringtech/go-collections/constraints"

// Product multiplies every element of input together, returning the result
// exact in T (it mirrors a sum, but with multiplication). For example the
// product of {2, 3, 4} is 24.
//
// The second return value is false — and the result the zero value of T — when
// input is empty (the empty product is undefined here by the (result, ok)
// convention rather than silently 1) or when input contains a non-finite value
// (NaN or ±Inf), which would make the result undefined.
func Product[T constraints.Numeric](input []T) (T, bool) {
	var zero T
	if len(input) == 0 {
		return zero, false
	}
	product := T(1)
	for _, v := range input {
		if nonFinite(float64(v)) {
			return zero, false
		}
		product *= v
	}
	return product, true
}

// Range returns the spread of input — its maximum minus its minimum — exact in
// T. For example the range of {3, 1, 4, 1, 5} is 4 (5 - 1). It is computed in a
// single pass.
//
// The second return value is false — and the result the zero value of T — when
// input is empty or contains a non-finite value (NaN or ±Inf), which poisons
// the min/max ordering the range depends on.
func Range[T constraints.Numeric](input []T) (T, bool) {
	var zero T
	if len(input) == 0 {
		return zero, false
	}
	lo, hi := input[0], input[0]
	for _, v := range input {
		if nonFinite(float64(v)) {
			return zero, false
		}
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	return hi - lo, true
}

// CumulativeSum returns the running prefix sums of input as a new slice of the
// same length, exact in T: out[i] is the sum of input[0..i]. For example the
// cumulative sum of {3, 1, 4} is {3, 4, 8}. The caller's slice is never
// mutated, and empty or nil input yields an empty slice.
//
// Unlike the single-value reductions, CumulativeSum has no ok flag to report a
// non-finite input, so any NaN or ±Inf simply propagates into that prefix and
// every prefix after it, per IEEE-754.
func CumulativeSum[T constraints.Numeric](input []T) []T {
	out := make([]T, len(input))
	var running T
	for i, v := range input {
		running += v
		out[i] = running
	}
	return out
}
