package heaps

import "iter"

// FromSeq builds a Binary heap from the values produced by seq, ordered by less.
// The values are heapified in O(n) after collection. It is the inbound
// counterpart to the All iterator.
func FromSeq[T any](less LessFunc[T], seq iter.Seq[T]) *Binary[T] {
	var values []T
	for value := range seq {
		values = append(values, value)
	}
	return New(less, values...)
}
