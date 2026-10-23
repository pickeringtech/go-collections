package lists_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// TestConcurrentArray_InPlaceStackAndQueueOpsReuseStorage verifies that the
// concurrent arrays' in-place stack and queue operations work on the array's
// own storage instead of rebuilding it: once capacity exists, a PushInPlace +
// PopInPlace pair and an EnqueueInPlace + DequeueInPlace pair allocate nothing.
// Routing them through copying slice helpers would make every push or pop O(n)
// and building or draining a list O(n²); popping the last element used to drop
// the backing array, so the next push reallocated.
func TestConcurrentArray_InPlaceStackAndQueueOpsReuseStorage(t *testing.T) {
	type args struct {
		held int // elements left in the list after warming up its capacity
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "an empty list that has held elements", args: args{held: 0}},
		{name: "a list holding elements", args: args{held: 10}},
	}
	factories := []struct {
		name string
		make func() lists.MutableList[int]
	}{
		{"ConcurrentArray", func() lists.MutableList[int] { return lists.NewConcurrentArray[int]() }},
		{"ConcurrentRWArray", func() lists.MutableList[int] { return lists.NewConcurrentRWArray[int]() }},
	}
	for _, tt := range tests {
		for _, f := range factories {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				l := f.make()
				for i := 0; i < 1024; i++ {
					l.PushInPlace(i)
				}
				for l.Length() > tt.args.held {
					l.PopInPlace() // drain from the back, so the capacity is kept
				}

				pushPop := testing.AllocsPerRun(100, func() {
					l.PushInPlace(1)
					l.PopInPlace()
				})
				if pushPop != 0 {
					t.Errorf("PushInPlace+PopInPlace allocated %v times per pair, want 0", pushPop)
				}
				enqueueDequeue := testing.AllocsPerRun(100, func() {
					l.EnqueueInPlace(1)
					l.DequeueInPlace()
				})
				if enqueueDequeue != 0 {
					t.Errorf("EnqueueInPlace+DequeueInPlace allocated %v times per pair, want 0", enqueueDequeue)
				}
				if l.Length() != tt.args.held {
					t.Errorf("Length() = %d after balanced pairs, want %d", l.Length(), tt.args.held)
				}
			})
		}
	}
}

// TestConcurrentArray_InPlaceRemovalKeepsOrder pins the observable behaviour of
// the in-place removals that now shift the array's own storage.
func TestConcurrentArray_InPlaceRemovalKeepsOrder(t *testing.T) {
	type args struct {
		remove func(l lists.MutableList[int]) (int, bool)
	}
	tests := []struct {
		name     string
		args     args
		wantGot  int
		wantOK   bool
		wantList []int
	}{
		{
			name:     "PopInPlace removes the last element",
			args:     args{remove: func(l lists.MutableList[int]) (int, bool) { return l.PopInPlace() }},
			wantGot:  4,
			wantOK:   true,
			wantList: []int{1, 2, 3},
		},
		{
			name:     "DequeueInPlace removes the first element",
			args:     args{remove: func(l lists.MutableList[int]) (int, bool) { return l.DequeueInPlace() }},
			wantGot:  1,
			wantOK:   true,
			wantList: []int{2, 3, 4},
		},
		{
			name:     "RemoveAtInPlace removes an interior element",
			args:     args{remove: func(l lists.MutableList[int]) (int, bool) { return l.RemoveAtInPlace(1) }},
			wantGot:  2,
			wantOK:   true,
			wantList: []int{1, 3, 4},
		},
		{
			name: "RemoveInPlace removes the first deeply-equal element",
			args: args{remove: func(l lists.MutableList[int]) (int, bool) {
				return 3, l.RemoveInPlace(3)
			}},
			wantGot:  3,
			wantOK:   true,
			wantList: []int{1, 2, 4},
		},
	}
	factories := []struct {
		name string
		make func(elements ...int) lists.MutableList[int]
	}{
		{"ConcurrentArray", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentArray(e...) }},
		{"ConcurrentRWArray", func(e ...int) lists.MutableList[int] { return lists.NewConcurrentRWArray(e...) }},
	}
	for _, tt := range tests {
		for _, f := range factories {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				l := f.make(1, 2, 3, 4)
				got, ok := tt.args.remove(l)
				if got != tt.wantGot || ok != tt.wantOK {
					t.Errorf("removal = (%d, %t), want (%d, %t)", got, ok, tt.wantGot, tt.wantOK)
				}
				l.PushInPlace(9) // reuses the freed slot; must not resurrect anything
				want := append(append([]int{}, tt.wantList...), 9)
				if !intSlicesMatch(l.AsSlice(), want) {
					t.Errorf("AsSlice() after removal and push = %v, want %v", l.AsSlice(), want)
				}
			})
		}
	}
}

// intSlicesMatch reports whether a and b hold the same elements in order.
func intSlicesMatch(a, b []int) bool {
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
