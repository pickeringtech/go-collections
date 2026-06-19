package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestPercentileOfScore(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}

	t.Run("ranks a score within the sample", func(t *testing.T) {
		cases := map[string]struct {
			score float64
			want  float64
		}{
			"below all":      {0, 0},
			"at minimum":     {1, 20},
			"middle":         {3, 60},
			"between values": {3.5, 60},
			"at maximum":     {5, 100},
			"above all":      {6, 100},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				got, ok := stats.PercentileOfScore(data, tc.score)
				if !ok || !approxEqual(got, tc.want) {
					t.Fatalf("PercentileOfScore(_, %v) = %v, %v; want %v, true", tc.score, got, ok, tc.want)
				}
			})
		}
	})

	t.Run("inverts Percentile at sample points", func(t *testing.T) {
		// The p-th percentile, ranked back, is >= p for every sample point.
		for p := 0.0; p <= 100; p += 25 {
			v, ok := stats.Percentile(data, p)
			if !ok {
				t.Fatalf("Percentile(_, %v) not ok", p)
			}
			rank, ok := stats.PercentileOfScore(data, v)
			if !ok || rank < p {
				t.Errorf("PercentileOfScore(_, %v) = %v, want >= %v", v, rank, p)
			}
		}
	})

	t.Run("unsorted input is handled", func(t *testing.T) {
		got, ok := stats.PercentileOfScore([]int{5, 1, 4, 2, 3}, 3)
		if !ok || !approxEqual(got, 60) {
			t.Fatalf("PercentileOfScore = %v, %v; want 60, true", got, ok)
		}
	})

	t.Run("rejects undefined inputs", func(t *testing.T) {
		if _, ok := stats.PercentileOfScore([]float64{}, 1); ok {
			t.Errorf("empty input reported ok")
		}
		if _, ok := stats.PercentileOfScore(data, math.NaN()); ok {
			t.Errorf("NaN score reported ok")
		}
		if _, ok := stats.PercentileOfScore(data, math.Inf(1)); ok {
			t.Errorf("Inf score reported ok")
		}
		if _, ok := stats.PercentileOfScore([]float64{1, math.NaN(), 3}, 2); ok {
			t.Errorf("NaN in input reported ok")
		}
	})
}
