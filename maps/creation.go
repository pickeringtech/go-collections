package maps

// FromKeys constructs a new map with each of the input keys and sets a default value for each of them.
func FromKeys[K comparable, V any](keys []K, defaultVal V) map[K]V {
	result := map[K]V{}
	for _, k := range keys {
		result[k] = defaultVal
	}
	return result
}
