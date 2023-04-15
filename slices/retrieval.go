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
