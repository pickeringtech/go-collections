package slices

import "github.com/pickeringtech/go-collections/constraints"

// ReductionFunc is a function that can be used to reduce a slice of values to a single value.
type ReductionFunc[I, O any] func(accum O, currVal I) O

// TotalReducer is a reduction function that can be used to sum up a slice of numeric values.
func TotalReducer[T constraints.Numeric](accum, currVal T) T {
	return accum + currVal
}

// NewCountOccurrencesReducer returns a reduction function that can be used to count the number of occurrences of a
// number of given values in a slice.
func NewCountOccurrencesReducer[I comparable, O constraints.Numeric](toCount []I) ReductionFunc[I, O] {
	return func(accum O, currVal I) O {
		var toAdd O
		for _, el := range toCount {
			if currVal == el {
				toAdd++
			}
		}
		return accum + toAdd
	}
}

// Reduce iterates over each element of the input, applying the provided reduction function, producing a single value
// which is the result of the reduction function.  If the input is empty or nil, the output will be nil.
func Reduce[I, O any](input []I, fn ReductionFunc[I, O]) O {
	var accumulator O
	for _, el := range input {
		accumulator = fn(accumulator, el)
	}
	return accumulator
}
