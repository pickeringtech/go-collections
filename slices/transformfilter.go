package slices

// FilterFunc is a function which can be used to filter a slice. It receives an element from the slice and returns true
// if that element should be included in the resulting slice.
type FilterFunc[T any] func(T) bool

// Filter returns a new slice containing only the elements of the input slice for which the provided function returns
// true. If the input is empty or nil, the output is an initialised, non-nil empty slice.
func Filter[T any](input []T, fn FilterFunc[T]) []T {
	output := []T{}
	for _, element := range input {
		if fn(element) {
			output = append(output, element)
		}
	}
	return output
}
