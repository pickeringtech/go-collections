package multimaps_test

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/multimaps"
)

// point is a pointer target whose distinct instances are reflect.DeepEqual to
// one another, so it tells "remove this binding" apart from "remove the first
// deeply-equal binding".
type point struct{ X int }

// identical reports whether a and b are the same value: == for comparable
// values (pointer identity for pointers), with NaN treated as identical to NaN.
func identical(a, b any) bool {
	if a == b {
		return true
	}
	af, aIsFloat := a.(float64)
	bf, bIsFloat := b.(float64)
	return aIsFloat && bIsFloat && math.IsNaN(af) && math.IsNaN(bf)
}

// TestListMultimapFilterInPlaceConcurrentMatchesSibling verifies that, with no
// concurrent write during the predicate window, the concurrent list-backed
// multimaps' FilterInPlace keeps exactly the bindings the predicate accepted —
// the same result as ListMultimap. Removal used to delete "one deeply-equal
// binding per rejected entry", which never removed NaN, removed the wrong one
// of two deeply-equal pointers, and broke order-dependent predicates.
func TestListMultimapFilterInPlaceConcurrentMatchesSibling(t *testing.T) {
	p1, p2 := &point{X: 1}, &point{X: 1}
	type args struct {
		values  []any
		newKeep func() func(key string, value any) bool
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			name: "a rejected NaN value is removed",
			args: args{values: []any{math.NaN(), 1.0}, newKeep: func() func(string, any) bool {
				return func(_ string, v any) bool { return !identical(v, math.NaN()) }
			}},
			want: []any{1.0},
		},
		{
			name: "the rejected binding is removed, not a deeply-equal sibling",
			args: args{values: []any{p1, p2}, newKeep: func() func(string, any) bool {
				return func(_ string, v any) bool { return v == any(p1) }
			}},
			want: []any{p1},
		},
		{
			name: "an order-dependent predicate keeps the first occurrence",
			args: args{values: []any{1, 2, 1}, newKeep: func() func(string, any) bool {
				seen := map[any]bool{}
				return func(_ string, v any) bool {
					if seen[v] {
						return false
					}
					seen[v] = true
					return true
				}
			}},
			want: []any{1, 2},
		},
		{
			name: "rejecting every value drops the key",
			args: args{values: []any{math.NaN(), p1}, newKeep: func() func(string, any) bool {
				return func(string, any) bool { return false }
			}},
			want: []any{},
		},
	}
	factories := []struct {
		name string
		make func() multimaps.MutableMultimap[string, any]
	}{
		{"ListMultimap", func() multimaps.MutableMultimap[string, any] { return multimaps.NewListMultimap[string, any]() }},
		{"ConcurrentListMultimap", func() multimaps.MutableMultimap[string, any] {
			return multimaps.NewConcurrentListMultimap[string, any]()
		}},
		{"ConcurrentRWListMultimap", func() multimaps.MutableMultimap[string, any] {
			return multimaps.NewConcurrentRWListMultimap[string, any]()
		}},
	}
	for _, tt := range tests {
		for _, f := range factories {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				m := f.make()
				m.PutAllInPlace("k", tt.args.values...)
				m.FilterInPlace(tt.args.newKeep())

				got := m.Get("k")
				if len(got) != len(tt.want) {
					t.Fatalf("FilterInPlace() left %v, want %v", got, tt.want)
				}
				for i := range got {
					if !identical(got[i], tt.want[i]) {
						t.Fatalf("FilterInPlace() left %v, want %v (binding %d differs)", got, tt.want, i)
					}
				}
				if m.ContainsKey("k") != (len(tt.want) > 0) {
					t.Errorf("ContainsKey(\"k\") = %t after FilterInPlace() left %v", m.ContainsKey("k"), got)
				}
			})
		}
	}
}

