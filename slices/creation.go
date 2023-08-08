package slices

import "github.com/pickeringtech/go-collections/constraints"

// GeneratorFunc is a function which can be used to generate an element in a slice. It receives the index in the slice
// that the returned element will have when it is placed into the resulting slice.
type GeneratorFunc[T any] func(index int) T

// Generate simplifies creating a slice of any type using a generator function. The generator function is called n times
// and the result of that function is used as an element in the resulting slice. This can be really helpful when
// templating a slice of items, in which only small parts of the items differ.
func Generate[T any](n int, fn GeneratorFunc[T]) []T {
	var results []T
	for i := 0; i < n; i++ {
		results = append(results, fn(i))
	}
	return results
}

// NumericIdentityGenerator is a GeneratorFunc which returns the index as the element. This is useful when you want to
// create a slice of numbers, where the numbers are the index of the element in the slice.
func NumericIdentityGenerator[T constraints.Numeric](index int) T {
	return T(index)
}
