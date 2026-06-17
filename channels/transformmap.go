package channels

import "context"

// MapFunc transforms a single input value of type I into an output value of type O.
type MapFunc[I, O any] func(I) O

// Map takes an input channel, transforms each of its entries using the MapFunc until the input channel is closed.
// The results are output to the receive-only channel returned from this function.
//
// The supplied context governs the lifetime of the transforming goroutine: when ctx is cancelled the goroutine
// stops reading from input, closes the output channel, and returns. This reclaims the goroutine deterministically
// rather than leaking it until the input drains, and unblocks it even when a downstream consumer has stalled.
func Map[I, O any](ctx context.Context, input <-chan I, fn MapFunc[I, O]) <-chan O {
	output := make(chan O)
	go func() {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-input:
				if !ok {
					return
				}
				if !send(ctx, output, fn(val)) {
					return
				}
			}
		}
	}()
	return output
}
