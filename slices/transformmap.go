package slices

// MapFunc is a function which can be used to map a slice. It receives an element from the slice and returns the result
// of the mapping function.  The result of the mapping function is placed into the resulting slice.
type MapFunc[I, O any] func(I) O

// Map iterates over each element of the input, applying the provided mapping function, producing a new slice with the
// outputs of the mapping function.  If the input is empty or nil, the output will be nil.
func Map[I, O any](input []I, fun MapFunc[I, O]) []O {
	var output []O
	for _, element := range input {
		output = append(output, fun(element))
	}
	return output
}
