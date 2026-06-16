package dicts_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

func ExampleMap() {
	d := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)
	// Swap each entry's key and value, doubling the value.
	swapped := dicts.Map(d, func(k string, v int) (int, string) {
		return v, k
	})
	fmt.Println(swapped.Get(1, "?"))
	fmt.Println(swapped.Get(2, "?"))

	// Output:
	// a true
	// b true
}

func ExampleReduce() {
	d := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
		dicts.Pair[string, int]{Key: "c", Value: 3},
	)
	sum := dicts.Reduce(d, 0, func(acc int, _ string, v int) int {
		return acc + v
	})
	fmt.Printf("%d\n", sum)

	// Output:
	// 6
}

func TestMap(t *testing.T) {
	d := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
		dicts.Pair[string, int]{Key: "c", Value: 3},
	)
	// Map keys to their uppercase form and values to their square.
	got := dicts.Map(d, func(k string, v int) (string, int) {
		return k + k, v * v
	})

	want := map[string]int{"aa": 1, "bb": 4, "cc": 9}
	gotMap := got.AsMap()
	if len(gotMap) != len(want) {
		t.Fatalf("Map() length = %d, want %d", len(gotMap), len(want))
	}
	for k, v := range want {
		if gv, ok := gotMap[k]; !ok || gv != v {
			t.Errorf("Map()[%q] = (%d, %v), want %d", k, gv, ok, v)
		}
	}
}

func TestMap_EmptyInputYieldsNonNilEmptyDict(t *testing.T) {
	d := dicts.NewHash[string, int]()
	got := dicts.Map(d, func(k string, v int) (int, int) {
		return v, v
	})
	if got == nil {
		t.Fatal("Map() returned nil, want non-nil empty Dict")
	}
	if !got.IsEmpty() {
		t.Errorf("Map() = %v, want empty", got.AsMap())
	}
}

func TestMap_CollidingKeysCollapse(t *testing.T) {
	d := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)
	// Both entries map to the same output key; one survives.
	got := dicts.Map(d, func(k string, v int) (string, int) {
		return "same", v
	})
	if got.Length() != 1 {
		t.Errorf("Map() length = %d, want 1", got.Length())
	}
}

func TestReduce(t *testing.T) {
	d := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
		dicts.Pair[string, int]{Key: "c", Value: 3},
	)

	t.Run("sums values onto init", func(t *testing.T) {
		got := dicts.Reduce(d, 100, func(acc int, _ string, v int) int {
			return acc + v
		})
		if got != 106 {
			t.Errorf("Reduce() = %d, want 106", got)
		}
	})

	t.Run("collects keys into a different accumulator type", func(t *testing.T) {
		got := dicts.Reduce(d, []string{}, func(acc []string, k string, _ int) []string {
			return append(acc, k)
		})
		sort.Strings(got)
		want := []string{"a", "b", "c"}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("Reduce() = %v, want %v", got, want)
		}
	})

	t.Run("empty input returns init unchanged", func(t *testing.T) {
		empty := dicts.NewHash[string, int]()
		got := dicts.Reduce(empty, 42, func(acc int, _ string, v int) int {
			return acc + v
		})
		if got != 42 {
			t.Errorf("Reduce() = %d, want 42", got)
		}
	})
}
