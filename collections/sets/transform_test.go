package sets_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

func ExampleMap() {
	s := sets.NewHash(1, 2, 3)
	doubled := sets.Map(s, func(n int) int {
		return n * 2
	})
	got := doubled.AsSlice()
	sort.Ints(got)
	fmt.Printf("%v\n", got)

	// Output:
	// [2 4 6]
}

func ExampleReduce() {
	s := sets.NewHash(1, 2, 3, 4)
	sum := sets.Reduce(s, 0, func(acc, n int) int {
		return acc + n
	})
	fmt.Printf("%d\n", sum)

	// Output:
	// 10
}

func TestMap(t *testing.T) {
	s := sets.NewHash("a", "ab", "abc")
	got := sets.Map(s, func(str string) int {
		return len(str)
	})

	want := []int{1, 2, 3}
	gotSlice := got.AsSlice()
	sort.Ints(gotSlice)
	if fmt.Sprint(gotSlice) != fmt.Sprint(want) {
		t.Errorf("Map() = %v, want %v", gotSlice, want)
	}
}

func TestMap_CollapsesDuplicateResults(t *testing.T) {
	s := sets.NewHash("a", "bb", "cc", "d")
	// Lengths collide: {"a","d"}->1, {"bb","cc"}->2, so the result has 2 elements.
	got := sets.Map(s, func(str string) int {
		return len(str)
	})
	if got.Length() != 2 {
		t.Errorf("Map() length = %d, want 2", got.Length())
	}
}

func TestMap_EmptyInputYieldsNonNilEmptySet(t *testing.T) {
	s := sets.NewHash[int]()
	got := sets.Map(s, func(n int) int { return n * 2 })
	if got == nil {
		t.Fatal("Map() returned nil, want non-nil empty Set")
	}
	if !got.IsEmpty() {
		t.Errorf("Map() = %v, want empty", got.AsSlice())
	}
}

func ExampleMapSorted() {
	s := sets.NewHash("a", "ab", "abc")
	// Map to lengths (an Ordered type); the result iterates in ascending order
	// without a post-hoc sort.
	byLen := sets.MapSorted(s, func(str string) int {
		return len(str)
	})
	fmt.Printf("%v\n", byLen.AsSlice())

	// Output:
	// [1 2 3]
}

func TestMapSorted_PreservesAscendingOrder(t *testing.T) {
	// Unordered Hash input mapped to ordered elements must come out sorted.
	s := sets.NewHash(3, 1, 2)
	got := sets.MapSorted(s, func(n int) int {
		return n * 10
	})
	want := []int{10, 20, 30}
	// AsSlice on a TreeSet-backed SortedSet is already ascending; assert directly.
	if fmt.Sprint(got.AsSlice()) != fmt.Sprint(want) {
		t.Errorf("MapSorted() = %v, want %v", got.AsSlice(), want)
	}
}

func TestMapSorted_CollapsesDuplicateResults(t *testing.T) {
	s := sets.NewHash("a", "bb", "cc", "d")
	got := sets.MapSorted(s, func(str string) int {
		return len(str)
	})
	if got.Length() != 2 {
		t.Errorf("MapSorted() length = %d, want 2", got.Length())
	}
}

func TestMapSorted_EmptyInputYieldsNonNilEmptySortedSet(t *testing.T) {
	s := sets.NewHash[int]()
	got := sets.MapSorted(s, func(n int) int { return n * 2 })
	if got == nil {
		t.Fatal("MapSorted() returned nil, want non-nil empty SortedSet")
	}
	if !got.IsEmpty() {
		t.Errorf("MapSorted() = %v, want empty", got.AsSlice())
	}
}

func TestReduce(t *testing.T) {
	s := sets.NewHash(1, 2, 3, 4)

	t.Run("sums elements onto init", func(t *testing.T) {
		got := sets.Reduce(s, 100, func(acc, n int) int {
			return acc + n
		})
		if got != 110 {
			t.Errorf("Reduce() = %d, want 110", got)
		}
	})

	t.Run("collects into a different accumulator type", func(t *testing.T) {
		got := sets.Reduce(s, []int{}, func(acc []int, n int) []int {
			return append(acc, n*n)
		})
		sort.Ints(got)
		want := []int{1, 4, 9, 16}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("Reduce() = %v, want %v", got, want)
		}
	})

	t.Run("empty input returns init unchanged", func(t *testing.T) {
		empty := sets.NewHash[int]()
		got := sets.Reduce(empty, 42, func(acc, n int) int {
			return acc + n
		})
		if got != 42 {
			t.Errorf("Reduce() = %d, want 42", got)
		}
	})
}
