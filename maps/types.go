package maps

// Entry is a key-value pair.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}
