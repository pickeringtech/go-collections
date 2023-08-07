package channels_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/channels"
	"reflect"
	"testing"
)

func ExampleReduce() {
	input := channels.FromSlice([]int{1, 2, 3, 4, 5})

	// Creates a new pipeline which totals the input channel.
	pipeline := channels.NewPipeline[int, int](input, func(input <-chan int) <-chan int {
		return channels.Reduce(input, func(accumulator int, element int) int {
			return accumulator + element
		})
	})

	// Capture results in a slice.
	results := pipeline.CollectAsSlice()

	fmt.Printf("Results: %v", results)
	// Output: Results: [15]
}

func TestReduce(t *testing.T) {
	type args[I any, O any] struct {
		input <-chan I
		fn    channels.ReduceFunc[I, O]
	}
	type testCase[I any, O any] struct {
		name string
		args args[I, O]
		want []O
	}
	tests := []testCase[int, int]{
		{
			name: "totals correctly",
			args: args[int, int]{
				input: channels.FromSlice([]int{1, 2, 3, 4, 5}),
				fn:    func(a int, b int) int { return a + b },
			},
			want: []int{15},
		},
		{
			name: "empty input results in zero output",
			args: args[int, int]{
				input: channels.FromSlice([]int{}),
				fn:    func(a int, b int) int { return a + b },
			},
			want: []int{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCh := channels.Reduce(tt.args.input, tt.args.fn)
			got := channels.CollectAsSlice(gotCh)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}
