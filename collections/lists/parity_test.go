package lists_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// mutableListConstructor builds a MutableList[int] seeded with the given
// elements. Every concrete implementation provides one so the parity tests can
// run against all of them uniformly.
type mutableListConstructor struct {
	name string
	make func(...int) lists.MutableList[int]
}

func allMutableListConstructors() []mutableListConstructor {
	return []mutableListConstructor{
		{"Array", func(e ...int) lists.MutableList[int] { return lists.NewArray(e...) }},
		{"Linked", func(e ...int) lists.MutableList[int] { return lists.NewLinked(e...) }},
		{"DoublyLinked", func(e ...int) lists.MutableList[int] { return lists.NewDoublyLinked(e...) }},
		{"ConcurrentArray", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentArray(e...) }},
		{"ConcurrentLinked", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentLinked(e...) }},
		{"ConcurrentDoublyLinked", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentDoublyLinked(e...) }},
		{"ConcurrentRWArray", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentRWArray(e...) }},
		{"ConcurrentRWLinked", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentRWLinked(e...) }},
		{"ConcurrentRWDoublyLinked", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentRWDoublyLinked(e...) }},
	}
}

func TestMutableList_IsEmpty(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			if !ctor.make().IsEmpty() {
				t.Errorf("expected freshly created empty list to be empty")
			}
			if ctor.make(1).IsEmpty() {
				t.Errorf("expected non-empty list to report not empty")
			}
		})
	}
}

func TestMutableList_Clear(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			l := ctor.make(1, 2, 3)
			l.Clear()
			if l.Length() != 0 {
				t.Errorf("expected length 0 after Clear, got %d", l.Length())
			}
			if !l.IsEmpty() {
				t.Errorf("expected IsEmpty true after Clear")
			}
			// The list remains usable after clearing.
			l.PushInPlace(9)
			if got := l.AsSlice(); !reflect.DeepEqual(got, []int{9}) {
				t.Errorf("expected [9] after push following Clear, got %v", got)
			}
		})
	}
}

func TestMutableList_RemoveAt(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			l := ctor.make(10, 20, 30)

			got := l.RemoveAt(1)
			if !reflect.DeepEqual(got, []int{10, 30}) {
				t.Errorf("RemoveAt(1) = %v, want [10 30]", got)
			}
			// Immutable: the receiver is untouched.
			if l.Length() != 3 {
				t.Errorf("RemoveAt must not mutate receiver, length = %d", l.Length())
			}

			// Out of bounds leaves the elements unchanged.
			for _, idx := range []int{-1, 3, 99} {
				if got := l.RemoveAt(idx); !reflect.DeepEqual(got, []int{10, 20, 30}) {
					t.Errorf("RemoveAt(%d) = %v, want [10 20 30]", idx, got)
				}
			}

			// The returned slice is independent of the list, including on the
			// out-of-bounds path: mutating it must not affect the receiver.
			oob := l.RemoveAt(99)
			if len(oob) > 0 {
				oob[0] = -1
			}
			if got := l.RemoveAt(99); !reflect.DeepEqual(got, []int{10, 20, 30}) {
				t.Errorf("RemoveAt result leaked backing array; list now %v", got)
			}
		})
	}
}

func TestMutableList_Remove(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			l := ctor.make(10, 20, 20, 30)

			// Removes the first matching value only.
			got := l.Remove(20)
			if !reflect.DeepEqual(got, []int{10, 20, 30}) {
				t.Errorf("Remove(20) = %v, want [10 20 30]", got)
			}
			if l.Length() != 4 {
				t.Errorf("Remove must not mutate receiver, length = %d", l.Length())
			}

			// Absent value leaves the elements unchanged.
			if got := l.Remove(99); !reflect.DeepEqual(got, []int{10, 20, 20, 30}) {
				t.Errorf("Remove(99) = %v, want [10 20 20 30]", got)
			}
		})
	}
}

