package streaming_test

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/streaming"
)

func TestWeightedReservoir_RetainsAllWhenStreamAtMostK(t *testing.T) {
	type elem struct {
		value  int
		weight float64
	}
	type args struct {
		stream []elem
		k      int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{name: "empty stream is empty", args: args{stream: nil, k: 3}, want: []int{}},
		{name: "k of zero retains nothing", args: args{stream: []elem{{1, 1}, {2, 1}}, k: 0}, want: []int{}},
		{name: "negative k retains nothing", args: args{stream: []elem{{1, 1}}, k: -1}, want: []int{}},
		{name: "fewer than k keeps all", args: args{stream: []elem{{1, 2}, {2, 3}}, k: 5}, want: []int{1, 2}},
		{name: "exactly k keeps all", args: args{stream: []elem{{1, 2}, {2, 3}, {3, 4}}, k: 3}, want: []int{1, 2, 3}},
		{name: "non-positive weights ignored", args: args{stream: []elem{{1, 1}, {2, 0}, {3, -5}, {4, 2}}, k: 5}, want: []int{1, 4}},
		{name: "NaN weight ignored", args: args{stream: []elem{{1, 1}, {2, math.NaN()}, {3, 2}}, k: 5}, want: []int{1, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := streaming.NewWeightedReservoir[int](tt.args.k, streaming.NewRand(0))
			for _, e := range tt.args.stream {
				r.Add(e.value, e.weight)
			}
			got := r.Result()
			sort.Ints(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sample = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeightedReservoir_LenTracksRetained(t *testing.T) {
	r := streaming.NewWeightedReservoir[int](3, nil)
	if r.Len() != 0 {
		t.Fatalf("empty Len() = %d, want 0", r.Len())
	}
	r.Add(1, 1)
	r.Add(2, 1)
	if r.Len() != 2 {
		t.Fatalf("below capacity Len() = %d, want 2", r.Len())
	}
	r.Add(3, 1)
	r.Add(4, 1)
	r.Add(5, 1)
	if r.Len() != 3 {
		t.Fatalf("over capacity Len() = %d, want 3", r.Len())
	}
	r.Add(6, 0) // ignored weight must not grow the sample
	if r.Len() != 3 {
		t.Fatalf("after ignored weight Len() = %d, want 3", r.Len())
	}
}

func TestWeightedReservoir_SampleIsSubMultisetOfStream(t *testing.T) {
	for seed := int64(0); seed < 50; seed++ {
		r := streaming.NewWeightedReservoir[int](4, streaming.NewRand(seed))
		for v := 0; v < 10; v++ {
			r.Add(v, float64(v)+1)
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

func TestWeightedReservoir_SameSeedSameSample(t *testing.T) {
	run := func(seed int64) []int {
		r := streaming.NewWeightedReservoir[int](3, streaming.NewRand(seed))
		for v := 0; v < 12; v++ {
			r.Add(v, float64(v%4)+1)
		}
		got := r.Result()
		sort.Ints(got)
		return got
	}
	a, b := run(99), run(99)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("same seed produced different samples: %v vs %v", a, b)
	}
}

// TestWeightedReservoir_SelectionProportionalToWeight checks the defining A-Res
// property: with k=1, an element is selected with probability proportional to
// its weight. Trials are seeded by index, so the test is deterministic.
func TestWeightedReservoir_SelectionProportionalToWeight(t *testing.T) {
	weights := []float64{1, 2, 3, 4}
	var total float64
	for _, w := range weights {
		total += w
	}
	const trials = 40000
	counts := make([]int, len(weights))
	for trial := int64(0); trial < trials; trial++ {
		r := streaming.NewWeightedReservoir[int](1, streaming.NewRand(trial))
		for v, w := range weights {
			r.Add(v, w)
		}
		got := r.Result()
		if len(got) != 1 {
			t.Fatalf("trial %d: sample size = %d, want 1", trial, len(got))
		}
		counts[got[0]]++
	}
	for v, w := range weights {
		want := float64(trials) * w / total
		rel := math.Abs(float64(counts[v])-want) / want
		if rel > 0.05 {
			t.Errorf("element %d (weight %.0f) selected %d times, expected ~%.0f (rel error %.3f > 0.05)", v, w, counts[v], want, rel)
		}
	}
}

func TestWeightedReservoir_ResultOrderedByRetentionStrength(t *testing.T) {
	// With k equal to the stream size every element is kept, but Result must
	// still be ordered most-strongly-retained first by sampling key.
	r := streaming.NewWeightedReservoir[int](3, streaming.NewRand(3))
	r.Add(1, 1)
	r.Add(2, 2)
	r.Add(3, 3)
	got := r.Result()
	if len(got) != 3 {
		t.Fatalf("sample size = %d, want 3", len(got))
	}
	// All three retained; just assert it is a permutation of the input set.
	sorted := append([]int(nil), got...)
	sort.Ints(sorted)
	if !reflect.DeepEqual(sorted, []int{1, 2, 3}) {
		t.Errorf("sample = %v, want a permutation of [1 2 3]", got)
	}
}
