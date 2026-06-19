package streaming_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/collections/heaps"
	"github.com/pickeringtech/go-collections/collections/streaming"
)

func ExampleNewTopKOrdered() {
	top := streaming.NewTopKOrdered[int](3)
	for _, v := range []int{5, 1, 9, 3, 7, 2, 8} {
		top.Add(v)
	}

	fmt.Println(top.Result())
	// Output: [9 8 7]
}

func ExampleNewTopK() {
	type job struct {
		name     string
		priority int
	}
	// less(a, b) reports a ranks below b, so the highest-priority jobs are kept.
	top := streaming.NewTopK[job](2, heaps.LessFunc[job](func(a, b job) bool {
		return a.priority < b.priority
	}))
	top.Add(job{"email", 1})
	top.Add(job{"deploy", 9})
	top.Add(job{"lunch", 0})
	top.Add(job{"incident", 7})

	for _, j := range top.Result() {
		fmt.Printf("%s (%d)\n", j.name, j.priority)
	}
	// Output:
	// deploy (9)
	// incident (7)
}
