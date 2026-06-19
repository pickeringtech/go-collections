package stats_test

import (
	"math"
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
		if got := bins[len(bins)-1].Max; got != 2 {
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
				if bins, ok := stats.Histogram(tc.input, tc.bins); ok {
					t.Errorf("ok = true (bins %+v), want false", bins)
				}
			})
		}
	})
}
