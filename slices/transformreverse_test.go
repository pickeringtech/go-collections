package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/slices"
	"reflect"
	"testing"
)

func ExampleReverse() {
	a := []int{1, 2, 3, 4, 5}
	b := slices.Reverse(a)
	for _, element := range b {
		fmt.Printf("element: %v\n", element)
	}

	// Output:
	// element: 5
	// element: 4
	// element: 3
	// element: 2
	// element: 1
}

func TestReverse(t *testing.T) {
	type args[T any] struct {
		input []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "reverses the input",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
			},
			want: []int{5, 4, 3, 2, 1},
		},
		{
			name: "nil input results in nil output",
			args: args[int]{
				input: nil,
			},
			want: nil,
		},
		{
			name: "empty input results in empty output",
			args: args[int]{
				input: []int{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Reverse(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
