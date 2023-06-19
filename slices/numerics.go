package slices

import "github.com/pickeringtech/go-collectionutil/constraints"

// NumericSlice represents a slice of numeric values. This type exposes some mathematical operations that can be
// performed on such a slice.
type NumericSlice[T constraints.Numeric] []T

func (n NumericSlice[T]) Avg() float64 {
	return Avg(n)
}

func (n NumericSlice[T]) Max() T {
	return Max(n)
}

func (n NumericSlice[T]) Min() T {
	return Min(n)
}

func (n NumericSlice[T]) Sum() T {
	return Sum(n)
}

// Avg calculates the average of the input, returning the result.  Empty or nil input results in zero.
func Avg[T constraints.Numeric](input []T) float64 {
	var total T
	for _, element := range input {
		total += element
	}
	if total == 0 {
		return 0
	}
	return float64(total) / float64(len(input))
}

// Max finds the maximum value in the input, returning the result.  Empty or nil input results in zero.
func Max[T constraints.Ordered](input []T) T {
	var result T
	for _, element := range input {
		if element > result {
			result = element
		}
	}
	return result
}

// Min finds the minimum value in the input, returning the result.  Empty or nil input results in max int value.
func Min[T constraints.Ordered](input []T) T {
	var result T
	if len(input) > 0 {
		result = input[0]
	}
	for _, element := range input {
		if element < result {
			result = element
		}
	}
	return result
}

// Sum adds up each element of the input slice, returning the total result.  Empty or nil input results in zero.
func Sum[T constraints.Numeric](input []T) T {
	var result T
	for _, element := range input {
		result += element
	}
	return result
}
