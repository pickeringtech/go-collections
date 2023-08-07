package slices

type FilterFunc[T any] func(T) bool

func Filter[T any](input []T, fn FilterFunc[T]) []T {
	var output []T
	for _, element := range input {
		if fn(element) {
			output = append(output, element)
		}
	}
	return output
}

func FilterInPlace[T any](input []T, fn FilterFunc[T]) {
	n := 0
	for _, x := range input {
		if fn(x) {
			input[n] = x
			n++
		}
	}
	input = input[:n]
}
