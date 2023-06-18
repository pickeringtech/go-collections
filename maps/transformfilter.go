package maps

type FilterFunc[K comparable, V any] func(key K, value V) bool

// Filter applies the provided FilterFunc to each entry in the input map, including the entry in the result map if the
// FilterFunc returns true.
func Filter[K comparable, V any](input map[K]V, fn FilterFunc[K, V]) map[K]V {
	result := map[K]V{}
	for key, value := range input {
		if fn(key, value) {
			result[key] = value
		}
	}
	return result
}
