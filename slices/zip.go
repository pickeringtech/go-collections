package slices

// Pair holds two values of independent types, produced by Zip when combining
// two slices element-wise. First comes from the first input, Second from the
// second.
type Pair[A, B any] struct {
	First  A
	Second B
}

// ZipFunc combines an element from each of two slices into a single result. It
// is the per-element callback used by ZipWith.
type ZipFunc[A, B, O any] func(A, B) O

// Zip combines two slices element-wise into a slice of Pairs: the element at
// index i of a is paired with the element at index i of b. When the inputs have
// unequal lengths the result is truncated to the shorter input, so no Pair ever
// holds a missing element. If either input is empty or nil the output is an
// initialised, non-nil empty slice. The inputs are never mutated.
func Zip[A, B any](a []A, b []B) []Pair[A, B] {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	output := make([]Pair[A, B], 0, n)
	for i := 0; i < n; i++ {
		output = append(output, Pair[A, B]{First: a[i], Second: b[i]})
	}
	return output
}

// ZipWith combines two slices element-wise using fun, which receives the
// elements at the same index from each input and returns the combined result.
// As with Zip, the result is truncated to the shorter input when the lengths
// differ, and an empty or nil input yields an initialised, non-nil empty slice.
// The inputs are never mutated.
func ZipWith[A, B, O any](a []A, b []B, fun ZipFunc[A, B, O]) []O {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	output := make([]O, 0, n)
	for i := 0; i < n; i++ {
		output = append(output, fun(a[i], b[i]))
	}
	return output
}
