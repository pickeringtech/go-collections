package lists

import "iter"

// Every list iterator ranges over a snapshot of the list taken when ranging
// BEGINS (not when All/Values/Backward is called): an iterator stored and ranged
// later sees the list as it is then, and the list may be modified freely during
// the range — by the loop body itself or, for the concurrent types, by another
// goroutine — without affecting the elements yielded. The snapshot functions
// below are called once per range, inside the returned iterator.

// seqAll returns an iterator over index/value pairs of snapshot(), front to back.
func seqAll[T any](snapshot func() []T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, element := range snapshot() {
			if !yield(i, element) {
				return
			}
		}
	}
}

// seqValues returns an iterator over the values of snapshot(), front to back.
func seqValues[T any](snapshot func() []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, element := range snapshot() {
			if !yield(element) {
				return
			}
		}
	}
}

// seqBackward returns an iterator over index/value pairs of snapshot(), back to
// front; indices still count from the front.
func seqBackward[T any](snapshot func() []T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		elements := snapshot()
		for i := len(elements) - 1; i >= 0; i-- {
			if !yield(i, elements[i]) {
				return
			}
		}
	}
}

// FromSeq builds an Array from the values produced by seq, preserving their
// order. It is the inbound counterpart to the Values iterator.
func FromSeq[T any](seq iter.Seq[T]) *Array[T] {
	var elements []T
	for element := range seq {
		elements = append(elements, element)
	}
	return NewArray(elements...)
}

// All returns an iterator over index/value pairs, front to back.
func (a *Array[T]) All() iter.Seq2[int, T] { return seqAll(a.AsSlice) }

// Values returns an iterator over values, front to back.
func (a *Array[T]) Values() iter.Seq[T] { return seqValues(a.AsSlice) }

// Backward returns an iterator over index/value pairs, back to front.
func (a *Array[T]) Backward() iter.Seq2[int, T] { return seqBackward(a.AsSlice) }

// All returns an iterator over index/value pairs, front to back.
func (l *Linked[T]) All() iter.Seq2[int, T] { return seqAll(l.AsSlice) }

// Values returns an iterator over values, front to back.
func (l *Linked[T]) Values() iter.Seq[T] { return seqValues(l.AsSlice) }

// Backward returns an iterator over index/value pairs, back to front.
func (l *Linked[T]) Backward() iter.Seq2[int, T] { return seqBackward(l.AsSlice) }

// All returns an iterator over index/value pairs, front to back.
func (dl *DoublyLinked[T]) All() iter.Seq2[int, T] { return seqAll(dl.AsSlice) }

// Values returns an iterator over values, front to back.
func (dl *DoublyLinked[T]) Values() iter.Seq[T] { return seqValues(dl.AsSlice) }

// Backward returns an iterator over index/value pairs, back to front.
func (dl *DoublyLinked[T]) Backward() iter.Seq2[int, T] { return seqBackward(dl.AsSlice) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock when ranging begins. It is safe for concurrent use.
func (a *ConcurrentArray[T]) All() iter.Seq2[int, T] { return seqAll(a.AsSlice) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock when ranging begins. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Values() iter.Seq[T] { return seqValues(a.AsSlice) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock when ranging begins. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Backward() iter.Seq2[int, T] { return seqBackward(a.AsSlice) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock when ranging begins. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) All() iter.Seq2[int, T] { return seqAll(a.AsSlice) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock when ranging begins. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Values() iter.Seq[T] { return seqValues(a.AsSlice) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock when ranging begins. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Backward() iter.Seq2[int, T] { return seqBackward(a.AsSlice) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentLinked[T]) All() iter.Seq2[int, T] { return seqAll(cl.AsSlice) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentLinked[T]) Values() iter.Seq[T] { return seqValues(cl.AsSlice) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentLinked[T]) Backward() iter.Seq2[int, T] { return seqBackward(cl.AsSlice) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentRWLinked[T]) All() iter.Seq2[int, T] { return seqAll(cl.AsSlice) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentRWLinked[T]) Values() iter.Seq[T] { return seqValues(cl.AsSlice) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentRWLinked[T]) Backward() iter.Seq2[int, T] { return seqBackward(cl.AsSlice) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentDoublyLinked[T]) All() iter.Seq2[int, T] { return seqAll(cl.AsSlice) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentDoublyLinked[T]) Values() iter.Seq[T] { return seqValues(cl.AsSlice) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentDoublyLinked[T]) Backward() iter.Seq2[int, T] { return seqBackward(cl.AsSlice) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentRWDoublyLinked[T]) All() iter.Seq2[int, T] { return seqAll(cl.AsSlice) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentRWDoublyLinked[T]) Values() iter.Seq[T] { return seqValues(cl.AsSlice) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock when ranging begins. It is safe for concurrent use.
func (cl *ConcurrentRWDoublyLinked[T]) Backward() iter.Seq2[int, T] { return seqBackward(cl.AsSlice) }
