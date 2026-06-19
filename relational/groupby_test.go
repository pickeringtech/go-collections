package relational_test

import (
	"iter"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/relational"
)

func sliceSeq[T any](values []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

func parity(n int) string {
	if n%2 == 0 {
		return "even"
	}
	return "odd"
}

func TestGroupBy(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want map[string][]int
	}{
		{
			name: "nil input yields non-nil empty map",
			args: args{input: nil},
			want: map[string][]int{},
		},
		{
			name: "empty input yields non-nil empty map",
			args: args{input: []int{}},
			want: map[string][]int{},
		},
		{
			name: "groups preserve first-seen order",
			args: args{input: []int{1, 2, 3, 4, 5}},
			want: map[string][]int{
				"odd":  {1, 3, 5},
				"even": {2, 4},
			},
		},
		{
			name: "single group",
			args: args{input: []int{2, 4, 6}},
			want: map[string][]int{"even": {2, 4, 6}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relational.GroupBy(tt.args.input, parity)
			if got == nil {
				t.Fatalf("GroupBy returned nil map")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupBySeq(t *testing.T) {
	type args struct {
		seq iter.Seq[int]
	}
	tests := []struct {
		name string
		args args
		want map[string][]int
	}{
		{
			name: "nil seq yields non-nil empty map",
			args: args{seq: nil},
			want: map[string][]int{},
		},
		{
			name: "empty seq yields non-nil empty map",
			args: args{seq: sliceSeq([]int{})},
			want: map[string][]int{},
		},
		{
			name: "groups preserve first-seen order",
			args: args{seq: sliceSeq([]int{1, 2, 3, 4, 5})},
			want: map[string][]int{
				"odd":  {1, 3, 5},
				"even": {2, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relational.GroupBySeq(tt.args.seq, parity)
			if got == nil {
				t.Fatalf("GroupBySeq returned nil map")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupBySeq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountBy(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			name: "nil input yields non-nil empty map",
			args: args{input: nil},
			want: map[string]int{},
		},
		{
			name: "empty input yields non-nil empty map",
			args: args{input: []int{}},
			want: map[string]int{},
		},
		{
			name: "counts per group",
			args: args{input: []int{1, 2, 3, 4, 5}},
			want: map[string]int{"odd": 3, "even": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relational.CountBy(tt.args.input, parity)
			if got == nil {
				t.Fatalf("CountBy returned nil map")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGroupByDoesNotMutateInput guards the non-mutation contract.
func TestGroupByDoesNotMutateInput(t *testing.T) {
	input := []int{1, 2, 3}
	snapshot := []int{1, 2, 3}
	_ = relational.GroupBy(input, parity)
	if !reflect.DeepEqual(input, snapshot) {
		t.Errorf("GroupBy mutated input: %v, want %v", input, snapshot)
	}
}
