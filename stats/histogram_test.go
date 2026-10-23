package stats_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestHistogram(t *testing.T) {
	t.Run("buckets values into equal-width bins", func(t *testing.T) {
		// span 4 over 2 bins => width 2: [1,3) and [3,5].
		bins, ok := stats.Histogram([]float64{1, 2, 3, 4, 5}, 2)
		if !ok {
			t.Fatalf("ok = false, want true")
		}
		if len(bins) != 2 {
			t.Fatalf("got %d bins, want 2", len(bins))
		}
		if bins[0].Min != 1 || bins[0].Max != 3 || bins[0].Count != 2 {
			t.Errorf("bin0 = %+v, want {1 3 2}", bins[0])
		}
		// 3 and 4 fall in the second bin; 5 (the max) is folded into it too.
		if bins[1].Min != 3 || bins[1].Max != 5 || bins[1].Count != 3 {
			t.Errorf("bin1 = %+v, want {3 5 3}", bins[1])
		}
	})

	t.Run("counts always sum to len(input)", func(t *testing.T) {
		input := []int{5, 1, 9, 3, 7, 2, 8, 4, 6, 10, 5, 5}
		bins, ok := stats.Histogram(input, 4)
		if !ok {
			t.Fatalf("ok = false, want true")
		}
		total := 0
		for _, b := range bins {
			total += b.Count
		}
		if total != len(input) {
			t.Errorf("counts sum to %d, want %d", total, len(input))
		}
	})

	t.Run("last bin upper bound is the exact maximum", func(t *testing.T) {
		bins, _ := stats.Histogram([]float64{0, 1, 2}, 3)
		got := bins[len(bins)-1].Max
		if got != 2 {
			t.Errorf("last bin Max = %v, want exactly 2", got)
		}
	})

	t.Run("rejects undefined inputs", func(t *testing.T) {
		cases := map[string]struct {
			input []float64
			bins  int
		}{
			"empty":              {nil, 4},
			"zero bins":          {[]float64{1, 2, 3}, 0},
			"negative bins":      {[]float64{1, 2, 3}, -1},
			"constant (no span)": {[]float64{7, 7, 7}, 3},
			"NaN":                {[]float64{1, math.NaN(), 3}, 3},
			"Inf":                {[]float64{1, math.Inf(1), 3}, 3},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				bins, ok := stats.Histogram(tc.input, tc.bins)
				if ok {
					t.Errorf("ok = true (bins %+v), want false", bins)
				}
			})
		}
	})
}

// recount tallies input against the bins' reported edges — [Min, Max) for every
// bin but the last, which is [Min, Max] — so a test can check the counts agree
// with the edges a caller sees.
func recount(input []float64, bins []stats.Bin) []int {
	counts := make([]int, len(bins))
	for _, x := range input {
		for i, b := range bins {
			last := i == len(bins)-1
			if x >= b.Min && (x < b.Max || (last && x <= b.Max)) {
				counts[i]++
				break
			}
		}
	}
	return counts
}

// TestHistogramCountsMatchReportedEdges pins the counts to the reported
// half-open edges: computing a bin as ⌊(x−min)/width⌋ disagrees with the edge
// min+i·width when the two round differently, so a value could be counted in a
// bin whose reported interval excludes it.
func TestHistogramCountsMatchReportedEdges(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		bins  int
	}{
		{name: "value on a rounded inner edge", input: []float64{1.6, 3, 4.4}, bins: 2},
		{name: "decimal steps", input: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7}, bins: 6},
		{name: "many bins over decimal data", input: []float64{0.3, 0.6, 0.9, 1.2, 1.5, 1.8, 2.1}, bins: 7},
		{name: "value just below a rounded-up inner edge", input: []float64{0, 3.9, 5.2}, bins: 4},
		{name: "another value below a rounded-up edge", input: []float64{0, 3.1, 6.2}, bins: 6},
		{name: "integers", input: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, bins: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bins, ok := stats.Histogram(tt.input, tt.bins)
			if !ok {
				t.Fatalf("ok = false, want true")
			}
			want := recount(tt.input, bins)
			for i, b := range bins {
				if b.Count != want[i] {
					t.Errorf("bin %d = %+v, but %d values lie in its reported edges", i, b, want[i])
				}
			}
		})
	}
}

// TestHistogramOverflowingSpan checks a finite range whose width max−min
// overflows float64 is still bucketed (rather than panicking on an Inf width).
func TestHistogramOverflowingSpan(t *testing.T) {
	maxF := math.MaxFloat64
	tests := []struct {
		name  string
		input []float64
		bins  int
		want  []stats.Bin
	}{
		{name: "two bins", input: []float64{-maxF, maxF}, bins: 2, want: []stats.Bin{{Min: -maxF, Max: 0, Count: 1}, {Min: 0, Max: maxF, Count: 1}}},
		{name: "one bin", input: []float64{-maxF, 0, maxF}, bins: 1, want: []stats.Bin{{Min: -maxF, Max: maxF, Count: 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.Histogram(tt.input, tt.bins)
			if !ok || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Histogram() = %+v, %v; want %+v, true", got, ok, tt.want)
			}
		})
	}

	t.Run("many bins have finite ascending edges that match the counts", func(t *testing.T) {
		input := []float64{-maxF, -1e308, -1, 1, 1e308, maxF}
		bins, ok := stats.Histogram(input, 7)
		if !ok {
			t.Fatalf("ok = false, want true")
		}
		want := recount(input, bins)
		for i, b := range bins {
			if math.IsInf(b.Min, 0) || math.IsInf(b.Max, 0) || b.Min > b.Max {
				t.Errorf("bin %d = %+v, want finite ascending edges", i, b)
			}
			if b.Count != want[i] {
				t.Errorf("bin %d = %+v, but %d values lie in its reported edges", i, b, want[i])
			}
		}
	})
}

// FuzzHistogram asserts that for any finite input Histogram never panics, its
// counts sum to len(input), its edges ascend, and every count agrees with the
// bin's reported edges.
func FuzzHistogram(f *testing.F) {
	f.Add(1.6, 3.0, 4.4, uint8(2))
	f.Add(0.0, 3.9, 5.2, uint8(4))
	f.Add(-math.MaxFloat64, 0.0, math.MaxFloat64, uint8(3))
	f.Fuzz(func(t *testing.T, a, b, c float64, rawBins uint8) {
		input := []float64{a, b, c}
		bins, ok := stats.Histogram(input, int(rawBins))
		if !ok {
			return
		}
		want := recount(input, bins)
		total := 0
		for i, bin := range bins {
			total += bin.Count
			if bin.Min > bin.Max || bin.Count != want[i] {
				t.Fatalf("bin %d = %+v of %+v disagrees with its edges (recount %d)", i, bin, bins, want[i])
			}
		}
		if total != len(input) {
			t.Fatalf("counts sum to %d, want %d", total, len(input))
		}
	})
}
