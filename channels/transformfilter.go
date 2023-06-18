package channels

type FilterFunc[T any] func(element T) bool

func Filter[T any](input <-chan T, fn FilterFunc[T]) <-chan T {
	output := make(chan T)
	go func() {
		for element := range input {
			if fn(element) {
				output <- element
			}
		}
		close(output)
	}()
	return output
}
