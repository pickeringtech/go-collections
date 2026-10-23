package lists_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// filterFactory builds a MutableList[float64] seeded with the given elements.
type filterFactory struct {
	name string
	make func(elements ...float64) lists.MutableList[float64]
}

// filterParityFactories returns the non-concurrent lists alongside every
// concurrent list (including the circular constructors), so FilterInPlace can
// be checked for parity between each concurrent type and its sibling.
func filterParityFactories() []filterFactory {
	return []filterFactory{
		{"Array", func(e ...float64) lists.MutableList[float64] { return lists.NewArray(e...) }},
		{"Linked", func(e ...float64) lists.MutableList[float64] { return lists.NewLinked(e...) }},
		{"DoublyLinked", func(e ...float64) lists.MutableList[float64] { return lists.NewDoublyLinked(e...) }},
		{"ConcurrentArray", func(e ...float64) lists.MutableList[float64] { return lists.NewConcurrentArray(e...) }},
		{"ConcurrentRWArray", func(e ...float64) lists.MutableList[float64] { return lists.NewConcurrentRWArray(e...) }},
		{"ConcurrentLinked", func(e ...float64) lists.MutableList[float64] { return lists.NewConcurrentLinked(e...) }},
		{"ConcurrentLinkedCircular", func(e ...float64) lists.MutableList[float64] {
			return lists.NewConcurrentLinkedCircular(e...)
		}},
		{"ConcurrentRWLinked", func(e ...float64) lists.MutableList[float64] { return lists.NewConcurrentRWLinked(e...) }},
		{"ConcurrentRWLinkedCircular", func(e ...float64) lists.MutableList[float64] {
			return lists.NewConcurrentRWLinkedCircular(e...)
		}},
		{"ConcurrentDoublyLinked", func(e ...float64) lists.MutableList[float64] {
			return lists.NewConcurrentDoublyLinked(e...)
		}},
		{"ConcurrentDoublyLinkedCircular", func(e ...float64) lists.MutableList[float64] {
			return lists.NewConcurrentDoublyLinkedCircular(e...)
		}},
		{"ConcurrentRWDoublyLinked", func(e ...float64) lists.MutableList[float64] {
			return lists.NewConcurrentRWDoublyLinked(e...)
		}},
		{"ConcurrentRWDoublyLinkedCircular", func(e ...float64) lists.MutableList[float64] {
			return lists.NewConcurrentRWDoublyLinkedCircular(e...)
		}},
	}
}

// sameFloats reports whether a and b hold the same elements in the same order,
// treating NaN as equal to NaN and nil as equal to empty.
func sameFloats(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		bothNaN := math.IsNaN(a[i]) && math.IsNaN(b[i])
		if a[i] != b[i] && !bothNaN {
			return false
		}
	}
	return true
}

// TestListFilterInPlaceConcurrentMatchesSibling verifies that, with no
// concurrent write during the predicate window, every concurrent list's
// FilterInPlace keeps exactly the elements the predicate accepted, in their
// original positions — the same result as Array, Linked and DoublyLinked — and
// calls the predicate once per element, in order. Removal used to be a multiset
// diff that deleted the first deeply-equal occurrence of each rejected value,
// which never removed NaN and broke order-dependent predicates.
func TestListFilterInPlaceConcurrentMatchesSibling(t *testing.T) {
	type args struct {
		elements []float64
		newKeep  func() func(element float64) bool
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "an order-dependent predicate keeps the first occurrence",
			args: args{elements: []float64{1, 2, 1}, newKeep: func() func(float64) bool {
				seen := map[float64]bool{}
				return func(e float64) bool {
					if seen[e] {
						return false
					}
					seen[e] = true
					return true
				}
			}},
			want: []float64{1, 2},
		},
		{
			name: "a rejected NaN is removed",
			args: args{elements: []float64{math.NaN(), 1}, newKeep: func() func(float64) bool {
				return func(e float64) bool { return !math.IsNaN(e) }
			}},
			want: []float64{1},
		},
		{
			name: "rejecting every element empties the list",
			args: args{elements: []float64{math.NaN(), 1, 2}, newKeep: func() func(float64) bool {
				return func(float64) bool { return false }
			}},
			want: []float64{},
		},
		{
			name: "keeping every element leaves the list unchanged",
			args: args{elements: []float64{3, math.NaN(), 1}, newKeep: func() func(float64) bool {
				return func(float64) bool { return true }
			}},
			want: []float64{3, math.NaN(), 1},
		},
		{
			name: "an empty list stays empty",
			args: args{elements: nil, newKeep: func() func(float64) bool {
				return func(float64) bool { return false }
			}},
			want: []float64{},
		},
	}
	for _, tt := range tests {
		for _, f := range filterParityFactories() {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				l := f.make(tt.args.elements...)
				keep := tt.args.newKeep()
				var calls []float64
				l.FilterInPlace(func(e float64) bool {
					calls = append(calls, e)
					return keep(e)
				})

				got := l.AsSlice()
				if !sameFloats(got, tt.want) {
					t.Errorf("FilterInPlace() = %v, want %v", got, tt.want)
				}
				if !sameFloats(calls, tt.args.elements) {
					t.Errorf("FilterInPlace() predicate calls = %v, want %v", calls, tt.args.elements)
				}
			})
		}
	}
}

