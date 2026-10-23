package lists_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// TestMutableList_IteratorsSnapshotWhenRangingBegins pins the Iterable contract:
// an iterator taken before a mutation and ranged after it sees the mutation, and
// a loop body that modifies the list does not disturb the elements yielded.
func TestMutableList_IteratorsSnapshotWhenRangingBegins(t *testing.T) {
	for _, c := range everyMutableListConstructor() {
		t.Run(c.name, func(t *testing.T) {
			l := c.make(1, 2)
			all, values, backward := l.All(), l.Values(), l.Backward()
			l.PushInPlace(3)

			var gotAll, gotValues, gotBackward []int
			for _, v := range all {
				gotAll = append(gotAll, v)
			}
			for v := range values {
				gotValues = append(gotValues, v)
			}
			for _, v := range backward {
				gotBackward = append(gotBackward, v)
			}
			if !reflect.DeepEqual(gotAll, []int{1, 2, 3}) {
				t.Errorf("All() ranged after PushInPlace = %v, want [1 2 3]", gotAll)
			}
			if !reflect.DeepEqual(gotValues, []int{1, 2, 3}) {
				t.Errorf("Values() ranged after PushInPlace = %v, want [1 2 3]", gotValues)
			}
			if !reflect.DeepEqual(gotBackward, []int{3, 2, 1}) {
				t.Errorf("Backward() ranged after PushInPlace = %v, want [3 2 1]", gotBackward)
			}

			var drained []int
			for v := range l.Values() {
				drained = append(drained, v)
				l.PopInPlace()
			}
			if !reflect.DeepEqual(drained, []int{1, 2, 3}) {
				t.Errorf("Values() while popping in the loop body = %v, want [1 2 3]", drained)
			}
			if !l.IsEmpty() {
				t.Errorf("list after popping once per element = %v, want empty", l.AsSlice())
			}
		})
	}
}

// TestArray_ForEachNeverVisitsARemovedSlot guards the Array's live walk: popping
// from inside the callback must not surface the vacated (zeroed) slot as a
// phantom zero element.
func TestArray_ForEachNeverVisitsARemovedSlot(t *testing.T) {
	tests := []struct {
		name string
		walk func(a *lists.Array[int], visit func(int))
	}{
		{
			name: "ForEach",
			walk: func(a *lists.Array[int], visit func(int)) {
				a.ForEach(func(v int) { visit(v) })
			},
		},
		{
			name: "ForEachWithIndex",
			walk: func(a *lists.Array[int], visit func(int)) {
				a.ForEachWithIndex(func(_ int, v int) { visit(v) })
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := lists.NewArray(1, 2, 3, 4)
			var visited []int
			tt.walk(a, func(v int) {
				visited = append(visited, v)
				a.PopInPlace()
			})
			want := []int{1, 2}
			if !reflect.DeepEqual(visited, want) {
				t.Errorf("visited while popping = %v, want %v", visited, want)
			}
		})
	}
}

// TestArray_InPlaceEndOperationsDoNotCopy guards against the in-place stack and
// queue operations degrading to O(n) by copying the backing array: once the
// array has capacity, a push/pop or enqueue/dequeue round-trip allocates nothing.
func TestArray_InPlaceEndOperationsDoNotCopy(t *testing.T) {
	tests := []struct {
		name      string
		roundTrip func(a *lists.Array[int])
	}{
		{
			name: "PushInPlace then PopInPlace",
			roundTrip: func(a *lists.Array[int]) {
				a.PushInPlace(1)
				a.PopInPlace()
			},
		},
		{
			name: "EnqueueInPlace then DequeueInPlace",
			roundTrip: func(a *lists.Array[int]) {
				a.EnqueueInPlace(1)
				a.DequeueInPlace()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := lists.NewArray(make([]int, 1000)...)
			a.PushInPlace(0) // grow once so spare capacity exists
			a.PopInPlace()
			got := testing.AllocsPerRun(100, func() { tt.roundTrip(a) })
			if got != 0 {
				t.Errorf("allocations per round-trip = %v, want 0", got)
			}
		})
	}
}
