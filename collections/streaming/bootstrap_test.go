package streaming_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/streaming"
)

func TestBootstrap_LengthAndSubMultiset(t *testing.T) {
	type args struct {
		input []int
		seed  int64
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
	}{
		{name: "nil input yields empty resample", args: args{input: nil, seed: 1}, wantLen: 0},
		{name: "empty input yields empty resample", args: args{input: []int{}, seed: 1}, wantLen: 0},
		{name: "single element", args: args{input: []int{42}, seed: 1}, wantLen: 1},
		{name: "five elements", args: args{input: []int{1, 2, 3, 4, 5}, seed: 2}, wantLen: 5},
		{name: "duplicates in input", args: args{input: []int{7, 7, 7}, seed: 3}, wantLen: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := streaming.Bootstrap(tt.args.input, streaming.NewRand(tt.args.seed))
			if got == nil {
				t.Fatal("Bootstrap() returned nil, want non-nil slice")
			}
			if len(got) != tt.wantLen {
				t.Errorf("len(Bootstrap()) = %d, want %d", len(got), tt.wantLen)
			}
			// Every resampled element must come from the input set.
			set := make(map[int]bool)
			for _, v := range tt.args.input {
				set[v] = true
			}
			for _, v := range got {
				if !set[v] {
					t.Errorf("Bootstrap() produced %d not present in input", v)
				}
			}
		})
	}
}

func TestBootstrap_SingleElementRepeats(t *testing.T) {
	got := streaming.Bootstrap([]int{99}, streaming.NewRand(5))
	want := []int{99}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Bootstrap() = %v, want %v", got, want)
	}
}

func TestBootstrap_NilRngIsDeterministic(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	first := streaming.Bootstrap(input, nil)
	second := streaming.Bootstrap(input, nil)
	if !reflect.DeepEqual(first, second) {
		t.Errorf("Bootstrap() with nil rng not deterministic: %v vs %v", first, second)
	}
}

func TestBootstrap_DoesNotMutateInput(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	snapshot := make([]int, len(input))
	copy(snapshot, input)
	_ = streaming.Bootstrap(input, streaming.NewRand(1))
	if !reflect.DeepEqual(input, snapshot) {
		t.Errorf("Bootstrap() mutated input: got %v, want %v", input, snapshot)
	}
}

func TestBootstrapN_CountAndShape(t *testing.T) {
	type args struct {
		input []int
		count int
		seed  int64
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantInner int
	}{
		{name: "negative count yields empty", args: args{input: []int{1, 2, 3}, count: -1, seed: 1}, wantCount: 0, wantInner: 0},
		{name: "zero count yields empty", args: args{input: []int{1, 2, 3}, count: 0, seed: 1}, wantCount: 0, wantInner: 0},
		{name: "three resamples of three", args: args{input: []int{1, 2, 3}, count: 3, seed: 1}, wantCount: 3, wantInner: 3},
		{name: "resamples of empty input", args: args{input: nil, count: 2, seed: 1}, wantCount: 2, wantInner: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := streaming.BootstrapN(tt.args.input, tt.args.count, streaming.NewRand(tt.args.seed))
			if got == nil {
				t.Fatal("BootstrapN() returned nil, want non-nil slice")
			}
			if len(got) != tt.wantCount {
				t.Fatalf("len(BootstrapN()) = %d, want %d", len(got), tt.wantCount)
			}
			for i, r := range got {
				if r == nil {
					t.Errorf("BootstrapN()[%d] is nil, want non-nil slice", i)
				}
				if len(r) != tt.wantInner {
					t.Errorf("len(BootstrapN()[%d]) = %d, want %d", i, len(r), tt.wantInner)
				}
			}
		})
	}
}

func TestBootstrapN_ResamplesAreIndependent(t *testing.T) {
	// With a non-trivial input and a fixed seed, two successive resamples in the
	// batch should differ (they consume distinct draws from the shared rng).
	resamples := streaming.BootstrapN([]int{1, 2, 3, 4, 5, 6, 7, 8}, 2, streaming.NewRand(99))
	if reflect.DeepEqual(resamples[0], resamples[1]) {
		t.Errorf("BootstrapN() produced identical successive resamples %v", resamples[0])
	}
}

func TestBootstrapN_NilRngIsDeterministic(t *testing.T) {
	input := []int{1, 2, 3, 4}
	first := streaming.BootstrapN(input, 3, nil)
	second := streaming.BootstrapN(input, 3, nil)
	if !reflect.DeepEqual(first, second) {
		t.Errorf("BootstrapN() with nil rng not deterministic: %v vs %v", first, second)
	}
}
