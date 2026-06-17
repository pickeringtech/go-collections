package deques

import "iter"

// RingBuffer is a ring-buffer-backed implementation of MutableDeque. It gives
// O(1) push and pop at both ends. It can be unbounded (growing on demand) or
// bounded to a fixed capacity with an OverflowPolicy. It is not safe for
// concurrent use; choose ConcurrentRingBuffer or ConcurrentRWRingBuffer for
// that.
//
// Elements are stored in buf with the front at index head and the remaining
// size-1 elements following at successive indices modulo len(buf), wrapping
// around the end. The logical element at index i lives at buf[(head+i)%len(buf)].
//
// Zero value: a RingBuffer{} is a valid, empty, unbounded deque ready for use;
// ensureCapacity grows the buffer from zero on the first push. NewRingBuffer (or
// NewBoundedRingBuffer) remains the recommended constructor, especially to seed
// elements or set a bounded capacity and overflow policy.
type RingBuffer[T any] struct {
	buf      []T
	head     int
	size     int
	bounded  bool
	capacity int
	policy   OverflowPolicy
}

// NewRingBuffer creates an unbounded RingBuffer seeded with the given elements,
// preserving their order (elements[0] becomes the front). The deque grows on
// demand as elements are added.
func NewRingBuffer[T any](elements ...T) *RingBuffer[T] {
	buf := make([]T, len(elements))
	copy(buf, elements)
	return &RingBuffer[T]{
		buf:      buf,
		size:     len(elements),
		capacity: Unbounded,
	}
}

// NewBoundedRingBuffer creates a bounded RingBuffer with the given capacity and
// overflow policy, seeded with the given elements (elements[0] becomes the
// front). A negative capacity is treated as zero. If more elements are supplied
// than the capacity holds, the policy decides which survive: OverwriteOldest
// keeps the last capacity elements, RejectWhenFull keeps the first capacity.
func NewBoundedRingBuffer[T any](capacity int, policy OverflowPolicy, elements ...T) *RingBuffer[T] {
	if capacity < 0 {
		capacity = 0
	}
	rb := &RingBuffer[T]{
		buf:      make([]T, capacity),
		bounded:  true,
		capacity: capacity,
		policy:   policy,
	}
	for _, element := range elements {
		rb.PushBackInPlace(element)
	}
	return rb
}

// Interface guards
var _ Deque[int] = &RingBuffer[int]{}
var _ MutableDeque[int] = &RingBuffer[int]{}

// physical returns the index in buf of the logical element at index i. It must
// only be called when size > 0 (so len(buf) > 0).
func (rb *RingBuffer[T]) physical(i int) int {
	return (rb.head + i) % len(rb.buf)
}

// ensureCapacity grows an unbounded buffer when it is full, re-laying the
// elements out contiguously from index 0. It is a no-op for bounded buffers,
// which are pre-allocated to their fixed capacity.
func (rb *RingBuffer[T]) ensureCapacity() {
	if rb.size < len(rb.buf) {
		return
	}
	newCap := len(rb.buf) * 2
	if newCap == 0 {
		newCap = 1
	}
	newBuf := make([]T, newCap)
	for i := 0; i < rb.size; i++ {
		newBuf[i] = rb.buf[rb.physical(i)]
	}
	rb.buf = newBuf
	rb.head = 0
}

// clone returns an independent deep copy of the receiver, sharing no backing
// storage. It underpins every immutable operation.
func (rb *RingBuffer[T]) clone() *RingBuffer[T] {
	newBuf := make([]T, len(rb.buf))
	copy(newBuf, rb.buf)
	return &RingBuffer[T]{
		buf:      newBuf,
		head:     rb.head,
		size:     rb.size,
		bounded:  rb.bounded,
		capacity: rb.capacity,
		policy:   rb.policy,
	}
}

// Length returns the number of elements currently held.
func (rb *RingBuffer[T]) Length() int {
	return rb.size
}

// IsEmpty returns true if the deque holds no elements.
func (rb *RingBuffer[T]) IsEmpty() bool {
	return rb.size == 0
}

// IsFull returns true if the deque is bounded and at capacity. Unbounded deques
// always return false.
func (rb *RingBuffer[T]) IsFull() bool {
	return rb.bounded && rb.size == rb.capacity
}

// Capacity returns the bounded deque's maximum size, or Unbounded (-1) for an
// unbounded deque.
func (rb *RingBuffer[T]) Capacity() int {
	if rb.bounded {
		return rb.capacity
	}
	return Unbounded
}

// PeekFront returns the front element and true, or the zero value and false if
// the deque is empty.
func (rb *RingBuffer[T]) PeekFront() (T, bool) {
	if rb.size == 0 {
		var zero T
		return zero, false
	}
	return rb.buf[rb.head], true
}

// PeekBack returns the back element and true, or the zero value and false if the
// deque is empty.
func (rb *RingBuffer[T]) PeekBack() (T, bool) {
	if rb.size == 0 {
		var zero T
		return zero, false
	}
	return rb.buf[rb.physical(rb.size-1)], true
}

