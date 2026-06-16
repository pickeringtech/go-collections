package deques

import (
	"iter"
	"sync"
)

// ConcurrentRingBuffer is a ring-buffer-backed implementation of MutableDeque
// that is safe for concurrent use. Every operation is guarded by a sync.Mutex.
// Operating on it yields another ConcurrentRingBuffer, so thread-safe in means
// thread-safe out.
type ConcurrentRingBuffer[T any] struct {
	inner *RingBuffer[T]
	lock  *sync.Mutex
}

// NewConcurrentRingBuffer creates an unbounded, thread-safe RingBuffer seeded
// with the given elements, preserving their order.
func NewConcurrentRingBuffer[T any](elements ...T) *ConcurrentRingBuffer[T] {
	return &ConcurrentRingBuffer[T]{
		inner: NewRingBuffer[T](elements...),
		lock:  &sync.Mutex{},
	}
}

// NewBoundedConcurrentRingBuffer creates a bounded, thread-safe RingBuffer with
// the given capacity and overflow policy, seeded with the given elements. See
// NewBoundedRingBuffer for how seed elements beyond the capacity are handled.
func NewBoundedConcurrentRingBuffer[T any](capacity int, policy OverflowPolicy, elements ...T) *ConcurrentRingBuffer[T] {
	return &ConcurrentRingBuffer[T]{
		inner: NewBoundedRingBuffer[T](capacity, policy, elements...),
		lock:  &sync.Mutex{},
	}
}

// Interface guards
var _ Deque[int] = &ConcurrentRingBuffer[int]{}
var _ MutableDeque[int] = &ConcurrentRingBuffer[int]{}

// wrap builds a new ConcurrentRingBuffer around an inner buffer with a fresh
// lock, so the result is independent of the receiver's lock.
func wrapConcurrent[T any](inner *RingBuffer[T]) *ConcurrentRingBuffer[T] {
	return &ConcurrentRingBuffer[T]{inner: inner, lock: &sync.Mutex{}}
}

// snapshot returns an independent front-to-back copy of the elements, taken
// under the lock, for safe lock-free iteration afterwards.
func (c *ConcurrentRingBuffer[T]) snapshot() []T {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.AsSlice()
}

// Length returns the number of elements currently held. It is safe for
// concurrent use.
func (c *ConcurrentRingBuffer[T]) Length() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Length()
}

// IsEmpty returns true if the deque holds no elements. It is safe for concurrent
// use.
func (c *ConcurrentRingBuffer[T]) IsEmpty() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.IsEmpty()
}

// IsFull returns true if the deque is bounded and at capacity. It is safe for
// concurrent use.
func (c *ConcurrentRingBuffer[T]) IsFull() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.IsFull()
}

// Capacity returns the bounded capacity, or Unbounded (-1). It is safe for
// concurrent use.
func (c *ConcurrentRingBuffer[T]) Capacity() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.Capacity()
}

// PeekFront returns the front element and true, or the zero value and false if
// empty. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) PeekFront() (T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.PeekFront()
}

// PeekBack returns the back element and true, or the zero value and false if
// empty. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) PeekBack() (T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.PeekBack()
}

// PushFrontInPlace adds element at the front, reporting acceptance, modifying the
// receiver. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) PushFrontInPlace(element T) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.PushFrontInPlace(element)
}

// PushBackInPlace adds element at the back, reporting acceptance, modifying the
// receiver. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) PushBackInPlace(element T) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.PushBackInPlace(element)
}

// PopFrontInPlace removes and returns the front element, modifying the receiver.
// It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) PopFrontInPlace() (T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.PopFrontInPlace()
}

// PopBackInPlace removes and returns the back element, modifying the receiver. It
// is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) PopBackInPlace() (T, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.PopBackInPlace()
}

// PushFront returns a new ConcurrentRingBuffer with element added at the front,
// without modifying the receiver. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) PushFront(element T) Deque[T] {
	c.lock.Lock()
	defer c.lock.Unlock()
	clone := c.inner.clone()
	clone.PushFrontInPlace(element)
	return wrapConcurrent(clone)
}

// PushBack returns a new ConcurrentRingBuffer with element added at the back,
// without modifying the receiver. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) PushBack(element T) Deque[T] {
	c.lock.Lock()
	defer c.lock.Unlock()
	clone := c.inner.clone()
	clone.PushBackInPlace(element)
	return wrapConcurrent(clone)
}

// PopFront returns the front element, whether one was present, and a new
// ConcurrentRingBuffer with that element removed, without modifying the receiver.
// It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) PopFront() (T, bool, Deque[T]) {
	c.lock.Lock()
	defer c.lock.Unlock()
	clone := c.inner.clone()
	element, ok := clone.PopFrontInPlace()
	return element, ok, wrapConcurrent(clone)
}

// PopBack returns the back element, whether one was present, and a new
// ConcurrentRingBuffer with that element removed, without modifying the receiver.
// It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) PopBack() (T, bool, Deque[T]) {
	c.lock.Lock()
	defer c.lock.Unlock()
	clone := c.inner.clone()
	element, ok := clone.PopBackInPlace()
	return element, ok, wrapConcurrent(clone)
}

// Clear removes all elements from the deque. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.Clear()
}

// AsSlice returns a new, non-nil front-to-back slice copy. It is safe for
// concurrent use.
func (c *ConcurrentRingBuffer[T]) AsSlice() []T {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.inner.AsSlice()
}

// ForEach calls fn once for each element, front to back, over a snapshot taken
// under the lock. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) ForEach(fn EachFunc[T]) {
	for _, element := range c.snapshot() {
		fn(element)
	}
}

// ForEachWithIndex calls fn once for each element, front to back, over a snapshot
// taken under the lock. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	for idx, element := range c.snapshot() {
		fn(idx, element)
	}
}

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) All() iter.Seq2[int, T] {
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
// under the lock. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) Values() iter.Seq[T] {
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
// snapshot taken under the lock. It is safe for concurrent use.
func (c *ConcurrentRingBuffer[T]) Backward() iter.Seq2[int, T] {
	snapshot := c.snapshot()
	return func(yield func(int, T) bool) {
		for i := len(snapshot) - 1; i >= 0; i-- {
			if !yield(i, snapshot[i]) {
				return
			}
		}
	}
}
