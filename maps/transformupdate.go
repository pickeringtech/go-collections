package maps

// Update copies the input map, and adds all the entries from the update map to it. If any keys in the update map
// collide with those in the input map, the value in the input map is overwritten with that from the update map.
func Update[K comparable, V any](input map[K]V, update map[K]V) map[K]V {
	result := Copy(input)
	for key, value := range update {
		result[key] = value
	}
	return result
}
