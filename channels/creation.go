package channels

import "github.com/pickeringtech/go-collectionutil/maps"

func FromSlice[T any](input []T) <-chan T {
	output := make(chan T)
	go func() {
		for _, el := range input {
			output <- el
		}
		close(output)
	}()
	return output
}

func FromMap[K comparable, V any](input map[K]V) <-chan maps.Entry[K, V] {
	output := make(chan maps.Entry[K, V])
	go func() {
		for key, val := range input {
			output <- maps.Entry[K, V]{
				Key:   key,
				Value: val,
			}
		}
		close(output)
	}()
	return output
}
