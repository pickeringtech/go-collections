package slices_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/slices"
)

// Example_quickStart is the runnable twin of the package godoc Quick Start. Keep
// the two in sync: `go test` compiles and output-checks this, which is what
// guarantees the documented entry-point API actually exists and behaves as shown.
func Example_quickStart() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Each operation is a standalone function that takes a slice and returns a
	// new one, so compose them by nesting the calls.
	evens := slices.Filter(numbers, func(n int) bool { return n%2 == 0 })
	squares := slices.Map(evens, func(n int) int { return n * n })
	sum := slices.Reduce(squares, func(acc, n int) int { return acc + n })

	fmt.Println(sum)
	// Output: 220
}
