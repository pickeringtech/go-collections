package slices_test

import (
	"github.com/pickeringtech/go-collectionutil/constraints"
	"github.com/pickeringtech/go-collectionutil/slices"
	"reflect"
	"testing"
)

func TestAscendingSortFunc(t *testing.T) {
	type args[T constraints.Ordered] struct {
		a T
		b T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "a < b == true",
			args: args[int]{
				a: 0,
				b: 1,
			},
			want: true,
		},
		{
			name: "a > b == false",
			args: args[int]{
				a: 1,
				b: 0,
			},
			want: false,
		},
		{
			name: "a == b == false",
			args: args[int]{
				a: 0,
				b: 0,
			},
			want: false,
		},
		{
			name: "(a < b) < 0",
			args: args[int]{
				a: -2,
				b: -1,
			},
			want: true,
		},
		{
			name: "(a > b) < 0",
			args: args[int]{
				a: -1,
				b: -2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.AscendingSortFunc(tt.args.a, tt.args.b)
			if got != tt.want {
				t.Errorf("AscendingSortFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDescendingSortFunc(t *testing.T) {
	type args[T constraints.Ordered] struct {
		a T
		b T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "a < b == false",
			args: args[int]{
				a: 0,
				b: 1,
			},
			want: false,
		},
		{
			name: "a > b == true",
			args: args[int]{
				a: 1,
				b: 0,
			},
			want: true,
		},
		{
			name: "a == b == true",
			args: args[int]{
				a: 0,
				b: 0,
			},
			want: false,
		},
		{
			name: "(a < b) < 0",
			args: args[int]{
				a: -2,
				b: -1,
			},
			want: false,
		},
		{
			name: "(a > b) < 0",
			args: args[int]{
				a: -1,
				b: -2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.DescendingSortFunc(tt.args.a, tt.args.b)
			if got != tt.want {
				t.Errorf("DescendingSortFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSort(t *testing.T) {
	type args[T any] struct {
		input []T
		fun   slices.SortFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts numbers ascending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
				fun:   slices.AscendingSortFunc[int],
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "sorts numbers descending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
				fun:   slices.DescendingSortFunc[int],
			},
			want: []int{9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name: "handles nil input",
			args: args[int]{
				input: nil,
				fun:   slices.DescendingSortFunc[int],
			},
			want: nil,
		},
		{
			name: "handles empty input",
			args: args[int]{
				input: []int{},
				fun:   slices.DescendingSortFunc[int],
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.args.input
			orig := append(input[:0:0], input...)
			got := slices.Sort(input, tt.args.fun)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(orig, input) {
				t.Errorf("Sort() changed input - no changes expected")
			}
		})
	}
}

func TestSortInPlace(t *testing.T) {
	type args[T any] struct {
		input []T
		fun   slices.SortFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts numbers ascending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
				fun:   slices.AscendingSortFunc[int],
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "sorts numbers descending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
				fun:   slices.DescendingSortFunc[int],
			},
			want: []int{9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name: "handles nil input",
			args: args[int]{
				input: nil,
				fun:   slices.DescendingSortFunc[int],
			},
			want: nil,
		},
		{
			name: "handles empty input",
			args: args[int]{
				input: []int{},
				fun:   slices.DescendingSortFunc[int],
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.args.input
			slices.SortInPlace(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(input, tt.want) {
				t.Errorf("SortInPlace() = %v, want %v", input, tt.want)
			}
		})
	}
}

func TestSortOrderedAsc(t *testing.T) {
	type args[T constraints.Ordered] struct {
		input []T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts numbers ascending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "handles nil input",
			args: args[int]{
				input: nil,
			},
			want: nil,
		},
		{
			name: "handles empty input",
			args: args[int]{
				input: []int{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.args.input
			orig := append(input[:0:0], input...)
			got := slices.SortOrderedAsc(input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortOrderedAsc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(orig, input) {
				t.Errorf("SortOrderedAsc() changed input - no changes expected")
			}
		})
	}
}

func TestSortOrderedAscInPlace(t *testing.T) {
	type args[T constraints.Ordered] struct {
		input []T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts numbers ascending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "handles nil input",
			args: args[int]{
				input: nil,
			},
			want: nil,
		},
		{
			name: "handles empty input",
			args: args[int]{
				input: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slices.SortOrderedAscInPlace(tt.args.input)
			if !reflect.DeepEqual(tt.args.input, tt.want) {
				t.Errorf("SortAscInPlace() = %v, want %v", tt.args.input, tt.want)
			}
		})
	}
}
