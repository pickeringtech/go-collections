package dicts_test

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// filterParityFactories builds every MutableDict implementation, plain and
// concurrent, so FilterInPlace can be checked for parity between each
// concurrent type and its non-concurrent sibling.
func filterParityFactories() []struct {
	name string
	make func(entries ...dicts.Pair[string, any]) dicts.MutableDict[string, any]
} {
	return []struct {
		name string
		make func(entries ...dicts.Pair[string, any]) dicts.MutableDict[string, any]
	}{
		{"Hash", func(e ...dicts.Pair[string, any]) dicts.MutableDict[string, any] { return dicts.NewHash(e...) }},
		{"ConcurrentHash", func(e ...dicts.Pair[string, any]) dicts.MutableDict[string, any] {
			return dicts.NewConcurrentHash(e...)
		}},
		{"ConcurrentHashRW", func(e ...dicts.Pair[string, any]) dicts.MutableDict[string, any] {
			return dicts.NewConcurrentHashRW(e...)
		}},
		{"Tree", func(e ...dicts.Pair[string, any]) dicts.MutableDict[string, any] { return dicts.NewTree(e...) }},
		{"ConcurrentTree", func(e ...dicts.Pair[string, any]) dicts.MutableDict[string, any] {
			return dicts.NewConcurrentTree(e...)
		}},
		{"ConcurrentTreeRW", func(e ...dicts.Pair[string, any]) dicts.MutableDict[string, any] {
			return dicts.NewConcurrentTreeRW(e...)
		}},
	}
}

// TestDictFilterInPlaceConcurrentMatchesSibling verifies that, with no
// concurrent write during the predicate window, every concurrent dictionary's
// FilterInPlace removes exactly the keys the predicate rejected — the same
// result as Hash and Tree. Values that are not reflect.DeepEqual to themselves
// (NaN, arrays and structs holding NaN, non-nil funcs) used to survive on the
// concurrent types because removal was conditioned on DeepEqual against the
// snapshot.
func TestDictFilterInPlaceConcurrentMatchesSibling(t *testing.T) {
	type args struct {
		entries []dicts.Pair[string, any]
		keep    func(key string, value any) bool
	}
	keepOnlyKeep := func(key string, _ any) bool { return key == "keep" }
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "a rejected NaN value is removed",
			args: args{entries: []dicts.Pair[string, any]{{Key: "nan", Value: math.NaN()}, {Key: "keep", Value: 1}}, keep: keepOnlyKeep},
			want: []string{"keep"},
		},
		{
			name: "a rejected array holding NaN is removed",
			args: args{entries: []dicts.Pair[string, any]{{Key: "array", Value: [2]float64{1, math.NaN()}}, {Key: "keep", Value: 1}}, keep: keepOnlyKeep},
			want: []string{"keep"},
		},
		{
			name: "a rejected struct holding NaN is removed",
			args: args{entries: []dicts.Pair[string, any]{{Key: "struct", Value: struct{ X float64 }{math.NaN()}}, {Key: "keep", Value: 1}}, keep: keepOnlyKeep},
			want: []string{"keep"},
		},
		{
			name: "a rejected non-nil func value is removed",
			args: args{entries: []dicts.Pair[string, any]{{Key: "func", Value: func() {}}, {Key: "keep", Value: 1}}, keep: keepOnlyKeep},
			want: []string{"keep"},
		},
		{
			name: "rejecting every entry empties the dictionary",
			args: args{entries: []dicts.Pair[string, any]{{Key: "a", Value: math.NaN()}, {Key: "b", Value: 2}}, keep: func(string, any) bool { return false }},
			want: []string{},
		},
		{
			name: "keeping every entry leaves the dictionary unchanged",
			args: args{entries: []dicts.Pair[string, any]{{Key: "a", Value: math.NaN()}, {Key: "b", Value: 2}}, keep: func(string, any) bool { return true }},
			want: []string{"a", "b"},
		},
		{
			name: "an empty dictionary stays empty",
			args: args{entries: nil, keep: keepOnlyKeep},
			want: []string{},
		},
	}
	for _, tt := range tests {
		for _, f := range filterParityFactories() {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				d := f.make(tt.args.entries...)
				d.FilterInPlace(tt.args.keep)

				got := d.Keys()
				sort.Strings(got)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FilterInPlace() left keys %v, want %v", got, tt.want)
				}
			})
		}
	}
}

