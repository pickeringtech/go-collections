package slices

import (
	"fmt"
	"strings"
)

// Concatenate joins two slices together, with inputA being joined with inputB following its last element.
func Concatenate[T any](inputA, inputB []T) []T {
	return append(inputA, inputB...)
}

// Copy duplicates the entries within the input into a new slice which is returned. If the input is empty or nil, the
// output is an initialised, non-nil empty slice.
func Copy[T any](input []T) []T {
	return append([]T{}, input...)
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
	input = Copy(input)
	return append(input[:index], input[index+1:]...)
}

// Fill sets every element in the input slice to the specified value, returning the resulting slice.
func Fill[T any](input []T, value T) []T {
	return FillFromTo[T](input, value, 0, len(input))
}

// FillFrom sets every element in the input slice after the specified index to the specified value, returning the resulting
// slice.
func FillFrom[T any](input []T, value T, fromIndex int) []T {
	return FillFromTo[T](input, value, fromIndex, len(input))
}

// FillFromTo sets every element in the input slice after the specified index and before an upper boundary index to the
// specified value.  The upper boundary is exclusive (i.e. will not be set to the new value - every element before it
// will).  The resulting slice is returned. If the input is empty or nil, the output is an initialised, non-nil empty
// slice.
func FillFromTo[T any](input []T, value T, fromIndex, toIndex int) []T {
	if len(input) == 0 {
		return []T{}
	}
	if fromIndex < 0 || toIndex > len(input) || fromIndex > toIndex {
		return input
	}
	inputCpy := Copy(input)
	for i := fromIndex; i < toIndex; i++ {
		inputCpy[i] = value
	}
	return inputCpy
}

// FillTo sets every element in the input slice up until the specified index (exclusive) to the specified value.  The
// resulting slice is returned.
func FillTo[T any](input []T, value T, toIndex int) []T {
	return FillFromTo[T](input, value, 0, toIndex)
}

// Insert adds the specified elements to the input slice at the specified index, returning the resulting slice. The
// index may range over 0 <= startIdx <= len(input): a startIdx equal to the length appends the elements to the end
// (so inserting into an empty slice yields just the elements). An out-of-range index (startIdx < 0 || startIdx >
// len(input)) leaves the input unchanged.
func Insert[T any](input []T, startIdx int, elements ...T) []T {
	if startIdx < 0 || startIdx > len(input) {
		return input
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
		sb.WriteString(fmt.Sprintf("%v", element))
		if idx != len(input)-1 {
			sb.WriteString(separator)
		}
	}
	return sb.String()
}

// Pop removes the last element from the input slice, returning it as well as the new, smaller slice.
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

	newSlice := input[:len(input)-1]
	return lastElement, true, newSlice
}

// PopFront removes the first element from the input slice, returning it as well as the new, smaller slice.
func PopFront[T any](input []T) (T, bool, []T) {
	var firstElement T
	if len(input) == 0 {
		return firstElement, false, nil
	}
	firstElement = input[0]
	if len(input) == 1 {
		return firstElement, true, nil
	}
	newSlice := input[1:]
	return firstElement, true, newSlice
}

// Push adds new elements to the end of the input slice.
func Push[T any](input []T, newElements ...T) []T {
	return append(input, newElements...)
}

// PushCopy returns a new slice containing the elements of input followed by
// newElements. Unlike Push, the result never aliases input's backing array
// (Push appends into input's spare capacity when present), and it is built with
// a single allocation sized to hold every element - avoiding the extra
// reallocation of copying and then pushing separately. The result is non-nil
// even when both input and newElements are empty.
func PushCopy[T any](input []T, newElements ...T) []T {
	out := make([]T, 0, len(input)+len(newElements))
	out = append(out, input...)
	return append(out, newElements...)
}

// PushFront adds the new elements to the front of the input slice.
func PushFront[T any](input []T, newElements ...T) []T {
	return append(newElements, input...)
}
