package collections

import (
	"iter"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/heaps"
)

func intSeq(values ...int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

func pairSeq(keys []string, values []int) iter.Seq2[string, int] {
	return func(yield func(string, int) bool) {
		for i := range keys {
			if !yield(keys[i], values[i]) {
				return
			}
		}
	}
}

func TestListFromSeq(t *testing.T) {
	list := ListFromSeq(intSeq(1, 2, 3))
	if got := list.AsSlice(); !equalInts(got, []int{1, 2, 3}) {
		t.Errorf("ListFromSeq = %v, want [1 2 3]", got)
	}
}

func TestDictFromSeq2(t *testing.T) {
	d := DictFromSeq2(pairSeq([]string{"a", "b"}, []int{1, 2}))
	if d.Length() != 2 {
		t.Errorf("DictFromSeq2 Length() = %d, want 2", d.Length())
	}
	if v, _ := d.Get("a", 0); v != 1 {
		t.Errorf(`DictFromSeq2 Get("a") = %d, want 1`, v)
	}
}

func TestSetFromSeq(t *testing.T) {
	s := SetFromSeq(intSeq(1, 1, 2, 3))
	if s.Length() != 3 {
		t.Errorf("SetFromSeq Length() = %d, want 3", s.Length())
	}
}

func TestDequeFromSeq(t *testing.T) {
	d := DequeFromSeq(intSeq(1, 2, 3))
	if got := d.AsSlice(); !equalInts(got, []int{1, 2, 3}) {
		t.Errorf("DequeFromSeq = %v, want [1 2 3]", got)
	}
}

func TestHeapFromSeq(t *testing.T) {
	h := HeapFromSeq(heaps.Min[int], intSeq(5, 3, 8, 1))
	if got := h.AsSortedSlice(); !equalInts(got, []int{1, 3, 5, 8}) {
		t.Errorf("HeapFromSeq = %v, want [1 3 5 8]", got)
	}
}

func TestListMultimapFromSeq2(t *testing.T) {
	m := ListMultimapFromSeq2(pairSeq([]string{"a", "a", "b"}, []int{1, 2, 3}))
	if got := m.Get("a"); !equalInts(got, []int{1, 2}) {
		t.Errorf(`ListMultimapFromSeq2 Get("a") = %v, want [1 2]`, got)
	}
}

func TestSetMultimapFromSeq2(t *testing.T) {
	m := SetMultimapFromSeq2(pairSeq([]string{"a", "a", "a"}, []int{1, 1, 2}))
	values := m.Get("a")
	sort.Ints(values)
	if !equalInts(values, []int{1, 2}) {
		t.Errorf(`SetMultimapFromSeq2 Get("a") = %v, want [1 2]`, values)
	}
}

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
