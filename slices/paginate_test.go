package slices

import (
	"reflect"
	"testing"
)

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
			name: "provides nil results when page index is less than 0",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				pageIndex: -1,
				pageSize:  10,
			},
			want: nil,
		},
		{
			name: "provides nil results when page is entirely beyond slice length",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				pageIndex: 1,
				pageSize:  10,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Paginate(tt.args.slice, tt.args.pageIndex, tt.args.pageSize)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Paginate() = %v, want %v", got, tt.want)
			}
		})
	}
}
