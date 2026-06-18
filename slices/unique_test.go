package slices_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/slices"
)

func ExampleUnique() {
	input := []int{3, 1, 3, 2, 1}
	output := slices.Unique(input)
	fmt.Printf("Output: %v\n", output)

	// Output: Output: [3 1 2]
}

func ExampleUniqueBy() {
	type person struct {
		name string
		dept string
	}
	people := []person{
		{"alice", "eng"},
		{"bob", "eng"},
		{"carol", "sales"},
	}
	output := slices.UniqueBy(people, func(p person) string {
		return p.dept
	})
	fmt.Printf("Output: %v\n", output)

	// Output: Output: [{alice eng} {carol sales}]
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "drops later duplicates preserving order",
			input: []int{3, 1, 3, 2, 1},
			want:  []int{3, 1, 2},
		},
		{
			name:  "already unique is unchanged",
			input: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
		{
			name:  "all duplicates collapse to one",
			input: []int{7, 7, 7},
			want:  []int{7},
		},
		{
			name:  "nil input yields non-nil empty output",
			input: nil,
			want:  []int{},
		},
		{
			name:  "empty input yields non-nil empty output",
			input: []int{},
			want:  []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Unique(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueBy(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		keyFn slices.KeyFunc[string, int]
		want  []string
	}{
		{
			name:  "keeps first element per key",
			input: []string{"a", "bb", "cc", "ddd"},
			keyFn: func(s string) int { return len(s) },
			want:  []string{"a", "bb", "ddd"},
		},
		{
			name:  "already unique by key is unchanged",
			input: []string{"a", "bb", "ccc"},
			keyFn: func(s string) int { return len(s) },
			want:  []string{"a", "bb", "ccc"},
		},
		{
			name:  "nil input yields non-nil empty output",
			input: nil,
			keyFn: func(s string) int { return len(s) },
			want:  []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.UniqueBy(tt.input, tt.keyFn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestUnique_DoesNotMutateInput guards the no-mutation contract.
func TestUnique_DoesNotMutateInput(t *testing.T) {
	input := []int{1, 1, 2}
	_ = slices.Unique(input)
	if !reflect.DeepEqual(input, []int{1, 1, 2}) {
		t.Errorf("Unique mutated the input: %v", input)
	}
}
