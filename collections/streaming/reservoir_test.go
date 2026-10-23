package streaming_test

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/streaming"
)

// drainReservoir feeds values into a fresh seeded Reservoir of size k and
// returns its sorted sample, so results are comparable independent of internal
// slot order.
func drainReservoir(values []int, k int, seed int64) []int {
	r := streaming.NewReservoir[int](k, streaming.NewRand(seed))
	for _, v := range values {
		r.Add(v)
	}
	got := r.Result()
	sort.Ints(got)
	return got
}

func TestReservoir_RetainsAllWhenStreamAtMostK(t *testing.T) {
	type args struct {
		values []int
		k      int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{name: "nil stream is empty", args: args{values: nil, k: 3}, want: []int{}},
		{name: "empty stream is empty", args: args{values: []int{}, k: 3}, want: []int{}},
		{name: "k of zero retains nothing", args: args{values: []int{5, 1, 9}, k: 0}, want: []int{}},
		{name: "negative k retains nothing", args: args{values: []int{5, 1, 9}, k: -2}, want: []int{}},
		{name: "fewer elements than k keeps all", args: args{values: []int{3, 1, 2}, k: 5}, want: []int{1, 2, 3}},
		{name: "exactly k keeps all", args: args{values: []int{3, 1, 2}, k: 3}, want: []int{1, 2, 3}},
		{name: "single element", args: args{values: []int{42}, k: 3}, want: []int{42}},
		{name: "duplicates kept below capacity", args: args{values: []int{5, 5, 5}, k: 4}, want: []int{5, 5, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := drainReservoir(tt.args.values, tt.args.k, 0)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sample = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReservoir_LenTracksRetained(t *testing.T) {
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
			r := streaming.NewReservoir[int](tt.args.k, nil)
			for _, v := range tt.args.values {
				r.Add(v)
			}
			if got := r.Len(); got != tt.want {
				t.Errorf("Len() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestReservoir_SampleIsSubMultisetOfStream(t *testing.T) {
	stream := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for seed := int64(0); seed < 50; seed++ {
		r := streaming.NewReservoir[int](4, streaming.NewRand(seed))
		for _, v := range stream {
			r.Add(v)
		}
		got := r.Result()
		if len(got) != 4 {
			t.Fatalf("seed %d: sample size = %d, want 4", seed, len(got))
		}
		seen := make(map[int]bool)
		for _, v := range got {
			if v < 0 || v > 9 {
				t.Fatalf("seed %d: sampled %d not in stream", seed, v)
			}
			if seen[v] {
				t.Fatalf("seed %d: sampled %d twice (without-replacement violated)", seed, v)
			}
			seen[v] = true
		}
	}
}

func TestReservoir_SameSeedSameSample(t *testing.T) {
	stream := []int{10, 20, 30, 40, 50, 60, 70, 80}
	a := drainReservoir(stream, 3, 1234)
	b := drainReservoir(stream, 3, 1234)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("same seed produced different samples: %v vs %v", a, b)
	}
	c := drainReservoir(stream, 3, 5678)
	if reflect.DeepEqual(a, c) {
		t.Errorf("different seeds produced identical samples %v — sampling not seed-driven", a)
	}
}

func TestReservoir_NilRNGIsDeterministicDefault(t *testing.T) {
	stream := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	withNil := streaming.NewReservoir[int](3, nil)
	withSeed0 := streaming.NewReservoir[int](3, streaming.NewRand(0))
	for _, v := range stream {
		withNil.Add(v)
		withSeed0.Add(v)
	}
	gotNil, gotSeed0 := withNil.Result(), withSeed0.Result()
	if !reflect.DeepEqual(gotNil, gotSeed0) {
		t.Errorf("nil rng = %v, NewRand(0) = %v — nil default not equivalent to seed 0", gotNil, gotSeed0)
	}
}

func TestReservoir_ResultIsAnIndependentCopy(t *testing.T) {
	r := streaming.NewReservoir[int](3, streaming.NewRand(7))
	for _, v := range []int{1, 2, 3, 4, 5, 6} {
		r.Add(v)
	}
	first := r.Result()
	for i := range first {
		first[i] = -1
	}
	second := r.Result()
	for _, v := range second {
		if v == -1 {
			t.Fatalf("mutating Result() output corrupted reservoir state: %v", second)
		}
	}
}

// TestReservoir_UniformDistribution checks the defining property of Algorithm R:
// after the whole stream, every element is in the sample with probability k/n.
// Trials are seeded by index, so the test is fully deterministic.
func TestReservoir_UniformDistribution(t *testing.T) {
	const n, k, trials = 10, 3, 20000
	counts := make([]int, n)
	for trial := int64(0); trial < trials; trial++ {
		r := streaming.NewReservoir[int](k, streaming.NewRand(trial))
		for v := 0; v < n; v++ {
			r.Add(v)
		}
		for _, v := range r.Result() {
			counts[v]++
		}
	}
	want := float64(trials) * float64(k) / float64(n) // expected count per element
	for v, got := range counts {
		rel := math.Abs(float64(got)-want) / want
		if rel > 0.05 {
			t.Errorf("element %d sampled %d times, expected ~%.0f (rel error %.3f > 0.05)", v, got, want, rel)
		}
	}
}
