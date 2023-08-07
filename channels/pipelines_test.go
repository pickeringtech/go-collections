package channels_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/channels"
	"reflect"
	"strconv"
	"testing"
)

func ExamplePipeline_CollectAsSlice() {
	input := channels.FromSlice([]int{1, 2, 5, 4, 3})

	// Creates a new pipeline which totals and then stringifies the input channel.
	pipeline := channels.NewPipeline[int, string](input, func(input <-chan int) <-chan string {
		reducer := channels.Reduce(input, func(accumulator int, element int) int {
			return accumulator + element
		})

		stringifier := channels.Map[int, string](reducer, func(element int) string {
			return strconv.Itoa(element)
		})

		return stringifier
	})

	// Capture results in a slice.
	results := pipeline.CollectAsSlice()

	fmt.Printf("Results: %v", results)
	// Output: Results: [15]
}

func TestPipeline_CollectAsSlice(t *testing.T) {
	type testCase[I any, O any] struct {
		name string
		p    *channels.Pipeline[I, O]
		want []O
	}
	tests := []testCase[string, int]{
		{
			name: "correctly maps elements through pipeline",
			p: channels.NewPipeline[string, int](channels.FromSlice([]string{"one", "two", "three", "four", "five"}), func(input <-chan string) <-chan int {
				return channels.Map[string, int](input, func(element string) int {
					return len(element)
				})
			}),
			want: []int{3, 3, 5, 4, 4},
		},
		{
			name: "empty input produces nil output",
			p: channels.NewPipeline[string, int](channels.FromSlice([]string{}), func(input <-chan string) <-chan int {
				return channels.Map[string, int](input, func(element string) int {
					return len(element)
				})
			}),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.CollectAsSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectAsSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
