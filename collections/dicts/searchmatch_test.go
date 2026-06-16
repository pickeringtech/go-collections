package dicts_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// dictBuilders enumerates every concrete Dict implementation so the shared
// Searchable contract (AllMatch/AnyMatch/NoneMatch) is exercised uniformly.
var dictBuilders = []struct {
	name  string
	build func(entries ...dicts.Pair[string, int]) dicts.Dict[string, int]
}{
	{"Hash", func(e ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewHash(e...) }},
	{"Tree", func(e ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewTree(e...) }},
	{"ConcurrentHash", func(e ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewConcurrentHash(e...) }},
	{"ConcurrentHashRW", func(e ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewConcurrentHashRW(e...) }},
}

func pairs(entries ...dicts.Pair[string, int]) []dicts.Pair[string, int] { return entries }

func p(key string, value int) dicts.Pair[string, int] {
	return dicts.Pair[string, int]{Key: key, Value: value}
}

func TestDict_AllMatch(t *testing.T) {
	valuePositive := func(_ string, value int) bool { return value > 0 }

	cases := []struct {
		name    string
		entries []dicts.Pair[string, int]
		fn      func(string, int) bool
		want    bool
	}{
		{"all values positive", pairs(p("a", 1), p("b", 2)), valuePositive, true},
		{"one value non-positive", pairs(p("a", 1), p("b", 0)), valuePositive, false},
		{"empty is vacuously true", nil, valuePositive, true},
	}

	for _, builder := range dictBuilders {
		for _, tc := range cases {
			t.Run(builder.name+"/"+tc.name, func(t *testing.T) {
				dict := builder.build(tc.entries...)
				if got := dict.AllMatch(tc.fn); got != tc.want {
					t.Errorf("AllMatch() = %v, want %v", got, tc.want)
				}
			})
		}
	}
}

func TestDict_AnyMatch(t *testing.T) {
	keyIsB := func(key string, _ int) bool { return key == "b" }

	cases := []struct {
		name    string
		entries []dicts.Pair[string, int]
		fn      func(string, int) bool
		want    bool
	}{
		{"a matching key exists", pairs(p("a", 1), p("b", 2)), keyIsB, true},
		{"no matching key", pairs(p("a", 1), p("c", 3)), keyIsB, false},
		{"empty is false", nil, keyIsB, false},
	}

	for _, builder := range dictBuilders {
		for _, tc := range cases {
			t.Run(builder.name+"/"+tc.name, func(t *testing.T) {
				dict := builder.build(tc.entries...)
				if got := dict.AnyMatch(tc.fn); got != tc.want {
					t.Errorf("AnyMatch() = %v, want %v", got, tc.want)
				}
			})
		}
	}
}

func TestDict_NoneMatch(t *testing.T) {
	valueIsThree := func(_ string, value int) bool { return value == 3 }

	cases := []struct {
		name    string
		entries []dicts.Pair[string, int]
		fn      func(string, int) bool
		want    bool
	}{
		{"no value matches", pairs(p("a", 1), p("b", 2)), valueIsThree, true},
		{"a value matches", pairs(p("a", 1), p("b", 3)), valueIsThree, false},
		{"empty is vacuously true", nil, valueIsThree, true},
	}

	for _, builder := range dictBuilders {
		for _, tc := range cases {
			t.Run(builder.name+"/"+tc.name, func(t *testing.T) {
				dict := builder.build(tc.entries...)
				if got := dict.NoneMatch(tc.fn); got != tc.want {
					t.Errorf("NoneMatch() = %v, want %v", got, tc.want)
				}
			})
		}
	}
}

// TestDict_NoneMatch_NegatesAnyMatch asserts the All/Any/None trio stays
// internally consistent across every Dict implementation.
func TestDict_NoneMatch_NegatesAnyMatch(t *testing.T) {
	valueEven := func(_ string, value int) bool { return value%2 == 0 }
	inputs := [][]dicts.Pair[string, int]{
		nil,
		pairs(p("a", 1), p("b", 3)),
		pairs(p("a", 2), p("b", 4)),
		pairs(p("a", 1), p("b", 2)),
	}

	for _, builder := range dictBuilders {
		for _, entries := range inputs {
			dict := builder.build(entries...)
			if none, any := dict.NoneMatch(valueEven), dict.AnyMatch(valueEven); none == any {
				t.Errorf("%s: NoneMatch()=%v should be the negation of AnyMatch()=%v for %v",
					builder.name, none, any, entries)
			}
		}
	}
}
