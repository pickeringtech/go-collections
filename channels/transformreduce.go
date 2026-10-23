package channels

import "context"

// ReduceFunc is a function which takes an accumulator and an input element, and returns the new accumulator value.
type ReduceFunc[I, O any] func(accumulator O, element I) O

// Reduce reads all elements from the input channel and reduces them to a single value using the given ReduceFunc.
// The accumulated value is emitted on the output channel once the input channel is closed.
//
// The supplied context governs the lifetime of the reducing goroutine: when ctx is cancelled the goroutine
// abandons the partial accumulation, closes the output channel without emitting, and returns. As in Map,
// cancellation takes priority over a ready input, so fn never runs on an element taken after cancellation.
func Reduce[I, O any](ctx context.Context, input <-chan I, fn ReduceFunc[I, O]) <-chan O {
	output := make(chan O)
	go func() {
		defer close(output)
		var accumulator O
		for {
			element, ok := receive(ctx, input)
			if !ok {
				// The input closed - or ctx was cancelled, in which case send
				// declines and the partial accumulation is abandoned unsent.
				send(ctx, output, accumulator)
				return
			}
			accumulator = fn(accumulator, element)
		}
	}()
	return output
}
