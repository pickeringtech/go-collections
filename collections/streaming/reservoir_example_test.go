package streaming_test

import (
	"fmt"
	"sort"

	"github.com/pickeringtech/go-collections/collections/streaming"
)

func ExampleNewReservoir() {
	// A fixed seed makes the sample reproducible; pass nil for the default.
	r := streaming.NewReservoir[int](3, streaming.NewRand(42))
	for v := 0; v < 100; v++ {
		r.Add(v)
	}

	sample := r.Result()
	sort.Ints(sample) // sample order is unspecified; sort for a stable print
	fmt.Println(sample)
	// Output: [9 15 91]
}

func ExampleNewWeightedReservoir() {
	// Heavier elements are proportionally more likely to be retained.
	r := streaming.NewWeightedReservoir[string](2, streaming.NewRand(7))
	r.Add("rare", 1)
	r.Add("common", 50)
	r.Add("frequent", 30)
	r.Add("scarce", 2)

	sample := r.Result()
	sort.Strings(sample)
	fmt.Println(sample)
	// Output: [common frequent]
}
