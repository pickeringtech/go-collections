package heaps

import (
	"iter"
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
	"github.com/pickeringtech/go-collections/constraints"
)

// ConcurrentRWBinary is a thread-safe binary-heap priority queue guarded by a
// read-write mutex. Read and immutable operations take a read lock, so many
// readers proceed concurrently; only the in-place mutators take the exclusive
// write lock. It suits read-heavy (Peek-heavy) workloads.
//
// Immutable operations (Push, PushMany, Pop) return another ConcurrentRWBinary,
// honouring the thread-safe-in / thread-safe-out contract.
//
// Zero value: always construct with NewConcurrentRW (or NewConcurrentRWMin /
// NewConcurrentRWMax). The embedded mutex is a value, so a bare
// &ConcurrentRWBinary{} is at least lock-safe, but its inner heap is nil until
// the constructor runs, so any operation — reads included — dereferences a nil
// pointer and panics.
//
// ConcurrentRWBinary must not be copied after first use; copying after construction
// produces an independent lock over shared backing data, which breaks the
// thread-safety contract. go vet reports any such copy.
type ConcurrentRWBinary[T any] struct {
	_     nocopy.NoCopy
	inner *Binary[T]
	lock  sync.RWMutex
}

// NewConcurrentRW creates a read-write-locked Binary heap ordered by the given
// comparator, seeded with the supplied values (heapified in O(n)).
func NewConcurrentRW[T any](less LessFunc[T], values ...T) *ConcurrentRWBinary[T] {
	return &ConcurrentRWBinary[T]{
		inner: New(less, values...),
	}
}

// NewConcurrentRWMin creates a read-write-locked min-heap over an ordered type.
func NewConcurrentRWMin[T constraints.Ordered](values ...T) *ConcurrentRWBinary[T] {
	return NewConcurrentRW(Min[T], values...)
}

// NewConcurrentRWMax creates a read-write-locked max-heap over an ordered type.
func NewConcurrentRWMax[T constraints.Ordered](values ...T) *ConcurrentRWBinary[T] {
	return NewConcurrentRW(Max[T], values...)
}

// Interface guards to ensure ConcurrentRWBinary implements the required interfaces.
var _ Heap[int] = &ConcurrentRWBinary[int]{}
var _ MutableHeap[int] = &ConcurrentRWBinary[int]{}

// wrap adapts a plain Binary into a thread-safe ConcurrentRWBinary, preserving
// the concurrent return contract.
func (c *ConcurrentRWBinary[T]) wrap(b *Binary[T]) Heap[T] {
	return &ConcurrentRWBinary[T]{inner: b}
}

// Peek returns the highest-priority element without removing it.
func (c *ConcurrentRWBinary[T]) Peek() (T, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.inner.Peek()
}

// Length returns the number of elements in the heap.
func (c *ConcurrentRWBinary[T]) Length() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.inner.Length()
}

// IsEmpty returns true if the heap contains no elements.
func (c *ConcurrentRWBinary[T]) IsEmpty() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.inner.IsEmpty()
}

// ForEach executes the given function for each element in heap-array order. fn
// is invoked after the lock is released, against a point-in-time snapshot taken
// under the lock, so fn may safely call back into the collection.
func (c *ConcurrentRWBinary[T]) ForEach(fn func(element T)) {
	c.lock.RLock()
	snapshot := c.inner.AsSlice()
	c.lock.RUnlock()

	for _, element := range snapshot {
		fn(element)
	}
}

// All returns an iterator over a snapshot of the elements in heap-array order.
func (c *ConcurrentRWBinary[T]) All() iter.Seq[T] {
	c.lock.RLock()
	defer c.lock.RUnlock()

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
func (c *ConcurrentRWBinary[T]) Drain() iter.Seq[T] {
	c.lock.RLock()
	defer c.lock.RUnlock()

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
func (c *ConcurrentRWBinary[T]) AsSlice() []T {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.inner.AsSlice()
}

// AsSortedSlice returns the elements in priority order, without modifying the
// receiver.
func (c *ConcurrentRWBinary[T]) AsSortedSlice() []T {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.inner.AsSortedSlice()
}

// Push returns a new thread-safe heap with the given element added.
func (c *ConcurrentRWBinary[T]) Push(element T) Heap[T] {
	c.lock.RLock()
	defer c.lock.RUnlock()

	clone := c.inner.clone()
	clone.PushInPlace(element)
	return c.wrap(clone)
}

// PushMany returns a new thread-safe heap with all given elements added.
func (c *ConcurrentRWBinary[T]) PushMany(elements ...T) Heap[T] {
	c.lock.RLock()
	defer c.lock.RUnlock()

	clone := c.inner.clone()
	clone.PushManyInPlace(elements...)
	return c.wrap(clone)
}

// Pop removes the highest-priority element, returning it, whether the heap was
// non-empty, and a new thread-safe heap. The receiver is left unchanged.
func (c *ConcurrentRWBinary[T]) Pop() (T, bool, Heap[T]) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	clone := c.inner.clone()
	element, ok := clone.PopInPlace()
	return element, ok, c.wrap(clone)
}

// PushInPlace adds the given element to the heap.
func (c *ConcurrentRWBinary[T]) PushInPlace(element T) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.inner.PushInPlace(element)
}

// PushManyInPlace adds all given elements to the heap.
func (c *ConcurrentRWBinary[T]) PushManyInPlace(elements ...T) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.inner.PushManyInPlace(elements...)
}

// PopInPlace removes and returns the highest-priority element, reporting
// whether the heap was non-empty.
func (c *ConcurrentRWBinary[T]) PopInPlace() (T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.inner.PopInPlace()
}
