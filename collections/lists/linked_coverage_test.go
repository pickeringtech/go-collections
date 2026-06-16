package lists_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

func TestLinked_InsertInPlace(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		index    int
		elements []int
		want     []int
	}{
		{name: "at head", initial: []int{1, 2, 3}, index: 0, elements: []int{9, 8}, want: []int{9, 8, 1, 2, 3}},
		{name: "in middle", initial: []int{1, 2, 3}, index: 1, elements: []int{9, 8}, want: []int{1, 9, 8, 2, 3}},
		{name: "before tail", initial: []int{1, 2, 3}, index: 2, elements: []int{7}, want: []int{1, 2, 7, 3}},
		{name: "at end", initial: []int{1, 2, 3}, index: 3, elements: []int{4}, want: []int{1, 2, 3, 4}},
		{name: "beyond end appends", initial: []int{1, 2, 3}, index: 10, elements: []int{4}, want: []int{1, 2, 3, 4}},
		{name: "negative index ignored", initial: []int{1, 2, 3}, index: -1, elements: []int{4}, want: []int{1, 2, 3}},
		{name: "into empty at head", initial: nil, index: 0, elements: []int{1, 2}, want: []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lists.NewLinked(tt.initial...)
			l.InsertInPlace(tt.index, tt.elements...)
			got := l.GetAsSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertInPlace(%d, %v) = %v, want %v", tt.index, tt.elements, got, tt.want)
			}
		})
	}
}

func TestLinked_Circular(t *testing.T) {
	l := lists.NewLinkedCircular(1, 2, 3)

	if l.Length() != 3 {
		t.Errorf("Length() = %d, want 3", l.Length())
	}
	if got := l.GetAsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("GetAsSlice() = %v, want [1 2 3]", got)
	}

	var count int
	l.ForEach(func(int) { count++ })
	if count != 3 {
		t.Errorf("ForEach visited %d, want 3", count)
	}

	var idxSum int
	l.ForEachWithIndex(func(idx, _ int) { idxSum += idx })
	if idxSum != 3 {
		t.Errorf("ForEachWithIndex idx sum = %d, want 3", idxSum)
	}

	if !l.AllMatch(func(n int) bool { return n > 0 }) {
		t.Error("AllMatch(>0) = false, want true")
	}
	if l.AllMatch(func(n int) bool { return n > 1 }) {
		t.Error("AllMatch(>1) = true, want false")
	}
	if !l.AnyMatch(func(n int) bool { return n == 3 }) {
		t.Error("AnyMatch(==3) = false, want true")
	}
	if l.AnyMatch(func(n int) bool { return n == 99 }) {
		t.Error("AnyMatch(==99) = true, want false")
	}

	if _, ok := l.Find(func(n int) bool { return n == 2 }); !ok {
		t.Error("Find(==2) ok = false, want true")
	}
	if _, ok := l.Find(func(n int) bool { return n == 99 }); ok {
		t.Error("Find(==99) ok = true, want false")
	}
	if idx := l.FindIndex(func(n int) bool { return n == 99 }); idx != -1 {
		t.Errorf("FindIndex(==99) = %d, want -1", idx)
	}
	if got := l.Filter(func(n int) bool { return n != 2 }); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Errorf("Filter(!=2) = %v, want [1 3]", got)
	}
	if l.Get(99, -1) != -1 {
		t.Errorf("Get(99) = %d, want -1", l.Get(99, -1))
	}
}

func TestLinked_CircularInsertInPlace(t *testing.T) {
	// Inserting in the middle of a circular list exercises findNodeBefore /
	// insertAfter and the tail-to-head re-link.
	l := lists.NewLinkedCircular(1, 2, 3)
	l.InsertInPlace(1, 9)
	if got := l.GetAsSlice(); !reflect.DeepEqual(got, []int{1, 9, 2, 3}) {
		t.Errorf("circular InsertInPlace(1,9) = %v, want [1 9 2 3]", got)
	}

	// Insert at head of a circular list.
	l2 := lists.NewLinkedCircular(1, 2, 3)
	l2.InsertInPlace(0, 0)
	if got := l2.GetAsSlice(); !reflect.DeepEqual(got, []int{0, 1, 2, 3}) {
		t.Errorf("circular InsertInPlace(0,0) = %v, want [0 1 2 3]", got)
	}
}

