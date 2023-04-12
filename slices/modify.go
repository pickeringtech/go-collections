package slices

// Concatenate joins two slices together, with inputA being joined with inputB following its last element.
func Concatenate[T any](inputA, inputB []T) []T {
	return append(inputA, inputB...)
}

// Copy duplicates the entries within the input into a new slice which is returned.
func Copy[T any](input []T) []T {
	return append([]T(nil), input...)
}

// Delete removes the element at the given index from the provided input slice, returning the resulting slice.
func Delete[T any](input []T, index int) []T {
	inputLen := len(input)
	if index >= inputLen {
		return input
	}
	if index < 0 {
		return input
	}
	return append(input[:index], input[index+1:]...)
}

// Pop removes the last element from the input slice, returning it as well as the new, smaller slice.
func Pop[T any](input []T) (lastElement T, newSlice []T) {
	if len(input) == 0 {
		return
	}
	lastElement = input[len(input)-1]
	newSlice = input[:len(input)-1]
	return
}

// Push adds new elements to the end of the input slice.
func Push[T any](input []T, newElements ...T) []T {
	return append(input, newElements...)
}

// PushFront adds the new elements to the front of the input slice.
func PushFront[T any](input []T, newElements ...T) []T {
	return append(newElements, input...)
}
