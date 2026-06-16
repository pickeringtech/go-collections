package heaps_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/heaps"
	"github.com/pickeringtech/go-collections/slices"
)

// sortedCopy returns the values sorted ascending as a non-nil slice, serving as
// the min-heap priority oracle.
func sortedCopy(values []int) []int {
	out := make([]int, len(values))
	copy(out, values)
	sort.Ints(out)
	return out
}

// reversed returns a non-nil reverse of the given slice, serving as the
// max-heap priority oracle when applied to the ascending sort.
func reversed(values []int) []int {
	out := make([]int, len(values))
	for i, v := range values {
		out[len(values)-1-i] = v
	}
	return out
}

// drainInPlace repeatedly PopInPlace-s a mutable heap into a slice.
func drainInPlace(h heaps.MutableHeap[int]) []int {
	out := make([]int, 0, h.Length())
	for {
		v, ok := h.PopInPlace()
		if !ok {
			return out
		}
		out = append(out, v)
	}
}

func TestBinary_DrainsInPriorityOrder(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "nil input drains empty", args: args{values: nil}},
		{name: "empty input drains empty", args: args{values: []int{}}},
		{name: "single element", args: args{values: []int{42}}},
		{name: "already ascending", args: args{values: []int{1, 2, 3, 4, 5}}},
		{name: "already descending", args: args{values: []int{5, 4, 3, 2, 1}}},
		{name: "shuffled", args: args{values: []int{3, 1, 4, 1, 5, 9, 2, 6}}},
		{name: "duplicates collapse to multiset", args: args{values: []int{3, 1, 3, 1, 2, 2}}},
		{name: "negatives and zero", args: args{values: []int{-5, 3, -1, 0, 8, -8}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minWant := sortedCopy(tt.args.values)
			maxWant := reversed(minWant)

			// AsSortedSlice is the non-mutating priority drain.
			if got := heaps.NewMin(tt.args.values...).AsSortedSlice(); !reflect.DeepEqual(got, minWant) {
				t.Errorf("NewMin AsSortedSlice() = %v, want %v", got, minWant)
			}
			if got := heaps.NewMax(tt.args.values...).AsSortedSlice(); !reflect.DeepEqual(got, maxWant) {
				t.Errorf("NewMax AsSortedSlice() = %v, want %v", got, maxWant)
			}

			// PopInPlace must agree with AsSortedSlice.
			if got := drainInPlace(heaps.NewMin(tt.args.values...)); !reflect.DeepEqual(got, minWant) {
				t.Errorf("NewMin drainInPlace = %v, want %v", got, minWant)
			}
			if got := drainInPlace(heaps.NewMax(tt.args.values...)); !reflect.DeepEqual(got, maxWant) {
				t.Errorf("NewMax drainInPlace = %v, want %v", got, maxWant)
			}

			// Drain yields the same priority order as AsSortedSlice.
			gotDrain := make([]int, 0, len(minWant))
			for v := range heaps.NewMin(tt.args.values...).Drain() {
				gotDrain = append(gotDrain, v)
			}
			if !reflect.DeepEqual(gotDrain, minWant) {
				t.Errorf("NewMin Drain() = %v, want %v", gotDrain, minWant)
			}
		})
	}
}

func TestBinary_PushInPlaceBuildsSameOrderAsHeapify(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "nil input", args: args{values: nil}},
		{name: "single", args: args{values: []int{7}}},
		{name: "many", args: args{values: []int{9, 2, 7, 1, 8, 3, 6, 4, 5}}},
		{name: "duplicates", args: args{values: []int{2, 2, 1, 1, 3, 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := sortedCopy(tt.args.values)

			incremental := heaps.NewMin[int]()
			for _, v := range tt.args.values {
				incremental.PushInPlace(v)
			}
			if got := incremental.AsSortedSlice(); !reflect.DeepEqual(got, want) {
				t.Errorf("incremental PushInPlace AsSortedSlice() = %v, want %v", got, want)
			}

			batch := heaps.NewMin[int]()
			batch.PushManyInPlace(tt.args.values...)
			if got := batch.AsSortedSlice(); !reflect.DeepEqual(got, want) {
				t.Errorf("PushManyInPlace AsSortedSlice() = %v, want %v", got, want)
			}
		})
	}
}

func TestBinary_PeekAndEmptyState(t *testing.T) {
	empty := heaps.NewMin[int]()
	if v, ok := empty.Peek(); ok || v != 0 {
		t.Errorf("empty Peek() = (%d, %v), want (0, false)", v, ok)
	}
	if v, ok := empty.PopInPlace(); ok || v != 0 {
		t.Errorf("empty PopInPlace() = (%d, %v), want (0, false)", v, ok)
	}
	if !empty.IsEmpty() {
		t.Error("empty IsEmpty() = false, want true")
	}
	if empty.Length() != 0 {
		t.Errorf("empty Length() = %d, want 0", empty.Length())
	}
	if got := empty.AsSlice(); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("empty AsSlice() = %v, want []int{}", got)
	}
	if got := empty.AsSortedSlice(); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("empty AsSortedSlice() = %v, want []int{}", got)
	}

	h := heaps.NewMin(4, 2, 6)
	if v, ok := h.Peek(); !ok || v != 2 {
		t.Errorf("Peek() = (%d, %v), want (2, true)", v, ok)
	}
	if h.IsEmpty() {
		t.Error("non-empty IsEmpty() = true, want false")
	}
	if h.Length() != 3 {
		t.Errorf("Length() = %d, want 3", h.Length())
	}
}

