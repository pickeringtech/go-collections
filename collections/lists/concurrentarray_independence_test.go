package lists_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

func ascending(a, b int) bool { return a < b }

// TestConcurrentArray_ImmutableOpsIndependentOfReceiver verifies that slices
// returned by immutable operations do not alias the receiver's backing array.
// Pop/Dequeue return sub-slices, so a later in-place mutation of the receiver
// (e.g. SortInPlace, which sorts the backing array in place) must not change
// the previously returned slice.
func TestConcurrentArray_ImmutableOpsIndependentOfReceiver(t *testing.T) {
	t.Run("Pop", func(t *testing.T) {
		a := lists.NewConcurrentArray(3, 1, 2)
		_, _, popped := a.Pop()

		a.SortInPlace(ascending) // mutates the receiver's backing array

		if want := []int{3, 1}; !reflect.DeepEqual(popped, want) {
			t.Errorf("Pop result = %v, want %v (returned slice aliased the receiver)", popped, want)
		}
	})

	t.Run("Dequeue", func(t *testing.T) {
		a := lists.NewConcurrentArray(3, 1, 2)
		_, _, dequeued := a.Dequeue()

		a.SortInPlace(ascending)

		if want := []int{1, 2}; !reflect.DeepEqual(dequeued, want) {
			t.Errorf("Dequeue result = %v, want %v (returned slice aliased the receiver)", dequeued, want)
		}
	})

	t.Run("Insert does not mutate receiver", func(t *testing.T) {
		a := lists.NewConcurrentArray(1, 2, 3)
		inserted := a.Insert(1, 99)

		if want := []int{1, 99, 2, 3}; !reflect.DeepEqual(inserted, want) {
			t.Errorf("Insert result = %v, want %v", inserted, want)
		}
		if got := a.AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
			t.Errorf("receiver mutated by immutable Insert: %v, want [1 2 3]", got)
		}
	})

	t.Run("Push result is independent", func(t *testing.T) {
		a := lists.NewConcurrentArray(1, 2, 3)
		pushed := a.Push(9)
		pushed[0] = 999

		got, _ := a.Get(0, -1)
		if got != 1 {
			t.Errorf("receiver index 0 = %d, want 1 (Push result aliased the receiver)", got)
		}
	})
}

// TestConcurrentRWArray_ImmutableOpsIndependentOfReceiver mirrors the above for
// the read-write mutex variant.
func TestConcurrentRWArray_ImmutableOpsIndependentOfReceiver(t *testing.T) {
	t.Run("Pop", func(t *testing.T) {
		a := lists.NewConcurrentRWArray(3, 1, 2)
		_, _, popped := a.Pop()

		a.SortInPlace(ascending)

		if want := []int{3, 1}; !reflect.DeepEqual(popped, want) {
			t.Errorf("Pop result = %v, want %v (returned slice aliased the receiver)", popped, want)
		}
	})

	t.Run("Dequeue", func(t *testing.T) {
		a := lists.NewConcurrentRWArray(3, 1, 2)
		_, _, dequeued := a.Dequeue()

		a.SortInPlace(ascending)

		if want := []int{1, 2}; !reflect.DeepEqual(dequeued, want) {
			t.Errorf("Dequeue result = %v, want %v (returned slice aliased the receiver)", dequeued, want)
		}
	})

	t.Run("Insert does not mutate receiver", func(t *testing.T) {
		a := lists.NewConcurrentRWArray(1, 2, 3)
		inserted := a.Insert(1, 99)

		if want := []int{1, 99, 2, 3}; !reflect.DeepEqual(inserted, want) {
			t.Errorf("Insert result = %v, want %v", inserted, want)
		}
		if got := a.AsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
			t.Errorf("receiver mutated by immutable Insert: %v, want [1 2 3]", got)
		}
	})
}
