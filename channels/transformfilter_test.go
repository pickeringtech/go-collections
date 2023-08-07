package channels

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleFilter() {
	input := FromSlice([]string{"hello", "everyone", "world", "goodness", "gracious"})
	output := Filter(input, func(element string) bool {
		return len(element) > 5
	})

	// Capture results in a slice.
	results := CollectAsSlice(output)

	// Print results.
	fmt.Printf("Results: %v", results)
	// Output: Results: [everyone goodness gracious]
}

func TestFilter(t *testing.T) {
	type args[T any] struct {
		input <-chan T
		fn    FilterFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []string
	}
	tests := []testCase[string]{
		{
			name: "filters out words with 5 characters or less",
			args: args[string]{
				input: FromSlice([]string{"hello", "everyone", "world", "goodness", "gracious"}),
				fn: func(element string) bool {
					return len(element) > 5
				},
			},
			want: []string{"everyone", "goodness", "gracious"},
		},
		{
			name: "empty input provides nil output",
			args: args[string]{
				input: FromSlice([]string{}),
				fn: func(element string) bool {
					return len(element) > 5
				},
			},
			want: nil,
		},
		{
			name: "nil input provides nil output",
			args: args[string]{
				input: FromSlice[string](nil),
				fn: func(element string) bool {
					return len(element) > 5
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := Filter(tt.args.input, tt.args.fn)
			got := CollectAsSlice(output)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
