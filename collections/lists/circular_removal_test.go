package lists_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// drainedResult records one removal: the value returned and whether the list
// reported an element as present.
type drainedResult struct {
	value int
	ok    bool
}

// TestCircularList_DrainEmptiesList drains each circular list from both ends and
// asserts every element comes out exactly once, followed by (0, false). Draining
// is bounded so a removal that never shrinks the list fails rather than hangs.
func TestCircularList_DrainEmptiesList(t *testing.T) {
	type args struct {
		initial []int
		remove  func(l lists.MutableList[int]) (int, bool)
	}
	dequeue := func(l lists.MutableList[int]) (int, bool) { return l.DequeueInPlace() }
	pop := func(l lists.MutableList[int]) (int, bool) { return l.PopInPlace() }
	tests := []struct {
		name string
		args args
		want []drainedResult
	}{
		{
			name: "dequeue from a single-element list removes it",
			args: args{initial: []int{42}, remove: dequeue},
			want: []drainedResult{{42, true}, {0, false}},
		},
		{
			name: "pop from a single-element list removes it",
			args: args{initial: []int{42}, remove: pop},
			want: []drainedResult{{42, true}, {0, false}},
		},
		{
			name: "dequeue drains front to back",
			args: args{initial: []int{1, 2, 3}, remove: dequeue},
			want: []drainedResult{{1, true}, {2, true}, {3, true}, {0, false}},
		},
		{
			name: "pop drains back to front",
			args: args{initial: []int{1, 2, 3}, remove: pop},
			want: []drainedResult{{3, true}, {2, true}, {1, true}, {0, false}},
		},
		{
			name: "removing from an empty list reports not found",
			args: args{initial: nil, remove: dequeue},
			want: []drainedResult{{0, false}},
		},
	}
	for _, ctor := range allCircularConstructors() {
		for _, tt := range tests {
			t.Run(ctor.name+"/"+tt.name, func(t *testing.T) {
				l := ctor.make(tt.args.initial...)
				got := make([]drainedResult, 0, len(tt.want))
				for range tt.want {
					value, ok := tt.args.remove(l)
					got = append(got, drainedResult{value, ok})
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("removals = %v, want %v", got, tt.want)
				}
				if l.Length() != 0 {
					t.Errorf("Length() after draining = %d, want 0", l.Length())
				}
				if !l.IsEmpty() {
					t.Errorf("IsEmpty() after draining = false, want true")
				}
			})
		}
	}
}

// TestCircularList_SingleElementRemovalPaths removes the only element of a
// one-element circular list through every removal path, then asserts the list
// is empty and still usable: a fresh push must be the sole element and must
// itself be removable.
func TestCircularList_SingleElementRemovalPaths(t *testing.T) {
	tests := []struct {
		name   string
		remove func(l lists.MutableList[int])
	}{
		{name: "DequeueInPlace", remove: func(l lists.MutableList[int]) { l.DequeueInPlace() }},
		{name: "PopInPlace", remove: func(l lists.MutableList[int]) { l.PopInPlace() }},
		{name: "RemoveAtInPlace", remove: func(l lists.MutableList[int]) { l.RemoveAtInPlace(0) }},
		{name: "RemoveInPlace", remove: func(l lists.MutableList[int]) { l.RemoveInPlace(42) }},
		{name: "FilterInPlace", remove: func(l lists.MutableList[int]) { l.FilterInPlace(func(int) bool { return false }) }},
		{name: "Clear", remove: func(l lists.MutableList[int]) { l.Clear() }},
	}
	for _, ctor := range allCircularConstructors() {
		for _, tt := range tests {
			t.Run(ctor.name+"/"+tt.name, func(t *testing.T) {
				l := ctor.make(42)
				tt.remove(l)
				got := l.AsSlice()
				if !reflect.DeepEqual(got, []int{}) {
					t.Fatalf("AsSlice() after %s = %v, want []", tt.name, got)
				}

				l.PushInPlace(7)
				got = l.AsSlice()
				if !reflect.DeepEqual(got, []int{7}) {
					t.Fatalf("AsSlice() after re-push = %v, want [7]", got)
				}
				value, ok := l.DequeueInPlace()
				if value != 7 || !ok {
					t.Errorf("DequeueInPlace() after re-push = (%d, %t), want (7, true)", value, ok)
				}
				if l.Length() != 0 {
					t.Errorf("Length() after final dequeue = %d, want 0", l.Length())
				}
			})
		}
	}
}
