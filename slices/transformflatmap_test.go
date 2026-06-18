package slices_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/slices"
)

func ExampleFlatMap() {
	input := []int{1, 2, 3}
	output := slices.FlatMap(input, func(n int) []int {
		return []int{n, n * 10}
	})
	fmt.Printf("Output: %v\n", output)

	// Output: Output: [1 10 2 20 3 30]
}

func TestFlatMap(t *testing.T) {
	type args struct {
		input []int
		fun   slices.FlatMapFunc[int, int]
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "expands each element",
			args: args{
				input: []int{1, 2, 3},
				fun:   func(n int) []int { return []int{n, n * 10} },
			},
			want: []int{1, 10, 2, 20, 3, 30},
		},
		{
			name: "empty result slices drop elements",
			args: args{
				input: []int{1, 2, 3, 4},
				fun: func(n int) []int {
					if n%2 == 0 {
						return []int{n}
					}
					return nil
				},
			},
			want: []int{2, 4},
		},
		{
			name: "nil input yields non-nil empty output",
			args: args{
				input: nil,
				fun:   func(n int) []int { return []int{n} },
			},
			want: []int{},
		},
		{
			name: "every call returning nothing yields non-nil empty output",
			args: args{
				input: []int{1, 2, 3},
				fun:   func(n int) []int { return nil },
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.FlatMap(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlatMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFlatMap_ChangesType exercises the I != O path.
func TestFlatMap_ChangesType(t *testing.T) {
	input := []string{"ab", "c"}
	got := slices.FlatMap(input, func(s string) []rune {
		return []rune(s)
	})
	want := []rune{'a', 'b', 'c'}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FlatMap() = %v, want %v", got, want)
	}
}
