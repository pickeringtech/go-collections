package sets_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// setBuilders enumerates every concrete Set implementation so the shared
// Searchable contract (here, NoneMatch) is exercised uniformly across them.
var setBuilders = []struct {
	name  string
	build func(elements ...int) sets.Set[int]
}{
	{"Hash", func(e ...int) sets.Set[int] { return sets.NewHash(e...) }},
	{"ConcurrentHash", func(e ...int) sets.Set[int] { return sets.NewConcurrentHash(e...) }},
	{"ConcurrentHashRW", func(e ...int) sets.Set[int] { return sets.NewConcurrentHashRW(e...) }},
}

func TestSet_NoneMatch(t *testing.T) {
	isThree := func(i int) bool { return i == 3 }

	cases := []struct {
		name     string
		elements []int
		fn       func(int) bool
		want     bool
	}{
		{"no element matches", []int{1, 2, 4, 5}, isThree, true},
		{"some element matches", []int{1, 2, 3, 4}, isThree, false},
		{"empty is vacuously true", nil, isThree, true},
	}

	for _, builder := range setBuilders {
		for _, tc := range cases {
			t.Run(builder.name+"/"+tc.name, func(t *testing.T) {
				set := builder.build(tc.elements...)
				if got := set.NoneMatch(tc.fn); got != tc.want {
					t.Errorf("NoneMatch() = %v, want %v", got, tc.want)
				}
			})
		}
	}
}

// TestSet_NoneMatch_NegatesAnyMatch asserts NoneMatch is the exact negation of
// AnyMatch, keeping the All/Any/None trio internally consistent.
func TestSet_NoneMatch_NegatesAnyMatch(t *testing.T) {
	isEven := func(i int) bool { return i%2 == 0 }
	inputs := [][]int{nil, {1, 3, 5}, {2, 4, 6}, {1, 2, 3}}

	for _, builder := range setBuilders {
		for _, elements := range inputs {
			set := builder.build(elements...)
			if none, any := set.NoneMatch(isEven), set.AnyMatch(isEven); none == any {
				t.Errorf("%s: NoneMatch()=%v should be the negation of AnyMatch()=%v for %v",
					builder.name, none, any, elements)
			}
		}
	}
}
