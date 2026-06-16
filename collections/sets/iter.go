package sets

import "iter"

// seqOf returns an iterator over the given elements.
func seqOf[T any](elements []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, element := range elements {
			if !yield(element) {
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

// All returns an iterator over the elements, over a snapshot taken under the
// lock. It is safe for concurrent use. Iteration order is unspecified.
func (ch *ConcurrentHash[T]) All() iter.Seq[T] { return seqOf(ch.AsSlice()) }

// All returns an iterator over the elements, over a snapshot taken under the
// lock. It is safe for concurrent use. Iteration order is unspecified.
func (ch *ConcurrentHashRW[T]) All() iter.Seq[T] { return seqOf(ch.AsSlice()) }
