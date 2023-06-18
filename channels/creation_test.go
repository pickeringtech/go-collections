package channels_test

import (
	"github.com/pickeringtech/go-collectionutil/channels"
	"github.com/pickeringtech/go-collectionutil/maps"
	"github.com/pickeringtech/go-collectionutil/slices"
	"reflect"
	"testing"
)

func TestFromSlice(t *testing.T) {
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
			name: "converts slice to channel and reads back consistently",
			args: args[int]{
				input: []int{1, 10, 5, 19, -1},
			},
			want: []int{1, 10, 5, 19, -1},
		},
		{
			name: "empty input provides nil output",
			args: args[int]{
				input: []int{},
			},
			want: nil,
		},
		{
			name: "nil input provides nil output",
			args: args[int]{
				input: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := channels.FromSlice(tt.args.input)
			got := channels.CollectAsSlice(output)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromSlice() = %v, want %v", output, tt.want)
			}
		})
	}
}

func TestFromMap(t *testing.T) {
	type args[K comparable, V any] struct {
		input map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want []maps.Entry[K, V]
	}
	tests := []testCase[int, string]{
		{
			name: "",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					5:  "five",
					2:  "two",
					-1: "negative one",
				},
			},
			want: []maps.Entry[int, string]{
				{
					Key:   -1,
					Value: "negative one",
				},
				{
					Key:   1,
					Value: "one",
				},
				{
					Key:   2,
					Value: "two",
				},
				{
					Key:   5,
					Value: "five",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := channels.FromMap(tt.args.input)
			got := channels.CollectAsSlice(output)
			got = slices.SortByOrderedField(got, slices.AscendingSortFunc[int], func(m maps.Entry[int, string]) int {
				return m.Key
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
