package lists_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// listBuilders enumerates every concrete List implementation, linear and
// circular, so the shared Searchable contract (NoneMatch, AllMatch and
// FindIndex) is exercised uniformly across them. The circular lists matter: a
// search walk there must stop only once it wraps back to the head.
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
	{"LinkedCircular", func(e ...int) lists.List[int] { return lists.NewLinkedCircular(e...) }},
	{"DoublyLinkedCircular", func(e ...int) lists.List[int] { return lists.NewDoublyLinkedCircular(e...) }},
	{"ConcurrentLinkedCircular", func(e ...int) lists.List[int] { return lists.NewConcurrentLinkedCircular(e...) }},
	{"ConcurrentDoublyLinkedCircular", func(e ...int) lists.List[int] {
		return lists.NewConcurrentDoublyLinkedCircular(e...)
	}},
	{"ConcurrentRWLinkedCircular", func(e ...int) lists.List[int] { return lists.NewConcurrentRWLinkedCircular(e...) }},
	{"ConcurrentRWDoublyLinkedCircular", func(e ...int) lists.List[int] {
		return lists.NewConcurrentRWDoublyLinkedCircular(e...)
	}},
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

// TestList_AllMatch pins that every implementation agrees AllMatch is vacuously
// true for an empty list. The slice-backed and concurrent lists delegate to
// slices.AllMatch, which once reported false for empty input while the linked
// lists reported true.
func TestList_AllMatch(t *testing.T) {
	isOdd := func(i int) bool { return i%2 == 1 }

	cases := []struct {
		name     string
		elements []int
		fn       func(int) bool
		want     bool
	}{
		{"every element matches", []int{1, 3, 5}, isOdd, true},
		{"some element fails", []int{1, 2, 3}, isOdd, false},
		{"empty is vacuously true", nil, isOdd, true},
	}

	for _, builder := range listBuilders {
		for _, tc := range cases {
			t.Run(builder.name+"/"+tc.name, func(t *testing.T) {
				list := builder.build(tc.elements...)
				got := list.AllMatch(tc.fn)
				if got != tc.want {
					t.Errorf("AllMatch() = %v, want %v", got, tc.want)
				}
			})
		}
	}
}

// TestList_FindIndex pins that every implementation reports the index of the
// first matching element wherever it sits — including past the head of a
// circular list — and -1 when nothing matches.
func TestList_FindIndex(t *testing.T) {
	type args struct {
		elements []int
		fn       func(int) bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "a match at the head is index zero", args: args{elements: []int{4, 5, 6}, fn: func(i int) bool { return i == 4 }}, want: 0},
		{name: "a match past the head reports its index", args: args{elements: []int{4, 5, 6}, fn: func(i int) bool { return i == 5 }}, want: 1},
		{name: "a match at the tail reports the last index", args: args{elements: []int{4, 5, 6}, fn: func(i int) bool { return i == 6 }}, want: 2},
		{name: "the first of several matches wins", args: args{elements: []int{4, 5, 6}, fn: func(i int) bool { return i > 4 }}, want: 1},
		{name: "no match reports minus one", args: args{elements: []int{4, 5, 6}, fn: func(i int) bool { return i > 6 }}, want: -1},
		{name: "an empty list reports minus one", args: args{elements: nil, fn: func(int) bool { return true }}, want: -1},
	}
	for _, builder := range listBuilders {
		for _, tt := range tests {
			t.Run(builder.name+"/"+tt.name, func(t *testing.T) {
				got := builder.build(tt.args.elements...).FindIndex(tt.args.fn)
				if got != tt.want {
					t.Errorf("FindIndex() = %d, want %d", got, tt.want)
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
