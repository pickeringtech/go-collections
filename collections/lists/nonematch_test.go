package lists_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// listBuilders enumerates every concrete List implementation so the shared
// Searchable contract (here, NoneMatch) is exercised uniformly across them.
var listBuilders = []struct {
	name  string
	build func(elements ...int) lists.List[int]
}{
	{"Array", func(e ...int) lists.List[int] { return lists.NewArray(e...) }},
	{"ConcurrentArray", func(e ...int) lists.List[int] { return lists.NewConcurrentArray(e...) }},
	{"ConcurrentRWArray", func(e ...int) lists.List[int] { return lists.NewConcurrentRWArray(e...) }},
	{"Linked", func(e ...int) lists.List[int] { return lists.NewLinked(e...) }},
	{"DoublyLinked", func(e ...int) lists.List[int] { return lists.NewDoublyLinked(e...) }},
	{"ConcurrentLinked", func(e ...int) lists.List[int] { return lists.NewConcurrentLinked(e...) }},
	{"ConcurrentDoublyLinked", func(e ...int) lists.List[int] { return lists.NewConcurrentDoublyLinked(e...) }},
	{"ConcurrentRWLinked", func(e ...int) lists.List[int] { return lists.NewConcurrentRWLinked(e...) }},
	{"ConcurrentRWDoublyLinked", func(e ...int) lists.List[int] { return lists.NewConcurrentRWDoublyLinked(e...) }},
}

func TestList_NoneMatch(t *testing.T) {
	isThree := func(i int) bool { return i == 3 }

	cases := []struct {
		name     string
		elements []int
		fn       func(int) bool
		want     bool
	}{
		{"no element matches", []int{1, 2, 4, 5}, isThree, true},
		{"some element matches", []int{1, 2, 3, 4}, isThree, false},
		{"every element matches", []int{3, 3, 3}, isThree, false},
		{"empty is vacuously true", nil, isThree, true},
	}

	for _, builder := range listBuilders {
		for _, tc := range cases {
			t.Run(builder.name+"/"+tc.name, func(t *testing.T) {
				list := builder.build(tc.elements...)
				if got := list.NoneMatch(tc.fn); got != tc.want {
					t.Errorf("NoneMatch() = %v, want %v", got, tc.want)
				}
			})
		}
	}
}

// TestList_NoneMatch_NegatesAnyMatch asserts the All/Any/None trio stays
// internally consistent: NoneMatch must be the exact negation of AnyMatch.
func TestList_NoneMatch_NegatesAnyMatch(t *testing.T) {
	isEven := func(i int) bool { return i%2 == 0 }
	inputs := [][]int{nil, {1, 3, 5}, {2, 4, 6}, {1, 2, 3}}

	for _, builder := range listBuilders {
		for _, elements := range inputs {
			list := builder.build(elements...)
			if none, any := list.NoneMatch(isEven), list.AnyMatch(isEven); none == any {
				t.Errorf("%s: NoneMatch()=%v should be the negation of AnyMatch()=%v for %v",
					builder.name, none, any, elements)
			}
		}
	}
}