// TestDictFilterInPlaceConcurrentPredicateCalls verifies that a concurrent
// dictionary's FilterInPlace calls the predicate exactly once per entry, as its
// non-concurrent sibling does, and — for the sorted types — in ascending key
// order.
func TestDictFilterInPlaceConcurrentPredicateCalls(t *testing.T) {
	entries := []dicts.Pair[string, any]{{Key: "c", Value: 3}, {Key: "a", Value: 1}, {Key: "b", Value: 2}}
	tests := []struct {
		name   string
		sorted bool
	}{
		{name: "Hash", sorted: false},
		{name: "ConcurrentHash", sorted: false},
		{name: "ConcurrentHashRW", sorted: false},
		{name: "Tree", sorted: true},
		{name: "ConcurrentTree", sorted: true},
		{name: "ConcurrentTreeRW", sorted: true},
	}
	factories := map[string]func(entries ...dicts.Pair[string, any]) dicts.MutableDict[string, any]{}
	for _, f := range filterParityFactories() {
		factories[f.name] = f.make
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := factories[tt.name](entries...)
			var calls []string
			d.FilterInPlace(func(key string, _ any) bool {
				calls = append(calls, key)
				return key != "b"
			})

			if !tt.sorted {
				sort.Strings(calls)
			}
			want := []string{"a", "b", "c"}
			if !reflect.DeepEqual(calls, want) {
				t.Errorf("FilterInPlace() predicate calls = %v, want %v", calls, want)
			}
		})
	}
}

