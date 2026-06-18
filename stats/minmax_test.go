package stats_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestMinMax(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		wantMin int
		wantMax int
		wantOK  bool
	}{
		{name: "spread", input: []int{3, 1, 4, 1, 5}, wantMin: 1, wantMax: 5, wantOK: true},
		{name: "single value", input: []int{9}, wantMin: 9, wantMax: 9, wantOK: true},
		{name: "all equal", input: []int{2, 2, 2}, wantMin: 2, wantMax: 2, wantOK: true},
		{name: "negatives", input: []int{-5, 0, 5}, wantMin: -5, wantMax: 5, wantOK: true},
		{name: "empty input", input: []int{}, wantOK: false},
		{name: "nil input", input: nil, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lo, hi, ok := stats.MinMax(tt.input)
			if ok != tt.wantOK || (ok && (lo != tt.wantMin || hi != tt.wantMax)) {
				t.Fatalf("MinMax(%v) = (%d, %d, %v), want (%d, %d, %v)", tt.input, lo, hi, ok, tt.wantMin, tt.wantMax, tt.wantOK)
			}
		})
	}
}

func TestMinMaxOnStrings(t *testing.T) {
	lo, hi, ok := stats.MinMax([]string{"banana", "apple", "cherry"})
	if !ok || lo != "apple" || hi != "cherry" {
		t.Fatalf(`MinMax(strings) = (%q, %q, %v), want ("apple", "cherry", true)`, lo, hi, ok)
	}
}

func TestArgMin(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		want   int
		wantOK bool
	}{
		{name: "minimum in middle", input: []int{3, 1, 4}, want: 1, wantOK: true},
		{name: "first occurrence on ties", input: []int{2, 1, 1, 0, 0}, want: 3, wantOK: true},
		{name: "single value", input: []int{9}, want: 0, wantOK: true},
		{name: "empty input", input: []int{}, wantOK: false},
		{name: "nil input", input: nil, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.ArgMin(tt.input)
			if ok != tt.wantOK || (ok && got != tt.want) {
				t.Fatalf("ArgMin(%v) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestArgMax(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		want   int
		wantOK bool
	}{
		{name: "maximum in middle", input: []int{1, 5, 4}, want: 1, wantOK: true},
		{name: "first occurrence on ties", input: []int{2, 9, 9, 3}, want: 1, wantOK: true},
		{name: "single value", input: []int{9}, want: 0, wantOK: true},
		{name: "empty input", input: []int{}, wantOK: false},
		{name: "nil input", input: nil, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.ArgMax(tt.input)
			if ok != tt.wantOK || (ok && got != tt.want) {
				t.Fatalf("ArgMax(%v) = (%d, %v), want (%d, %v)", tt.input, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}
