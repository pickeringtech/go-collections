package stats

// Mode returns the value (or values) that occur most frequently in input. When
// several values tie for the highest frequency they are all returned, in order
// of first appearance in input; when every value is unique they are therefore
// all returned, each having frequency one. The returned slice is freshly
// allocated and the caller's slice is never mutated.
//
// Mode works on any comparable T, not just numbers. The second return value is
// false — and the result nil — only when input is empty.
//
// Caveat for floating-point input: NaN is never equal to itself, so multiple
// NaN elements are never aggregated into a single count and a NaN can never be
// reported as the mode in the way a caller might expect. Pre-filter NaN if that
// matters.
func Mode[T comparable](input []T) ([]T, bool) {
	if len(input) == 0 {
		return nil, false
	}
	counts := make(map[T]int, len(input))
	var order []T
	for _, v := range input {
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
