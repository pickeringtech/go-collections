package maps

// MapFunc is a function that takes a key and value and returns a new key and value.
type MapFunc[K comparable, V any, OK comparable, OV any] func(key K, value V) (OK, OV)

// Map takes each entry in the input map, transforming them using the provided mapping function, building a new map to
// output. It does not modify the input map, rather creating a new map which is returned. An entry's key and value can
// both be changed - transformation is not limited to the value.
func Map[K comparable, V any, OK comparable, OV any](input map[K]V, fn MapFunc[K, V, OK, OV]) map[OK]OV {
	results := map[OK]OV{}
	for key, value := range input {
		outputKey, outputVal := fn(key, value)
		results[outputKey] = outputVal
	}
	return results
}
