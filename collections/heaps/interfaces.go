package heaps

import (
	"iter"

	"github.com/pickeringtech/go-collections/constraints"
)

// LessFunc reports whether a has strictly higher priority than b — that is,
// whether a should leave the heap before b. It matches the lists Sort
// convention of a func(T, T) bool less-than comparator.
//
// For a min-heap over an ordered type the comparator is `a < b`; for a max-heap
// it is `a > b`. See [Min] and [Max] for ready-made comparators.
type LessFunc[T any] func(a, b T) bool

// Min is a [LessFunc] producing a min-heap over any ordered type: the smallest
// element has the highest priority and leaves the heap first.
func Min[T constraints.Ordered](a, b T) bool {
	return a < b
}

// Max is a [LessFunc] producing a max-heap over any ordered type: the largest
// element has the highest priority and leaves the heap first.
func Max[T constraints.Ordered](a, b T) bool {
	return a > b
}

// Indexable provides priority-aware access and size reporting for a heap.
type Indexable[T any] interface {
	// Peek returns the highest-priority element without removing it.
	// Returns the element and true if the heap is non-empty; the zero value
	// and false otherwise.
	Peek() (T, bool)

	// Length returns the number of elements in the heap.
	Length() int

	// IsEmpty returns true if the heap contains no elements.
	IsEmpty() bool
}

// Iterable provides iteration over a heap in unspecified (heap-array) order.
type Iterable[T any] interface {
	// ForEach executes the given function for each element. The iteration
	// order is the internal heap-array order, which is unspecified beyond the
	// heap invariant; use AsSortedSlice or Drain for priority order.
	ForEach(fn func(element T))

	// All returns an iterator over the elements in unspecified (heap-array)
	// order, without modifying the heap.
	All() iter.Seq[T]
}

// Drainable provides priority-ordered iteration over a heap.
type Drainable[T any] interface {
	// Drain returns an iterator that yields the elements in priority order
	// (highest priority first), without modifying the heap. It operates on an
	// internal copy, so the receiver is left intact.
	Drain() iter.Seq[T]
}

// Convertible provides conversion of a heap to a slice.
type Convertible[T any] interface {
	// AsSlice returns a copy of the elements in unspecified (heap-array) order.
	AsSlice() []T

	// AsSortedSlice returns the elements in priority order (highest priority
	// first), without modifying the heap.
	AsSortedSlice() []T
}

// Pushable provides insertion that returns a new heap, leaving the receiver
// unchanged.
type Pushable[T any] interface {
	// Push returns a new heap with the given element added.
	Push(element T) Heap[T]

	// PushMany returns a new heap with all given elements added.
	PushMany(elements ...T) Heap[T]
}

// MutablePushable provides in-place insertion.
type MutablePushable[T any] interface {
	// PushInPlace adds the given element to the heap.
	PushInPlace(element T)

	// PushManyInPlace adds all given elements to the heap.
	PushManyInPlace(elements ...T)
}

// Poppable provides removal of the highest-priority element that returns a new
// heap, leaving the receiver unchanged.
type Poppable[T any] interface {
	// Pop removes the highest-priority element, returning it, whether the heap
	// was non-empty, and the resulting heap. When the heap is empty it returns
	// the zero value, false, and a heap equal to the receiver.
	Pop() (element T, ok bool, rest Heap[T])
}

// MutablePoppable provides in-place removal of the highest-priority element.
type MutablePoppable[T any] interface {
	// PopInPlace removes and returns the highest-priority element, reporting
	// whether the heap was non-empty.
	PopInPlace() (element T, ok bool)
}

// Heap is the immutable priority-queue interface: mutating operations return a
// new heap and never modify the receiver.
type Heap[T any] interface {
	Indexable[T]
	Iterable[T]
	Drainable[T]
	Convertible[T]
	Pushable[T]
	Poppable[T]
}

// MutableHeap extends Heap with in-place mutation operations.
type MutableHeap[T any] interface {
	Heap[T]
	MutablePushable[T]
	MutablePoppable[T]
}
