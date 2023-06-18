package channels

import (
	"github.com/pickeringtech/go-collectionutil/maps"
)

func CollectAsSlice[T any](input <-chan T) []T {
	var results []T
	for el := range input {
		results = append(results, el)
	}
	return results
}

type MapBuilderFunc[I any, OK comparable, OV any] func(input I) maps.Entry[OK, OV]

func CollectAsMap[I any, OK comparable, OV any](input <-chan I, fn MapBuilderFunc[I, OK, OV]) map[OK]OV {
	results := map[OK]OV{}
	for el := range input {
		entry := fn(el)
		results[entry.Key] = entry.Value
	}
	return results
}
