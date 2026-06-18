package slices

import (
	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/stats"
)

// NumericSlice represents a slice of numeric values. It is a thin ergonomic
// accessor: every method below delegates to the canonical implementation (the
// numeric summaries live in the stats package, the ordering reductions in
// slices) and holds no logic of its own.
type NumericSlice[T constraints.Numeric] []T

// Mean returns the arithmetic mean of the slice and whether it is defined.
// Empty or nil input yields (0, false). It delegates to stats.Mean.
func (n NumericSlice[T]) Mean() (float64, bool) {
	return stats.Mean(n)
}

// Sum returns the total of the slice and whether it is defined. Empty or nil
// input yields (zero value, false). It delegates to stats.Sum.
func (n NumericSlice[T]) Sum() (T, bool) {
	return stats.Sum(n)
}

// Max returns the maximum value of the slice and whether it is defined. Empty
// or nil input yields (zero value, false). It delegates to slices.Max.
func (n NumericSlice[T]) Max() (T, bool) {
	return Max(n)
}

// Min returns the minimum value of the slice and whether it is defined. Empty
// or nil input yields (zero value, false). It delegates to slices.Min.
func (n NumericSlice[T]) Min() (T, bool) {
	return Min(n)
}

// Max finds the maximum value in the input. The boolean result is false when the
// input is empty or nil, in which case the T result is the zero value.
//
// Max is an ordering reduction (it works on any constraints.Ordered, strings
// included, and matches the standard library's slices.Max), so it lives here in
// slices rather than in stats.
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
//
// Min is an ordering reduction (it works on any constraints.Ordered, strings
// included, and matches the standard library's slices.Min), so it lives here in
// slices rather than in stats.
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
