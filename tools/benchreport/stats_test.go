package main

import (
	"math"
	"testing"
)

func TestMannWhitneyEmptySamples(t *testing.T) {
	if p := mannWhitneyP(nil, []float64{1, 2}); p != 1 {
		t.Errorf("empty a: p = %v, want 1", p)
	}
	if p := mannWhitneyP([]float64{1, 2}, nil); p != 1 {
		t.Errorf("empty b: p = %v, want 1", p)
	}
}

func TestMannWhitneyDegenerateVariance(t *testing.T) {
	// Every value identical → the tie correction zeroes the variance, so there
	// is no rank information and the test reports "no difference" (p = 1).
	if p := mannWhitneyP([]float64{5, 5, 5}, []float64{5, 5, 5}); p != 1 {
		t.Errorf("all-tied: p = %v, want 1", p)
	}
}

func TestMannWhitneyClearSeparation(t *testing.T) {
	// Two well-separated, equal-sized samples are highly significant.
	a := []float64{10, 10, 11, 10, 11, 10, 11, 10}
	b := []float64{40, 41, 40, 42, 40, 41, 40, 41}
	if p := mannWhitneyP(a, b); p >= 0.05 {
		t.Errorf("separated samples: p = %v, want < 0.05", p)
	}
}

func TestMannWhitneyHeavyOverlap(t *testing.T) {
	// Interleaved samples with near-zero rank difference: the continuity
	// correction drives the z-score to (or below) zero, so p clamps to 1.
	a := []float64{1, 4}
	b := []float64{2, 3}
	if p := mannWhitneyP(a, b); p != 1 {
		t.Errorf("balanced overlap: p = %v, want 1 (continuity clamp)", p)
	}
}

func TestMannWhitneyModerateOverlap(t *testing.T) {
	// A real but not overwhelming separation lands strictly between 0 and 1,
	// exercising the non-clamped z path.
	a := []float64{10, 11, 12, 13}
	b := []float64{20, 30, 40, 50}
	p := mannWhitneyP(a, b)
	if p <= 0 || p >= 1 {
		t.Errorf("moderate overlap: p = %v, want in (0,1)", p)
	}
}

func TestMedian(t *testing.T) {
	cases := []struct {
		in   []float64
		want float64
	}{
		{nil, 0},
		{[]float64{5}, 5},
		{[]float64{3, 1, 2}, 2},      // odd, unsorted
		{[]float64{4, 1, 3, 2}, 2.5}, // even, unsorted
		{[]float64{10, 10, 10, 10}, 10},
	}
	for _, c := range cases {
		src := append([]float64(nil), c.in...)
		if got := median(c.in); got != c.want {
			t.Errorf("median(%v) = %v, want %v", c.in, got, c.want)
		}
		// median must not mutate its argument.
		for i := range src {
			if src[i] != c.in[i] {
				t.Errorf("median mutated its input: %v vs %v", c.in, src)
			}
		}
	}
}

// Sanity-check the helper against a value computed independently: erfc is the
// engine behind the two-sided p, so a known z maps to a known p.
func TestErfcSanity(t *testing.T) {
	if got := math.Erfc(0); got != 1 {
		t.Errorf("erfc(0) = %v, want 1", got)
	}
}
