package maps_test

import (
	"github.com/pickeringtech/go-collectionutil/maps"
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	type args[K comparable, V any] struct {
		input map[K]V
		fn    maps.FilterFunc[K, V]
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[int, string]{
		{
			name: "filters out negative or zero keys",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					0:  "zero",
					10: "ten",
				},
				fn: func(key int, value string) bool {
					return key > 0
				},
			},
			want: map[int]string{
				1:  "one",
				10: "ten",
			},
		},
		{
			name: "empty input provides empty output",
			args: args[int, string]{
				input: map[int]string{},
				fn: func(key int, value string) bool {
					return key > 0
				},
			},
			want: map[int]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origInput := maps.Copy(tt.args.input)
			got := maps.Filter(tt.args.input, tt.args.fn)
			if !reflect.DeepEqual(tt.args.input, origInput) {
				t.Errorf("Filter() changed input - wanted %v, got %v", origInput, got)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
