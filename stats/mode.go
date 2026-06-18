package stats

// Mode returns the value (or values) that occur most frequently in input. When
// several values tie for the highest frequency they are all returned, in order
// of first appearance in input; when every value is unique they are therefore
// all returned, each having frequency one. The returned slice is freshly
// allocated and the caller's slice is never mutated.
//
// Mode works on any comparable T, not just numbers. The second return value is
// false — and the result nil — when input is empty, or, for floating-point
// element types, when input contains a non-finite value (NaN or ±Inf). NaN is
// rejected because it never compares equal to itself, so it cannot be counted
// coherently; ±Inf is rejected alongside it to match the package's uniform
// non-finite policy for the ok-returning reductions.
func Mode[T comparable](input []T) ([]T, bool) {
	if len(input) == 0 {
		return nil, false
	}
	counts := make(map[T]int, len(input))
	var order []T
	for _, v := range input {
		if nonFiniteComparable(v) {
			return nil, false
		}
		if counts[v] == 0 {
			order = append(order, v)
		}
		counts[v]++
	}
	best := 0
	for _, v := range order {
		if counts[v] > best {
			best = counts[v]
		}
	}
	var modes []T
	for _, v := range order {
		if counts[v] == best {
			modes = append(modes, v)
		}
	}
	return modes, true
}

// nonFiniteComparable reports whether a comparable value is a non-finite
// floating-point number (NaN or ±Inf). Only float32/float64 dynamic types can
// be non-finite; every other comparable type (ints, strings, structs, …) is
// always finite, so the check is a cheap type assertion that returns false.
func nonFiniteComparable[T comparable](v T) bool {
	switch f := any(v).(type) {
	case float64:
		return nonFinite(f)
	case float32:
		return nonFinite(float64(f))
	default:
		return false
	}
}
