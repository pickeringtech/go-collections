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
			got := Concatenate(tt.args.inputA, tt.args.inputB)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Concatenate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "duplicates the input into a new slice",
			args: args{
				input: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "nil input provides nil output",
			args: args{
				input: nil,
			},
			want: nil,
		},
		{
			name: "empty input provides nil output",
			args: args{
				input: []int{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Copy(tt.args.input)
			tt.args.input = append(tt.args.input, 45)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
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
			got := Delete(tt.args.input, tt.args.index)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPush(t *testing.T) {
	type args struct {
		input       []int
		newElements []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "pushes the new elements to the end of the input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{4, 5, 6},
			},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "nil input slice results in only the new elements",
			args: args{
				input:       nil,
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "empty input slice results in only the new elements",
			args: args{
				input:       []int{},
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "nil new elements results in only original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "empty new elements results in only original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Push(tt.args.input, tt.args.newElements...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Push() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPushFront(t *testing.T) {
	type args struct {
		input       []int
		newElements []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "adds the new elements to the front of the input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6, 1, 2, 3},
		},
		{
			name: "nil new elements results in original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "empty new elements results in original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "nil input slice results in only new elements",
			args: args{
				input:       nil,
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "empty input slice results in only new elements",
			args: args{
				input:       []int{},
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PushFront(tt.args.input, tt.args.newElements...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PushFront() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPop(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantSli []int
	}{
		{
			name: "gets the last element from the input returning the smaller slice",
			args: args{
				input: []int{1, 2, 3, 4},
			},
			want:    4,
			wantSli: []int{1, 2, 3},
		},
		{
			name: "nil input provides zero value for type and nil resulting slice",
			args: args{
				input: nil,
			},
			want:    0,
			wantSli: nil,
		},
		{
			name: "empty input provides zero value for type and nil resulting slice",
			args: args{
				input: []int{},
			},
			want:    0,
			wantSli: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Pop(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pop() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.wantSli) {
				t.Errorf("Pop() got1 = %v, want %v", got1, tt.wantSli)
			}
		})
	}
}
