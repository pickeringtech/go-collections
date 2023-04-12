package slices

import (
	"reflect"
	"testing"
)

func TestConcatenate(t *testing.T) {
	type args struct {
		inputA []int
		inputB []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "inputs are joined consecutively",
			args: args{
				inputA: []int{1, 2, 3},
				inputB: []int{4, 5, 6},
			},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "nil first input means result is second input",
			args: args{
				inputA: nil,
				inputB: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "empty first input means result is second input",
			args: args{
				inputA: []int{},
				inputB: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "nil second input means result is first input",
			args: args{
				inputA: []int{1, 2, 3},
				inputB: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "empty second input means result is first input",
			args: args{
				inputA: []int{1, 2, 3},
				inputB: []int{},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Concatenate(tt.args.inputA, tt.args.inputB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Concatenate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		input []int
		index int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "removes the element at the specified index",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 2,
			},
			want: []int{1, 2, 4},
		},
		{
			name: "removes the element at the last index",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 3,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "removes the zeroth element",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 0,
			},
			want: []int{2, 3, 4},
		},
		{
			name: "if index is beyond range the slice is not modified",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 4,
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "if index is below zero the slice is not modified",
			args: args{
				input: []int{1, 2, 3, 4},
				index: -1,
			},
			want: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Delete(tt.args.input, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