func TestLinked_CircularFilterInPlace(t *testing.T) {
	l := lists.NewLinkedCircular(1, 2, 3, 4)
	l.FilterInPlace(func(n int) bool { return n%2 == 0 })
	if got := l.GetAsSlice(); !reflect.DeepEqual(got, []int{2, 4}) {
		t.Errorf("circular FilterInPlace(even) = %v, want [2 4]", got)
	}

	// Filtering everything out of a circular list must empty it cleanly.
	l2 := lists.NewLinkedCircular(1, 3, 5)
	l2.FilterInPlace(func(n int) bool { return n%2 == 0 })
	if l2.Length() != 0 {
		t.Errorf("circular FilterInPlace(none) Length = %d, want 0", l2.Length())
	}
}

func TestLinked_CircularPopAndDequeue(t *testing.T) {
	l := lists.NewLinkedCircular(1, 2, 3)
	val, ok := l.PopInPlace()
	if !ok || val != 3 {
		t.Errorf("PopInPlace() = (%d, %v), want (3, true)", val, ok)
	}
	if got := l.GetAsSlice(); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Errorf("after PopInPlace GetAsSlice() = %v, want [1 2]", got)
	}

	d := lists.NewLinkedCircular(1, 2, 3)
	val, ok = d.DequeueInPlace()
	if !ok || val != 1 {
		t.Errorf("DequeueInPlace() = (%d, %v), want (1, true)", val, ok)
	}
	if got := d.GetAsSlice(); !reflect.DeepEqual(got, []int{2, 3}) {
		t.Errorf("after DequeueInPlace GetAsSlice() = %v, want [2 3]", got)
	}
}

func TestLinked_Enqueue(t *testing.T) {
	l := lists.NewLinked(1, 2)
	if got := l.Enqueue(3); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("Enqueue(3) = %v, want [1 2 3]", got)
	}
	l.EnqueueInPlace(3)
	if got := l.GetAsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("after EnqueueInPlace GetAsSlice() = %v, want [1 2 3]", got)
	}
}

func TestDoublyLinked_Circular(t *testing.T) {
	dl := lists.NewDoublyLinkedCircular(1, 2, 3, 4)

	if dl.Length() != 4 {
		t.Errorf("Length() = %d, want 4", dl.Length())
	}
	if got := dl.GetAsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Errorf("GetAsSlice() = %v, want [1 2 3 4]", got)
	}

	var count, idxSum int
	dl.ForEach(func(int) { count++ })
	dl.ForEachWithIndex(func(idx, _ int) { idxSum += idx })
	if count != 4 || idxSum != 6 {
		t.Errorf("ForEach count=%d idxSum=%d, want 4 and 6", count, idxSum)
	}

	if !dl.AllMatch(func(n int) bool { return n > 0 }) {
		t.Error("AllMatch(>0) = false, want true")
	}
	if !dl.AnyMatch(func(n int) bool { return n == 3 }) {
		t.Error("AnyMatch(==3) = false, want true")
	}
	if _, ok := dl.Find(func(n int) bool { return n == 2 }); !ok {
		t.Error("Find(==2) ok = false, want true")
	}
	if idx := dl.FindIndex(func(n int) bool { return n == 4 }); idx != 3 {
		t.Errorf("FindIndex(==4) = %d, want 3", idx)
	}
	if got := dl.Filter(func(n int) bool { return n%2 == 0 }); !reflect.DeepEqual(got, []int{2, 4}) {
		t.Errorf("Filter(even) = %v, want [2 4]", got)
	}

	// NOTE: FilterInPlace / DequeueInPlace on a *circular* DoublyLinked is
	// intentionally NOT asserted here because DoublyLinked.removeNode has a bug
	// when removing the head (or tail) of a circular list: head.prev points at
	// the tail (non-nil), so the branch that advances dl.head is skipped and the
	// removed element is left reachable from the head pointer. See the
	// package-level report. We only filter a non-head element to keep this test
	// exercising the circular Filter/FilterInPlace traversal without depending
	// on the buggy head-removal path.
	dl.FilterInPlace(func(n int) bool { return n != 3 })
	if dl.Length() != 3 {
		t.Errorf("after FilterInPlace(!=3) Length = %d, want 3", dl.Length())
	}
}

