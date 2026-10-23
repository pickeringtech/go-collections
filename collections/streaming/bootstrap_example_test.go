package streaming_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/collections/streaming"
)

func ExampleBootstrap() {
	input := []int{1, 2, 3, 4, 5}
	// A fixed seed makes the resample reproducible.
	resample := streaming.Bootstrap(input, streaming.NewRand(42))

	fmt.Println(len(resample))
	fmt.Println(resample)
	// Output:
	// 5
	// [4 2 4 3 5]
}

func ExampleBootstrapN() {
	input := []float64{10, 20, 30}
	resamples := streaming.BootstrapN(input, 3, streaming.NewRand(7))

	for _, r := range resamples {
		fmt.Println(r)
	}
	// Output:
	// [30 10 30]
	// [10 10 10]
	// [30 10 20]
}