func TestMutableList_RemoveAtInPlace(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			// Remove from the middle.
			l := ctor.make(10, 20, 30)
			value, ok := l.RemoveAtInPlace(1)
			if !ok || value != 20 {
				t.Errorf("RemoveAtInPlace(1) = (%d, %v), want (20, true)", value, ok)
			}
			if got := l.AsSlice(); !reflect.DeepEqual(got, []int{10, 30}) {
				t.Errorf("after RemoveAtInPlace(1) = %v, want [10 30]", got)
			}

			// Remove the head and the tail.
			l = ctor.make(10, 20, 30)
			if v, ok := l.RemoveAtInPlace(0); !ok || v != 10 {
				t.Errorf("RemoveAtInPlace(0) = (%d, %v), want (10, true)", v, ok)
			}
			if v, ok := l.RemoveAtInPlace(l.Length() - 1); !ok || v != 30 {
				t.Errorf("RemoveAtInPlace(last) = (%d, %v), want (30, true)", v, ok)
			}
			if got := l.AsSlice(); !reflect.DeepEqual(got, []int{20}) {
				t.Errorf("after head+tail removal = %v, want [20]", got)
			}

			// Out of bounds returns the zero value and false without mutating.
			l = ctor.make(10, 20, 30)
			for _, idx := range []int{-1, 3, 99} {
				if v, ok := l.RemoveAtInPlace(idx); ok || v != 0 {
					t.Errorf("RemoveAtInPlace(%d) = (%d, %v), want (0, false)", idx, v, ok)
				}
			}
			if l.Length() != 3 {
				t.Errorf("out-of-bounds RemoveAtInPlace must not mutate, length = %d", l.Length())
			}
		})
	}
}

func TestMutableList_RemoveInPlace(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			l := ctor.make(10, 20, 20, 30)

			if !l.RemoveInPlace(20) {
				t.Errorf("RemoveInPlace(20) = false, want true")
			}
			if got := l.AsSlice(); !reflect.DeepEqual(got, []int{10, 20, 30}) {
				t.Errorf("after RemoveInPlace(20) = %v, want [10 20 30]", got)
			}

			if l.RemoveInPlace(99) {
				t.Errorf("RemoveInPlace(99) = true, want false")
			}
			if got := l.AsSlice(); !reflect.DeepEqual(got, []int{10, 20, 30}) {
				t.Errorf("after absent RemoveInPlace = %v, want [10 20 30]", got)
			}

			// Removing the only element empties the list.
			single := ctor.make(7)
			if !single.RemoveInPlace(7) {
				t.Errorf("RemoveInPlace(7) on single-element list = false, want true")
			}
			if !single.IsEmpty() {
				t.Errorf("expected empty list after removing only element")
			}
		})
	}
}

// TestMutableList_InsertParity pins down the shared Insert/InsertInPlace
// contract (issue #78): every implementation accepts 0 <= index <= Length(),
// treats index == Length() as an append (so inserting into an empty list yields
// just the elements), and leaves the list unchanged for an out-of-range index.
func TestMutableList_InsertParity(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			// Insert is immutable: it returns the resulting slice without
			// mutating the receiver. InsertInPlace performs the same logical
			// edit on the receiver. Both must agree for every index below.
			type insertCase struct {
				name  string
				seed  []int
				index int
				elems []int
				want  []int
			}
			cases := []insertCase{
				{"at head", []int{1, 2, 3}, 0, []int{8, 9}, []int{8, 9, 1, 2, 3}},
				{"in middle", []int{1, 2, 3}, 1, []int{8, 9}, []int{1, 8, 9, 2, 3}},
				{"index equal to length appends", []int{1, 2, 3}, 3, []int{8, 9}, []int{1, 2, 3, 8, 9}},
				{"into empty list appends", nil, 0, []int{8, 9}, []int{8, 9}},
				{"index beyond length unchanged", []int{1, 2, 3}, 99, []int{8, 9}, []int{1, 2, 3}},
				{"negative index unchanged", []int{1, 2, 3}, -1, []int{8, 9}, []int{1, 2, 3}},
			}

			for _, tc := range cases {
				t.Run("Insert/"+tc.name, func(t *testing.T) {
					l := ctor.make(tc.seed...)
					got := l.Insert(tc.index, tc.elems...)
					if !equalInts(got, tc.want) {
						t.Errorf("Insert(%d, %v) = %v, want %v", tc.index, tc.elems, got, tc.want)
					}
					// Immutable: the receiver must be untouched.
					if rcv := l.AsSlice(); !equalInts(rcv, tc.seed) {
						t.Errorf("Insert mutated receiver: got %v, want %v", rcv, tc.seed)
					}
				})

				t.Run("InsertInPlace/"+tc.name, func(t *testing.T) {
					l := ctor.make(tc.seed...)
					l.InsertInPlace(tc.index, tc.elems...)
					if got := l.AsSlice(); !equalInts(got, tc.want) {
						t.Errorf("InsertInPlace(%d, %v) = %v, want %v", tc.index, tc.elems, got, tc.want)
					}
				})
			}
		})
	}
}

