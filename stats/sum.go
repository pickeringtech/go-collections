package stats

import "github.com/pickeringtech/go-collections/constraints"

// Sum adds up every element of input, returning the total together with an ok
// flag. Empty or nil input yields (zero value, false) — the empty sum is
// undefined under the library's (result, ok) idiom rather than a silent zero.
//
// Sum is an exact-in-T reduction: the total accumulates in T, so an integer sum
// is exact (no float round-off) but can overflow T's range on very large
// inputs; widen T (e.g. accumulate as int64/float64) if that is a concern. For
// the same reason Sum does not use Kahan compensation — that is a float64
// technique and would defeat integer exactness. Non-finite float inputs (NaN,
// ±Inf) propagate through the total per IEEE arithmetic; it is the float64
// summaries (Mean, …) that reject them, because there the statistic itself
// would be undefined.
func Sum[T constraints.Numeric](input []T) (T, bool) {
	if len(input) == 0 {
		var zero T
		return zero, false
	}
	var total T
	for _, element := range input {
		total += element
	}
	return total, true
}
