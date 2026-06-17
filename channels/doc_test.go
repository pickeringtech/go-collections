package channels_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/channels"
)

// Example_quickStart is the runnable twin of the package godoc Quick Start. Keep
// the two in sync: `go test` compiles and output-checks this, which is what
// guarantees the documented entry-point API actually exists and behaves as shown.
func Example_quickStart() {
	// Feed numbers into a channel.
	input := channels.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// Build a pipeline: square every number, then keep the even squares.
	pipeline := channels.NewPipeline[int, int](input, func(in <-chan int) <-chan int {
		squares := channels.Map(in, func(n int) int { return n * n })
		return channels.Filter(squares, func(n int) bool { return n%2 == 0 })
	})

	// CollectAsSlice drains the pipeline once the input channel is closed.
	results := pipeline.CollectAsSlice()

	fmt.Println(results)
	// Output: [4 16 36 64 100]
}
