package maps

// Clear removes every key-value pair from the input map, modifying the input map.
func Clear[K comparable, V any](input map[K]V) {
	for key := range input {
		delete(input, key)
	}
}

// ContainsValue searches through the input map for the given value. If the value is found, a truthy bool is returned.
// Otherwise, a falsy bool is returned.
func ContainsValue[K, V comparable](input map[K]V, value V) bool {
	for _, val := range input {
		if val == value {
			return true
		}
	}
	return false
}

// Copy creates a new map in memory which is identical to the input map.
func Copy[K comparable, V any](input map[K]V) map[K]V {
	newMap := map[K]V{}
	for key, val := range input {
		newMap[key] = val
	}
	return newMap
}

// GetMany attempts to find many entries in the map, returning their values in a slice. If a value does not exist, it
// is simply omitted from the output - no default value is inserted.
func GetMany[K comparable, V any](input map[K]V, keys []K) []V {
	var results []V
	for _, key := range keys {
		val, ok := input[key]
		if !ok {
			continue
		}
		results = append(results, val)
	}
	return results
}

// GetManyOrDefault attempts to find many entries in the map, returning their values in a slice. If a value does not
// exist, the default value specified is included in the output instead of the targeted value.
func GetManyOrDefault[K comparable, V any](input map[K]V, keys []K, defaultVal V) []V {
	var results []V
	for _, key := range keys {
		val, ok := input[key]
		if !ok {
			val = defaultVal
		}
		results = append(results, val)
	}
	return results
}

// GetOrDefault attempts to find the key within the input map. If the key can be found, its associated value is returned.
// if the key cannot be found, the 'orElse' value is returned.
func GetOrDefault[K comparable, V any](input map[K]V, key K, orElse V) V {
	val, ok := input[key]
	if !ok {
		return orElse
	}
	return val
}

// Items returns a slice of map entries representing each entry in the input map.
func Items[K comparable, V any](input map[K]V) []Entry[K, V] {
	var results []Entry[K, V]
	for key, val := range input {
		results = append(results, Entry[K, V]{
			Key:   key,
			Value: val,
		})
	}
	return results
}

// Keys provides a slice of all the keys of the input map.
func Keys[K comparable, V any](input map[K]V) []K {
	var results []K
	for key := range input {
		results = append(results, key)
	}
	return results
}

// Values returns a slice of all the values of the input map.
func Values[K comparable, V any](input map[K]V) []V {
	var results []V
	for _, val := range input {
		results = append(results, val)
	}
	return results
}
