package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/slices"
	"reflect"
	"testing"
)

func ExampleReduce() {
	a := []int{1, 2, 3, 4, 5}
	b := slices.Reduce(a, slices.TotalReducer[int])
	fmt.Printf("total: %v\n", b)

	// Output:
	// total: 15
}

func TestReduce(t *testing.T) {
	type args[I any, O any] struct {
		input []I
		fn    slices.ReductionFunc[I, O]
	}
	type testCase[I any, O any] struct {
		name string
		args args[I, O]
		want O
	}
	tests := []testCase[int, int]{
		{
			name: "reduce can add up all inputs",
			args: args[int, int]{
				input: []int{1, 2, 3, 4, 5},
				fn:    slices.TotalReducer[int],
			},
			want: 15,
		},
		{
			name: "empty input provides zero value output",
			args: args[int, int]{
				input: []int{},
				fn:    slices.TotalReducer[int],
			},
			want: 0,
		},
		{
			name: "nil input provides zero value output",
			args: args[int, int]{
				input: []int{},
				fn:    slices.TotalReducer[int],
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Reduce(tt.args.input, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleNewCountOccurrencesReducer() {
	a := []int{1, 2, 3, 4, 5, 3, 2, 2, 5, 4, 1}
	b := slices.Reduce(a, slices.NewCountOccurrencesReducer[int, int]([]int{1, 2, 3}))
	fmt.Printf("occurrences: %v\n", b)

	// Output:
	// occurrences: 7
}

func TestReduce_CountOccurrences(t *testing.T) {
	type args[I any, O any] struct {
		input []I
	}
	type testCase[I any, O any] struct {
		name string
		args args[I, O]
		want O
	}
	tests := []testCase[int, int]{
		{
			name: "reduce can add up all inputs",
			args: args[int, int]{
				input: []int{1, 2, 3, 4, 5, 3, 2, 2, 5, 4, 1},
			},
			want: 7,
		},
		{
			name: "empty input provides zero value output",
			args: args[int, int]{
				input: []int{},
			},
			want: 0,
		},
		{
			name: "nil input provides zero value output",
			args: args[int, int]{
				input: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Reduce[int, int](tt.args.input, slices.NewCountOccurrencesReducer[int, int]([]int{1, 2, 3}))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}
