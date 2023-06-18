package constraints_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/constraints"
)

type complexValue[T constraints.Complex] struct {
	value T
}

func ExampleComplex() {
	a := complexValue[complex64]{0}
	b := complexValue[complex128]{0}
	fmt.Printf("%v, %v", a, b)
	// Output: {(0+0i)}, {(0+0i)}
}
