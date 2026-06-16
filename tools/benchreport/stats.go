package main

import (
	"math"
	"sort"
)

// mannWhitneyP returns the two-sided p-value of the Mann–Whitney U test for the
// null hypothesis that samples a and b are drawn from the same distribution,
// computed with the normal approximation plus tie and continuity corrections.
//
// This is the same rank-based significance test benchstat applies to decide
// whether a delta is real or noise. Reimplementing it (≈ standard library only)
// keeps the trend tooling dependency-free and a pure function of the committed
// samples — the regression check never needs to shell out to benchstat.
//
// It returns 1 (no evidence of a difference) when either sample is empty or the
// combined values are degenerate (zero rank variance), so callers can treat a
// high p-value uniformly as "not significant".
func mannWhitneyP(a, b []float64) float64 {
	n1, n2 := len(a), len(b)
	if n1 == 0 || n2 == 0 {
		return 1
	}

	ranks, tieTerm := rankCombined(a, b)
	// Sum the ranks assigned to the a-sample (the first n1 combined values).
	var r1 float64
	for i := 0; i < n1; i++ {
		r1 += ranks[i]
	}

	u1 := r1 - float64(n1)*(float64(n1)+1)/2
	muU := float64(n1) * float64(n2) / 2
	n := float64(n1 + n2)

	// Variance of U under H0, corrected for ties: the Σ(t³−t) term shrinks the
	// variance when many values are equal (common for ns-scale benchmarks).
	sigma2 := (float64(n1) * float64(n2) / 12) * ((n + 1) - tieTerm/(n*(n-1)))
	if sigma2 <= 0 {
		return 1
	}
	sigma := math.Sqrt(sigma2)

	// Continuity correction: pull |U−μ| toward μ by 0.5 before standardizing.
	z := (math.Abs(u1-muU) - 0.5) / sigma
	if z < 0 {
		z = 0
	}
	// Two-sided p = 2·Φ(−z) = erfc(z/√2).
	return math.Erfc(z / math.Sqrt2)
}

// rankCombined ranks the concatenation a‖b in ascending order, assigning the
// average rank to tied values, and returns each element's rank (a's elements
// first, then b's) together with the tie-correction term Σ(tᵢ³ − tᵢ) summed over
// groups of tᵢ tied values.
func rankCombined(a, b []float64) (ranks []float64, tieTerm float64) {
	n := len(a) + len(b)
	type item struct {
		v   float64
		idx int
	}
	items := make([]item, 0, n)
	for i, v := range a {
		items = append(items, item{v, i})
	}
	for i, v := range b {
		items = append(items, item{v, len(a) + i})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].v < items[j].v })

	ranks = make([]float64, n)
	for i := 0; i < n; {
		j := i
		for j < n && items[j].v == items[i].v {
			j++
		}
		// items[i:j] are tied and occupy 1-based ranks i+1 … j; their shared
		// rank is the average of those, ((i+1)+j)/2.
		avg := float64(i+1+j) / 2
		for k := i; k < j; k++ {
			ranks[items[k].idx] = avg
		}
		t := float64(j - i)
		tieTerm += t*t*t - t
		i = j
	}
	return ranks, tieTerm
}

// median returns the median of xs (mean of the two middle values for an even
// count). It does not mutate xs. Returns 0 for an empty slice.
func median(xs []float64) float64 {
	n := len(xs)
	if n == 0 {
		return 0
	}
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	if n%2 == 1 {
		return s[n/2]
	}
	return (s[n/2-1] + s[n/2]) / 2
}
