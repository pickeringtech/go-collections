package slices

import "github.com/pickeringtech/go-collectionutil/constraints"

type ReductionFunc[I, O any] func(accum O, currVal I) O

func ReductionTotalFunc[T constraints.Numeric](accum, currVal T) T {
	return accum + currVal
}

func NewReductionCountOccurrencesFunc[I comparable, O constraints.Numeric](toCount []I) ReductionFunc[I, O] {
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

func Reduce[I, O any](input []I, fn ReductionFunc[I, O]) O {
	var accumulator O
	for _, el := range input {
		accumulator = fn(accumulator, el)
	}
	return accumulator
}
