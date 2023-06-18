package channels

type ReduceFunc[I, O any] func(accumulator O, element I) O

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
