package slices

// Reverse reorders the input such that all the elements are reversed.  The last element becomes the first.  The second
// from last element becomes the second element, etc.
func Reverse[T any](input []T) []T {
	inputCpy := Copy(input)
	for left, right := 0, len(input)-1; left < right; left, right = left+1, right-1 {
		inputCpy[left], inputCpy[right] = inputCpy[right], inputCpy[left]
	}
	return inputCpy
}
