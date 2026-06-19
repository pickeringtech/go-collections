package streaming

import (
	"github.com/pickeringtech/go-collections/collections/heaps"
	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/slices"
)

// TopK retains the k highest-ranked elements seen across an unbounded stream,
// using O(k) memory regardless of how many elements are fed in. Rank is decided
// by a [heaps.LessFunc]: less(a, b) reports that a ranks below b, so the
// elements TopK keeps are the k for which less reports no greater element — the
// "largest" under the comparator.
//
// Internally TopK is a size-bounded min-heap (the lowest-ranked retained
// element sits at the top), so Add is O(log k): a new element is discarded in
// O(1) when it cannot beat the current minimum, and otherwise displaces it in
// O(log k). Ties never displace an incumbent — when a new element ranks equal
// to the retained minimum, the element already held is kept.
//
// TopK is exact, not approximate: Result is precisely the k highest-ranked
// elements of the stream. It is single-threaded; see the package documentation
// on thread safety.
type TopK[T any] struct {
	k    int
	less heaps.LessFunc[T]
	heap *heaps.Binary[T]
}

// NewTopK creates a TopK that retains the k highest-ranked elements according
// to less, where less(a, b) reports whether a ranks below b. For k <= 0 the
// resulting TopK retains nothing and Result is always empty.
//
// To rank by an ordered type's natural order, use [NewTopKOrdered].
func NewTopK[T any](k int, less heaps.LessFunc[T]) *TopK[T] {
	return &TopK[T]{
		k:    k,
		less: less,
		heap: heaps.New(less),
	}
}

// NewTopKOrdered creates a TopK over an ordered type that retains the k largest
// elements by natural order. For k <= 0 it retains nothing.
func NewTopKOrdered[T constraints.Ordered](k int) *TopK[T] {
	return NewTopK[T](k, heaps.Min[T])
}

// Add feeds one element into the stream. It runs in O(log k): the element is
// kept if fewer than k elements are held, or if it ranks strictly above the
// current lowest-ranked retained element (which it then displaces); otherwise
// it is discarded. Add is a no-op when k <= 0.
func (t *TopK[T]) Add(element T) {
	if t.k <= 0 {
		return
	}
	if t.heap.Length() < t.k {
		t.heap.PushInPlace(element)
		return
	}
	lowest, _ := t.heap.Peek()
	if t.less(lowest, element) {
		t.heap.PopInPlace()
		t.heap.PushInPlace(element)
	}
}

// Result returns the retained elements highest-ranked first, as a non-nil
// slice. It does not modify the TopK, so Add may continue afterwards. Duplicates
// are retained, so the length is min(k, number of elements fed). Result is
// O(k log k).
func (t *TopK[T]) Result() []T {
	return slices.Reverse(t.heap.AsSortedSlice())
}

// Len returns the number of elements currently retained, between 0 and
// max(k, 0) — it is always 0 when k <= 0.
func (t *TopK[T]) Len() int {
	return t.heap.Length()
}
