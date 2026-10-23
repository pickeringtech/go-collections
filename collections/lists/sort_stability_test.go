package lists_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// keyedRecord is sorted by key alone; seq records its original position so a
// sort that reorders equal keys is observable.
type keyedRecord struct {
	key int
	seq int
}

// keyedRecords builds n records keyed i%3, so every key is shared by many
// records and the input is far beyond the small-slice insertion-sort cutoff
// where an unstable sort can still happen to preserve order.
func keyedRecords(n int) []keyedRecord {
	records := make([]keyedRecord, n)
	for i := range records {
		records[i] = keyedRecord{key: i % 3, seq: i}
	}
	return records
}

// stableSortOracle is the stable-sorted form of keyedRecords(n): grouped by key
// ascending, each group in its original seq order.
func stableSortOracle(n int) []keyedRecord {
	want := make([]keyedRecord, 0, n)
	for key := 0; key < 3; key++ {
		for i := key; i < n; i += 3 {
			want = append(want, keyedRecord{key: key, seq: i})
		}
	}
	return want
}

func TestList_SortIsStable(t *testing.T) {
	byKey := func(a, b keyedRecord) bool { return a.key < b.key }
	constructors := []struct {
		name string
		make func(...keyedRecord) lists.MutableList[keyedRecord]
	}{
		{"Array", func(e ...keyedRecord) lists.MutableList[keyedRecord] { return lists.NewArray(e...) }},
		{"Linked", func(e ...keyedRecord) lists.MutableList[keyedRecord] { return lists.NewLinked(e...) }},
		{"LinkedCircular", func(e ...keyedRecord) lists.MutableList[keyedRecord] { return lists.NewLinkedCircular(e...) }},
		{"DoublyLinked", func(e ...keyedRecord) lists.MutableList[keyedRecord] { return lists.NewDoublyLinked(e...) }},
		{"DoublyLinkedCircular", func(e ...keyedRecord) lists.MutableList[keyedRecord] {
			return lists.NewDoublyLinkedCircular(e...)
		}},
		{"ConcurrentArray", func(e ...keyedRecord) lists.MutableList[keyedRecord] { return lists.NewConcurrentArray(e...) }},
		{"ConcurrentRWArray", func(e ...keyedRecord) lists.MutableList[keyedRecord] {
			return lists.NewConcurrentRWArray(e...)
		}},
		{"ConcurrentLinked", func(e ...keyedRecord) lists.MutableList[keyedRecord] { return lists.NewConcurrentLinked(e...) }},
		{"ConcurrentRWLinked", func(e ...keyedRecord) lists.MutableList[keyedRecord] {
			return lists.NewConcurrentRWLinked(e...)
		}},
		{"ConcurrentDoublyLinked", func(e ...keyedRecord) lists.MutableList[keyedRecord] {
			return lists.NewConcurrentDoublyLinked(e...)
		}},
		{"ConcurrentRWDoublyLinked", func(e ...keyedRecord) lists.MutableList[keyedRecord] {
			return lists.NewConcurrentRWDoublyLinked(e...)
		}},
	}
	tests := []struct {
		name string
		sort func(l lists.MutableList[keyedRecord]) []keyedRecord
	}{
		{
			name: "Sort keeps equal keys in their original order",
			sort: func(l lists.MutableList[keyedRecord]) []keyedRecord { return l.Sort(byKey).AsSlice() },
		},
		{
			name: "SortInPlace keeps equal keys in their original order",
			sort: func(l lists.MutableList[keyedRecord]) []keyedRecord {
				l.SortInPlace(byKey)
				return l.AsSlice()
			},
		},
	}
	const n = 40
	want := stableSortOracle(n)
	for _, ctor := range constructors {
		for _, tt := range tests {
			t.Run(ctor.name+"/"+tt.name, func(t *testing.T) {
				got := tt.sort(ctor.make(keyedRecords(n)...))
				if !reflect.DeepEqual(got, want) {
					t.Errorf("sorted = %v, want %v", got, want)
				}
			})
		}
	}
}
