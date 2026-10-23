package slices

import (
	"fmt"
	"strings"
)

// Concatenate joins two slices together, with inputA being joined with inputB following its last element. It returns a
// new slice without modifying or aliasing either input. If both inputs are empty, the result is a non-nil empty slice.
func Concatenate[T any](inputA, inputB []T) []T {
	output := make([]T, 0, len(inputA)+len(inputB))
	output = append(output, inputA...)
	return append(output, inputB...)
}

// Copy duplicates the entries within the input into a new slice which is returned. If the input is empty or nil, the
// output is an initialised, non-nil empty slice.
func Copy[T any](input []T) []T {
	return append([]T{}, input...)
}

// Delete returns a new slice holding the elements of the input except the one at the given index. The input is not
// modified, and the result never shares its backing array - even when the index is out of range (index < 0 || index >=
// len(input)), in which case the result is a copy of the input.
func Delete[T any](input []T, index int) []T {
	inputLen := len(input)
	if index >= inputLen {
		return Copy(input)
	}
	if index < 0 {
		return Copy(input)
	}
	input = Copy(input)
	return append(input[:index], input[index+1:]...)
}

// Fill returns a copy of the input slice with every element set to the specified value. The input is not modified.
func Fill[T any](input []T, value T) []T {
	return FillFromTo[T](input, value, 0, len(input))
}

// FillFrom returns a copy of the input slice with every element from the specified index (inclusive) to the end set to
// the specified value. The input is not modified; see FillFromTo for the out-of-range behaviour.
func FillFrom[T any](input []T, value T, fromIndex int) []T {
	return FillFromTo[T](input, value, fromIndex, len(input))
}

// FillFromTo returns a copy of the input slice with every element from fromIndex (inclusive) up to toIndex (exclusive)
// set to the specified value. The input is not modified, and the result never shares its backing array. An invalid
// range (fromIndex < 0, toIndex > len(input) or fromIndex > toIndex) sets nothing, so the result is an unchanged copy of
// the input. If the input is empty or nil, the output is an initialised, non-nil empty slice.
func FillFromTo[T any](input []T, value T, fromIndex, toIndex int) []T {
	if len(input) == 0 {
		return []T{}
	}
	if fromIndex < 0 || toIndex > len(input) || fromIndex > toIndex {
		return Copy(input)
	}
	inputCpy := Copy(input)
	for i := fromIndex; i < toIndex; i++ {
		inputCpy[i] = value
	}
	return inputCpy
}

// FillTo returns a copy of the input slice with every element up until the specified index (exclusive) set to the
// specified value. The input is not modified; see FillFromTo for the out-of-range behaviour.
func FillTo[T any](input []T, value T, toIndex int) []T {
	return FillFromTo[T](input, value, 0, toIndex)
}

// Insert returns a new slice holding the input with the specified elements inserted at the specified index. The
// index may range over 0 <= startIdx <= len(input): a startIdx equal to the length appends the elements to the end
// (so inserting into an empty slice yields just the elements). An out-of-range index (startIdx < 0 || startIdx >
// len(input)) inserts nothing, so the result is an unchanged copy of the input. Neither the input nor elements is
// modified, and the result never shares a backing array with either.
func Insert[T any](input []T, startIdx int, elements ...T) []T {
	if startIdx < 0 || startIdx > len(input) {
		return Copy(input)
	}
	// Build a fresh slice so neither the input nor the caller-provided elements
	// (whose backing array is shared via the variadic argument) are mutated.
	output := make([]T, 0, len(input)+len(elements))
	output = append(output, input[:startIdx]...)
	output = append(output, elements...)
	output = append(output, input[startIdx:]...)
	return output
}

// JoinToString creates a new string by stringifying each of the elements within the input, and placing the separator
// between them in the resulting string.
func JoinToString[T any](input []T, separator string) string {
	var sb strings.Builder
	for idx, element := range input {
		fmt.Fprintf(&sb, "%v", element)
		if idx != len(input)-1 {
			sb.WriteString(separator)
		}
	}
	return sb.String()
}

// Pop removes the last element from the input slice, returning it, whether there was one, and a new, smaller slice
// holding the remaining elements. The input is not modified, and the returned slice never shares its backing array, so
// writing to or appending to it cannot affect the input.
func Pop[T any](input []T) (T, bool, []T) {
	var lastElement T
	if len(input) == 0 {
		return lastElement, false, nil
	}

	lastElement = input[len(input)-1]

	// If there's only one element, return a nil slice.
	if len(input) == 1 {
		return lastElement, true, nil
	}

	return lastElement, true, Copy(input[:len(input)-1])
}

// PopFront removes the first element from the input slice, returning it, whether there was one, and a new, smaller
// slice holding the remaining elements. The input is not modified, and the returned slice never shares its backing
// array, so writing to or appending to it cannot affect the input.
func PopFront[T any](input []T) (T, bool, []T) {
	var firstElement T
	if len(input) == 0 {
		return firstElement, false, nil
	}
	firstElement = input[0]
	if len(input) == 1 {
		return firstElement, true, nil
	}
	return firstElement, true, Copy(input[1:])
}

// Push adds new elements to the end of the input slice in a new slice without modifying or aliasing the input. If the
// input and new elements are empty, it returns a non-nil empty slice.
func Push[T any](input []T, newElements ...T) []T {
	return PushCopy(input, newElements...)
}

// PushCopy returns a new slice containing the elements of input followed by newElements. Like Push, it never aliases
// the input's backing array and returns a non-nil empty slice when both inputs are empty. Its name emphasizes that the
// input is copied into the result.
func PushCopy[T any](input []T, newElements ...T) []T {
	out := make([]T, 0, len(input)+len(newElements))
	out = append(out, input...)
	return append(out, newElements...)
}

// PushFront adds new elements to the front of the input slice in a new slice without modifying or aliasing either input.
// If the input and new elements are empty, it returns a non-nil empty slice.
func PushFront[T any](input []T, newElements ...T) []T {
	output := make([]T, 0, len(input)+len(newElements))
	output = append(output, newElements...)
	return append(output, input...)
}
