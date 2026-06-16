package lists

// EachFunc is the callback invoked for each element during iteration.
type EachFunc[T any] func(element T)

// IndexedEachFunc is the callback invoked for each element during iteration,
// receiving both the element's index and its value.
type IndexedEachFunc[T any] func(idx int, element T)
