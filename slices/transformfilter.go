package slices

type FilterFunc[T any] func(T) bool

func Filter[T any](input []T, fun FilterFunc[T]) []T {
	var output []T
	for _, element := range input {
		if fun(element) {
			output = append(output, element)
		}
	}
	return output
}
