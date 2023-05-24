package slices

// AllMatch tests each element of the input with the provided function.  If all the elements, when passed through the
// function result in a truthy boolean value, true is returned from this function.  Otherwise, false is returned.
func AllMatch[T any](input []T, fun FindFunc[T]) bool {
	for _, element := range input {
		if !fun(element) {
			return false
		}
	}
	return true
}

// AnyMatch tests each element of the input with the provided function.  If any of the elements, when passed through the
// function result in a truthy boolean value, a match is found and true is returned from this function.  Otherwise, false
// is returned.
func AnyMatch[T any](input []T, fun FindFunc[T]) bool {
	for _, element := range input {
		if fun(element) {
			return true
		}
	}
	return false
}

type FindFunc[T any] func(T) bool

// Find tests each element of the input with the provided function.  If the function returns true, the selected element
// is returned, along with a boolean truthy value.
func Find[T any](input []T, fun FindFunc[T]) (result T, ok bool) {
	for _, element := range input {
		if fun(element) {
			return element, true
		}
	}
	return
}

// FindIndex tests each element of input with the provided testing function, and returns the index of the first element
// that satisfies the testing function.  If no matches are found, -1 is returned.
func FindIndex[T any](input []T, fun FindFunc[T]) int {
	for idx, element := range input {
		if fun(element) {
			return idx
		}
	}
	return -1
}

// FindLast tests each element of the input with the provided function, starting from the end and working background.
// If the function returns true, the selected element is returned, along with a boolean truthy value.
func FindLast[T any](input []T, fun FindFunc[T]) (result T, ok bool) {
	for i := len(input) - 1; i >= 0; i-- {
		element := input[i]
		if fun(element) {
			return element, true
		}
	}
	return
}

// FindLastIndex tests each element of input with the provided testing function, and returns the index of the last element
// that satisfies the testing function.  If no matches are found, -1 is returned.
func FindLastIndex[T any](input []T, fun FindFunc[T]) int {
	for i := len(input) - 1; i >= 0; i-- {
		element := input[i]
		if fun(element) {
			return i
		}
	}
	return -1
}

// First provides the first element of the input slice.  If there is no possible element to return, a boolean false value
// is provided as the ok named return value.
func First[T any](input []T) (result T, ok bool) {
	if len(input) == 0 {
		return
	}
	result = input[0]
	ok = true
	return
}

// Includes determines whether the input slice contains the specified value.  If it does, a truthy boolean is returned.
// Otherwise, a falsy boolean is returned.
func Includes[T comparable](input []T, value T) bool {
	for _, element := range input {
		if element == value {
			return true
		}
	}
	return false
}

// IndexOf returns the first index at which a given element can be found in the slice or -1 if it is not present.
func IndexOf[T comparable](input []T, value T) int {
	for idx, element := range input {
		if element == value {
			return idx
		}
	}
	return -1
}

// SubSlice provides a new slice containing the entries between the two indexes of the input slice (from is inclusive,
// to is exclusive).
func SubSlice[T any](input []T, fromIndex, toIndex int) []T {
	l := len(input)
	if l == 0 {
		return nil
	}
	// If the range is before the start of the slice, return nil...
	if fromIndex < 0 && toIndex < 0 {
		return nil
	}
	// If the range is after the end of the slice, return nil...
	if fromIndex > l && toIndex > l {
		return nil
	}
	// If the range is backward (i.e. fromIndex > toIndex), return nil...
	if fromIndex > toIndex {
		return nil
	}
	if toIndex >= l {
		toIndex = l
	}
	if fromIndex < 0 {
		fromIndex = 0
	}
	return input[fromIndex:toIndex]
}
