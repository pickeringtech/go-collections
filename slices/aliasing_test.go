package slices_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/slices"
)

// TestResultsDoNotAliasInput pins the non-mutating contract for every function
// that once handed back a view of its input (a sub-slice, or the input itself on
// an out-of-range or no-op path): the result must own fresh storage, so writing
// to it or appending to it can never reach the caller's backing array. The input
// is built with spare capacity, the condition under which an append to an
// aliasing result silently overwrites the caller's data beyond len.
func TestResultsDoNotAliasInput(t *testing.T) {
	pop := func(input []int) []int {
		_, _, rest := slices.Pop(input)
		return rest
	}
	popFront := func(input []int) []int {
		_, _, rest := slices.PopFront(input)
		return rest
	}
	tests := []struct {
		name string
		fn   func(input []int) []int
	}{
		{name: "Pop", fn: pop},
		{name: "PopFront", fn: popFront},
		{name: "SubSlice within range", fn: func(in []int) []int { return slices.SubSlice(in, 1, 3) }},
		{name: "SubSlice whole slice", fn: func(in []int) []int { return slices.SubSlice(in, 0, 4) }},
		{name: "SubSlice clamped range", fn: func(in []int) []int { return slices.SubSlice(in, -5, 10) }},
		{name: "SubSlice empty range", fn: func(in []int) []int { return slices.SubSlice(in, 2, 2) }},
		{name: "Paginate full page", fn: func(in []int) []int { return slices.Paginate(in, 0, 2) }},
		{name: "Paginate partial last page", fn: func(in []int) []int { return slices.Paginate(in, 1, 3) }},
		{name: "Delete in range", fn: func(in []int) []int { return slices.Delete(in, 1) }},
		{name: "Delete index too high", fn: func(in []int) []int { return slices.Delete(in, 4) }},
		{name: "Delete negative index", fn: func(in []int) []int { return slices.Delete(in, -1) }},
		{name: "Fill", fn: func(in []int) []int { return slices.Fill(in, 9) }},
		{name: "FillFrom negative index", fn: func(in []int) []int { return slices.FillFrom(in, 9, -1) }},
		{name: "FillTo index too high", fn: func(in []int) []int { return slices.FillTo(in, 9, 5) }},
		{name: "FillFromTo in range", fn: func(in []int) []int { return slices.FillFromTo(in, 9, 1, 3) }},
		{name: "FillFromTo backward range", fn: func(in []int) []int { return slices.FillFromTo(in, 9, 3, 1) }},
		{name: "Insert in range", fn: func(in []int) []int { return slices.Insert(in, 2, 9) }},
		{name: "Insert index too high", fn: func(in []int) []int { return slices.Insert(in, 5, 9) }},
		{name: "Insert negative index", fn: func(in []int) []int { return slices.Insert(in, -1, 9) }},
		{name: "Insert with no elements", fn: func(in []int) []int { return slices.Insert(in, 2) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]int, 4, 8)
			copy(input, []int{1, 2, 3, 4})
			backing := input[:cap(input)]
			backingBefore := slices.Copy(backing)

			got := tt.fn(input)
			if !reflect.DeepEqual(backing, backingBefore) {
				t.Fatalf("%s modified its input: backing array = %v, want %v", tt.name, backing, backingBefore)
			}

			for i := range got {
				got[i] = -1
			}
			_ = append(got, -2, -2, -2, -2)
			if !reflect.DeepEqual(backing, backingBefore) {
				t.Errorf("%s result aliases its input: backing array = %v after writing to the result, want %v",
					tt.name, backing, backingBefore)
			}
		})
	}
}
