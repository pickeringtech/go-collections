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
// the input channel is closed. If the channel closes without delivering any element, the output is an initialised,
// non-nil empty slice.
func CollectAsSlice[T any](input <-chan T) []T {
	results := []T{}
	for el := range input {
		results = append(results, el)
	}
	return results
}

// CollectNAsSlice reads at most howMany elements from the input channel and returns them as a slice, in the order they
// were received. It blocks until it has howMany elements or the input channel closes, whichever comes first; a howMany
// of zero or less reads nothing. If no element is collected, the output is an initialised, non-nil empty slice.
//
// CollectNAsSlice neither drains nor cancels the input once it stops reading, so any producer still writing to it - the
// FromSlice, Map or other stages upstream - stays blocked on its next send, and its goroutine leaks unless something
// releases it. Following Go's explicit-cancellation pipeline pattern, the caller owns that teardown: build the upstream
// stages with a cancellable context and cancel it once CollectNAsSlice returns.
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	doubled := channels.Map(ctx, channels.FromSlice(ctx, values), double)
//	firstThree := channels.CollectNAsSlice(doubled, 3)
//	cancel() // release the stages still blocked on the elements not collected
func CollectNAsSlice[T any](input <-chan T, howMany int) []T {
	results := []T{}

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
