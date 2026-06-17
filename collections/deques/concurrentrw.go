package deques

import (
	"iter"
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// ConcurrentRWRingBuffer is a ring-buffer-backed implementation of MutableDeque
// that is safe for concurrent use, guarded by a sync.RWMutex. Read-only
// operations take a read lock so they can proceed concurrently; mutations take
// the full lock. Favour it over ConcurrentRingBuffer for read-heavy workloads.
//
// Zero value: always construct with NewConcurrentRWRingBuffer (or
// NewBoundedConcurrentRWRingBuffer). Unlike the plain RingBuffer, whose zero
// value is a valid empty deque, the concurrent wrapper holds its RingBuffer by
// pointer. The embedded mutex is a value, so a bare &ConcurrentRWRingBuffer{} is
// at least lock-safe, but its inner buffer is nil until the constructor runs, so
// any operation — reads included — dereferences a nil pointer and panics.
//
// ConcurrentRWRingBuffer must not be copied after first use; copying after construction
// produces an independent lock over shared backing data, which breaks the
// thread-safety contract. go vet reports any such copy.
type ConcurrentRWRingBuffer[T any] struct {
	_     nocopy.NoCopy
	inner *RingBuffer[T]
	lock  sync.RWMutex
}

// NewConcurrentRWRingBuffer creates an unbounded, thread-safe RingBuffer
// optimised for concurrent reads, seeded with the given elements.
func NewConcurrentRWRingBuffer[T any](elements ...T) *ConcurrentRWRingBuffer[T] {
	return &ConcurrentRWRingBuffer[T]{
		inner: NewRingBuffer[T](elements...),
	}
}

// NewBoundedConcurrentRWRingBuffer creates a bounded, thread-safe RingBuffer
// optimised for concurrent reads, with the given capacity and overflow policy,
// seeded with the given elements. See NewBoundedRingBuffer for how seed elements
// beyond the capacity are handled.
func NewBoundedConcurrentRWRingBuffer[T any](capacity int, policy OverflowPolicy, elements ...T) *ConcurrentRWRingBuffer[T] {
	return &ConcurrentRWRingBuffer[T]{
		inner: NewBoundedRingBuffer[T](capacity, policy, elements...),
	}
}

// Interface guards
var _ Deque[int] = &ConcurrentRWRingBuffer[int]{}
var _ MutableDeque[int] = &ConcurrentRWRingBuffer[int]{}

// wrapConcurrentRW builds a new ConcurrentRWRingBuffer around an inner buffer
// with a fresh lock, so the result is independent of the receiver's lock.
func wrapConcurrentRW[T any](inner *RingBuffer[T]) *ConcurrentRWRingBuffer[T] {
	return &ConcurrentRWRingBuffer[T]{inner: inner}
}

// snapshot returns an independent front-to-back copy of the elements, taken
// under a read lock, for safe lock-free iteration afterwards.
func (c *ConcurrentRWRingBuffer[T]) snapshot() []T {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.AsSlice()
}

// Length returns the number of elements currently held. It is safe for
// concurrent use.
func (c *ConcurrentRWRingBuffer[T]) Length() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.Length()
}

// IsEmpty returns true if the deque holds no elements. It is safe for concurrent
// use.
func (c *ConcurrentRWRingBuffer[T]) IsEmpty() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.IsEmpty()
}

// IsFull returns true if the deque is bounded and at capacity. It is safe for
// concurrent use.
func (c *ConcurrentRWRingBuffer[T]) IsFull() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.IsFull()
}

// Capacity returns the bounded capacity, or Unbounded (-1). It is safe for
// concurrent use.
func (c *ConcurrentRWRingBuffer[T]) Capacity() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.Capacity()
}

// PeekFront returns the front element and true, or the zero value and false if
// empty. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) PeekFront() (T, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.PeekFront()
}

// PeekBack returns the back element and true, or the zero value and false if
// empty. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) PeekBack() (T, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.PeekBack()
}

// PushFrontInPlace adds element at the front, reporting acceptance, modifying the
// receiver. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) PushFrontInPlace(element T) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.PushFrontInPlace(element)
}

// PushBackInPlace adds element at the back, reporting acceptance, modifying the
// receiver. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) PushBackInPlace(element T) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.PushBackInPlace(element)
}

// PopFrontInPlace removes and returns the front element, modifying the receiver.
// It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) PopFrontInPlace() (T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.PopFrontInPlace()
}

// PopBackInPlace removes and returns the back element, modifying the receiver. It
// is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) PopBackInPlace() (T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.PopBackInPlace()
}

// PushFront returns a new ConcurrentRWRingBuffer with element added at the front,
// without modifying the receiver. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) PushFront(element T) Deque[T] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	clone := c.inner.clone()
	clone.PushFrontInPlace(element)
	return wrapConcurrentRW(clone)
}

// PushBack returns a new ConcurrentRWRingBuffer with element added at the back,
// without modifying the receiver. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) PushBack(element T) Deque[T] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	clone := c.inner.clone()
	clone.PushBackInPlace(element)
	return wrapConcurrentRW(clone)
}

// PopFront returns the front element, whether one was present, and a new
// ConcurrentRWRingBuffer with that element removed, without modifying the
// receiver. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) PopFront() (T, bool, Deque[T]) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	clone := c.inner.clone()
	element, ok := clone.PopFrontInPlace()
	return element, ok, wrapConcurrentRW(clone)
}

// PopBack returns the back element, whether one was present, and a new
// ConcurrentRWRingBuffer with that element removed, without modifying the
// receiver. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) PopBack() (T, bool, Deque[T]) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	clone := c.inner.clone()
	element, ok := clone.PopBackInPlace()
	return element, ok, wrapConcurrentRW(clone)
}

// Clear removes all elements from the deque. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.Clear()
}

// AsSlice returns a new, non-nil front-to-back slice copy. It is safe for
// concurrent use.
func (c *ConcurrentRWRingBuffer[T]) AsSlice() []T {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.inner.AsSlice()
}

// ForEach calls fn once for each element, front to back, over a snapshot taken
// under a read lock. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) ForEach(fn EachFunc[T]) {
	for _, element := range c.snapshot() {
		fn(element)
	}
}

// ForEachWithIndex calls fn once for each element, front to back, over a snapshot
// taken under a read lock. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	for idx, element := range c.snapshot() {
		fn(idx, element)
	}
}

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under a read lock. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) All() iter.Seq2[int, T] {
	snapshot := c.snapshot()
	return func(yield func(int, T) bool) {
		for i, element := range snapshot {
			if !yield(i, element) {
				return
			}
		}
	}
}

// Values returns an iterator over values, front to back, over a snapshot taken
// under a read lock. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) Values() iter.Seq[T] {
	snapshot := c.snapshot()
	return func(yield func(T) bool) {
		for _, element := range snapshot {
			if !yield(element) {
				return
			}
		}
	}
}

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under a read lock. It is safe for concurrent use.
func (c *ConcurrentRWRingBuffer[T]) Backward() iter.Seq2[int, T] {
	snapshot := c.snapshot()
	return func(yield func(int, T) bool) {
		for i := len(snapshot) - 1; i >= 0; i-- {
			if !yield(i, snapshot[i]) {
				return
			}
		}
	}
}
