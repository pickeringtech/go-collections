package maps_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/maps"
	"github.com/pickeringtech/go-collectionutil/slices"
	"reflect"
	"testing"
)

func ExampleClear() {
	input := map[int]string{
		1:  "one",
		-1: "negative one",
		10: "ten",
	}
	maps.Clear(input)
	fmt.Printf("%v", input)
	// Output: map[]
}

func TestClear(t *testing.T) {
	type args[K comparable, V any] struct {
		input map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
	}
	tests := []testCase[int, string]{
		{
			name: "clears the provided input map as expected",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maps.Clear(tt.args.input)
			if len(tt.args.input) > 0 {
				t.Errorf("clear did not remove all entries in the input map: %v", tt.args.input)
			}
		})
	}
}

func ExampleContainsValue() {
	input := map[int]string{
		1:  "one",
		-1: "negative one",
		10: "ten",
	}
	value := "negative one"
	result := maps.ContainsValue(input, value)
	fmt.Printf("result: %v", result)
	// Output: result: true
}

func TestContainsValue(t *testing.T) {
	type args[K comparable, V comparable] struct {
		input map[K]V
		value V
	}
	type testCase[K comparable, V comparable] struct {
		name string
		args args[K, V]
		want bool
	}
	tests := []testCase[int, string]{
		{
			name: "finds the values as expected",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
				},
				value: "negative one",
			},
			want: true,
		},
		{
			name: "does not find missing value",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
				},
				value: "negative two",
			},
			want: false,
		},
		{
			name: "empty input results in no find",
			args: args[int, string]{
				input: map[int]string{},
				value: "negative two",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.ContainsValue(tt.args.input, tt.args.value)
			if got != tt.want {
				t.Errorf("ContainsValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleCopy() {
	input := map[int]string{
		1: "one",
	}
	output := maps.Copy(input)
	maps.Clear(input)
	fmt.Printf("%v", output)
	// Output: map[1:one]
}

func TestCopy(t *testing.T) {
	type args[K comparable, V any] struct {
		input map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[int, string]{
		{
			name: "copies whole map in memory",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
				},
			},
			want: map[int]string{
				1:  "one",
				-1: "negative one",
				10: "ten",
			},
		},
		{
			name: "empty input provides empty output",
			args: args[int, string]{
				input: map[int]string{},
			},
			want: map[int]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.Copy(tt.args.input)
			maps.Clear(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleGetMany() {
	input := map[int]string{
		1:  "one",
		-1: "negative one",
		10: "ten",
	}

	results := maps.GetMany(input, []int{1, 100, 10})
	fmt.Printf("results: %v", results)
	// Output: results: [one ten]
}

func TestGetMany(t *testing.T) {
	type args[K comparable, V any] struct {
		input map[K]V
		keys  []K
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want []V
	}
	tests := []testCase[int, string]{
		{
			name: "retrieves many values and omits a value if it cannot be found",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
				},
				keys: []int{1, 100, 10},
			},
			want: []string{"one", "ten"},
		},
		{
			name: "empty input provides nil output",
			args: args[int, string]{
				input: map[int]string{},
				keys:  []int{1, 100, 10},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.GetMany(tt.args.input, tt.args.keys)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleGetManyOrDefault() {
	input := map[int]string{
		1:  "one",
		-1: "negative one",
		10: "ten",
	}

	results := maps.GetManyOrDefault(input, []int{1, 100, 10}, "missing")
	fmt.Printf("results: %v", results)
	// Output: results: [one missing ten]
}

func TestGetManyOrDefault(t *testing.T) {
	type args[K comparable, V any] struct {
		input      map[K]V
		keys       []K
		defaultVal V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want []V
	}
	tests := []testCase[int, string]{
		{
			name: "retrieves many values and uses the default if an entry cannot be found",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
				},
				keys:       []int{1, 100, 10},
				defaultVal: "missing",
			},
			want: []string{"one", "missing", "ten"},
		},
		{
			name: "empty input provides slice of default values",
			args: args[int, string]{
				input:      map[int]string{},
				keys:       []int{1, 100, 10},
				defaultVal: "missing",
			},
			want: []string{"missing", "missing", "missing"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.GetManyOrDefault(tt.args.input, tt.args.keys, tt.args.defaultVal)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetManyOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleGetOrDefault() {
	input := map[int]string{
		1:  "one",
		-1: "negative one",
		10: "ten",
	}

	result := maps.GetOrDefault(input, -1, "missing")
	fmt.Printf("result: %v", result)
	// Output: result: negative one
}

func TestGetOrDefault(t *testing.T) {
	type args[K comparable, V any] struct {
		input  map[K]V
		key    K
		orElse V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want V
	}
	tests := []testCase[int, string]{
		{
			name: "finds expected value",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
				},
				key:    -1,
				orElse: "missing",
			},
			want: "negative one",
		},
		{
			name: "search failure provides default value",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
				},
				key:    -2,
				orElse: "missing",
			},
			want: "missing",
		},
		{
			name: "empty input provides default value",
			args: args[int, string]{
				input:  map[int]string{},
				key:    -2,
				orElse: "missing",
			},
			want: "missing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.GetOrDefault(tt.args.input, tt.args.key, tt.args.orElse)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleItems() {
	input := map[int]string{
		1: "one",
	}

	results := maps.Items(input)
	fmt.Printf("results: %v", results)
	// Output: results: [{1 one}]
}

func TestItems(t *testing.T) {
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
			name: "provides each entry in the map as an element in the output",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
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
					Key:   10,
					Value: "ten",
				},
			},
		},
		{
			name: "empty input provides nil output",
			args: args[int, string]{
				input: map[int]string{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.Items(tt.args.input)
			got = slices.SortByOrderedField(got, slices.AscendingSortFunc[int], func(e maps.Entry[int, string]) int {
				return e.Key
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Items() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleKeys() {
	input := map[int]string{
		1: "one",
	}

	results := maps.Keys(input)
	fmt.Printf("results: %v", results)
	// Output: results: [1]
}

func TestKeys(t *testing.T) {
	type args[K comparable, V any] struct {
		input map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want []K
	}
	tests := []testCase[int, string]{
		{
			name: "provides the keys as a slice",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
				},
			},
			want: []int{-1, 1, 10},
		},
		{
			name: "empty input provides nil output",
			args: args[int, string]{
				input: map[int]string{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.Keys(tt.args.input)
			got = slices.SortOrderedAsc(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleValues() {
	input := map[int]string{
		1: "one",
	}

	results := maps.Values(input)
	fmt.Printf("results: %v", results)
	// Output: results: [one]
}

func TestValues(t *testing.T) {
	type args[K comparable, V any] struct {
		input map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want []V
	}
	tests := []testCase[int, string]{
		{
			name: "provides the values as a slice",
			args: args[int, string]{
				input: map[int]string{
					1:  "one",
					-1: "negative one",
					10: "ten",
				},
			},
			want: []string{"negative one", "one", "ten"},
		},
		{
			name: "empty input provides nil output",
			args: args[int, string]{
				input: map[int]string{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.Values(tt.args.input)
			got = slices.SortOrderedAsc(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}
