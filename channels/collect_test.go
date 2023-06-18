package channels_test

import (
	"github.com/pickeringtech/go-collectionutil/channels"
	"github.com/pickeringtech/go-collectionutil/maps"
	"reflect"
	"testing"
)

func TestCollectAsSlice(t *testing.T) {
	type args[T any] struct {
		input <-chan T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "correctly collect as a slice",
			args: args[int]{
				input: channels.FromSlice([]int{1, 19, 21, 3, -1, 100}),
			},
			want: []int{1, 19, 21, 3, -1, 100},
		},
		{
			name: "empty input provides nil output",
			args: args[int]{
				input: channels.FromSlice([]int{}),
			},
			want: nil,
		},
		{
			name: "nil input provides nil output",
			args: args[int]{
				input: channels.FromSlice[int](nil),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := channels.CollectAsSlice(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectAsSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollectAsMap(t *testing.T) {
	type args[I any, OK comparable, OV any] struct {
		input <-chan I
		fn    channels.MapBuilderFunc[I, OK, OV]
	}
	type testCase[I any, OK comparable, OV any] struct {
		name string
		args args[I, OK, OV]
		want map[OK]OV
	}
	tests := []testCase[string, string, int]{
		{
			name: "converts and collect elements into a map as expected",
			args: args[string, string, int]{
				input: channels.FromSlice([]string{"hello", "generous", "and", "glorious", "world"}),
				fn: func(input string) maps.Entry[string, int] {
					return maps.Entry[string, int]{
						Key:   input,
						Value: len(input),
					}
				},
			},
			want: map[string]int{
				"hello":    5,
				"generous": 8,
				"and":      3,
				"glorious": 8,
				"world":    5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := channels.CollectAsMap(tt.args.input, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectAsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
