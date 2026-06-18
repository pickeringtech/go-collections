package stats

import "github.com/pickeringtech/go-collections/constraints"

// Clamp constrains value to the closed interval [lo, hi]: it returns lo when
// value is below lo, hi when value is above hi, and value otherwise. It works
// on any constraints.Ordered type (including strings), and returns its result
// directly since clamping is always well-defined for a finite interval.
//
// Callers must pass lo <= hi; a reversed interval (lo > hi) has no sensible
// clamp and the result is left unspecified. For floating-point value a NaN is
// returned unchanged, since NaN compares false against both bounds.
func Clamp[T constraints.Ordered](value, lo, hi T) T {
	if value < lo {
		return lo
	}
	if value > hi {
		return hi
	}
	return value
}

// ClampAll clamps every element of input to [lo, hi] using Clamp, returning a
// new slice of the same length; the caller's slice is never mutated and empty
// or nil input yields an empty slice. See Clamp for the lo <= hi requirement
// and the NaN behaviour.
func ClampAll[T constraints.Ordered](input []T, lo, hi T) []T {
	out := make([]T, len(input))
	for i, v := range input {
		out[i] = Clamp(v, lo, hi)
	}
	return out
}
