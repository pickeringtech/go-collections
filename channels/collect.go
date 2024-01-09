package channels

import (
	"github.com/pickeringtech/go-collections/maps"
)

// BuildMapFromEntries takes a slice of maps.Entry and returns a map built from those entries.
func BuildMapFromEntries[K comparable, V any](entries []maps.Entry[K, V]) map[K]V {
	results := map[K]V{}
	for _, entry := range entries {
		results[entry.Key] = entry.Value
	}
	return results
}

// CollectAsSlice reads all elements from the input channel and returns them as a slice. This function will block until
// the input channel is closed.
func CollectAsSlice[T any](input <-chan T) []T {
	var results []T
	for el := range input {
		results = append(results, el)
	}
	return results
}

// CollectNAsSlice reads all elements from the input channel and returns them as a slice. This function will block until
// the input channel is closed.
func CollectNAsSlice[T any](input <-chan T, howMany int) []T {
	var results []T

	for i := 0; i < howMany; i++ {
		el, ok := <-input
		if !ok {
			break
		}
		results = append(results, el)
	}
	return results
}

// MapBuilderFunc is a function which takes an input element and returns a maps.Entry, which is used to build a map.
type MapBuilderFunc[I any, OK comparable, OV any] func(input I) maps.Entry[OK, OV]

// CollectAsMap reads all elements from the input channel and returns them as a map. This function will block until the
// input channel is closed.
func CollectAsMap[I any, OK comparable, OV any](input <-chan I, fn MapBuilderFunc[I, OK, OV]) map[OK]OV {
	results := map[OK]OV{}

	for el := range input {
		entry := fn(el)
		results[entry.Key] = entry.Value
	}
	return results
}
