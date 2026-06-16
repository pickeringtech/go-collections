package heaps

import (
	"iter"

	"github.com/pickeringtech/go-collections/constraints"
)

// Binary is a generic binary-heap priority queue backed by a slice. The element
// ordering is governed by a LessFunc comparator, so the heap works for any type
// — not just ordered primitives — and can be configured as a min-heap, a
// max-heap, or anything in between (e.g. ordering tasks by a priority field).
//
// Push and PopInPlace run in O(log n); Peek, Length and IsEmpty in O(1).
// Building from a slice via New is O(n) thanks to bottom-up heapify.
//
// Binary is single-threaded; use ConcurrentBinary or ConcurrentRWBinary for
// concurrent access.
//
// Example usage:
//
//	// A min-heap of ints.
//	pq := heaps.NewMin(5, 1, 3, 2, 4)
//	top, _ := pq.Peek()   // 1
//
//	// A max-heap, or a custom comparator over a struct.
//	tasks := heaps.New(func(a, b Task) bool { return a.Priority > b.Priority })
type Binary[T any] struct {
	data []T
	less LessFunc[T]
}

// New creates a Binary heap ordered by the given comparator, seeded with the
// supplied values. The values are heapified in O(n).
//
// Example:
//
//	// Order tasks so the highest Priority leaves the heap first.
//	pq := heaps.New(func(a, b Task) bool { return a.Priority > b.Priority }, tasks...)
func New[T any](less LessFunc[T], values ...T) *Binary[T] {
	h := &Binary[T]{
		data: make([]T, len(values)),
		less: less,
	}
	copy(h.data, values)
	h.heapify()
	return h
}

// NewMin creates a min-heap over an ordered type: the smallest element has the
// highest priority. Values are heapified in O(n).
func NewMin[T constraints.Ordered](values ...T) *Binary[T] {
	return New(Min[T], values...)
}

// NewMax creates a max-heap over an ordered type: the largest element has the
// highest priority. Values are heapified in O(n).
func NewMax[T constraints.Ordered](values ...T) *Binary[T] {
	return New(Max[T], values...)
}

// Interface guards to ensure Binary implements the required interfaces.
var _ Heap[int] = &Binary[int]{}
var _ MutableHeap[int] = &Binary[int]{}

// heapify establishes the heap invariant over all elements in O(n) by sifting
// down every internal node from the last parent up to the root.
func (h *Binary[T]) heapify() {
	for i := len(h.data)/2 - 1; i >= 0; i-- {
		h.siftDown(i)
	}
}

// siftDown restores the heap invariant for the subtree rooted at index i by
// repeatedly swapping the node with its highest-priority child.
func (h *Binary[T]) siftDown(i int) {
	n := len(h.data)
	for {
		best := i
		left := 2*i + 1
		right := 2*i + 2
		if left < n && h.less(h.data[left], h.data[best]) {
			best = left
		}
		if right < n && h.less(h.data[right], h.data[best]) {
			best = right
		}
		if best == i {
			return
		}
		h.data[i], h.data[best] = h.data[best], h.data[i]
		i = best
	}
}

// siftUp restores the heap invariant for the node at index i by repeatedly
// swapping it with its parent while it has higher priority.
func (h *Binary[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !h.less(h.data[i], h.data[parent]) {
			return
		}
		h.data[i], h.data[parent] = h.data[parent], h.data[i]
		i = parent
	}
}

// clone returns an independent copy of the heap, sharing the comparator.
func (h *Binary[T]) clone() *Binary[T] {
	cp := make([]T, len(h.data))
	copy(cp, h.data)
	return &Binary[T]{data: cp, less: h.less}
}

// Peek returns the highest-priority element without removing it.
func (h *Binary[T]) Peek() (T, bool) {
	if len(h.data) == 0 {
		var zero T
		return zero, false
	}
	return h.data[0], true
}

// Length returns the number of elements in the heap.
func (h *Binary[T]) Length() int {
	return len(h.data)
}

// IsEmpty returns true if the heap contains no elements.
func (h *Binary[T]) IsEmpty() bool {
	return len(h.data) == 0
}

// ForEach executes the given function for each element in heap-array order.
func (h *Binary[T]) ForEach(fn func(element T)) {
	for _, element := range h.data {
		fn(element)
	}
}

// All returns an iterator over the elements in heap-array order.
func (h *Binary[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, element := range h.data {
			if !yield(element) {
				return
			}
		}
	}
}

// Drain returns an iterator yielding the elements in priority order without
// modifying the receiver.
func (h *Binary[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		clone := h.clone()
		for {
			element, ok := clone.PopInPlace()
			if !ok {
				return
			}
			if !yield(element) {
				return
			}
		}
	}
}

// AsSlice returns a copy of the elements in heap-array order.
func (h *Binary[T]) AsSlice() []T {
	result := make([]T, len(h.data))
	copy(result, h.data)
	return result
}

// AsSortedSlice returns the elements in priority order without modifying the
// receiver.
func (h *Binary[T]) AsSortedSlice() []T {
	clone := h.clone()
	result := make([]T, 0, clone.Length())
	for {
		element, ok := clone.PopInPlace()
		if !ok {
			return result
		}
		result = append(result, element)
	}
}

// Push returns a new heap with the given element added, leaving the receiver
// unchanged.
func (h *Binary[T]) Push(element T) Heap[T] {
	clone := h.clone()
	clone.PushInPlace(element)
	return clone
}

// PushMany returns a new heap with all given elements added, leaving the
// receiver unchanged.
func (h *Binary[T]) PushMany(elements ...T) Heap[T] {
	clone := h.clone()
	clone.PushManyInPlace(elements...)
	return clone
}

// Pop removes the highest-priority element, returning it, whether the heap was
// non-empty, and the resulting heap. The receiver is left unchanged.
func (h *Binary[T]) Pop() (T, bool, Heap[T]) {
	clone := h.clone()
	element, ok := clone.PopInPlace()
	return element, ok, clone
}

// PushInPlace adds the given element to the heap.
func (h *Binary[T]) PushInPlace(element T) {
	h.data = append(h.data, element)
	h.siftUp(len(h.data) - 1)
}

// PushManyInPlace adds all given elements to the heap.
func (h *Binary[T]) PushManyInPlace(elements ...T) {
	for _, element := range elements {
		h.PushInPlace(element)
	}
}

// PopInPlace removes and returns the highest-priority element, reporting
// whether the heap was non-empty.
func (h *Binary[T]) PopInPlace() (T, bool) {
	n := len(h.data)
	if n == 0 {
		var zero T
		return zero, false
	}
	top := h.data[0]
	last := n - 1
	h.data[0] = h.data[last]
	var zero T
	h.data[last] = zero // avoid retaining a reference in the freed slot
	h.data = h.data[:last]
	if len(h.data) > 0 {
		h.siftDown(0)
	}
	return top, true
}
