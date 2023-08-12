package channels_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/channels"
	"reflect"
	"testing"
)

func ExampleMap() {
	input := channels.FromSlice([]string{"one", "two", "three", "four", "five"})
	output := channels.Map(input, func(s string) int {
		return len(s)
	})

	// Capture results in a slice.
	results := channels.CollectAsSlice(output)

	// Print results.
	fmt.Printf("Results: %v", results)
	// Output: Results: [3 3 5 4 4]
}

func TestMap(t *testing.T) {
	type args[I any, O any] struct {
		input <-chan I
		fn    channels.MapFunc[I, O]
	}
	type testCase[I any, O any] struct {
		name string
		args args[I, O]
		want []O
	}
	tests := []testCase[string, int]{
		{
			name: "transforms elements correctly",
			args: args[string, int]{
				input: channels.FromSlice[string]([]string{"one", "two", "three", "four", "five"}),
				fn: func(s string) int {
					return len(s)
				},
			},
			want: []int{3, 3, 5, 4, 4},
		},
		{
			name: "empty input provides nil output",
			args: args[string, int]{
				input: channels.FromSlice[string]([]string{}),
				fn: func(s string) int {
					return len(s)
				},
			},
			want: nil,
		},
		{
			name: "nil input provides nil output",
			args: args[string, int]{
				input: channels.FromSlice[string](nil),
				fn: func(s string) int {
					return len(s)
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := channels.Map(tt.args.input, tt.args.fn)
			got := channels.CollectAsSlice(output)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
