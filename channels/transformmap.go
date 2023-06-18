package channels

type MapFunc[I, O any] func(I) O

// Map takes an input channel, transforms each of its entries using the MapFunc until the input channel is closed.
// The results are output to an output channel returned from this function.
func Map[I, O any](input <-chan I, fn MapFunc[I, O]) chan O {
	output := make(chan O)
	go func() {
		for val := range input {
			res := fn(val)
			output <- res
		}
		close(output)
	}()
	return output
}
