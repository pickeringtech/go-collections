package dicts

// EachFunc is a function that operates on a single element.
type EachFunc[T any] func(element T)

// EachFuncWithIndex is a function that operates on an element with its index.
type EachFuncWithIndex[T any] func(idx int, element T)

// Pair represents a key-value pair in a dictionary.
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// FilterFunc is a function that takes a key and value and returns true
// if the key-value pair should be included in the result.
type FilterFunc[K comparable, V any] func(key K, value V) bool

// MapFunc is a function that takes a key and value and returns a new key and value.
type MapFunc[K comparable, V any, OK comparable, OV any] func(key K, value V) (OK, OV)

// KeyFilterFunc is a function that takes a key and returns true
// if the key should be included in the result.
type KeyFilterFunc[K comparable] func(key K) bool

// ValueFilterFunc is a function that takes a value and returns true
// if the value should be included in the result.
type ValueFilterFunc[V any] func(value V) bool

// KeyValueFunc is a function that operates on a key-value pair.
type KeyValueFunc[K comparable, V any] func(key K, value V)

// KeyFunc is a function that operates on a key.
type KeyFunc[K comparable] func(key K)

// ValueFunc is a function that operates on a value.
type ValueFunc[V any] func(value V)
