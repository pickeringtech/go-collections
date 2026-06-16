package dicts_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

type dictConstructor struct {
	name string
	make func(...dicts.Pair[string, int]) dicts.Dict[string, int]
}

func allDictConstructors() []dictConstructor {
	return []dictConstructor{
		{"Hash", func(p ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewHash(p...) }},
		{"Tree", func(p ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewTree(p...) }},
		{"ConcurrentHash", func(p ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewConcurrentHash(p...) }},
		{"ConcurrentHashRW", func(p ...dicts.Pair[string, int]) dicts.Dict[string, int] { return dicts.NewConcurrentHashRW(p...) }},
	}
}

func seedPairs() []dicts.Pair[string, int] {
	return []dicts.Pair[string, int]{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c", Value: 3},
	}
}

func TestDict_All(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			d := ctor.make(seedPairs()...)

			got := map[string]int{}
			for k, v := range d.All() {
				got[k] = v
			}
			want := map[string]int{"a": 1, "b": 2, "c": 3}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("All() = %v, want %v", got, want)
			}
		})
	}
}

func TestDict_KeysSeq(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var got []string
			for k := range ctor.make(seedPairs()...).KeysSeq() {
				got = append(got, k)
			}
			sort.Strings(got)
			if want := []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
				t.Errorf("KeysSeq() = %v, want %v", got, want)
			}
		})
	}
}

func TestDict_ValuesSeq(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var got []int
			for v := range ctor.make(seedPairs()...).ValuesSeq() {
				got = append(got, v)
			}
			sort.Ints(got)
			if want := []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
				t.Errorf("ValuesSeq() = %v, want %v", got, want)
			}
		})
	}
}

func TestDict_All_Empty(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make().All() {
				count++
			}
			if count != 0 {
				t.Errorf("All() over empty dict yielded %d times, want 0", count)
			}
		})
	}
}

func TestDict_All_EarlyBreak(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make(seedPairs()...).All() {
				count++
				break
			}
			if count != 1 {
				t.Errorf("All() with early break yielded %d times, want 1", count)
			}
		})
	}
}

func TestTree_All_SortedOrder(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "charlie", Value: 3},
		dicts.Pair[string, int]{Key: "alice", Value: 1},
		dicts.Pair[string, int]{Key: "bob", Value: 2},
	)
	var keys []string
	for k := range tree.All() {
		keys = append(keys, k)
	}
	if want := []string{"alice", "bob", "charlie"}; !reflect.DeepEqual(keys, want) {
		t.Errorf("Tree.All() keys = %v, want sorted %v", keys, want)
	}
}

func TestFromSeq2(t *testing.T) {
	for _, ctor := range allDictConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			source := ctor.make(seedPairs()...)
			got := dicts.FromSeq2(source.All())
			if !reflect.DeepEqual(got.AsMap(), source.AsMap()) {
				t.Errorf("FromSeq2 round-trip = %v, want %v", got.AsMap(), source.AsMap())
			}
		})
	}
}

func TestFromSeq2_LastWins(t *testing.T) {
	seq := func(yield func(string, int) bool) {
		_ = yield("k", 1) && yield("k", 2)
	}
	got := dicts.FromSeq2(seq)
	if v, _ := got.Get("k", 0); v != 2 {
		t.Errorf("FromSeq2 duplicate key: got %d, want last value 2", v)
	}
}

func ExampleFromSeq2() {
	source := dicts.NewTree(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)
	d := dicts.FromSeq2(source.All())
	fmt.Println(d.Length())
	// Output: 2
}
