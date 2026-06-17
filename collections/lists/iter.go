package lists

import "iter"

// seqAll returns an iterator over index/value pairs of elements, front to back.
func seqAll[T any](elements []T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, element := range elements {
			if !yield(i, element) {
				return
			}
		}
	}
}

// seqValues returns an iterator over values of elements, front to back.
func seqValues[T any](elements []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, element := range elements {
			if !yield(element) {
				return
			}
		}
	}
}

// seqBackward returns an iterator over index/value pairs of elements, back to
// front; indices still count from the front.
func seqBackward[T any](elements []T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
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
func (a *Array[T]) All() iter.Seq2[int, T] { return seqAll(a.elements) }

// Values returns an iterator over values, front to back.
func (a *Array[T]) Values() iter.Seq[T] { return seqValues(a.elements) }

// Backward returns an iterator over index/value pairs, back to front.
func (a *Array[T]) Backward() iter.Seq2[int, T] { return seqBackward(a.elements) }

// All returns an iterator over index/value pairs, front to back.
func (l *Linked[T]) All() iter.Seq2[int, T] { return seqAll(l.AsSlice()) }

// Values returns an iterator over values, front to back.
func (l *Linked[T]) Values() iter.Seq[T] { return seqValues(l.AsSlice()) }

// Backward returns an iterator over index/value pairs, back to front.
func (l *Linked[T]) Backward() iter.Seq2[int, T] { return seqBackward(l.AsSlice()) }

// All returns an iterator over index/value pairs, front to back.
func (dl *DoublyLinked[T]) All() iter.Seq2[int, T] { return seqAll(dl.AsSlice()) }

// Values returns an iterator over values, front to back.
func (dl *DoublyLinked[T]) Values() iter.Seq[T] { return seqValues(dl.AsSlice()) }

// Backward returns an iterator over index/value pairs, back to front.
func (dl *DoublyLinked[T]) Backward() iter.Seq2[int, T] { return seqBackward(dl.AsSlice()) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock. It is safe for concurrent use.
func (a *ConcurrentArray[T]) All() iter.Seq2[int, T] { return seqAll(a.AsSlice()) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Values() iter.Seq[T] { return seqValues(a.AsSlice()) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Backward() iter.Seq2[int, T] { return seqBackward(a.AsSlice()) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) All() iter.Seq2[int, T] { return seqAll(a.AsSlice()) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Values() iter.Seq[T] { return seqValues(a.AsSlice()) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Backward() iter.Seq2[int, T] { return seqBackward(a.AsSlice()) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock. It is safe for concurrent use.
func (cl *ConcurrentLinked[T]) All() iter.Seq2[int, T] { return seqAll(cl.AsSlice()) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock. It is safe for concurrent use.
func (cl *ConcurrentLinked[T]) Values() iter.Seq[T] { return seqValues(cl.AsSlice()) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock. It is safe for concurrent use.
func (cl *ConcurrentLinked[T]) Backward() iter.Seq2[int, T] { return seqBackward(cl.AsSlice()) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock. It is safe for concurrent use.
func (cl *ConcurrentRWLinked[T]) All() iter.Seq2[int, T] { return seqAll(cl.AsSlice()) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock. It is safe for concurrent use.
func (cl *ConcurrentRWLinked[T]) Values() iter.Seq[T] { return seqValues(cl.AsSlice()) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock. It is safe for concurrent use.
func (cl *ConcurrentRWLinked[T]) Backward() iter.Seq2[int, T] { return seqBackward(cl.AsSlice()) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock. It is safe for concurrent use.
func (cl *ConcurrentDoublyLinked[T]) All() iter.Seq2[int, T] { return seqAll(cl.AsSlice()) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock. It is safe for concurrent use.
func (cl *ConcurrentDoublyLinked[T]) Values() iter.Seq[T] { return seqValues(cl.AsSlice()) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock. It is safe for concurrent use.
func (cl *ConcurrentDoublyLinked[T]) Backward() iter.Seq2[int, T] { return seqBackward(cl.AsSlice()) }

// All returns an iterator over index/value pairs, front to back, over a snapshot
// taken under the lock. It is safe for concurrent use.
func (cl *ConcurrentRWDoublyLinked[T]) All() iter.Seq2[int, T] { return seqAll(cl.AsSlice()) }

// Values returns an iterator over values, front to back, over a snapshot taken
// under the lock. It is safe for concurrent use.
func (cl *ConcurrentRWDoublyLinked[T]) Values() iter.Seq[T] { return seqValues(cl.AsSlice()) }

// Backward returns an iterator over index/value pairs, back to front, over a
// snapshot taken under the lock. It is safe for concurrent use.
func (cl *ConcurrentRWDoublyLinked[T]) Backward() iter.Seq2[int, T] { return seqBackward(cl.AsSlice()) }
