package slices

// Concatenate joins two slices together, with inputA being joined with inputB following its last element.
func Concatenate[T any](inputA, inputB []T) []T {
	return append(inputA, inputB...)
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
