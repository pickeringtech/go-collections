package maps_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/maps"
	"reflect"
	"testing"
)

func ExampleFromKeys() {
	keys := []int{1, 2, 3, 4, 5}
	defaultVal := "default"
	out := maps.FromKeys(keys, defaultVal)

	fmt.Printf("result: %v", out[2])
	// Output: result: default
}

func TestFromKeys(t *testing.T) {
	type args[K comparable, V any] struct {
		keys       []K
		defaultVal V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[int, string]{
		{
			name: "creates a map as expected",
			args: args[int, string]{
				keys:       []int{1, 2, 3, 4, -1, 10},
				defaultVal: "default",
			},
			want: map[int]string{
				1:  "default",
				2:  "default",
				3:  "default",
				4:  "default",
				-1: "default",
				10: "default",
			},
		},
		{
			name: "empty input creates empty output",
			args: args[int, string]{
				keys:       []int{},
				defaultVal: "default",
			},
			want: map[int]string{},
		},
		{
			name: "nil input creates empty output",
			args: args[int, string]{
				keys:       nil,
				defaultVal: "default",
			},
			want: map[int]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.FromKeys(tt.args.keys, tt.args.defaultVal)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
