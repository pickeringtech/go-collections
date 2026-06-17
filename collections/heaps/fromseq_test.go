package heaps_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/heaps"
)

func TestFromSeq(t *testing.T) {
	source := heaps.NewMin(5, 3, 8, 1)
	got := heaps.FromSeq(heaps.Min[int], source.All())
	if want := []int{1, 3, 5, 8}; !reflect.DeepEqual(got.AsSortedSlice(), want) {
		t.Errorf("FromSeq priority order = %v, want %v", got.AsSortedSlice(), want)
	}
}

func TestFromSeq_Empty(t *testing.T) {
	got := heaps.FromSeq(heaps.Min[int], heaps.NewMin[int]().All())
	if !got.IsEmpty() {
		t.Errorf("FromSeq over empty sequence should be empty")
	}
}

func TestFromSeq_RespectsComparator(t *testing.T) {
	source := heaps.NewMin(1, 2, 3)
	maxHeap := heaps.FromSeq(heaps.Max[int], source.All())
	if top, _ := maxHeap.Peek(); top != 3 {
		t.Errorf("FromSeq with Max comparator Peek() = %d, want 3", top)
	}
}

func ExampleFromSeq() {
	source := heaps.NewMax(5, 3, 8, 1)
	h := heaps.FromSeq(heaps.Min[int], source.All())
	fmt.Println(h.AsSortedSlice())
	// Output: [1 3 5 8]
}
