package sets

import "iter"

// elementsSeq returns an iterator over the elements that snapshot returns.
// snapshot is called each time a range over the iterator begins — never when
// the iterator is created — so a stored iterator always reflects the set as it
// is when it is ranged, and the loop body never runs while snapshot holds a
// lock.
func elementsSeq[T any](snapshot func() []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, element := range snapshot() {
			if !yield(element) {
				return
			}
		}
	}
}

// reverseElementsSeq is elementsSeq yielding the snapshot in reverse order.
func reverseElementsSeq[T any](snapshot func() []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		elements := snapshot()
		for i := len(elements) - 1; i >= 0; i-- {
			if !yield(elements[i]) {
				return
			}
		}
	}
}

// FromSeq builds a Hash set from the elements produced by seq. Duplicate
// elements are collapsed, as with any set. It is the inbound counterpart to the
// All iterator.
func FromSeq[T comparable](seq iter.Seq[T]) Hash[T] {
	h := make(Hash[T])
	for element := range seq {
		h[element] = struct{}{}
	}
	return h
}

// All returns an iterator over the elements. Iteration order is unspecified.
func (h Hash[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for element := range h {
			if !yield(element) {
				return
			}
		}
	}
}

// All returns an iterator over the elements. Each range over it snapshots the
// elements under the lock when it begins — not when All is called — and yields
// them after releasing it, so it is safe for concurrent use and the loop body
// may call back into the set. Iteration order is unspecified.
func (ch *ConcurrentHash[T]) All() iter.Seq[T] { return elementsSeq(ch.AsSlice) }

// All returns an iterator over the elements. Each range over it snapshots the
// elements under the read lock when it begins — not when All is called — and
// yields them after releasing it, so it is safe for concurrent use and the loop
// body may call back into the set. Iteration order is unspecified.
func (ch *ConcurrentHashRW[T]) All() iter.Seq[T] { return elementsSeq(ch.AsSlice) }
