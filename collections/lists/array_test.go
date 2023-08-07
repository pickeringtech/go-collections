package lists_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/collections/lists"
	"reflect"
	"testing"
)

func ExampleArray_AllMatch() {
	a := lists.NewArray(3, 4)
	match := a.AllMatch(func(a int) bool {
		return a > 2 && a < 5
	})
	fmt.Printf("Matches 1: %v\n", match)

	a = lists.NewArray(2, 3, 4)
	match = a.AllMatch(func(a int) bool {
		return a > 2 && a < 5
	})
	fmt.Printf("Matches 2: %v\n", match)

	// Output:
	// Matches 1: true
	// Matches 2: false
}

func TestArray_AllMatch(t *testing.T) {
	type args[T any] struct {
		fn func(T) bool
	}
	type testCase[T any] struct {
		name string
		a    lists.Array[T]
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "all matches",
			a:    []int{3, 4},
			args: args[int]{
				fn: func(a int) bool {
					return a > 2 && a < 5
				},
			},
			want: true,
		},
		{
			name: "do not all match",
			a:    []int{2, 3, 4},
			args: args[int]{
				fn: func(a int) bool {
					return a > 2 && a < 5
				},
			},
			want: false,
		},
		{
			name: "empty input provides true",
			a:    []int{},
			args: args[int]{
				fn: func(a int) bool {
					return a > 2 && a < 5
				},
			},
			want: false,
		},
		{
			name: "nil input provides true",
			a:    nil,
			args: args[int]{
				fn: func(a int) bool {
					return a > 2 && a < 5
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.AllMatch(tt.args.fn)
			if got != tt.want {
				t.Errorf("AllMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleArray_AnyMatch() {
	arr := lists.NewArray(4, 5, 3, 1, 2)

	match := arr.AnyMatch(func(a int) bool {
		return a == 3
	})
	fmt.Printf("Matches 1: %v\n", match)

	arr = lists.NewArray(4, 5, 1, 2)
	match = arr.AnyMatch(func(a int) bool {
		return a == 3
	})
	fmt.Printf("Matches 2: %v\n", match)

	// Output:
	// Matches 1: true
	// Matches 2: false
}

func TestArray_AnyMatch(t *testing.T) {
	type args[T any] struct {
		fn func(T) bool
	}
	type testCase[T any] struct {
		name string
		a    lists.Array[T]
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "matches with first element",
			a:    []int{3, 4, 5, 1, 2},
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: true,
		},
		{
			name: "matches with last element",
			a:    []int{4, 5, 1, 2, 3},
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: true,
		},
		{
			name: "matches with middle element",
			a:    []int{4, 5, 3, 1, 2},
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: true,
		},
		{
			name: "no match",
			a:    []int{4, 5, 1, 2},
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: false,
		},
		{
			name: "empty input provides false",
			a:    []int{},
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: false,
		},
		{
			name: "nil input provides false",
			a:    nil,
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.AnyMatch(tt.args.fn); got != tt.want {
				t.Errorf("AnyMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArray_Filter(t *testing.T) {
	type args[T any] struct {
		fn func(T) bool
	}
	type testCase[T any] struct {
		name string
		a    lists.Array[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "filters out values outside range",
			a:    []int{1, 2, 3, 4, 5},
			args: args[int]{
				fn: func(i int) bool {
					return i > 2 && i < 5
				},
			},
			want: []int{3, 4},
		},
		{
			name: "empty input provides nil output",
			a:    []int{},
			args: args[int]{
				fn: func(i int) bool {
					return i > 2 && i < 5
				},
			},
			want: nil,
		},
		{
			name: "nil input provides nil output",
			a:    nil,
			args: args[int]{
				fn: func(i int) bool {
					return i > 2 && i < 5
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Filter(tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
