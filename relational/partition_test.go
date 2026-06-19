package relational_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/relational"
)

func TestPartition(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	tests := []struct {
		name          string
		input         []int
		wantMatched   []int
		wantUnmatched []int
	}{
		{
			name:          "nil input yields two non-nil empty slices",
			input:         nil,
			wantMatched:   []int{},
			wantUnmatched: []int{},
		},
		{
			name:          "empty input yields two non-nil empty slices",
			input:         []int{},
			wantMatched:   []int{},
			wantUnmatched: []int{},
		},
		{
			name:          "splits preserving order in each half",
			input:         []int{1, 2, 3, 4, 5},
			wantMatched:   []int{2, 4},
			wantUnmatched: []int{1, 3, 5},
		},
		{
			name:          "all match",
			input:         []int{2, 4},
			wantMatched:   []int{2, 4},
			wantUnmatched: []int{},
		},
		{
			name:          "none match",
			input:         []int{1, 3},
			wantMatched:   []int{},
			wantUnmatched: []int{1, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched, unmatched := relational.Partition(tt.input, isEven)
			if matched == nil || unmatched == nil {
				t.Fatalf("Partition returned a nil slice: matched=%v unmatched=%v", matched, unmatched)
			}
			if !reflect.DeepEqual(matched, tt.wantMatched) {
				t.Errorf("matched = %v, want %v", matched, tt.wantMatched)
			}
			if !reflect.DeepEqual(unmatched, tt.wantUnmatched) {
				t.Errorf("unmatched = %v, want %v", unmatched, tt.wantUnmatched)
			}
		})
	}
}

func TestPartitionDoesNotMutateInput(t *testing.T) {
	input := []int{1, 2, 3}
	snapshot := []int{1, 2, 3}
	_, _ = relational.Partition(input, func(n int) bool { return n > 1 })
	if !reflect.DeepEqual(input, snapshot) {
		t.Errorf("Partition mutated input: %v, want %v", input, snapshot)
	}
}
