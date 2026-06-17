package constraints_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/constraints"
)

// sumOf and maxOf mirror the generic helpers shown in the package godoc Quick
// Start; defining them here is what proves the documented constraints actually
// constrain the type parameters as described.
func sumOf[T constraints.Numeric](numbers []T) T {
	var total T
	for _, n := range numbers {
		total += n
	}
	return total
}

func maxOf[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Example_quickStart is the runnable twin of the package godoc Quick Start. Keep
// the two in sync: `go test` compiles and output-checks this, which is what
// guarantees the documented constraints support the generic functions as shown.
func Example_quickStart() {
	intSum := sumOf([]int{1, 2, 3, 4, 5})

	fmt.Println(intSum, maxOf(10, 20), maxOf("apple", "banana"))
	// Output: 15 20 banana
}
