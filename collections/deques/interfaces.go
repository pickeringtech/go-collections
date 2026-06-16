package deques

import "iter"

// Indexable is implemented by deques that can report their size and capacity.
type Indexable[T any] interface {
	// Length returns the number of elements currently held.
	Length() int
	// IsEmpty returns true if the deque holds no elements.
	IsEmpty() bool
	// IsFull returns true if the deque is bounded and at capacity. Unbounded
	// deques always return false.
	IsFull() bool
	// Capacity returns the maximum number of elements a bounded deque can hold,
	// or Unbounded (-1) for an unbounded deque.
	Capacity() int
}

// Peekable is implemented by deques that can report the elements at their ends
// without removing them.
type Peekable[T any] interface {
	// PeekFront returns the front element and true, or the zero value and false
	// if the deque is empty.
	PeekFront() (T, bool)
	// PeekBack returns the back element and true, or the zero value and false if
	// the deque is empty.
	PeekBack() (T, bool)
}

// Insertable is implemented by deques that can produce a new deque with an
// element added at either end, without modifying the receiver.
type Insertable[T any] interface {
	// PushFront returns a new deque with element added at the front, without
	// modifying the receiver. For a full bounded deque the result depends on the
	// overflow policy: OverwriteOldest drops the back element, RejectWhenFull
	// returns an unchanged copy.
	PushFront(element T) Deque[T]
	// PushBack returns a new deque with element added at the back, without
	// modifying the receiver. For a full bounded deque the result depends on the
	// overflow policy: OverwriteOldest drops the front element, RejectWhenFull
	// returns an unchanged copy.
	PushBack(element T) Deque[T]
}

// MutableInsertable is implemented by deques that can add an element at either
// end in place, modifying the receiver.
type MutableInsertable[T any] interface {
	// PushFrontInPlace adds element at the front, reporting whether it was
	// accepted. It is always accepted for unbounded deques and for bounded deques
	// using OverwriteOldest; it is rejected (false, receiver unchanged) only when
	// a bounded RejectWhenFull deque is full.
	PushFrontInPlace(element T) bool
	// PushBackInPlace adds element at the back, reporting whether it was accepted.
	// See PushFrontInPlace for the acceptance rules.
	PushBackInPlace(element T) bool
}

// Removable is implemented by deques that can produce a new deque with an
// element removed from either end, without modifying the receiver.
type Removable[T any] interface {
	// PopFront returns the front element, whether one was present, and a new deque
	// with that element removed, without modifying the receiver.
	PopFront() (T, bool, Deque[T])
	// PopBack returns the back element, whether one was present, and a new deque
	// with that element removed, without modifying the receiver.
	PopBack() (T, bool, Deque[T])
}

// MutableRemovable is implemented by deques that can remove an element from
// either end in place, modifying the receiver.
type MutableRemovable[T any] interface {
	// PopFrontInPlace removes and returns the front element, reporting whether one
	// was present, modifying the receiver.
	PopFrontInPlace() (T, bool)
	// PopBackInPlace removes and returns the back element, reporting whether one
	// was present, modifying the receiver.
	PopBackInPlace() (T, bool)
	// Clear removes all elements from the deque.
	Clear()
}

// Convertible is implemented by deques that can be converted into a slice
// ordered front to back.
type Convertible[T any] interface {
	// AsSlice returns a new slice holding the elements in front-to-back order. It
	// is always non-nil, even when the deque is empty.
	AsSlice() []T
}

// Iterable is implemented by deques that can be iterated front to back, with or
// without element indices, via callbacks or range-over-func iterators.
type Iterable[T any] interface {
	// ForEach calls fn once for each element, front to back.
	ForEach(fn EachFunc[T])
	// ForEachWithIndex calls fn once for each element, front to back, passing the
	// element's index and value.
	ForEachWithIndex(fn IndexedEachFunc[T])
	// All returns an iterator over index/value pairs, front to back.
	All() iter.Seq2[int, T]
	// Values returns an iterator over values, front to back.
	Values() iter.Seq[T]
	// Backward returns an iterator over index/value pairs, back to front; indices
	// still count from the front (the back element has index Length()-1).
	Backward() iter.Seq2[int, T]
}

// Deque is the read-oriented double-ended queue interface, combining size and
// capacity reporting, peeking, immutable end insertion and removal, iteration,
// and conversion to a slice. Implementations may be unbounded (growing on
// demand) or bounded with a fixed capacity and an OverflowPolicy.
type Deque[T any] interface {
	Indexable[T]
	Peekable[T]
	Insertable[T]
	Removable[T]
	Convertible[T]
	Iterable[T]
}

// MutableDeque extends Deque with the in-place mutation operations for adding
// and removing elements at either end.
type MutableDeque[T any] interface {
	Deque[T]
	MutableInsertable[T]
	MutableRemovable[T]
}
