package collections_test

import (
	"reflect"
	"slices"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections"
	"github.com/pickeringtech/go-collections/collections/dicts"
	"github.com/pickeringtech/go-collections/collections/heaps"
	"github.com/pickeringtech/go-collections/collections/sets"
)

// These tests pin that the facade hands back the Mutable* interface for every
// mutable structure, as it already did for lists, deques, dicts, multimaps and
// LRU caches: the in-place operations must be reachable without dropping down
// to the subpackage constructors.

func TestFacade_HeapConstructorsReturnMutableHeap(t *testing.T) {
	less := heaps.Min[int]
	tests := []struct {
		name string
		make func(values ...int) heaps.MutableHeap[int]
		push int
		want []int
	}{
		{name: "NewHeap", make: func(v ...int) heaps.MutableHeap[int] { return collections.NewHeap(less, v...) }, push: 0, want: []int{0, 1, 3, 5}},
		{name: "NewMinHeap", make: collections.NewMinHeap[int], push: 0, want: []int{0, 1, 3, 5}},
		{name: "NewMaxHeap", make: collections.NewMaxHeap[int], push: 9, want: []int{9, 5, 3, 1}},
		{name: "NewConcurrentHeap", make: func(v ...int) heaps.MutableHeap[int] { return collections.NewConcurrentHeap(less, v...) }, push: 0, want: []int{0, 1, 3, 5}},
		{name: "NewConcurrentMinHeap", make: collections.NewConcurrentMinHeap[int], push: 0, want: []int{0, 1, 3, 5}},
		{name: "NewConcurrentMaxHeap", make: collections.NewConcurrentMaxHeap[int], push: 9, want: []int{9, 5, 3, 1}},
		{name: "NewConcurrentRWHeap", make: func(v ...int) heaps.MutableHeap[int] { return collections.NewConcurrentRWHeap(less, v...) }, push: 0, want: []int{0, 1, 3, 5}},
		{name: "NewConcurrentRWMinHeap", make: collections.NewConcurrentRWMinHeap[int], push: 0, want: []int{0, 1, 3, 5}},
		{name: "NewConcurrentRWMaxHeap", make: collections.NewConcurrentRWMaxHeap[int], push: 9, want: []int{9, 5, 3, 1}},
		{name: "HeapFromSeq", make: func(v ...int) heaps.MutableHeap[int] { return collections.HeapFromSeq(less, slices.Values(v)) }, push: 0, want: []int{0, 1, 3, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.make(5, 1, 3)
			h.PushInPlace(tt.push)
			got := []int{}
			for {
				value, ok := h.PopInPlace()
				if !ok {
					break
				}
				got = append(got, value)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PopInPlace drain = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFacade_SetConstructorsReturnMutableSet(t *testing.T) {
	tests := []struct {
		name string
		make func(elements ...int) sets.MutableSet[int]
	}{
		{name: "NewSet", make: collections.NewSet[int]},
		{name: "NewConcurrentSet", make: collections.NewConcurrentSet[int]},
		{name: "NewConcurrentRWSet", make: collections.NewConcurrentRWSet[int]},
		{name: "SetFromSeq", make: func(e ...int) sets.MutableSet[int] { return collections.SetFromSeq(slices.Values(e)) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.make(1, 2, 3)
			s.AddInPlace(4)
			s.RemoveInPlace(1)
			s.FilterInPlace(func(n int) bool { return n != 3 })
			got := s.AsSlice()
			sort.Ints(got)
			if !reflect.DeepEqual(got, []int{2, 4}) {
				t.Errorf("after in-place edits = %v, want [2 4]", got)
			}
		})
	}
}

func TestFacade_DictFromSeq2ReturnsMutableDict(t *testing.T) {
	tests := []struct {
		name string
		make func(entries map[string]int) dicts.MutableDict[string, int]
	}{
		{
			name: "DictFromSeq2",
			make: func(entries map[string]int) dicts.MutableDict[string, int] {
				return collections.DictFromSeq2(func(yield func(string, int) bool) {
					for key, value := range entries {
						if !yield(key, value) {
							return
						}
					}
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.make(map[string]int{"a": 1, "b": 2})
			d.PutInPlace("c", 3)
			d.RemoveInPlace("a")
			got := d.AsMap()
			want := map[string]int{"b": 2, "c": 3}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("after in-place edits = %v, want %v", got, want)
			}
		})
	}
}
