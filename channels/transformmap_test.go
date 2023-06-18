package channels

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	type args[I any, O any] struct {
		input <-chan I
		fn    MapFunc[I, O]
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
				input: FromSlice[string]([]string{"one", "two", "three", "four", "five"}),
				fn: func(s string) int {
					return len(s)
				},
			},
			want: []int{3, 3, 5, 4, 4},
		},
		{
			name: "empty input provides nil output",
			args: args[string, int]{
				input: FromSlice[string]([]string{}),
				fn: func(s string) int {
					return len(s)
				},
			},
			want: nil,
		},
		{
			name: "nil input provides nil output",
			args: args[string, int]{
				input: FromSlice[string](nil),
				fn: func(s string) int {
					return len(s)
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := Map(tt.args.input, tt.args.fn)
			got := CollectAsSlice(output)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
