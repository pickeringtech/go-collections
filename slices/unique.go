package slices

// KeyFunc derives a comparable key from an element. UniqueBy uses it to decide
// which elements are considered duplicates.
type KeyFunc[T any, K comparable] func(T) K

// Unique returns a new slice containing the first occurrence of each distinct
// element, preserving their original order. Later duplicates are dropped. If the
// input is empty or nil the output is an initialised, non-nil empty slice. The
// input is never mutated.
func Unique[T comparable](input []T) []T {
	output := []T{}
	seen := make(map[T]struct{}, len(input))
	for _, element := range input {
		_, ok := seen[element]
		if ok {
			continue
		}
		seen[element] = struct{}{}
		output = append(output, element)
	}
	return output
}

// UniqueBy returns a new slice containing the first element seen for each
// distinct key produced by keyFn, preserving the original order. Two elements
// collide when keyFn maps them to the same key; the first one wins and later
// ones are dropped. This is the order-preserving dedup for element types that
// are not themselves comparable, or when uniqueness is defined by a field. If
// the input is empty or nil the output is an initialised, non-nil empty slice.
// The input is never mutated.
func UniqueBy[T any, K comparable](input []T, keyFn KeyFunc[T, K]) []T {
	output := []T{}
	seen := make(map[K]struct{}, len(input))
	for _, element := range input {
		key := keyFn(element)
		_, ok := seen[key]
		if ok {
			continue
		}
		seen[key] = struct{}{}
		output = append(output, element)
	}
	return output
}
