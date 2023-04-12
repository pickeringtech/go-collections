package slices

func FindFirst[T any](input []T) (result T, ok bool) {
	if len(input) == 0 {
		return
	}
	result = input[0]
	ok = true
	return
}

type FindFunc[T any] func(T) bool

// FindAny tests each element of the input with the provided function.  If the function returns true, the selected element
// is returned, along with a boolean truthy value.
func FindAny[T any](input []T, fun FindFunc[T]) (result T, ok bool) {
	for _, element := range input {
		if fun(element) {
			return element, true
		}
	}
	return
}

// AnyMatch tests each element of the input with the provided function.  If any of the elements, when passed through the
// function result in a truthy boolean value, a match is found and true is returned from this function.  Otherwise false
// is returned.
func AnyMatch[T any](input []T, fun FindFunc[T]) bool {
	for _, element := range input {
		if fun(element) {
			return true
		}
	}
	return false
}
