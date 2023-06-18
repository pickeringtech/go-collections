package channels_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/channels"
	"strconv"
)

func ExamplePipeline_CollectAsSlice() {
	input := channels.FromSlice([]int{1, 2, 5, 4, 3})

	// Creates a new pipeline which totals and then stringifies the input channel.
	pipeline := channels.NewPipeline[int, string](input, func(input <-chan int) <-chan string {
		reducer := channels.Reduce(input, func(accumulator int, element int) int {
			return accumulator + element
		})

		stringifier := channels.Map[int, string](reducer, func(element int) string {
			return strconv.Itoa(element)
		})

		return stringifier
	})

	// Capture results in a slice.
	results := pipeline.CollectAsSlice()

	fmt.Printf("Results: %v", results)
	// Output: Results: [15]
}
