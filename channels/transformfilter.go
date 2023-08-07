package channels

// FilterFunc is a function which takes an input element and returns true if the element should be included in the output
// channel, or false if it should be excluded.
type FilterFunc[T any] func(element T) bool

// Filter reads all elements from the input channel and writes them to the output channel if the given FilterFunc returns
// true for that element.  This function will block until the input channel is closed.
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
