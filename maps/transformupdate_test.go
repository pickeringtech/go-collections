package maps_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/maps"
	"reflect"
	"testing"
)

func ExampleUpdate() {
	input := map[int]string{
		1:  "one",
		-1: "negative one",
		0:  "zero",
	}
	update := map[int]string{
		-1: "negative one",
		0:  "zero",
		10: "ten",
	}
	out := maps.Update(input, update)

	for k, v := range out {
		fmt.Printf("%v: %v\n", k, v)
	}

	// Unordered output:
	// 1: one
	// -1: negative one
	// 0: zero
	// 10: ten
}

func TestUpdate(t *testing.T) {
	type args[K comparable, V any] struct {
		input  map[K]V
		update map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[int, string]{
		{
			name: "flattens input and update into a single map",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					10: "ten",
				},
				update: map[int]string{
					-1: "negative one",
					0:  "zero",
				},
			},
			want: map[int]string{
				1:  "one",
				-1: "negative one",
				0:  "zero",
				10: "ten",
			},
		},
		{
			name: "conflicting keys in input and update overwrites input",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					10: "ten",
				},
				update: map[int]string{
					-1: "negative one",
					1:  "one million",
				},
			},
			want: map[int]string{
				1:  "one million",
				-1: "negative one",
				10: "ten",
			},
		},
		{
			name: "empty input provides only update map entries",
			args: args[int, string]{
				input: map[int]string{},
				update: map[int]string{
					-1: "negative one",
					0:  "zero",
				},
			},
			want: map[int]string{
				-1: "negative one",
				0:  "zero",
			},
		},
		{
			name: "empty update provides only input map entries",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					10: "ten",
				},
				update: map[int]string{},
			},
			want: map[int]string{
				1:  "one",
				10: "ten",
			},
		},
		{
			name: "empty input and update provides empty output",
			args: args[int, string]{
				input:  map[int]string{},
				update: map[int]string{},
			},
			want: map[int]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.Update(tt.args.input, tt.args.update)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
