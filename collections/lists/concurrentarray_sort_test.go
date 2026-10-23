package lists_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// keyedRecord is a sort input whose key has many ties; seq records the
// original position so a reordering of equal keys is observable.
type keyedRecord struct {
	key int
	seq int
}

// TestConcurrentArray_SortIsStable verifies that the concurrent arrays' Sort and
// SortInPlace keep elements with equal keys in their original relative order,
// as the List contract requires. They used the unstable slices.Sort and
// slices.SortInPlace, which reorder ties.
func TestConcurrentArray_SortIsStable(t *testing.T) {
	records := make([]keyedRecord, 0, 40)
	for i := 0; i < 40; i++ {
		records = append(records, keyedRecord{key: i % 3, seq: i})
	}
	byKey := func(a, b keyedRecord) bool { return a.key < b.key }
	want := append([]keyedRecord{}, records...)
	sort.SliceStable(want, func(i, j int) bool { return byKey(want[i], want[j]) })

	sorts := []struct {
		name string
		run  func(l lists.MutableList[keyedRecord]) []keyedRecord
	}{
		{"Sort", func(l lists.MutableList[keyedRecord]) []keyedRecord { return l.Sort(byKey).AsSlice() }},
		{"SortInPlace", func(l lists.MutableList[keyedRecord]) []keyedRecord {
			l.SortInPlace(byKey)
			return l.AsSlice()
		}},
	}
	factories := []struct {
		name string
		make func(elements ...keyedRecord) lists.MutableList[keyedRecord]
	}{
		{"ConcurrentArray", func(e ...keyedRecord) lists.MutableList[keyedRecord] { return lists.NewConcurrentArray(e...) }},
		{"ConcurrentRWArray", func(e ...keyedRecord) lists.MutableList[keyedRecord] {
			return lists.NewConcurrentRWArray(e...)
		}},
	}
	for _, s := range sorts {
		for _, f := range factories {
			t.Run(s.name+"/"+f.name, func(t *testing.T) {
				got := s.run(f.make(records...))
				if !reflect.DeepEqual(got, want) {
					t.Errorf("%s() = %v, want the stable order %v", s.name, got, want)
				}
			})
		}
	}
}