// TestConcurrentListFilterInPlacePreservesConcurrentWrite pins the issue #153
// contract for every mutating method: a write that lands while the predicate
// runs outside the lock must survive the apply phase. The predicate performs
// the racing write itself on its first call — it runs outside the lock, which
// makes the otherwise-racy window deterministic. Every case rejects 2 from
// [1 2 3], so a mutating method that failed to register its write would let the
// apply phase overwrite the list with the stale filtered snapshot [1 3].
func TestConcurrentListFilterInPlacePreservesConcurrentWrite(t *testing.T) {
	type args struct {
		write func(l lists.MutableList[int])
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "PushInPlace keeps the pushed element",
			args: args{write: func(l lists.MutableList[int]) { l.PushInPlace(4) }},
			want: []int{1, 3, 4},
		},
		{
			name: "EnqueueInPlace keeps the enqueued element",
			args: args{write: func(l lists.MutableList[int]) { l.EnqueueInPlace(4) }},
			want: []int{1, 3, 4},
		},
		{
			name: "InsertInPlace keeps the inserted element",
			args: args{write: func(l lists.MutableList[int]) { l.InsertInPlace(0, 9) }},
			want: []int{9, 1, 3},
		},
		{
			name: "PopInPlace stays popped",
			args: args{write: func(l lists.MutableList[int]) { l.PopInPlace() }},
			want: []int{1},
		},
		{
			name: "DequeueInPlace stays dequeued",
			args: args{write: func(l lists.MutableList[int]) { l.DequeueInPlace() }},
			want: []int{3},
		},
		{
			name: "RemoveAtInPlace stays removed",
			args: args{write: func(l lists.MutableList[int]) { l.RemoveAtInPlace(0) }},
			want: []int{3},
		},
		{
			name: "RemoveInPlace stays removed",
			args: args{write: func(l lists.MutableList[int]) { l.RemoveInPlace(3) }},
			want: []int{1},
		},
		{
			name: "Clear stays cleared",
			args: args{write: func(l lists.MutableList[int]) { l.Clear() }},
			want: []int{},
		},
		{
			name: "SortInPlace keeps the new order",
			args: args{write: func(l lists.MutableList[int]) { l.SortInPlace(func(a, b int) bool { return a > b }) }},
			want: []int{3, 1},
		},
		{
			name: "a nested FilterInPlace is applied too",
			args: args{write: func(l lists.MutableList[int]) { l.FilterInPlace(func(e int) bool { return e < 3 }) }},
			want: []int{1},
		},
	}
	for _, tt := range tests {
		for _, f := range concurrentWriteFactories() {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				l := f.make(1, 2, 3)
				wrote := false
				l.FilterInPlace(func(e int) bool {
					if !wrote {
						wrote = true
						tt.args.write(l) // racing write lands in the evaluation window
					}
					return e != 2
				})

				got := l.AsSlice()
				if got == nil {
					got = []int{} // the linked lists report an empty list as nil
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FilterInPlace() with a concurrent write = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

// TestConcurrentListFilterInPlacePreservesWritesThroughTwoMethods extends the
// issue #153 contract to two writes through different methods in the same
// predicate window. FilterInPlace detects a concurrent write by comparing a
// write counter before and after the predicate runs, so every mutating method
// must move that counter the same way: a method that moved it backwards would
// cancel a companion method's write (an ABA on the counter), the list would
// look untouched, and the apply phase would overwrite both writes with the
// stale filtered snapshot [1 3]. Each case pairs the method under test with a
// PushInPlace(4) companion (EnqueueInPlace(5) when PushInPlace is the method
// under test); a nested FilterInPlace covers the counter bump FilterInPlace
// makes itself.
func TestConcurrentListFilterInPlacePreservesWritesThroughTwoMethods(t *testing.T) {
	pushFour := func(l lists.MutableList[int]) { l.PushInPlace(4) }
	type args struct {
		write     func(l lists.MutableList[int])
		companion func(l lists.MutableList[int])
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "PushInPlace then EnqueueInPlace keeps both elements",
			args: args{
				write:     func(l lists.MutableList[int]) { l.PushInPlace(4) },
				companion: func(l lists.MutableList[int]) { l.EnqueueInPlace(5) },
			},
			want: []int{1, 3, 4, 5},
		},
		{
			name: "EnqueueInPlace then PushInPlace keeps both elements",
			args: args{write: func(l lists.MutableList[int]) { l.EnqueueInPlace(5) }, companion: pushFour},
			want: []int{1, 3, 5, 4},
		},
		{
			name: "InsertInPlace then PushInPlace keeps both elements",
			args: args{write: func(l lists.MutableList[int]) { l.InsertInPlace(0, 9) }, companion: pushFour},
			want: []int{9, 1, 3, 4},
		},
		{
			name: "PopInPlace then PushInPlace keeps the pop and the push",
			args: args{write: func(l lists.MutableList[int]) { l.PopInPlace() }, companion: pushFour},
			want: []int{1, 4},
		},
		{
			name: "DequeueInPlace then PushInPlace keeps the dequeue and the push",
			args: args{write: func(l lists.MutableList[int]) { l.DequeueInPlace() }, companion: pushFour},
			want: []int{3, 4},
		},
		{
			name: "RemoveAtInPlace then PushInPlace keeps the removal and the push",
			args: args{write: func(l lists.MutableList[int]) { l.RemoveAtInPlace(0) }, companion: pushFour},
			want: []int{3, 4},
		},
		{
			name: "RemoveInPlace then PushInPlace keeps the removal and the push",
			args: args{write: func(l lists.MutableList[int]) { l.RemoveInPlace(3) }, companion: pushFour},
			want: []int{1, 4},
		},
		{
			name: "Clear then PushInPlace keeps only the pushed element",
			args: args{write: func(l lists.MutableList[int]) { l.Clear() }, companion: pushFour},
			want: []int{4},
		},
		{
			name: "SortInPlace then PushInPlace keeps the new order and the push",
			args: args{
				write:     func(l lists.MutableList[int]) { l.SortInPlace(func(a, b int) bool { return a > b }) },
				companion: pushFour,
			},
			want: []int{3, 1, 4},
		},
		{
			name: "a nested FilterInPlace then PushInPlace keeps both writes",
			args: args{
				write:     func(l lists.MutableList[int]) { l.FilterInPlace(func(e int) bool { return e < 3 }) },
				companion: pushFour,
			},
			want: []int{1, 4},
		},
	}
	for _, tt := range tests {
		for _, f := range concurrentWriteFactories() {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				l := f.make(1, 2, 3)
				wrote := false
				l.FilterInPlace(func(e int) bool {
					if !wrote {
						wrote = true
						tt.args.write(l) // both racing writes land in the evaluation window
						tt.args.companion(l)
					}
					return e != 2
				})

				got := l.AsSlice()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FilterInPlace() with two concurrent writes = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

// concurrentWriteFactories builds every concurrent list, including the
// circular constructors.
func concurrentWriteFactories() []struct {
	name string
	make func(elements ...int) lists.MutableList[int]
} {
	return []struct {
		name string
		make func(elements ...int) lists.MutableList[int]
	}{
		{"ConcurrentArray", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentArray(e...) }},
		{"ConcurrentRWArray", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentRWArray(e...) }},
		{"ConcurrentLinked", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentLinked(e...) }},
		{"ConcurrentLinkedCircular", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentLinkedCircular(e...) }},
		{"ConcurrentRWLinked", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentRWLinked(e...) }},
		{"ConcurrentRWLinkedCircular", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentRWLinkedCircular(e...) }},
		{"ConcurrentDoublyLinked", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentDoublyLinked(e...) }},
		{"ConcurrentDoublyLinkedCircular", func(e ...int) lists.MutableList[int] {
			return lists.NewConcurrentDoublyLinkedCircular(e...)
		}},
		{"ConcurrentRWDoublyLinked", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentRWDoublyLinked(e...) }},
		{"ConcurrentRWDoublyLinkedCircular", func(e ...int) lists.MutableList[int] {
			return lists.NewConcurrentRWDoublyLinkedCircular(e...)
		}},
	}
}
