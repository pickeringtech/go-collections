package slices

import "github.com/pickeringtech/go-collectionutil/constraints"

// SortFunc is a function which compares the relative values of two elements of a slice, returning a boolean value
// indicating whether the positions of `a` and `b` should be switched.  Ascending is `a < b`, descending is `a > b`.
type SortFunc[T any] func(a, b T) bool

// AscendingSortFunc is a sort function which naturally orders an input slice in ascending order by the value of the elements.
func AscendingSortFunc[T constraints.Ordered](a, b T) bool {
	return a < b
}

// DescendingSortFunc is a sort function which naturally orders an input slice in descending order by the value of the elements.
func DescendingSortFunc[T constraints.Ordered](a, b T) bool {
	return a > b
}

// SortAsc orders the elements within the input slice in ascending order, using the provided function to determine the
// relative value of each element, and whether they should be before or after each other.
func SortAsc[T any](input []T, fun SortFunc[T]) []T {
	panic("implement me")
}

// SortAscInPlace orders the elements within the input slice in ascending order, using the provided function to determine the
// relative value of each element, and whether they should be before or after each other. The sort is performed on the
// input slice, with no copy being made.
func SortAscInPlace[T any](input []T, fun SortFunc[T]) {
	panic("implement me")
}

// SortDesc orders the elements within the input slice in descending order, using the provided function to determine the
// relative value of each element, and whether they should be before or after each other.
func SortDesc[T any](input []T, fun SortFunc[T]) []T {
	panic("implement me")
}

// SortDescInPlace orders the elements within the input slice in descending order, using the provided function to determine the
// relative value of each element, and whether they should be before or after each other.  The sort is performed on the
// input slice, with no copy being made.
func SortDescInPlace[T any](input []T, fun SortFunc[T]) {
	panic("implement me")
}

// SortOrderedAsc orders the elements within the input slice in ascending order, using their relative values to determine
// where within the slice they should be.
func SortOrderedAsc[T constraints.Ordered](input []T) []T {
	panic("implement me")
}

// SortOrderedAscInPlace orders the elements within the input slice in ascending order, using their relative values to determine
// where within the slice they should be.  The sort is performed on the input slice, with no copy being made.
func SortOrderedAscInPlace[T constraints.Ordered](input []T) {
	panic("implement me")
}

// SortOrderedDesc orders the elements within the input slice in descending order, using their relative values to determine
// where within the slice they should be.
func SortOrderedDesc[T constraints.Ordered](input []T) []T {
	panic("implement me")
}

// SortOrderedDescInPlace orders the elements within the input slice in descending order, using their relative values to determine
// where within the slice they should be.  The sort is performed on the input slice, with no copy being made.
func SortOrderedDescInPlace[T constraints.Ordered](input []T) {
	panic("implement me")
}
