package maps_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/maps"
	"reflect"
	"testing"
)

func ExampleMap() {
	input := map[int]string{
		1:  "one",
		2:  "two",
		5:  "five",
		-1: "negative one",
	}
	out := maps.Map(input, func(key int, value string) (string, int) {
		return value, key
	})

	for k, v := range out {
		fmt.Printf("%v %v\n", k, v)
	}
	// Unordered output:
	// one 1
	// two 2
	// five 5
	// negative one -1
}

func TestMap(t *testing.T) {
	type args[K comparable, V any, OK comparable, OV any] struct {
		input map[K]V
		fn    maps.MapFunc[K, V, OK, OV]
	}
	type testCase[K comparable, V any, OK comparable, OV any] struct {
		name string
		args args[K, V, OK, OV]
		want map[OK]OV
	}
	tests := []testCase[int, string, string, int]{
		{
			name: "transforms as expected",
			args: args[int, string, string, int]{
				input: map[int]string{
					1:  "one",
					2:  "two",
					5:  "five",
					-1: "negative one",
				},
				fn: func(key int, value string) (string, int) {
					return value, key
				},
			},
			want: map[string]int{
				"one":          1,
				"two":          2,
				"five":         5,
				"negative one": -1,
			},
		},
		{
			name: "empty input provides empty output",
			args: args[int, string, string, int]{
				input: map[int]string{},
				fn: func(key int, value string) (string, int) {
					return value, key
				},
			},
			want: map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.Map(tt.args.input, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
