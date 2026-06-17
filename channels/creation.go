package channels

import (
	"context"

	"github.com/pickeringtech/go-collections/maps"
)

// FromSlice converts a slice into a channel, writing them to the channel one-by-one. The channel will be closed after
// all elements have been read.
//
// The supplied context governs the lifetime of the producing goroutine: when ctx is cancelled the goroutine stops
// sending, closes the output channel, and returns, so a partially consumed source does not leak.
func FromSlice[T any](ctx context.Context, input []T) <-chan T {
	output := make(chan T)
	go func() {
		defer close(output)
		for _, el := range input {
			if !send(ctx, output, el) {
				return
			}
		}
	}()
	return output
}

// FromMap converts a map into a channel, writing the entries to the channel one-by-one. The channel will be closed
// after all entries have been read.
//
// The supplied context governs the lifetime of the producing goroutine: when ctx is cancelled the goroutine stops
// sending, closes the output channel, and returns, so a partially consumed source does not leak.
func FromMap[K comparable, V any](ctx context.Context, input map[K]V) <-chan maps.Entry[K, V] {
	output := make(chan maps.Entry[K, V])
	go func() {
		defer close(output)
		for key, val := range input {
			entry := maps.Entry[K, V]{
				Key:   key,
				Value: val,
			}
			if !send(ctx, output, entry) {
				return
			}
		}
	}()
	return output
}