// PushFrontInPlace adds element at the front, reporting whether it was accepted.
// See the MutableInsertable interface for the acceptance rules.
func (rb *RingBuffer[T]) PushFrontInPlace(element T) bool {
	if rb.bounded {
		if rb.capacity == 0 {
			return false
		}
		if rb.size == rb.capacity {
			if rb.policy == RejectWhenFull {
				return false
			}
			// OverwriteOldest: drop the back by claiming its slot for the new front.
			rb.head = (rb.head - 1 + len(rb.buf)) % len(rb.buf)
			rb.buf[rb.head] = element
			return true
		}
	} else {
		rb.ensureCapacity()
	}
	rb.head = (rb.head - 1 + len(rb.buf)) % len(rb.buf)
	rb.buf[rb.head] = element
	rb.size++
	return true
}

// PushBackInPlace adds element at the back, reporting whether it was accepted.
// See the MutableInsertable interface for the acceptance rules.
func (rb *RingBuffer[T]) PushBackInPlace(element T) bool {
	if rb.bounded {
		if rb.capacity == 0 {
			return false
		}
		if rb.size == rb.capacity {
			if rb.policy == RejectWhenFull {
				return false
			}
			// OverwriteOldest: when full the back slot coincides with the front, so
			// writing there and advancing head drops the old front.
			rb.buf[rb.head] = element
			rb.head = (rb.head + 1) % len(rb.buf)
			return true
		}
	} else {
		rb.ensureCapacity()
	}
	rb.buf[rb.physical(rb.size)] = element
	rb.size++
	return true
}

// PopFrontInPlace removes and returns the front element, reporting whether one
// was present, modifying the receiver.
func (rb *RingBuffer[T]) PopFrontInPlace() (T, bool) {
	var zero T
	if rb.size == 0 {
		return zero, false
	}
	element := rb.buf[rb.head]
	rb.buf[rb.head] = zero
	rb.head = (rb.head + 1) % len(rb.buf)
	rb.size--
	return element, true
}

// PopBackInPlace removes and returns the back element, reporting whether one was
// present, modifying the receiver.
func (rb *RingBuffer[T]) PopBackInPlace() (T, bool) {
	var zero T
	if rb.size == 0 {
		return zero, false
	}
	idx := rb.physical(rb.size - 1)
	element := rb.buf[idx]
	rb.buf[idx] = zero
	rb.size--
	return element, true
}

// PushFront returns a new deque with element added at the front, without
// modifying the receiver.
func (rb *RingBuffer[T]) PushFront(element T) Deque[T] {
	clone := rb.clone()
	clone.PushFrontInPlace(element)
	return clone
}

// PushBack returns a new deque with element added at the back, without modifying
// the receiver.
func (rb *RingBuffer[T]) PushBack(element T) Deque[T] {
	clone := rb.clone()
	clone.PushBackInPlace(element)
	return clone
}

// PopFront returns the front element, whether one was present, and a new deque
// with that element removed, without modifying the receiver.
func (rb *RingBuffer[T]) PopFront() (T, bool, Deque[T]) {
	clone := rb.clone()
	element, ok := clone.PopFrontInPlace()
	return element, ok, clone
}

// PopBack returns the back element, whether one was present, and a new deque
// with that element removed, without modifying the receiver.
func (rb *RingBuffer[T]) PopBack() (T, bool, Deque[T]) {
	clone := rb.clone()
	element, ok := clone.PopBackInPlace()
	return element, ok, clone
}

// Clear removes all elements from the deque, keeping a bounded deque's capacity.
func (rb *RingBuffer[T]) Clear() {
	if rb.bounded {
		rb.buf = make([]T, rb.capacity)
	} else {
		rb.buf = nil
	}
	rb.head = 0
	rb.size = 0
}

// AsSlice returns a new, non-nil slice holding the elements in front-to-back
// order, independent of the receiver's backing array.
func (rb *RingBuffer[T]) AsSlice() []T {
	out := make([]T, rb.size)
	for i := 0; i < rb.size; i++ {
		out[i] = rb.buf[rb.physical(i)]
	}
	return out
}

// ForEach calls fn once for each element, front to back.
func (rb *RingBuffer[T]) ForEach(fn EachFunc[T]) {
	for i := 0; i < rb.size; i++ {
		fn(rb.buf[rb.physical(i)])
	}
}

// ForEachWithIndex calls fn once for each element, front to back, passing the
// element's index and value.
func (rb *RingBuffer[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	for i := 0; i < rb.size; i++ {
		fn(i, rb.buf[rb.physical(i)])
	}
}

// All returns an iterator over index/value pairs, front to back.
func (rb *RingBuffer[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := 0; i < rb.size; i++ {
			if !yield(i, rb.buf[rb.physical(i)]) {
				return
			}
		}
	}
}

// Values returns an iterator over values, front to back.
func (rb *RingBuffer[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < rb.size; i++ {
			if !yield(rb.buf[rb.physical(i)]) {
				return
			}
		}
	}
}

// Backward returns an iterator over index/value pairs, back to front; indices
// still count from the front.
func (rb *RingBuffer[T]) Backward() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := rb.size - 1; i >= 0; i-- {
			if !yield(i, rb.buf[rb.physical(i)]) {
				return
			}
		}
	}
}
