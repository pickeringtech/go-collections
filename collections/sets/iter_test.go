package sets_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

type setConstructor struct {
	name string
	make func(...int) sets.Set[int]
}

func allSetConstructors() []setConstructor {
	return []setConstructor{
		{"Hash", func(e ...int) sets.Set[int] { return sets.NewHash(e...) }},
		{"ConcurrentHash", func(e ...int) sets.Set[int] { return sets.NewConcurrentHash(e...) }},
		{"ConcurrentHashRW", func(e ...int) sets.Set[int] { return sets.NewConcurrentHashRW(e...) }},
	}
}

func TestSet_All(t *testing.T) {
	for _, ctor := range allSetConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var got []int
			for v := range ctor.make(3, 1, 2).All() {
				got = append(got, v)
			}
			sort.Ints(got)
			if want := []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
				t.Errorf("All() = %v, want %v", got, want)
			}
		})
	}
}

func TestSet_All_Empty(t *testing.T) {
	for _, ctor := range allSetConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make().All() {
				count++
			}
			if count != 0 {
				t.Errorf("All() over empty set yielded %d times, want 0", count)
			}
		})
	}
}

func TestSet_All_EarlyBreak(t *testing.T) {
	for _, ctor := range allSetConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make(1, 2, 3, 4, 5).All() {
				count++
				break
			}
			if count != 1 {
				t.Errorf("All() with early break yielded %d times, want 1", count)
			}
		})
	}
}

func TestFromSeq(t *testing.T) {
	for _, ctor := range allSetConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			source := ctor.make(1, 2, 3)
			got := sets.FromSeq(source.All())

			gotSlice := got.AsSlice()
			sort.Ints(gotSlice)
			if want := []int{1, 2, 3}; !reflect.DeepEqual(gotSlice, want) {
				t.Errorf("FromSeq round-trip = %v, want %v", gotSlice, want)
			}
		})
	}
}

func TestFromSeq_CollapsesDuplicates(t *testing.T) {
	seq := func(yield func(int) bool) {
		for _, v := range []int{1, 1, 2, 2, 3} {
			if !yield(v) {
				return
			}
		}
	}
	got := sets.FromSeq(seq)
	if got.Length() != 3 {
		t.Errorf("FromSeq with duplicates Length() = %d, want 3", got.Length())
	}
}

func ExampleHash_All() {
	set := sets.NewHash(42)
	for v := range set.All() {
		fmt.Println(v)
	}
	// Output: 42
}

func ExampleFromSeq() {
	source := sets.NewHash(1, 2, 3)
	set := sets.FromSeq(source.All())
	fmt.Println(set.Length())
	// Output: 3
}