// TestConcurrentDictFilterInPlacePreservesConcurrentWrite pins the issue #153
// contract for every mutating method: a write that lands while the predicate
// runs outside the lock must survive the apply phase. The predicate performs
// the racing write itself on its first call — it runs outside the lock, which
// makes the otherwise-racy window deterministic. Every case rejects "a" on its
// snapshot value, so a mutating method that failed to register its write would
// let the apply phase delete "a" unconditionally and lose an update to it.
func TestConcurrentDictFilterInPlacePreservesConcurrentWrite(t *testing.T) {
	type args struct {
		write func(d dicts.MutableDict[string, int])
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			name: "PutInPlace updating a rejected key keeps the new value",
			args: args{write: func(d dicts.MutableDict[string, int]) { d.PutInPlace("a", 10) }},
			want: map[string]int{"a": 10, "b": 2, "c": 3},
		},
		{
			name: "PutManyInPlace updating a rejected key keeps the new value",
			args: args{write: func(d dicts.MutableDict[string, int]) {
				d.PutManyInPlace(dicts.Pair[string, int]{Key: "a", Value: 10}, dicts.Pair[string, int]{Key: "d", Value: 4})
			}},
			want: map[string]int{"a": 10, "b": 2, "c": 3, "d": 4},
		},
		{
			name: "UpdateInPlace updating a rejected key keeps the new value",
			args: args{write: func(d dicts.MutableDict[string, int]) {
				d.UpdateInPlace("a", func(old int, _ bool) int { return old + 10 })
			}},
			want: map[string]int{"a": 11, "b": 2, "c": 3},
		},
		{
			name: "RemoveInPlace of a kept key stays removed",
			args: args{write: func(d dicts.MutableDict[string, int]) { d.RemoveInPlace("b") }},
			want: map[string]int{"c": 3},
		},
		{
			name: "RemoveManyInPlace of kept keys stays removed",
			args: args{write: func(d dicts.MutableDict[string, int]) { d.RemoveManyInPlace("b", "c") }},
			want: map[string]int{},
		},
		{
			name: "Clear stays cleared",
			args: args{write: func(d dicts.MutableDict[string, int]) { d.Clear() }},
			want: map[string]int{},
		},
		{
			name: "a nested FilterInPlace is applied too",
			args: args{write: func(d dicts.MutableDict[string, int]) {
				d.FilterInPlace(func(key string, _ int) bool { return key != "c" })
			}},
			want: map[string]int{"b": 2},
		},
	}
	for _, tt := range tests {
		for _, f := range concurrentWriteFactories() {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				d := f.make()
				wrote := false
				d.FilterInPlace(func(key string, _ int) bool {
					if !wrote {
						wrote = true
						tt.args.write(d) // racing write lands in the evaluation window
					}
					return key != "a"
				})

				got := d.AsMap()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FilterInPlace() with a concurrent write = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

// TestConcurrentDictFilterInPlacePreservesWritesThroughTwoMethods extends the
// issue #153 contract to two writes through different methods in the same
// predicate window. FilterInPlace detects a concurrent write by comparing a
// write counter before and after the predicate runs, so every mutating method
// must move that counter the same way: a method that moved it backwards would
// cancel a companion method's write (an ABA on the counter), the dictionary
// would look untouched, and the apply phase would delete the rejected key "a"
// unconditionally — losing the companion's update to it. Each case pairs the
// method under test with a PutInPlace("a", 10) companion (UpdateInPlace when
// PutInPlace is the method under test); a nested FilterInPlace covers the
// counter bump FilterInPlace makes itself.
func TestConcurrentDictFilterInPlacePreservesWritesThroughTwoMethods(t *testing.T) {
	putA := func(d dicts.MutableDict[string, int]) { d.PutInPlace("a", 10) }
	type args struct {
		write     func(d dicts.MutableDict[string, int])
		companion func(d dicts.MutableDict[string, int])
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			name: "PutInPlace then UpdateInPlace of the rejected key keeps both writes",
			args: args{
				write: func(d dicts.MutableDict[string, int]) { d.PutInPlace("d", 4) },
				companion: func(d dicts.MutableDict[string, int]) {
					d.UpdateInPlace("a", func(old int, _ bool) int { return old + 10 })
				},
			},
			want: map[string]int{"a": 11, "b": 2, "c": 3, "d": 4},
		},
		{
			name: "PutManyInPlace then PutInPlace of the rejected key keeps both writes",
			args: args{
				write:     func(d dicts.MutableDict[string, int]) { d.PutManyInPlace(dicts.Pair[string, int]{Key: "d", Value: 4}) },
				companion: putA,
			},
			want: map[string]int{"a": 10, "b": 2, "c": 3, "d": 4},
		},
		{
			name: "UpdateInPlace then PutInPlace of the rejected key keeps both writes",
			args: args{
				write: func(d dicts.MutableDict[string, int]) {
					d.UpdateInPlace("b", func(old int, _ bool) int { return old + 10 })
				},
				companion: putA,
			},
			want: map[string]int{"a": 10, "b": 12, "c": 3},
		},
		{
			name: "RemoveInPlace then PutInPlace of the rejected key keeps both writes",
			args: args{write: func(d dicts.MutableDict[string, int]) { d.RemoveInPlace("b") }, companion: putA},
			want: map[string]int{"a": 10, "c": 3},
		},
		{
			name: "RemoveManyInPlace then PutInPlace of the rejected key keeps both writes",
			args: args{write: func(d dicts.MutableDict[string, int]) { d.RemoveManyInPlace("b", "c") }, companion: putA},
			want: map[string]int{"a": 10},
		},
		{
			name: "Clear then PutInPlace of the rejected key keeps only the new value",
			args: args{write: func(d dicts.MutableDict[string, int]) { d.Clear() }, companion: putA},
			want: map[string]int{"a": 10},
		},
		{
			name: "a nested FilterInPlace then PutInPlace of the rejected key keeps both writes",
			args: args{
				write: func(d dicts.MutableDict[string, int]) {
					d.FilterInPlace(func(key string, _ int) bool { return key != "c" })
				},
				companion: putA,
			},
			want: map[string]int{"a": 10, "b": 2},
		},
	}
	for _, tt := range tests {
		for _, f := range concurrentWriteFactories() {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				d := f.make()
				wrote := false
				d.FilterInPlace(func(key string, _ int) bool {
					if !wrote {
						wrote = true
						tt.args.write(d) // both racing writes land in the evaluation window
						tt.args.companion(d)
					}
					return key != "a"
				})

				got := d.AsMap()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FilterInPlace() with two concurrent writes = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

// concurrentWriteFactories builds each concurrent dictionary seeded with
// {"a": 1, "b": 2, "c": 3}.
func concurrentWriteFactories() []struct {
	name string
	make func() dicts.MutableDict[string, int]
} {
	seed := []dicts.Pair[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}}
	return []struct {
		name string
		make func() dicts.MutableDict[string, int]
	}{
		{"ConcurrentHash", func() dicts.MutableDict[string, int] { return dicts.NewConcurrentHash(seed...) }},
		{"ConcurrentHashRW", func() dicts.MutableDict[string, int] { return dicts.NewConcurrentHashRW(seed...) }},
		{"ConcurrentTree", func() dicts.MutableDict[string, int] { return dicts.NewConcurrentTree(seed...) }},
		{"ConcurrentTreeRW", func() dicts.MutableDict[string, int] { return dicts.NewConcurrentTreeRW(seed...) }},
	}
}
