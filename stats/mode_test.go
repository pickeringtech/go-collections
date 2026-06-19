package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestMode(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		want   []int
		wantOK bool
	}{
		{name: "single clear mode", input: []int{1, 2, 2, 3}, want: []int{2}, wantOK: true},
		{name: "ties returned in first-appearance order", input: []int{3, 1, 3, 1, 2}, want: []int{3, 1}, wantOK: true},
		{name: "all unique are all modes", input: []int{4, 5, 6}, want: []int{4, 5, 6}, wantOK: true},
		{name: "single value", input: []int{9}, want: []int{9}, wantOK: true},
		{name: "empty input", input: []int{}, wantOK: false},
		{name: "nil input", input: nil, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.Mode(tt.input)
			if ok != tt.wantOK || (ok && !equalInts(got, tt.want)) {
				t.Fatalf("Mode(%v) = (%v, %v), want (%v, %v)", tt.input, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestModeOnStrings(t *testing.T) {
	got, ok := stats.Mode([]string{"a", "b", "a", "c", "b", "a"})
	if !ok || len(got) != 1 || got[0] != "a" {
		t.Fatalf(`Mode(strings) = (%v, %v), want (["a"], true)`, got, ok)
	}
}

func TestModeRejectsNonFinite(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		got, ok := stats.Mode([]float64{1, 1, bad, 2})
		if ok {
			t.Fatalf("Mode([1 1 %v 2]) = (%v, %v), want (nil, false)", bad, got, ok)
		}
	}
	// float32 inputs are checked too.
	_, ok := stats.Mode([]float32{1, float32(math.Inf(1)), 1})
	if ok {
		t.Fatalf("Mode(float32 with +Inf) should be ok=false")
	}
}

func TestModeFiniteFloatsStillWork(t *testing.T) {
	got, ok := stats.Mode([]float64{1.5, 2.5, 1.5})
	if !ok || len(got) != 1 || got[0] != 1.5 {
		t.Fatalf("Mode([1.5 2.5 1.5]) = (%v, %v), want ([1.5], true)", got, ok)
	}
}