// TestConcurrentMultimapFilterInPlacePreservesConcurrentWrite pins the issue
// #153 contract for every mutating method: a write that lands while the
// predicate runs outside the lock must survive the apply phase. The predicate
// performs the racing write itself on its first call — it runs outside the
// lock, which makes the otherwise-racy window deterministic. Every case rejects
// (b, 2), so a list-backed mutating method that failed to register its write
// would let the apply phase replace the contents with the stale filtered
// snapshot and lose it.
func TestConcurrentMultimapFilterInPlacePreservesConcurrentWrite(t *testing.T) {
	type args struct {
		write func(m multimaps.MutableMultimap[string, int])
	}
	tests := []struct {
		name string
		args args
		want map[string][]int
	}{
		{
			name: "PutInPlace to the rejected key keeps the new binding",
			args: args{write: func(m multimaps.MutableMultimap[string, int]) { m.PutInPlace("b", 3) }},
			want: map[string][]int{"a": {1}, "b": {3}},
		},
		{
			name: "PutAllInPlace to a new key keeps the new bindings",
			args: args{write: func(m multimaps.MutableMultimap[string, int]) { m.PutAllInPlace("c", 4, 5) }},
			want: map[string][]int{"a": {1}, "c": {4, 5}},
		},
		{
			name: "RemoveInPlace of a kept binding stays removed",
			args: args{write: func(m multimaps.MutableMultimap[string, int]) { m.RemoveInPlace("a", 1) }},
			want: map[string][]int{},
		},
		{
			name: "RemoveAllInPlace of a kept key stays removed",
			args: args{write: func(m multimaps.MutableMultimap[string, int]) { m.RemoveAllInPlace("a") }},
			want: map[string][]int{},
		},
		{
			name: "Clear stays cleared",
			args: args{write: func(m multimaps.MutableMultimap[string, int]) { m.Clear() }},
			want: map[string][]int{},
		},
		{
			name: "a nested FilterInPlace is applied too",
			args: args{write: func(m multimaps.MutableMultimap[string, int]) {
				m.FilterInPlace(func(key string, _ int) bool { return key == "b" })
			}},
			want: map[string][]int{},
		},
	}
	for _, tt := range tests {
		for _, f := range multimapReentrancyFactories() {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				m := f.make() // a->1, b->2
				wrote := false
				m.FilterInPlace(func(key string, _ int) bool {
					if !wrote {
						wrote = true
						tt.args.write(m) // racing write lands in the evaluation window
					}
					return key != "b"
				})

				got := m.AsMap()
				for _, values := range got {
					sort.Ints(values)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FilterInPlace() with a concurrent write = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

// TestConcurrentMultimapFilterInPlacePreservesWritesThroughTwoMethods extends
// the issue #153 contract to two writes through different methods in the same
// predicate window. The list-backed multimaps detect a concurrent write by
// comparing a write counter before and after the predicate runs, so every
// mutating method must move that counter the same way: a method that moved it
// backwards would cancel a companion method's write (an ABA on the counter),
// the multimap would look untouched, and the apply phase would replace the
// contents with the stale filtered snapshot {a: [1]}, losing both writes. Each
// case pairs the method under test with a PutInPlace("c", 9) companion
// (PutAllInPlace when PutInPlace is the method under test); a nested
// FilterInPlace covers the counter bump FilterInPlace makes itself.
func TestConcurrentMultimapFilterInPlacePreservesWritesThroughTwoMethods(t *testing.T) {
	putC := func(m multimaps.MutableMultimap[string, int]) { m.PutInPlace("c", 9) }
	type args struct {
		write     func(m multimaps.MutableMultimap[string, int])
		companion func(m multimaps.MutableMultimap[string, int])
	}
	tests := []struct {
		name string
		args args
		want map[string][]int
	}{
		{
			name: "PutInPlace then PutAllInPlace keeps both bindings",
			args: args{
				write:     func(m multimaps.MutableMultimap[string, int]) { m.PutInPlace("b", 3) },
				companion: func(m multimaps.MutableMultimap[string, int]) { m.PutAllInPlace("c", 9) },
			},
			want: map[string][]int{"a": {1}, "b": {3}, "c": {9}},
		},
		{
			name: "PutAllInPlace then PutInPlace keeps both writes",
			args: args{write: func(m multimaps.MutableMultimap[string, int]) { m.PutAllInPlace("d", 4, 5) }, companion: putC},
			want: map[string][]int{"a": {1}, "c": {9}, "d": {4, 5}},
		},
		{
			name: "RemoveInPlace then PutInPlace keeps the removal and the binding",
			args: args{write: func(m multimaps.MutableMultimap[string, int]) { m.RemoveInPlace("a", 1) }, companion: putC},
			want: map[string][]int{"c": {9}},
		},
		{
			name: "RemoveAllInPlace then PutInPlace keeps the removal and the binding",
			args: args{write: func(m multimaps.MutableMultimap[string, int]) { m.RemoveAllInPlace("a") }, companion: putC},
			want: map[string][]int{"c": {9}},
		},
		{
			name: "Clear then PutInPlace keeps only the new binding",
			args: args{write: func(m multimaps.MutableMultimap[string, int]) { m.Clear() }, companion: putC},
			want: map[string][]int{"c": {9}},
		},
		{
			name: "a nested FilterInPlace then PutInPlace keeps both writes",
			args: args{
				write: func(m multimaps.MutableMultimap[string, int]) {
					m.FilterInPlace(func(key string, _ int) bool { return key == "b" })
				},
				companion: putC,
			},
			want: map[string][]int{"c": {9}},
		},
	}
	for _, tt := range tests {
		for _, f := range multimapReentrancyFactories() {
			t.Run(tt.name+"/"+f.name, func(t *testing.T) {
				m := f.make() // a->1, b->2
				wrote := false
				m.FilterInPlace(func(key string, _ int) bool {
					if !wrote {
						wrote = true
						tt.args.write(m) // both racing writes land in the evaluation window
						tt.args.companion(m)
					}
					return key != "b"
				})

				got := m.AsMap()
				for _, values := range got {
					sort.Ints(values)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FilterInPlace() with two concurrent writes = %v, want %v", got, tt.want)
				}
			})
		}
	}
}
