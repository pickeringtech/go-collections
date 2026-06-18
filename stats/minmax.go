package stats

import "github.com/pickeringtech/go-collections/constraints"

// MinMax returns the smallest and largest element of input together, found in a
// single pass — cheaper than calling a separate Min and Max. It works on any
// constraints.Ordered type (including strings), and the results are returned
// exact in T.
//
// The third return value is false — and both results the zero value of T — when
// input is empty. Like the standard library's ordering reductions, MinMax does
// not give NaN a defined position, so a floating-point input containing NaN
// yields an unspecified min/max rather than an error.
func MinMax[T constraints.Ordered](input []T) (min, max T, ok bool) {
	if len(input) == 0 {
		var zero T
		return zero, zero, false
	}
	min, max = input[0], input[0]
	for _, v := range input {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max, true
}

// ArgMin returns the index of the smallest element of input. When several
// elements share the minimum value the lowest such index is returned. It works
// on any constraints.Ordered type.
//
// The second return value is false — and the index 0 — when input is empty. As
// with MinMax, NaN has no defined ordering, so a NaN-contaminated float input
// yields an unspecified index.
func ArgMin[T constraints.Ordered](input []T) (int, bool) {
	if len(input) == 0 {
		return 0, false
	}
	best := 0
	for i, v := range input {
		if v < input[best] {
			best = i
		}
	}
	return best, true
}

// ArgMax returns the index of the largest element of input. When several
// elements share the maximum value the lowest such index is returned. It works
// on any constraints.Ordered type.
//
// The second return value is false — and the index 0 — when input is empty. As
// with MinMax, NaN has no defined ordering, so a NaN-contaminated float input
// yields an unspecified index.
func ArgMax[T constraints.Ordered](input []T) (int, bool) {
	if len(input) == 0 {
		return 0, false
	}
	best := 0
	for i, v := range input {
		if v > input[best] {
			best = i
		}
	}
	return best, true
}
