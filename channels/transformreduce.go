package channels

import "context"

// ReduceFunc is a function which takes an accumulator and an input element, and returns the new accumulator value.
type ReduceFunc[I, O any] func(accumulator O, element I) O

// Reduce reads all elements from the input channel and reduces them to a single value using the given ReduceFunc.
// The accumulated value is emitted on the output channel once the input channel is closed.
//
// The supplied context governs the lifetime of the reducing goroutine: when ctx is cancelled the goroutine
// abandons the partial accumulation, closes the output channel without emitting, and returns.
func Reduce[I, O any](ctx context.Context, input <-chan I, fn ReduceFunc[I, O]) <-chan O {
	output := make(chan O)
	go func() {
		defer close(output)
		var accumulator O
		for {
			select {
			case <-ctx.Done():
				return
			case element, ok := <-input:
				if !ok {
					send(ctx, output, accumulator)
					return
				}
				accumulator = fn(accumulator, element)
			}
		}
	}()
	return output
}
