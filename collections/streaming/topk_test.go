package streaming_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/heaps"
	"github.com/pickeringtech/go-collections/collections/streaming"
)

// topKOracle returns the k largest of values, highest first, as a non-nil
// slice — the brute-force buffer-sort-slice answer streaming.TopK must match.
func topKOracle(values []int, k int) []int {
	if k <= 0 {
		return []int{}
	}
	sorted := make([]int, len(values))
	copy(sorted, values)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	if k < len(sorted) {
		sorted = sorted[:k]
	}
	out := make([]int, len(sorted))
	copy(out, sorted)
	return out
}

func TestTopK_ResultMatchesSortOracle(t *testing.T) {
	type args struct {
		values []int
		k      int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{name: "nil stream results in empty output", args: args{values: nil, k: 3}, want: []int{}},
		{name: "empty stream results in empty output", args: args{values: []int{}, k: 3}, want: []int{}},
		{name: "k of zero retains nothing", args: args{values: []int{5, 1, 9}, k: 0}, want: []int{}},
		{name: "negative k retains nothing", args: args{values: []int{5, 1, 9}, k: -2}, want: []int{}},
		{name: "fewer elements than k returns all sorted", args: args{values: []int{3, 1, 2}, k: 5}, want: []int{3, 2, 1}},
		{name: "exactly k elements", args: args{values: []int{3, 1, 2}, k: 3}, want: []int{3, 2, 1}},
		{name: "more elements than k", args: args{values: []int{5, 1, 9, 3, 7, 2, 8}, k: 3}, want: []int{9, 8, 7}},
		{name: "duplicates at the boundary", args: args{values: []int{3, 1, 2, 2}, k: 2}, want: []int{3, 2}},
		{name: "all equal", args: args{values: []int{5, 5, 5, 5}, k: 2}, want: []int{5, 5}},
		{name: "already ascending", args: args{values: []int{1, 2, 3, 4, 5}, k: 2}, want: []int{5, 4}},
		{name: "already descending", args: args{values: []int{5, 4, 3, 2, 1}, k: 2}, want: []int{5, 4}},
		{name: "single element", args: args{values: []int{42}, k: 3}, want: []int{42}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			top := streaming.NewTopKOrdered[int](tt.args.k)
			for _, v := range tt.args.values {
				top.Add(v)
			}

			got := top.Result()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Result() = %v, want %v", got, tt.want)
			}
			oracle := topKOracle(tt.args.values, tt.args.k)
			if !reflect.DeepEqual(got, oracle) {
				t.Errorf("Result() = %v, oracle = %v", got, oracle)
			}
		})
	}
}

func TestTopK_LenTracksRetained(t *testing.T) {
	type args struct {
		values []int
		k      int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "empty stream", args: args{values: nil, k: 3}, want: 0},
		{name: "below capacity", args: args{values: []int{1, 2}, k: 5}, want: 2},
		{name: "at capacity", args: args{values: []int{1, 2, 3, 4, 5}, k: 3}, want: 3},
		{name: "over capacity stays capped", args: args{values: []int{1, 2, 3, 4, 5, 6}, k: 3}, want: 3},
		{name: "k of zero stays empty", args: args{values: []int{1, 2, 3}, k: 0}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			top := streaming.NewTopKOrdered[int](tt.args.k)
			for _, v := range tt.args.values {
				top.Add(v)
			}

			got := top.Len()
			if got != tt.want {
				t.Errorf("Len() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTopK_ResultDoesNotMutateState(t *testing.T) {
	top := streaming.NewTopKOrdered[int](3)
	for _, v := range []int{5, 1, 9, 3, 7} {
		top.Add(v)
	}

	first := top.Result()
	second := top.Result()
	if !reflect.DeepEqual(first, second) {
		t.Errorf("Result() not idempotent: first %v, second %v", first, second)
	}

	// Adding after a Result still behaves correctly.
	top.Add(100)
	got := top.Result()
	want := []int{100, 9, 7}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Result() after further Add = %v, want %v", got, want)
	}
}

func TestTopK_CustomLessFunc(t *testing.T) {
	type job struct {
		name     string
		priority int
	}
	// less(a, b) reports a ranks below b; we keep the highest-priority jobs.
	byPriority := func(a, b job) bool { return a.priority < b.priority }

	top := streaming.NewTopK[job](2, heaps.LessFunc[job](byPriority))
	top.Add(job{"email", 1})
	top.Add(job{"deploy", 9})
	top.Add(job{"lunch", 0})
	top.Add(job{"incident", 7})

	got := top.Result()
	want := []job{{"deploy", 9}, {"incident", 7}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Result() = %v, want %v", got, want)
	}
}
