package slices

import (
	"fmt"
	"strings"
)

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
// will).  The resulting slice is returned.
func FillFromTo[T any](input []T, value T, fromIndex, toIndex int) []T {
	if len(input) == 0 {
		return nil
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

// Insert adds the specified elements to the input slice at the specified index, returning the resulting slice.
func Insert[T any](input []T, startIdx int, elements ...T) []T {
	if startIdx < 0 || startIdx >= len(input) {
		return nil
	}
	output := Copy(input)
	output = append(input[:startIdx], append(elements, input[startIdx:]...)...)
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
func Pop[T any](input []T) (lastElement T, newSlice []T) {
	if len(input) == 0 {
		return
	}
	lastElement = input[len(input)-1]
	newSlice = input[:len(input)-1]
	return
}

// PopFront removes the first element from the input slice, returning it as well as the new, smaller slice.
func PopFront[T any](input []T) (firstElement T, newSlice []T) {
	if len(input) == 0 {
		return
	}
	firstElement = input[0]
	newSlice = input[1:]
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
