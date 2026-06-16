package deques

import "iter"

// FromSeq builds an unbounded RingBuffer from the values produced by seq,
// preserving their order (the first value becomes the front). It is the inbound
// counterpart to the Values iterator.
func FromSeq[T any](seq iter.Seq[T]) *RingBuffer[T] {
	var elements []T
	for element := range seq {
		elements = append(elements, element)
	}
	return NewRingBuffer(elements...)
}
