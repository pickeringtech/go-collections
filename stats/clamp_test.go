package stats_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestClamp(t *testing.T) {
	tests := []struct {
		name          string
		value, lo, hi int
		want          int
	}{
		{name: "within range unchanged", value: 3, lo: 0, hi: 5, want: 3},
		{name: "below clamps to lo", value: -2, lo: 0, hi: 5, want: 0},
		{name: "above clamps to hi", value: 7, lo: 0, hi: 5, want: 5},
		{name: "at lower bound", value: 0, lo: 0, hi: 5, want: 0},
		{name: "at upper bound", value: 5, lo: 0, hi: 5, want: 5},
		{name: "degenerate interval", value: 9, lo: 4, hi: 4, want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stats.Clamp(tt.value, tt.lo, tt.hi)
			if got != tt.want {
				t.Fatalf("Clamp(%d, %d, %d) = %d, want %d", tt.value, tt.lo, tt.hi, got, tt.want)
			}
		})
	}
}

func TestClampOnStrings(t *testing.T) {
	got := stats.Clamp("m", "a", "f")
	if got != "f" {
		t.Fatalf(`Clamp("m","a","f") = %q, want "f"`, got)
	}
}

func TestClampAll(t *testing.T) {
	input := []int{-3, 0, 2, 9}
	got := stats.ClampAll(input, 0, 5)
	want := []int{0, 0, 2, 5}
	if !equalInts(got, want) {
		t.Fatalf("ClampAll(%v, 0, 5) = %v, want %v", input, got, want)
	}
	// input must not be mutated.
	for i, v := range []int{-3, 0, 2, 9} {
		if input[i] != v {
			t.Fatalf("ClampAll mutated input to %v", input)
		}
	}
}

func TestClampAllEmpty(t *testing.T) {
	got := stats.ClampAll([]int{}, 0, 5)
	if len(got) != 0 {
		t.Fatalf("ClampAll(empty) = %v, want empty", got)
	}
	got = stats.ClampAll[int](nil, 0, 5)
	if len(got) != 0 {
		t.Fatalf("ClampAll(nil) = %v, want empty", got)
	}
}
