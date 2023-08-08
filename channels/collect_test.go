package channels_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/channels"
	"github.com/pickeringtech/go-collections/maps"
	"reflect"
	"testing"
)

func TestBuildMapFromEntries(t *testing.T) {
	type args[K comparable, V any] struct {
		entries []maps.Entry[K, V]
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[string, int]{
		{
			name: "builds a map from the input entries",
			args: args[string, int]{
				entries: []maps.Entry[string, int]{
					{
						Key:   "hello",
						Value: 10,
					},
					{
						Key:   "world",
						Value: 20,
					},
				},
			},
			want: map[string]int{
				"hello": 10,
				"world": 20,
			},
		},
		{
			name: "empty input produces empty output",
			args: args[string, int]{
				entries: []maps.Entry[string, int]{},
			},
			want: map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := channels.BuildMapFromEntries(tt.args.entries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildMapFromEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleCollectAsSlice() {
	input := channels.FromSlice([]int{1, 2, 3, 4, 5})
	mapOut := channels.Map(input, func(element int) int {
		return element * 2
	})
	output := channels.CollectAsSlice(mapOut)
	fmt.Printf("result: %v", output)
	// Output: result: [2 4 6 8 10]
}

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

func ExampleCollectAsMap() {
	input := channels.FromSlice([]string{"hello", "generous", "and", "glorious", "world"})
	output := channels.CollectAsMap(input, func(element string) maps.Entry[string, int] {
		return maps.Entry[string, int]{
			Key:   element,
			Value: len(element),
		}
	})
	fmt.Printf("result: %v", output)
	// Output: result: map[and:3 generous:8 glorious:8 hello:5 world:5]
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