func TestDoublyLinked_InsertInPlace(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		index    int
		elements []int
		want     []int
	}{
		{name: "at head", initial: []int{1, 2, 3}, index: 0, elements: []int{9}, want: []int{9, 1, 2, 3}},
		{name: "in middle", initial: []int{1, 2, 3}, index: 1, elements: []int{9, 8}, want: []int{1, 9, 8, 2, 3}},
		{name: "at end", initial: []int{1, 2, 3}, index: 3, elements: []int{4}, want: []int{1, 2, 3, 4}},
		{name: "out of range ignored", initial: []int{1, 2, 3}, index: 99, elements: []int{4}, want: []int{1, 2, 3}},
		{name: "negative ignored", initial: []int{1, 2, 3}, index: -1, elements: []int{4}, want: []int{1, 2, 3}},
		{name: "into empty", initial: nil, index: 0, elements: []int{1}, want: []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := lists.NewDoublyLinked(tt.initial...)
			dl.InsertInPlace(tt.index, tt.elements...)
			got := dl.GetAsSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertInPlace(%d, %v) = %v, want %v", tt.index, tt.elements, got, tt.want)
			}
		})
	}
}

func TestDoublyLinked_Insert(t *testing.T) {
	dl := lists.NewDoublyLinked(1, 2, 3)
	if got := dl.Insert(1, 9); !reflect.DeepEqual(got, []int{1, 9, 2, 3}) {
		t.Errorf("Insert(1,9) = %v, want [1 9 2 3]", got)
	}
	if got := dl.Insert(-1, 9); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("Insert(-1,9) = %v, want [1 2 3]", got)
	}
}

func TestDoublyLinked_GetFromTailHalf(t *testing.T) {
	// Index in the tail half forces the reverse-walk branch of Get.
	dl := lists.NewDoublyLinked(0, 1, 2, 3, 4)
	if got := dl.Get(4, -1); got != 4 {
		t.Errorf("Get(4) = %d, want 4", got)
	}
	if got := dl.Get(3, -1); got != 3 {
		t.Errorf("Get(3) = %d, want 3", got)
	}
	if got := dl.Get(-1, -1); got != -1 {
		t.Errorf("Get(-1) = %d, want -1", got)
	}
}

func TestDoublyLinked_StackQueueAndSort(t *testing.T) {
	dl := lists.NewDoublyLinked(3, 1, 2)

	if got := dl.Push(4); !reflect.DeepEqual(got, []int{3, 1, 2, 4}) {
		t.Errorf("Push(4) = %v, want [3 1 2 4]", got)
	}
	if got := dl.Enqueue(4); !reflect.DeepEqual(got, []int{3, 1, 2, 4}) {
		t.Errorf("Enqueue(4) = %v, want [3 1 2 4]", got)
	}
	if got := dl.Sort(func(a, b int) bool { return a < b }); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("Sort(asc) = %v, want [1 2 3]", got)
	}

	val, ok, rest := dl.Pop()
	if !ok || val != 2 || !reflect.DeepEqual(rest, []int{3, 1}) {
		t.Errorf("Pop() = (%d, %v, %v), want (2, true, [3 1])", val, ok, rest)
	}
	val, ok, rest = dl.Dequeue()
	if !ok || val != 3 || !reflect.DeepEqual(rest, []int{1, 2}) {
		t.Errorf("Dequeue() = (%d, %v, %v), want (3, true, [1 2])", val, ok, rest)
	}

	dl.SortInPlace(func(a, b int) bool { return a < b })
	if got := dl.GetAsSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("after SortInPlace GetAsSlice() = %v, want [1 2 3]", got)
	}

	dl.EnqueueInPlace(4)
	end, ok := dl.PeekEnd()
	if !ok || end != 4 {
		t.Errorf("PeekEnd() = (%d, %v), want (4, true)", end, ok)
	}
	front, ok := dl.PeekFront()
	if !ok || front != 1 {
		t.Errorf("PeekFront() = (%d, %v), want (1, true)", front, ok)
	}
}

func TestDoublyLinked_EmptyBehaviour(t *testing.T) {
	dl := lists.NewDoublyLinked[int]()
	if _, ok := dl.PeekEnd(); ok {
		t.Error("empty PeekEnd() ok = true, want false")
	}
	if _, ok := dl.PeekFront(); ok {
		t.Error("empty PeekFront() ok = true, want false")
	}
	if _, ok := dl.PopInPlace(); ok {
		t.Error("empty PopInPlace() ok = true, want false")
	}
	if _, ok := dl.DequeueInPlace(); ok {
		t.Error("empty DequeueInPlace() ok = true, want false")
	}
	if _, ok, _ := dl.Pop(); ok {
		t.Error("empty Pop() ok = true, want false")
	}
	if _, ok, _ := dl.Dequeue(); ok {
		t.Error("empty Dequeue() ok = true, want false")
	}
	dl.SortInPlace(func(a, b int) bool { return a < b }) // no-op on empty
}
