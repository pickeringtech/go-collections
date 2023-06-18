package constraints_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/constraints"
)

func ExampleOrdered() {
	type orderedValue[T constraints.Ordered] struct {
		value T
	}

	a := orderedValue[float32]{0}
	b := orderedValue[float64]{1}
	c := orderedValue[int]{2}
	d := orderedValue[int8]{3}
	e := orderedValue[int16]{4}
	f := orderedValue[int32]{5}
	g := orderedValue[int64]{6}
	h := orderedValue[uint]{7}
	i := orderedValue[uint8]{8}
	j := orderedValue[uint16]{9}
	k := orderedValue[uint32]{10}
	l := orderedValue[uint64]{11}
	m := orderedValue[uintptr]{12}
	o := orderedValue[string]{"thirteen"}

	values := []any{a, b, c, d, e, f, g, h, i, j, k, l, m, o}

	fmt.Printf("%v", values)
	// Output: [{0} {1} {2} {3} {4} {5} {6} {7} {8} {9} {10} {11} {12} {thirteen}]
}
