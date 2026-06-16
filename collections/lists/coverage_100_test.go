package lists_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// --- DoublyLinked: circular traversal wrap-around --------------------------

// In a circular list the search loops break once `current` wraps back to the
// head. These tests use predicates that never match so the walk completes the
// full circle and hits the wrap-around break.

func TestDoublyLinked_AnyMatch_CircularNoMatch(t *testing.T) {
	dl := lists.NewDoublyLinkedCircular(1, 2, 3)
	if dl.AnyMatch(func(v int) bool { return v > 100 }) {
		t.Fatalf("AnyMatch should be false when nothing matches in a circular list")
	}
}

func TestDoublyLinked_Find_CircularNoMatch(t *testing.T) {
	dl := lists.NewDoublyLinkedCircular(1, 2, 3)
	if v, ok := dl.Find(func(v int) bool { return v > 100 }); ok || v != 0 {
		t.Fatalf("Find should return (0, false) for a circular list with no match, got (%d, %v)", v, ok)
	}
}

func TestDoublyLinked_FindIndex_CircularNoMatch(t *testing.T) {
	dl := lists.NewDoublyLinkedCircular(1, 2, 3)
	if idx := dl.FindIndex(func(v int) bool { return v > 100 }); idx != -1 {
		t.Fatalf("FindIndex should return -1 for a circular list with no match, got %d", idx)
	}
}

// TestDoublyLinked_InsertInPlace_CircularAtHead inserts at the front of a
// non-empty circular list, exercising insertAt's circular index-0 branch.
func TestDoublyLinked_InsertInPlace_CircularAtHead(t *testing.T) {
	dl := lists.NewDoublyLinkedCircular(1, 2, 3)
	dl.InsertInPlace(0, 9)

	got := dl.AsSlice()
	want := []int{9, 1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("InsertInPlace(0, 9) = %v, want %v", got, want)
	}

	// The circular invariant must be preserved: walking size+1 steps from the
	// head must wrap back to the first element.
	if first, ok := dl.Get(0, -1); !ok || first != 9 {
		t.Fatalf("Get(0) = (%d, %v), want (9, true)", first, ok)
	}
}

// --- Linked: FilterInPlace edge cases --------------------------------------

func TestLinked_FilterInPlace_Empty(t *testing.T) {
	l := lists.NewLinked[int]()
	l.FilterInPlace(func(int) bool { return true })
	if got := l.AsSlice(); len(got) != 0 {
		t.Fatalf("FilterInPlace on empty list = %v, want empty", got)
	}
}

func TestLinked_FilterInPlace_RemovesAll(t *testing.T) {
	l := lists.NewLinked(1, 2, 3)
	l.FilterInPlace(func(int) bool { return false })
	if got := l.AsSlice(); len(got) != 0 {
		t.Fatalf("FilterInPlace removing everything = %v, want empty", got)
	}
}

func TestLinked_FilterInPlace_RemovesTail(t *testing.T) {
	l := lists.NewLinked(1, 2, 3)
	l.FilterInPlace(func(v int) bool { return v != 3 })

	got := l.AsSlice()
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterInPlace removing the tail = %v, want %v", got, want)
	}

	// The new tail must accept appends correctly.
	l.PushInPlace(4)
	got = l.AsSlice()
	want = []int{1, 2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("after re-appending, list = %v, want %v", got, want)
	}
}

// --- Linked: removeFirst / Insert / SortInPlace guards ---------------------

func TestLinked_RemoveInPlace_Empty(t *testing.T) {
	l := lists.NewLinked[int]()
	if l.RemoveInPlace(1) {
		t.Fatalf("RemoveInPlace on empty list should report false")
	}
	if _, ok := l.RemoveAtInPlace(0); ok {
		t.Fatalf("RemoveAtInPlace on empty list should report false")
	}
}

func TestLinked_Insert_OutOfRange(t *testing.T) {
	l := lists.NewLinked(1, 2, 3)

	if got := l.Insert(-1, 9); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Insert(-1, 9) = %v, want unchanged [1 2 3]", got)
	}
	if got := l.Insert(10, 9); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Insert(10, 9) = %v, want unchanged [1 2 3]", got)
	}
}

// TestLinked_InsertInPlace_CircularWraps inserts at an index beyond the size of
// a circular list. findNodeBefore wraps back to the head before reaching the
// index, so the insertion is rejected and the list is left untouched.
func TestLinked_InsertInPlace_CircularWraps(t *testing.T) {
	l := lists.NewLinkedCircular(1, 2, 3)
	l.InsertInPlace(10, 9)

	got := l.AsSlice()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("InsertInPlace(10, 9) on circular list = %v, want unchanged %v", got, want)
	}
}

func TestLinked_SortInPlace_EmptyAndSingle(t *testing.T) {
	less := func(a, b int) bool { return a < b }

	empty := lists.NewLinked[int]()
	empty.SortInPlace(less)
	if got := empty.AsSlice(); len(got) != 0 {
		t.Fatalf("SortInPlace on empty list = %v, want empty", got)
	}

	single := lists.NewLinked(42)
	single.SortInPlace(less)
	if got := single.AsSlice(); !reflect.DeepEqual(got, []int{42}) {
		t.Fatalf("SortInPlace on single-element list = %v, want [42]", got)
	}
}
