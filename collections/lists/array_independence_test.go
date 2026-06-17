package lists_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// TestArray_ImmutableOpsIndependentOfReceiver mirrors the ConcurrentArray
// independence test for the non-concurrent Array: slices/lists returned by
// immutable operations must not alias the receiver's backing array. Pop/Dequeue
// return sub-slices, so a later in-place mutation of the receiver (e.g.
// SortInPlace, which sorts the backing array in place) must not change the
// previously returned result. ascending is defined in
// concurrentarray_independence_test.go.
func TestArray_ImmutableOpsIndependentOfReceiver(t *testing.T) {
	t.Run("Pop", func(t *testing.T) {
		a := lists.NewArray(3, 1, 2)
		_, _, popped := a.Pop()

		a.SortInPlace(ascending) // mutates the receiver's backing array

		if want := []int{3, 1}; !reflect.DeepEqual(popped.AsSlice(), want) {
			t.Errorf("Pop result = %v, want %v (returned slice aliased the receiver)", popped, want)
		}
	})

	t.Run("Dequeue", func(t *testing.T) {
		a := lists.NewArray(3, 1, 2)
		_, _, dequeued := a.Dequeue()

		a.SortInPlace(ascending)

		if want := []int{1, 2}; !reflect.DeepEqual(dequeued.AsSlice(), want) {
			t.Errorf("Dequeue result = %v, want %v (returned slice aliased the receiver)", dequeued, want)
		}
	})

	t.Run("Insert does not mutate receiver", func(t *testing.T) {
		a := lists.NewArray(1, 2, 3)
		inserted := a.Insert(1, 99)

		if want := []int{1, 99, 2, 3}; !reflect.DeepEqual(inserted.AsSlice(), want) {
			t.Errorf("Insert result = %v, want %v", inserted, want)
		}
		if got := a.AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
			t.Errorf("receiver mutated by immutable Insert: %v, want [1 2 3]", got)
		}
	})

	t.Run("Push does not mutate receiver", func(t *testing.T) {
		a := lists.NewArray(1, 2, 3)
		pushed := a.Push(9)

		if want := []int{1, 2, 3, 9}; !reflect.DeepEqual(pushed.AsSlice(), want) {
			t.Errorf("Push result = %v, want %v", pushed.AsSlice(), want)
		}
		if got := a.AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
			t.Errorf("receiver mutated by immutable Push: %v, want [1 2 3]", got)
		}
	})

	t.Run("Enqueue does not mutate receiver", func(t *testing.T) {
		a := lists.NewArray(1, 2, 3)
		enqueued := a.Enqueue(9)

		if want := []int{1, 2, 3, 9}; !reflect.DeepEqual(enqueued.AsSlice(), want) {
			t.Errorf("Enqueue result = %v, want %v", enqueued.AsSlice(), want)
		}
		if got := a.AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
			t.Errorf("receiver mutated by immutable Enqueue: %v, want [1 2 3]", got)
		}
	})

	// The subtests below exercise the shared-capacity aliasing path the issue
	// flagged: when the receiver's backing array has spare capacity, a plain
	// append (slices.Push) writes the new element into that spare slot, so the
	// returned List aliases the receiver. PopInPlace shrinks the length while
	// keeping the capacity, manufacturing that spare slot; a subsequent in-place
	// append must not corrupt the previously returned immutable result. These
	// would fail if Push/Enqueue used slices.Push instead of slices.PushCopy.
	t.Run("Push does not alias receiver spare capacity", func(t *testing.T) {
		a := lists.NewArray(1, 2, 3, 4)
		a.PopInPlace() // a.elements = [1 2 3] with spare capacity for one more

		pushed := a.Push(9)
		a.PushInPlace(7) // appends into the shared spare slot under a plain append

		got := pushed.AsSlice()
		want := []int{1, 2, 3, 9}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Push result = %v, want %v (returned slice aliased the receiver's spare capacity)", got, want)
		}
	})

	t.Run("Enqueue does not alias receiver spare capacity", func(t *testing.T) {
		a := lists.NewArray(1, 2, 3, 4)
		a.PopInPlace() // a.elements = [1 2 3] with spare capacity for one more

		enqueued := a.Enqueue(9)
		a.EnqueueInPlace(7) // appends into the shared spare slot under a plain append

		got := enqueued.AsSlice()
		want := []int{1, 2, 3, 9}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Enqueue result = %v, want %v (returned slice aliased the receiver's spare capacity)", got, want)
		}
	})
}

// TestNewArray_CopiesCallerSlice verifies that NewArray copies the caller's
// variadic backing array, so a later mutation of the source slice does not
// reach into the list's backing array.
func TestNewArray_CopiesCallerSlice(t *testing.T) {
	src := []int{1, 2, 3}
	a := lists.NewArray(src...)

	src[0] = 99

	if got := a.AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("list = %v, want [1 2 3] (constructor aliased the caller's slice)", got)
	}
}
