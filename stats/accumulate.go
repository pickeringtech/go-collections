package stats

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
)

// nonFinite reports whether x is NaN or ±Inf — values for which the means in
// this package are undefined and so are rejected.
func nonFinite(x float64) bool {
	return math.IsNaN(x) || math.IsInf(x, 0)
}

// maxAbs returns the largest absolute value in input, or 0 for empty input —
// the scale factor the overflow-safe computations divide out. A NaN anywhere
// makes the result NaN, and otherwise an infinity makes it +Inf, so a
// non-finite element always surfaces in the scale.
func maxAbs[T constraints.Numeric](input []T) float64 {
	var m float64
	for _, v := range input {
		a := math.Abs(float64(v))
		if a > m || a != a {
			m = a
		}
	}
	return m
}
