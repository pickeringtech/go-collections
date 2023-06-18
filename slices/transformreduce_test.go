package slices

import (
	"reflect"
	"testing"
)

func TestReduce(t *testing.T) {
	type args[I any, O any] struct {
		input []I
		fn    ReductionFunc[I, O]
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
				fn:    ReductionTotalFunc[int],
			},
			want: 15,
		},
		{
			name: "empty input provides zero value output",
			args: args[int, int]{
				input: []int{},
				fn:    ReductionTotalFunc[int],
			},
			want: 0,
		},
		{
			name: "nil input provides zero value output",
			args: args[int, int]{
				input: []int{},
				fn:    ReductionTotalFunc[int],
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reduce(tt.args.input, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
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
			got := Reduce[int, int](tt.args.input, NewReductionCountOccurrencesFunc[int, int]([]int{1, 2, 3}))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}
