package slices

// FlatMapFunc maps a single input element to zero or more output elements. It is
// the per-element callback used by FlatMap; returning an empty or nil slice
// drops the element from the result.
type FlatMapFunc[I, O any] func(I) []O

// FlatMap applies fun to every element of the input, in order, and concatenates
// the slices it returns into a single flat slice. It is the natural choice when
// each input element expands into zero or more output elements (a Map followed
// by a flatten). If the input is empty or nil, or every call returns nothing,
// the output is an initialised, non-nil empty slice. The input is never mutated.
func FlatMap[I, O any](input []I, fun FlatMapFunc[I, O]) []O {
	output := []O{}
	for _, element := range input {
		output = append(output, fun(element)...)
	}
	return output
}
