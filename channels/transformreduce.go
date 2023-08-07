package channels

// ReduceFunc is a function which takes an accumulator and an input element, and returns the new accumulator value.
type ReduceFunc[I, O any] func(accumulator O, element I) O

// Reduce reads all elements from the input channel and reduces them to a single value using the given ReduceFunc.
func Reduce[I, O any](input <-chan I, fn ReduceFunc[I, O]) <-chan O {
	output := make(chan O)
	go func() {
		var accumulator O
		for element := range input {
			accumulator = fn(accumulator, element)
		}
		output <- accumulator
		close(output)
	}()
	return output
}