func TestBinary_ImmutableOpsLeaveReceiverUnchanged(t *testing.T) {
	original := heaps.NewMin(5, 3, 8)

	pushed := original.Push(1)
	if v, _ := pushed.Peek(); v != 1 {
		t.Errorf("Push result Peek() = %d, want 1", v)
	}
	if original.Length() != 3 {
		t.Errorf("Push mutated receiver: Length() = %d, want 3", original.Length())
	}
	if v, _ := original.Peek(); v != 3 {
		t.Errorf("Push mutated receiver: Peek() = %d, want 3", v)
	}

	multi := original.PushMany(0, 9, -1)
	if v, _ := multi.Peek(); v != -1 {
		t.Errorf("PushMany result Peek() = %d, want -1", v)
	}
	if original.Length() != 3 {
		t.Errorf("PushMany mutated receiver: Length() = %d, want 3", original.Length())
	}

	v, ok, rest := original.Pop()
	if !ok || v != 3 {
		t.Errorf("Pop() = (%d, %v), want (3, true)", v, ok)
	}
	if rest.Length() != 2 {
		t.Errorf("Pop rest Length() = %d, want 2", rest.Length())
	}
	if original.Length() != 3 {
		t.Errorf("Pop mutated receiver: Length() = %d, want 3", original.Length())
	}

	// Pop on an empty heap returns the zero value, false, and an empty heap.
	emptyV, emptyOK, emptyRest := heaps.NewMin[int]().Pop()
	if emptyOK || emptyV != 0 || emptyRest.Length() != 0 {
		t.Errorf("empty Pop() = (%d, %v, len=%d), want (0, false, len=0)", emptyV, emptyOK, emptyRest.Length())
	}
}

func TestBinary_NewCopiesInputSlice(t *testing.T) {
	src := []int{5, 1, 3}
	h := heaps.NewMin(src...)
	src[0] = -100 // mutating the source must not affect the heap
	if v, _ := h.Peek(); v != 1 {
		t.Errorf("Peek() = %d, want 1 (heap aliased its input)", v)
	}
}

func TestBinary_ForEachAndAllVisitEveryElement(t *testing.T) {
	values := []int{3, 1, 4, 1, 5}
	want := sortedCopy(values)
	h := heaps.NewMin(values...)

	var each []int
	h.ForEach(func(e int) { each = append(each, e) })
	sort.Ints(each)
	if !reflect.DeepEqual(each, want) {
		t.Errorf("ForEach visited %v, want multiset %v", each, want)
	}

	var all []int
	for v := range h.All() {
		all = append(all, v)
	}
	sort.Ints(all)
	if !reflect.DeepEqual(all, want) {
		t.Errorf("All() visited %v, want multiset %v", all, want)
	}

	if got := h.AsSlice(); len(got) != len(values) {
		t.Errorf("AsSlice() length = %d, want %d", len(got), len(values))
	}
}

func TestBinary_IteratorEarlyBreak(t *testing.T) {
	h := heaps.NewMin(5, 1, 3, 2, 4)

	count := 0
	for range h.All() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("All() early break visited %d, want 1", count)
	}

	// Drain stops cleanly when the consumer breaks, leaving the heap intact.
	first := 0
	got := 0
	for v := range h.Drain() {
		first = v
		got++
		break
	}
	if got != 1 || first != 1 {
		t.Errorf("Drain() early break = (count %d, first %d), want (1, 1)", got, first)
	}
	if h.Length() != 5 {
		t.Errorf("Drain() mutated receiver: Length() = %d, want 5", h.Length())
	}
}

func TestBinary_CustomComparatorOnStruct(t *testing.T) {
	type task struct {
		name     string
		priority int
	}
	byPriorityDesc := func(a, b task) bool { return a.priority > b.priority }
	h := heaps.New(byPriorityDesc,
		task{"low", 1}, task{"high", 10}, task{"mid", 5})

	top, ok := h.Peek()
	if !ok || top.name != "high" {
		t.Errorf("Peek() = %+v, want the high-priority task", top)
	}
}

func benchmarkInputs() []struct {
	name string
	sli  []int
} {
	return []struct {
		name string
		sli  []int
	}{
		{name: "3", sli: []int{3, 1, 2}},
		{name: "10", sli: slices.Generate(10, func(i int) int { return 10 - i })},
		{name: "100", sli: slices.Generate(100, func(i int) int { return 100 - i })},
		{name: "1_000", sli: slices.Generate(1_000, func(i int) int { return 1_000 - i })},
		{name: "10_000", sli: slices.Generate(10_000, func(i int) int { return 10_000 - i })},
		{name: "100_000", sli: slices.Generate(100_000, func(i int) int { return 100_000 - i })},
		{name: "1_000_000", sli: slices.Generate(1_000_000, func(i int) int { return 1_000_000 - i })},
	}
}

func BenchmarkBinary_New(b *testing.B) {
	for _, bm := range benchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = heaps.NewMin(bm.sli...)
			}
		})
	}
}

func BenchmarkBinary_PushInPlace(b *testing.B) {
	for _, bm := range benchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				h := heaps.NewMin[int]()
				h.PushManyInPlace(bm.sli...)
			}
		})
	}
}

func BenchmarkBinary_PopInPlace(b *testing.B) {
	for _, bm := range benchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				h := heaps.NewMin(bm.sli...)
				for {
					if _, ok := h.PopInPlace(); !ok {
						break
					}
				}
			}
		})
	}
}
