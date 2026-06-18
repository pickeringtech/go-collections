package slices

import "github.com/pickeringtech/go-collections/constraints"

// NumericSlice represents a slice of numeric values. This type exposes some mathematical operations that can be
// performed on such a slice.
type NumericSlice[T constraints.Numeric] []T

// Avg calculates the average of the input, returning the result and whether it is
// defined. Empty or nil input yields (0, false).
func (n NumericSlice[T]) Avg() (float64, bool) {
	return Avg(n)
}

// Max finds the maximum value in the input, returning the result and whether it is
// defined. Empty or nil input yields (zero value, false).
func (n NumericSlice[T]) Max() (T, bool) {
	return Max(n)
}

// Min finds the minimum value in the input, returning the result and whether it is
// defined. Empty or nil input yields (zero value, false).
func (n NumericSlice[T]) Min() (T, bool) {
	return Min(n)
}

// Sum adds up each element of the input slice, returning the total and whether it is
// defined. Empty or nil input yields (zero value, false).
func (n NumericSlice[T]) Sum() (T, bool) {
	return Sum(n)
}

// Avg calculates the arithmetic mean of the input. The boolean result is false when
// the input is empty or nil (with a zero result), letting callers distinguish "no
// data" from a genuine zero mean.
//
// The elements are summed in T before the division, so averaging large integer
// inputs can overflow T's range and skew the result; widen T (e.g. accumulate as
// int64/float64) if that is a concern.
func Avg[T constraints.Numeric](input []T) (float64, bool) {
	if len(input) == 0 {
		return 0, false
	}
	var total T
	for _, element := range input {
		total += element
	}
	return float64(total) / float64(len(input)), true
}

// Max finds the maximum value in the input. The boolean result is false when the
// input is empty or nil, in which case the T result is the zero value.
func Max[T constraints.Ordered](input []T) (T, bool) {
	if len(input) == 0 {
		var zero T
		return zero, false
	}
	result := input[0]
	for _, element := range input[1:] {
		if element > result {
			result = element
		}
	}
	return result, true
}

// Min finds the minimum value in the input. The boolean result is false when the
// input is empty or nil, in which case the T result is the zero value.
func Min[T constraints.Ordered](input []T) (T, bool) {
	if len(input) == 0 {
		var zero T
		return zero, false
	}
	result := input[0]
	for _, element := range input[1:] {
		if element < result {
			result = element
		}
	}
	return result, true
}

// Sum adds up each element of the input slice, returning the total. The boolean
// result is false when the input is empty or nil, in which case the total is the
// zero value.
//
// The total accumulates in T, so summing large integer inputs can overflow T's
// range; widen T (e.g. accumulate as int64/float64) if that is a concern.
func Sum[T constraints.Numeric](input []T) (T, bool) {
	if len(input) == 0 {
		var zero T
		return zero, false
	}
	var result T
	for _, element := range input {
		result += element
	}
	return result, true
}
