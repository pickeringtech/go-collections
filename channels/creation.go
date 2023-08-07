package channels

import "github.com/pickeringtech/go-collectionutil/maps"

// FromSlice converts a slice into a channel, writing them to the channel one-by-one. The channel will be closed after
// all elements have been read.
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

// FromMap converts a map into a channel, writing the entries to the channel one-by-one. The channel will be closed
// after all entries have been read.
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
