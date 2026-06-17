package channels

import "context"

// FilterFunc is a function which takes an input element and returns true if the element should be included in the output
// channel, or false if it should be excluded.
type FilterFunc[T any] func(element T) bool

// Filter reads all elements from the input channel and writes them to the output channel if the given FilterFunc returns
// true for that element.  This function reads until the input channel is closed.
//
// The supplied context governs the lifetime of the filtering goroutine: when ctx is cancelled the goroutine stops
// reading from input, closes the output channel, and returns, so an abandoned pipeline is reclaimed deterministically.
func Filter[T any](ctx context.Context, input <-chan T, fn FilterFunc[T]) <-chan T {
	output := make(chan T)
	go func() {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return
			case element, ok := <-input:
				if !ok {
					return
				}
				if fn(element) {
					if !send(ctx, output, element) {
						return
					}
				}
			}
		}
	}()
	return output
}
