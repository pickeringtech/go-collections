package lists

import (
	"reflect"
	"testing"
)

// These tests pin down the fix for removeNode on circular doubly linked lists,
// where head/tail must be detected by node identity (prev/next never go nil).

func TestDoublyLinkedCircular_DequeueInPlaceRemovesHead(t *testing.T) {
	dl := NewDoublyLinkedCircular(1, 2, 3)

	value, ok := dl.DequeueInPlace()
	if !ok || value != 1 {
		t.Fatalf("DequeueInPlace() = (%d, %t), want (1, true)", value, ok)
	}
	if dl.Length() != 2 {
		t.Errorf("Length() = %d, want 2", dl.Length())
	}
	got := dl.GetAsSlice()
	want := []int{2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAsSlice() = %v, want %v", got, want)
	}
}

func TestDoublyLinkedCircular_FilterInPlaceRemovesHeadAndTail(t *testing.T) {
	dl := NewDoublyLinkedCircular(1, 2, 3, 4, 5)

	// Remove the first (head) and last (tail) elements.
	dl.FilterInPlace(func(v int) bool {
		return v != 1 && v != 5
	})

	if dl.Length() != 3 {
		t.Errorf("Length() = %d, want 3", dl.Length())
	}
	got := dl.GetAsSlice()
	want := []int{2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAsSlice() = %v, want %v", got, want)
	}
}

func TestDoublyLinkedCircular_RemoveDownToEmpty(t *testing.T) {
	dl := NewDoublyLinkedCircular(7, 8)

	first, ok := dl.DequeueInPlace()
	if !ok || first != 7 {
		t.Fatalf("first DequeueInPlace() = (%d, %t), want (7, true)", first, ok)
	}
	second, ok := dl.DequeueInPlace()
	if !ok || second != 8 {
		t.Fatalf("second DequeueInPlace() = (%d, %t), want (8, true)", second, ok)
	}
	if dl.Length() != 0 {
		t.Errorf("Length() = %d, want 0", dl.Length())
	}
	if got := dl.GetAsSlice(); len(got) != 0 {
		t.Errorf("GetAsSlice() = %v, want empty", got)
	}
}
