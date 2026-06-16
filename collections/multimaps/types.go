package multimaps

// Entry represents a single key-value association in a multimap.
//
// A multimap maps one key to many values, so the same key appears in multiple
// entries — one Entry per value bound to that key.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// EntryFunc is a function that operates on a single key-value entry.
type EntryFunc[K comparable, V any] func(key K, value V)

// KeyValuesFunc is a function that operates on a key together with all of the
// values currently bound to it.
type KeyValuesFunc[K comparable, V any] func(key K, values []V)

// KeyFunc is a function that operates on a key.
type KeyFunc[K comparable] func(key K)

// FilterFunc is a function that takes a key and value and returns true if the
// entry should be included in the result.
type FilterFunc[K comparable, V any] func(key K, value V) bool
