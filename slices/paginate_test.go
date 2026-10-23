package slices_test

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/pickeringtech/go-collections/slices"
)

func ExamplePaginate() {
	results := []string{"a", "b", "c", "d", "e"}

	first := slices.Paginate(results, 0, 2)
	last := slices.Paginate(results, 2, 2)

	fmt.Printf("first: %v, last: %v", first, last)
	// Output: first: [a b], last: [e]
}

func TestPaginate(t *testing.T) {
	type args[T any] struct {
		slice     []T
		pageIndex int
		pageSize  int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "provides first page of results",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				pageIndex: 0,
				pageSize:  5,
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "provides second page of results",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				pageIndex: 1,
				pageSize:  5,
			},
			want: []int{6, 7, 8, 9, 10},
		},
		{
			name: "provides all results when page size is larger than slice",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				pageIndex: 0,
				pageSize:  15,
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name: "provides as many results as possible when page exceeds slice length",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				pageIndex: 1,
				pageSize:  10,
			},
			want: []int{11, 12, 13},
		},
		{
			name: "provides non-nil empty results when page index is less than 0",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				pageIndex: -1,
				pageSize:  10,
			},
			want: []int{},
		},
		{
			name: "provides non-nil empty results when page size is zero",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5},
				pageIndex: 0,
				pageSize:  0,
			},
			want: []int{},
		},
		{
			name: "provides non-nil empty results when page size is negative",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5},
				pageIndex: 1,
				pageSize:  -2,
			},
			want: []int{},
		},
		{
			name: "provides non-nil empty results when page is entirely beyond slice length",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				pageIndex: 1,
				pageSize:  10,
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Paginate(tt.args.slice, tt.args.pageIndex, tt.args.pageSize)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Paginate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPaginate_LargeArgumentsDoNotOverflow pins that page arithmetic is
// overflow-safe. pageIndex*pageSize used to wrap around int for large
// arguments, either panicking on a negative offset or landing on a wrong page
// (an offset that wraps to exactly 0 returned the first page).
func TestPaginate_LargeArgumentsDoNotOverflow(t *testing.T) {
	// half*half == 2^IntSize, which wraps to exactly 0 on every architecture.
	half := 1 << (strconv.IntSize / 2)
	input := []int{1, 2, 3, 4, 5}
	tests := []struct {
		name      string
		pageIndex int
		pageSize  int
		want      []int
	}{
		{name: "huge page index is beyond the slice", pageIndex: math.MaxInt, pageSize: 2, want: []int{}},
		{name: "huge page size on a later page is beyond the slice", pageIndex: 2, pageSize: math.MaxInt, want: []int{}},
		{name: "offset that wraps negative is beyond the slice", pageIndex: math.MaxInt/2 + 1, pageSize: 2, want: []int{}},
		{name: "offset that wraps to zero is beyond the slice, not the first page", pageIndex: half, pageSize: half, want: []int{}},
		{name: "huge page size on the first page returns everything", pageIndex: 0, pageSize: math.MaxInt, want: []int{1, 2, 3, 4, 5}},
		{name: "last partial page is still found", pageIndex: 2, pageSize: 2, want: []int{5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Paginate(input, tt.pageIndex, tt.pageSize)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Paginate(_, %d, %d) = %v, want %v", tt.pageIndex, tt.pageSize, got, tt.want)
			}
		})
	}
}