// equalInts compares two int slices by content, treating nil and empty as
// equal. The list implementations legitimately differ on whether an unchanged
// result is returned as nil or a non-nil empty slice, so the parity tests
// compare contents rather than exact slice identity.
func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// circularConstructor builds a circular linked list seeded with the given
// elements. Removal must preserve the circular invariant.
type circularConstructor struct {
	name string
	make func(...int) lists.MutableList[int]
}

func allCircularConstructors() []circularConstructor {
	return []circularConstructor{
		{"LinkedCircular", func(e ...int) lists.MutableList[int] { return lists.NewLinkedCircular(e...) }},
		{"DoublyLinkedCircular", func(e ...int) lists.MutableList[int] { return lists.NewDoublyLinkedCircular(e...) }},
		{"ConcurrentLinkedCircular", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentLinkedCircular(e...) }},
		{"ConcurrentDoublyLinkedCircular", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentDoublyLinkedCircular(e...) }},
		{"ConcurrentRWLinkedCircular", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentRWLinkedCircular(e...) }},
		{"ConcurrentRWDoublyLinkedCircular", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentRWDoublyLinkedCircular(e...) }},
	}
}

func TestCircularList_RemovePreservesIntegrity(t *testing.T) {
	for _, ctor := range allCircularConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			// Remove the head.
			l := ctor.make(1, 2, 3)
			if _, ok := l.RemoveAtInPlace(0); !ok {
				t.Fatalf("RemoveAtInPlace(0) failed")
			}
			if got := l.AsSlice(); !reflect.DeepEqual(got, []int{2, 3}) {
				t.Errorf("after head removal = %v, want [2 3]", got)
			}

			// Remove the tail.
			l = ctor.make(1, 2, 3)
			if _, ok := l.RemoveAtInPlace(2); !ok {
				t.Fatalf("RemoveAtInPlace(2) failed")
			}
			if got := l.AsSlice(); !reflect.DeepEqual(got, []int{1, 2}) {
				t.Errorf("after tail removal = %v, want [1 2]", got)
			}

			// Remove a middle element by value.
			l = ctor.make(1, 2, 3)
			if !l.RemoveInPlace(2) {
				t.Fatalf("RemoveInPlace(2) failed")
			}
			if got := l.AsSlice(); !reflect.DeepEqual(got, []int{1, 3}) {
				t.Errorf("after middle removal = %v, want [1 3]", got)
			}
		})
	}
}

// TestRemove_DeepEqualSemantics confirms value-based removal works for [T any]
// element types that are not comparable with ==, using reflect.DeepEqual.
func TestRemove_DeepEqualSemantics(t *testing.T) {
	l := lists.NewArray([]int{1, 2}, []int{3, 4}, []int{5, 6})

	if !l.RemoveInPlace([]int{3, 4}) {
		t.Errorf("RemoveInPlace([3 4]) = false, want true")
	}
	got := l.AsSlice()
	want := [][]int{{1, 2}, {5, 6}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("after RemoveInPlace = %v, want %v", got, want)
	}
}

func ExampleComparable_Contains() {
	l := lists.NewComparable("apple", "banana", "cherry")
	fmt.Println(l.Contains("banana"))
	fmt.Println(l.Contains("durian"))
	fmt.Println(l.IndexOf("cherry"))
	// Output:
	// true
	// false
	// 2
}

func TestComparable_ContainsAndIndexOf(t *testing.T) {
	l := lists.NewComparable(10, 20, 30)

	if !l.Contains(20) {
		t.Errorf("Contains(20) = false, want true")
	}
	if l.Contains(99) {
		t.Errorf("Contains(99) = true, want false")
	}
	if got := l.IndexOf(30); got != 2 {
		t.Errorf("IndexOf(30) = %d, want 2", got)
	}
	if got := l.IndexOf(99); got != -1 {
		t.Errorf("IndexOf(99) = %d, want -1", got)
	}

	// The embedded list API is available and stays consistent with the queries.
	l.PushInPlace(40)
	if !l.Contains(40) {
		t.Errorf("Contains(40) = false after PushInPlace, want true")
	}
}

func TestComparable_WrapsAnyMutableList(t *testing.T) {
	// NewComparableFrom shares the wrapped list, including concurrent ones.
	backing := lists.NewConcurrentArray(1, 2, 3)
	l := lists.NewComparableFrom[int](backing)

	if !l.Contains(2) {
		t.Errorf("Contains(2) = false, want true")
	}
	// Mutations through the wrapper are visible on the backing list.
	l.RemoveInPlace(2)
	if backing.Length() != 2 {
		t.Errorf("expected backing length 2 after RemoveInPlace, got %d", backing.Length())
	}
	if l.Contains(2) {
		t.Errorf("Contains(2) = true after removal, want false")
	}
}
