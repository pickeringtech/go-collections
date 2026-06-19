package stats

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
)

// Dot computes the dot product of two numeric vectors a and b, returning the
// result as float64 together with an ok flag. The dot product is the sum of
// paired element-wise products: Σ aᵢ·bᵢ.
//
// It returns ok == false when the vectors are empty or have differing lengths —
// the dot product is undefined in those cases. Non-finite inputs (NaN/Inf)
// propagate to a non-finite result with ok == true, following the package's
// NaN/Inf policy for vector operations.
func Dot[T constraints.Numeric](a, b []T) (float64, bool) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	var sum float64
	for i := range a {
		sum += float64(a[i]) * float64(b[i])
	}
	return sum, true
}

// Norm computes the L2 (Euclidean) norm of a numeric vector, returning the
// result as float64 together with an ok flag. The L2 norm is the square root
// of the sum of squared elements: √(Σ xᵢ²).
//
// The implementation is numerically stable: it accumulates using math.Hypot
// in a chained fashion, avoiding catastrophic cancellation or overflow that
// would occur from naively squaring large values before summing.
//
// It returns ok == false for an empty input — the norm of an empty vector is
// undefined under the library's (result, ok) idiom. Non-finite inputs (NaN/Inf)
// propagate to a non-finite result with ok == true.
func Norm[T constraints.Numeric](input []T) (float64, bool) {
	if len(input) == 0 {
		return 0, false
	}
	var norm float64
	for _, v := range input {
		norm = math.Hypot(norm, float64(v))
	}
	return norm, true
}
