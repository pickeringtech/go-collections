package heaps

import (
	"iter"
	"sync"

	"github.com/pickeringtech/go-collections/constraints"
)

// ConcurrentBinary is a thread-safe binary-heap priority queue. Every operation
// is guarded by a single mutex. Immutable operations (Push, PushMany, Pop)
// return another ConcurrentBinary, so a thread-safe heap in yields a
// thread-safe heap out.
type ConcurrentBinary[T any] struct {
	inner *Binary[T]
	lock  sync.Mutex
}

// NewConcurrent creates a thread-safe Binary heap ordered by the given
// comparator, seeded with the supplied values (heapified in O(n)).
func NewConcurrent[T any](less LessFunc[T], values ...T) *ConcurrentBinary[T] {
	return &ConcurrentBinary[T]{
		inner: New(less, values...),
	}
}

// NewConcurrentMin creates a thread-safe min-heap over an ordered type.
func NewConcurrentMin[T constraints.Ordered](values ...T) *ConcurrentBinary[T] {
	return NewConcurrent(Min[T], values...)
}

// NewConcurrentMax creates a thread-safe max-heap over an ordered type.
func NewConcurrentMax[T constraints.Ordered](values ...T) *ConcurrentBinary[T] {
	return NewConcurrent(Max[T], values...)
}

// Interface guards to ensure ConcurrentBinary implements the required interfaces.
var _ Heap[int] = &ConcurrentBinary[int]{}
var _ MutableHeap[int] = &ConcurrentBinary[int]{}

// wrap adapts a plain Binary into a thread-safe ConcurrentBinary, preserving
// the concurrent return contract.
func (c *ConcurrentBinary[T]) wrap(b *Binary[T]) Heap[T] {
	return &ConcurrentBinary[T]{inner: b}
}

// Peek returns the highest-priority element without removing it.
func (c *ConcurrentBinary[T]) Peek() (T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.inner.Peek()
}

// Length returns the number of elements in the heap.
func (c *ConcurrentBinary[T]) Length() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.inner.Length()
}

// IsEmpty returns true if the heap contains no elements.
func (c *ConcurrentBinary[T]) IsEmpty() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.inner.IsEmpty()
}

// ForEach executes the given function for each element in heap-array order. fn
// is invoked after the lock is released, against a point-in-time snapshot taken
// under the lock, so fn may safely call back into the collection.
func (c *ConcurrentBinary[T]) ForEach(fn func(element T)) {
	c.lock.Lock()
	snapshot := c.inner.AsSlice()
	c.lock.Unlock()

	for _, element := range snapshot {
		fn(element)
	}
}

// All returns an iterator over a snapshot of the elements in heap-array order.
func (c *ConcurrentBinary[T]) All() iter.Seq[T] {
	c.lock.Lock()
	defer c.lock.Unlock()

	snapshot := c.inner.AsSlice()
	return func(yield func(T) bool) {
		for _, element := range snapshot {
			if !yield(element) {
				return
			}
		}
	}
}

// Drain returns an iterator over a snapshot of the elements in priority order,
// without modifying the receiver.
func (c *ConcurrentBinary[T]) Drain() iter.Seq[T] {
	c.lock.Lock()
	defer c.lock.Unlock()

	sorted := c.inner.AsSortedSlice()
	return func(yield func(T) bool) {
		for _, element := range sorted {
			if !yield(element) {
				return
			}
		}
	}
}

// AsSlice returns a copy of the elements in heap-array order.
func (c *ConcurrentBinary[T]) AsSlice() []T {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.inner.AsSlice()
}

// AsSortedSlice returns the elements in priority order, without modifying the
// receiver.
func (c *ConcurrentBinary[T]) AsSortedSlice() []T {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.inner.AsSortedSlice()
}

// Push returns a new thread-safe heap with the given element added.
func (c *ConcurrentBinary[T]) Push(element T) Heap[T] {
	c.lock.Lock()
	defer c.lock.Unlock()

	clone := c.inner.clone()
	clone.PushInPlace(element)
	return c.wrap(clone)
}

// PushMany returns a new thread-safe heap with all given elements added.
func (c *ConcurrentBinary[T]) PushMany(elements ...T) Heap[T] {
	c.lock.Lock()
	defer c.lock.Unlock()

	clone := c.inner.clone()
	clone.PushManyInPlace(elements...)
	return c.wrap(clone)
}

// Pop removes the highest-priority element, returning it, whether the heap was
// non-empty, and a new thread-safe heap. The receiver is left unchanged.
func (c *ConcurrentBinary[T]) Pop() (T, bool, Heap[T]) {
	c.lock.Lock()
	defer c.lock.Unlock()

	clone := c.inner.clone()
	element, ok := clone.PopInPlace()
	return element, ok, c.wrap(clone)
}

// PushInPlace adds the given element to the heap.
func (c *ConcurrentBinary[T]) PushInPlace(element T) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.inner.PushInPlace(element)
}

// PushManyInPlace adds all given elements to the heap.
func (c *ConcurrentBinary[T]) PushManyInPlace(elements ...T) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.inner.PushManyInPlace(elements...)
}

// PopInPlace removes and returns the highest-priority element, reporting
// whether the heap was non-empty.
func (c *ConcurrentBinary[T]) PopInPlace() (T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.inner.PopInPlace()
}
