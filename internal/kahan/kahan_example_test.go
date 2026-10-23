package kahan_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/internal/kahan"
)

func ExampleSum() {
	// A naive running total absorbs the 1 and 2 into 1e17 and loses them when
	// 1e17 is subtracted again; the compensated sum carries them through.
	var naive float64
	var total kahan.Sum
	for _, x := range []float64{1e17, 1, 2, -1e17} {
		naive += x
		total.Add(x)
	}
	fmt.Println(naive, total.Total())
	// Output: 0 3
}
