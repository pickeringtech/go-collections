package slices

import (
	"github.com/pickeringtech/go-collectionutil/constraints"
	"sort"
)

// SortFunc is a function which compares the relative values of two elements of a slice, returning a boolean value
// indicating whether the positions of `a` and `b` should be switched.  Ascending is `a < b`, descending is `a > b`.
type SortFunc[T any] func(a, b T) bool

// Ensure AscendingSortFunc implements SortFunc.
var _ SortFunc[int] = AscendingSortFunc[int]

// AscendingSortFunc is a sort function which naturally orders an input slice in ascending order by the value of the elements.
func AscendingSortFunc[T constraints.Ordered](a, b T) bool {
	return a < b
}

// Ensure DescendingSortFunc implements SortFunc.
var _ SortFunc[int] = DescendingSortFunc[int]

// DescendingSortFunc is a sort function which naturally orders an input slice in descending order by the value of the elements.
func DescendingSortFunc[T constraints.Ordered](a, b T) bool {
	return a > b
}

// Sort orders the elements within the input slice in order, using the provided function to determine the
// relative value of each element, and whether they should be before or after each other.
func Sort[T any](input []T, fun SortFunc[T]) []T {
	if len(input) == 0 {
		return nil
	}
	inputCopy := append([]T(nil), input...)
	sort.Slice(inputCopy, func(i, j int) bool {
		a, b := inputCopy[i], inputCopy[j]
		return fun(a, b)
	})
	return inputCopy
}

// SortFieldExtractorFunc is a function which extracts a field from an element of a slice, returning a value which can
// be compared to other values of the same type for the purposes of sorting. This can be used to sort a slice of structs
// by one of the struct member fields.
type SortFieldExtractorFunc[T any, S constraints.Ordered] func(T) S

// SortByOrderedField orders the elements within the input slice using the sort function, and using a field which is
// extracted from each element by the extractor function. Particularly useful when trying to sort a slice of structs
// by one of the struct member fields.
func SortByOrderedField[T any, S constraints.Ordered](input []T, fun SortFunc[S], extractor SortFieldExtractorFunc[T, S]) []T {
	if len(input) == 0 {
		return nil
	}
	inputCopy := append([]T(nil), input...)
	sort.Slice(inputCopy, func(i, j int) bool {
		a, b := extractor(inputCopy[i]), extractor(inputCopy[j])
		return fun(a, b)
	})
	return inputCopy
}

// SortInPlace orders the elements within the input slice in order, using the provided function to determine the
// relative value of each element, and whether they should be before or after each other. The sort is performed on the
// input slice, with no copy being made.
func SortInPlace[T any](input []T, fun SortFunc[T]) {
	if len(input) == 0 {
		return
	}
	sort.Slice(input, func(i, j int) bool {
		a, b := input[i], input[j]
		return fun(a, b)
	})
	return
}

// SortOrderedAsc orders the elements within the input slice in ascending order, using their relative values to determine
// where within the slice they should be.
func SortOrderedAsc[T constraints.Ordered](input []T) []T {
	return Sort[T](input, AscendingSortFunc[T])
}

// SortOrderedAscInPlace orders the elements within the input slice in ascending order, using their relative values to determine
// where within the slice they should be.  The sort is performed on the input slice, with no copy being made.
func SortOrderedAscInPlace[T constraints.Ordered](input []T) {
	SortInPlace[T](input, AscendingSortFunc[T])
}

// SortOrderedDesc orders the elements within the input slice in descending order, using their relative values to determine
// where within the slice they should be.
func SortOrderedDesc[T constraints.Ordered](input []T) []T {
	return Sort[T](input, DescendingSortFunc[T])
}

// SortOrderedDescInPlace orders the elements within the input slice in descending order, using their relative values to determine
// where within the slice they should be.  The sort is performed on the input slice, with no copy being made.
func SortOrderedDescInPlace[T constraints.Ordered](input []T) {
	SortInPlace[T](input, DescendingSortFunc[T])
}
